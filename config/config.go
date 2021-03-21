package config

import (
	"os"
	"time"
)

type Config struct {
	Port       string
	Server     ServerConfig
	Middleware MiddlewareConfig
}

func New() *Config {
	return &Config{
		Port:       os.Getenv("SERVER_PORT"),
		Server:     NewServerConfig(),
		Middleware: NewMiddlewareConfig(),
	}
}

type ServerConfig struct {
	Port              string
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

func NewServerConfig() ServerConfig {
	return ServerConfig{
		Port:              os.Getenv("SERVER_PORT"),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
		ShutdownTimeout:   10 * time.Second,
	}
}

type MiddlewareConfig struct {
	Timeout          time.Duration
	RequestSizeLimit int64
}

func NewMiddlewareConfig() MiddlewareConfig {
	return MiddlewareConfig{
		Timeout:          10 * time.Second,
		RequestSizeLimit: 1e+7, // 10MB
	}
}
