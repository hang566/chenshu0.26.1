package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Path string
}

type JWTConfig struct {
	Secret string
	Expiry int // 过期时间（小时）
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnvInt("SERVER_PORT", 8080),
		},
		Database: DatabaseConfig{
			Path: getEnvStr("DB_PATH", "./data/yunshu.db"),
		},
		JWT: JWTConfig{
			Secret: getEnvStr("JWT_SECRET", "yunshu-secret-key-change-in-production"),
			Expiry: getEnvInt("JWT_EXPIRY", 24),
		},
	}
}

func getEnvStr(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}
