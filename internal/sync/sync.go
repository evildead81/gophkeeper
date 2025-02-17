package sync

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/evildead81/gophkeeper/api/password"
	"github.com/evildead81/gophkeeper/pkg/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// StartSync подключается к серверу и слушает изменения паролей
func StartSync(token string) {
	cfg := config.LoadConfig()

	conn, err := grpc.Dial(cfg.ClientAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения к серверу: %v", err)
	}
	defer conn.Close()

	client := pb.NewPasswordServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	stream, err := client.SyncPasswords(ctx, &pb.SyncRequest{Token: token})
	if err != nil {
		log.Fatalf("Ошибка подписки: %v", err)
	}

	fmt.Println("Синхронизация запущена. Ожидание обновлений...")

	for {
		update, err := stream.Recv()
		if err != nil {
			log.Printf("Потеряно соединение, попытка переподключения...")
			time.Sleep(5 * time.Second)
			go StartSync(token)
			return
		}

		fmt.Printf("Изменение: %s | Сайт: %s | Логин: %s\n", update.Action, update.Site, update.Login)
	}
}
