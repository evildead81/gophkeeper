package main

import (
	"context"
	"log"
	"time"

	pb "github.com/evildead81/gophkeeper/api/auth"
	"github.com/evildead81/gophkeeper/pkg/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.LoadConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, cfg.ClientAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	resp, err := client.Login(ctx, &pb.LoginRequest{
		Username: "testuser",
		Password: "testpassword",
	})
	if err != nil {
		log.Fatalf("Ошибка входа: %v", err)
	}

	log.Printf("Токен получен: %s", resp.Token)
}
