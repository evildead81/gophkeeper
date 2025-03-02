package config

import (
	"os"
)

// Config структура для хранения настроек сервера
type Config struct {
	ServerAddr  string
	DatabaseDSN string
	EnableTLS   bool
	TLSCertFile string
	TLSKeyFile  string
	AESKey      string
}

// LoadConfig загружает конфигурацию из .env или переменных окружения
func LoadConfig() *Config {
	return &Config{
		ServerAddr:  getEnv("SERVER_ADDR", ":50051"),
		DatabaseDSN: getEnv("DATABASE_DSN", "postgres://user:password@localhost:5432/gophkeeper?sslmode=disable"),
		EnableTLS:   getEnvAsBool("ENABLE_TLS", false),
		TLSCertFile: getEnv("TLS_CERT_FILE", ""),
		TLSKeyFile:  getEnv("TLS_KEY_FILE", ""),
		AESKey:      getEnv("AES_KEY", "abcdefghijklmnopqrstuvwxyz123456"),
	}
}

// getEnv возвращает строковое значение из окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsBool получает переменную окружения как bool
func getEnvAsBool(name string, defaultValue bool) bool {
	valStr := getEnv(name, "")
	if valStr == "true" || valStr == "1" {
		return true
	}
	return defaultValue
}
