package config

import (
	"flag"
	"os"
)

// Config содержит настройки приложения
type Config struct {
	DatabaseDSN string
	ServerAddr  string
	ClientAddr  string
}

// LoadConfig загружает конфигурацию из флагов или переменных окружения
func LoadConfig() Config {
	dsn := flag.String("db_dsn", getEnv("DATABASE_DSN", "postgres://postgres:password@localhost:5432/gophkeeper?sslmode=disable"), "Database DSN")
	serverAddr := flag.String("server_addr", getEnv("SERVER_ADDR", ":50051"), "gRPC server address")
	clientAddr := flag.String("client_addr", getEnv("CLIENT_ADDR", "localhost:50051"), "gRPC client address")

	flag.Parse()

	return Config{
		DatabaseDSN: *dsn,
		ServerAddr:  *serverAddr,
		ClientAddr:  *clientAddr,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
