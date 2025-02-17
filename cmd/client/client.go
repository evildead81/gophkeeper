package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/evildead81/gophkeeper/api/auth"
	"github.com/evildead81/gophkeeper/internal/sync"
	"github.com/evildead81/gophkeeper/pkg/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.LoadConfig()

	token, err := os.ReadFile(".token")
	if err != nil {
		fmt.Println("Токен не найден. Нужно войти в систему.")

		conn, err := grpc.Dial(cfg.ClientAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Ошибка подключения: %v", err)
		}
		defer conn.Close()

		client := pb.NewAuthServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		fmt.Print("Введите логин: ")
		var username string
		fmt.Scanln(&username)

		fmt.Print("Введите пароль: ")
		var password string
		fmt.Scanln(&password)

		resp, err := client.Login(ctx, &pb.LoginRequest{
			Username: username,
			Password: password,
		})
		if err != nil {
			log.Fatalf("Ошибка входа: %v", err)
		}

		err = os.WriteFile(".token", []byte(resp.Token), 0600)
		if err != nil {
			log.Fatalf("Ошибка сохранения токена: %v", err)
		}

		fmt.Println("Вход выполнен! Токен сохранён.")
		token = []byte(resp.Token)
	}

	fmt.Println("Запускаем синхронизацию данных...")
	sync.StartSync(string(token))
}
