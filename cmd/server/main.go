package main

import (
	"log"
	"net"

	"github.com/evildead81/gophkeeper/internal/auth"
	"github.com/evildead81/gophkeeper/internal/password"
	"github.com/evildead81/gophkeeper/pkg/config"
	"github.com/evildead81/gophkeeper/pkg/db"

	pbAuth "github.com/evildead81/gophkeeper/api/auth"
	pbPassword "github.com/evildead81/gophkeeper/api/password"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	conn, err := db.Connect(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer conn.Close()

	grpcServer := grpc.NewServer()

	authService := auth.NewAuthService(conn)
	pbAuth.RegisterAuthServiceServer(grpcServer, authService)

	passwordService := password.NewPasswordService(conn)
	pbPassword.RegisterPasswordServiceServer(grpcServer, passwordService)

	listener, err := net.Listen("tcp", cfg.ServerAddr)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

	log.Printf("gRPC сервер запущен на %s", cfg.ServerAddr)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка работы gRPC сервера: %v", err)
	}
}
