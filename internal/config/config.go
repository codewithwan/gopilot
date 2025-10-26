package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Log      LogConfig
	Metrics  MetricsConfig
	Tracing  TracingConfig
}

type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type LogConfig struct {
	Level  string
	Format string
}

type MetricsConfig struct {
	Enabled bool
	Port    string
}

type TracingConfig struct {
	Enabled     bool
	ServiceName string
	Endpoint    string
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.readTimeout", "15s")
	viper.SetDefault("server.writeTimeout", "15s")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.dbname", "gopilot")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("jwt.secret", "your-secret-key-change-this")
	viper.SetDefault("jwt.expiration", "24h")
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("metrics.enabled", true)
	viper.SetDefault("metrics.port", "9090")
	viper.SetDefault("tracing.enabled", false)
	viper.SetDefault("tracing.serviceName", "gopilot")
	viper.SetDefault("tracing.endpoint", "")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var cfg Config
	cfg.Server.Port = viper.GetString("server.port")
	cfg.Server.Host = viper.GetString("server.host")
	cfg.Server.ReadTimeout = viper.GetDuration("server.readTimeout")
	cfg.Server.WriteTimeout = viper.GetDuration("server.writeTimeout")
	cfg.Database.Host = viper.GetString("database.host")
	cfg.Database.Port = viper.GetString("database.port")
	cfg.Database.User = viper.GetString("database.user")
	cfg.Database.Password = viper.GetString("database.password")
	cfg.Database.DBName = viper.GetString("database.dbname")
	cfg.Database.SSLMode = viper.GetString("database.sslmode")
	cfg.JWT.Secret = viper.GetString("jwt.secret")
	cfg.JWT.Expiration = viper.GetDuration("jwt.expiration")
	cfg.Log.Level = viper.GetString("log.level")
	cfg.Log.Format = viper.GetString("log.format")
	cfg.Metrics.Enabled = viper.GetBool("metrics.enabled")
	cfg.Metrics.Port = viper.GetString("metrics.port")
	cfg.Tracing.Enabled = viper.GetBool("tracing.enabled")
	cfg.Tracing.ServiceName = viper.GetString("tracing.serviceName")
	cfg.Tracing.Endpoint = viper.GetString("tracing.endpoint")

	return &cfg, nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}
