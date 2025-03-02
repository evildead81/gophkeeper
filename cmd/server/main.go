package main

import (
	"log"
	"net"

	"github.com/evildead81/gophkeeper/internal/repository"
	"github.com/evildead81/gophkeeper/internal/service"
	"github.com/evildead81/gophkeeper/pkg/config"
	"github.com/evildead81/gophkeeper/pkg/db"
	"github.com/evildead81/gophkeeper/pkg/session"

	pbAuth "github.com/evildead81/gophkeeper/api/auth"
	pbFile "github.com/evildead81/gophkeeper/api/file"
	pbSecureData "github.com/evildead81/gophkeeper/api/secure_data"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	cfg := config.LoadConfig()

	conn, err := db.Connect(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer conn.Close()

	sessionManager := session.NewSessionManager()

	userRepo := repository.NewUserRepository(conn)
	secureDataRepo := repository.NewSecureDataRepository(conn)
	fileRepo := repository.NewFileRepository(conn)

	authService := service.NewAuthService(userRepo, sessionManager)
	secureDataService := service.NewSecureDataService(secureDataRepo, userRepo, sessionManager)
	fileService := service.NewFileService(fileRepo, userRepo, sessionManager)

	var opts []grpc.ServerOption
	if cfg.EnableTLS {
		creds, err := credentials.NewServerTLSFromFile(cfg.TLSCertFile, cfg.TLSKeyFile)
		if err != nil {
			log.Fatalf("Ошибка загрузки TLS-сертификата: %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(opts...)

	pbAuth.RegisterAuthServiceServer(grpcServer, authService)
	pbSecureData.RegisterSecureDataServiceServer(grpcServer, secureDataService)
	pbFile.RegisterFileServiceServer(grpcServer, fileService)

	listener, err := net.Listen("tcp", cfg.ServerAddr)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

	log.Printf("gRPC сервер запущен на %s", cfg.ServerAddr)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка работы gRPC сервера: %v", err)
	}
}
