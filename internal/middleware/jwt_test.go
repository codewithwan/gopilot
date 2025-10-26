package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestGenerateToken(t *testing.T) {
	secret := "test-secret"
	jwtMiddleware := NewJWTMiddleware(secret)

	token, err := jwtMiddleware.GenerateToken(1, "testuser", 24*time.Hour)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Generated token is empty")
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	secret := "test-secret"
	jwtMiddleware := NewJWTMiddleware(secret)

	token, _ := jwtMiddleware.GenerateToken(1, "testuser", 24*time.Hour)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(jwtMiddleware.AuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		userID, _ := GetUserID(c)
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	secret := "test-secret"
	jwtMiddleware := NewJWTMiddleware(secret)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(jwtMiddleware.AuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	secret := "test-secret"
	jwtMiddleware := NewJWTMiddleware(secret)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(jwtMiddleware.AuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestGetUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("user_id", int64(123))

	userID, err := GetUserID(c)
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	if userID != 123 {
		t.Errorf("Expected user ID 123, got %d", userID)
	}
}

func TestGetUserID_Missing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	_, err := GetUserID(c)
	if err == nil {
		t.Error("Expected error for missing user_id, got nil")
	}
}
