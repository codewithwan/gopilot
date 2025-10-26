package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codewithwan/gopilot/internal/config"
	"github.com/codewithwan/gopilot/internal/handler"
	"github.com/codewithwan/gopilot/internal/middleware"
	"github.com/codewithwan/gopilot/internal/repository"
	"github.com/codewithwan/gopilot/internal/repository/db"
	"github.com/codewithwan/gopilot/internal/service"
	"github.com/codewithwan/gopilot/pkg/logger"
	"github.com/codewithwan/gopilot/pkg/metrics"
	"github.com/codewithwan/gopilot/pkg/tracing"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"

	_ "github.com/codewithwan/gopilot/docs"
)

// @title GoPilot API
// @version 1.0
// @description Production-ready REST API for managing todos
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://github.com/codewithwan/gopilot
// @contact.email support@gopilot.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger
	log, err := logger.New(cfg.Log.Level, cfg.Log.Format)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer log.Close()

	log.Info("Starting GoPilot service",
		zap.String("port", cfg.Server.Port),
		zap.String("log_level", cfg.Log.Level),
	)

	// Connect to database
	dbpool, err := pgxpool.New(context.Background(), cfg.Database.DSN())
	if err != nil {
		log.Error("Failed to connect to database", zap.Error(err))
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer dbpool.Close()

	if err := dbpool.Ping(context.Background()); err != nil {
		log.Error("Failed to ping database", zap.Error(err))
		return fmt.Errorf("failed to ping database: %w", err)
	}
	log.Info("Connected to database")

	// Initialize tracing if enabled
	if cfg.Tracing.Enabled && cfg.Tracing.Endpoint != "" {
		tp, tracerErr := tracing.InitTracer(cfg.Tracing.ServiceName, cfg.Tracing.Endpoint)
		if tracerErr != nil {
			log.Error("Failed to initialize tracer", zap.Error(tracerErr))
		} else {
			defer func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if shutdownErr := tp.Shutdown(ctx); shutdownErr != nil {
					log.Error("Failed to shutdown tracer", zap.Error(shutdownErr))
				}
			}()
			log.Info("Tracing enabled", zap.String("endpoint", cfg.Tracing.Endpoint))
		}
	}

	// Initialize repositories
	queries := db.New(dbpool)
	userRepo := repository.NewUserRepository(queries)
	todoRepo := repository.NewTodoRepository(queries)

	// Initialize JWT middleware
	jwtMiddleware := middleware.NewJWTMiddleware(cfg.JWT.Secret)

	// Initialize services
	authService := service.NewAuthService(userRepo, jwtMiddleware, cfg.JWT.Expiration, log.Logger)
	todoService := service.NewTodoService(todoRepo, log.Logger)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	todoHandler := handler.NewTodoHandler(todoService)

	// Set Gin mode
	if cfg.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.New()
	router.Use(gin.Recovery())

	// Add OpenTelemetry middleware if tracing is enabled
	if cfg.Tracing.Enabled {
		router.Use(otelgin.Middleware(cfg.Tracing.ServiceName))
	}

	// Add Prometheus metrics middleware if enabled
	if cfg.Metrics.Enabled {
		router.Use(metrics.PrometheusMiddleware())
		router.GET("/metrics", metrics.Handler())
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Todo routes (protected)
		todos := v1.Group("/todos")
		todos.Use(jwtMiddleware.AuthMiddleware())
		{
			todos.POST("", todoHandler.CreateTodo)
			todos.GET("", todoHandler.ListTodos)
			todos.GET("/:id", todoHandler.GetTodo)
			todos.PUT("/:id", todoHandler.UpdateTodo)
			todos.DELETE("/:id", todoHandler.DeleteTodo)
		}
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in a goroutine
	serverErr := make(chan error, 1)
	go func() {
		log.Info("Server started", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", zap.Error(err))
			serverErr <- err
		}
	}()

	// Wait for interrupt signal or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	var serverError error
	select {
	case <-quit:
		log.Info("Shutting down server...")
	case serverError = <-serverErr:
		log.Error("Server error, shutting down", zap.Error(serverError))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", zap.Error(err))
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Info("Server exited")
	if serverError != nil {
		return fmt.Errorf("server error: %w", serverError)
	}
	return nil
}
