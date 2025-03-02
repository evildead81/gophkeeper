package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	pbAuth "github.com/evildead81/gophkeeper/api/auth"
	pbFile "github.com/evildead81/gophkeeper/api/file"
	pbSecureData "github.com/evildead81/gophkeeper/api/secure_data"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	serverAddr  = flag.String("server", "localhost:50051", "Адрес gRPC сервера")
	enableTLS   = flag.Bool("tls", false, "Использовать TLS")
	tlsCertFile = flag.String("cert", "certs/server.crt", "Файл TLS-сертификата")
	command     = flag.String("cmd", "", "Команда (register, login, add-data, get-data, upload-file, download-file)")
	username    = flag.String("username", "", "Имя пользователя (для регистрации/логина)")
	password    = flag.String("password", "", "Пароль (для регистрации/логина)")
	token       = flag.String("token", "", "JWT токен")
	dataType    = flag.String("type", "", "Тип данных (password, text, binary, card)")
	data        = flag.String("data", "", "Данные")
	metaInfo    = flag.String("meta", "", "Метаинформация")
	filePath    = flag.String("file", "", "Путь к файлу")
	fileID      = flag.String("fileid", "", "ID файла для скачивания")
)

func main() {
	flag.Parse()

	var opts []grpc.DialOption
	if *enableTLS {
		creds, err := credentials.NewClientTLSFromFile(*tlsCertFile, "")
		if err != nil {
			log.Fatalf("Ошибка загрузки TLS-сертификата: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("Ошибка подключения к серверу: %v", err)
	}
	defer conn.Close()

	switch *command {
	case "register":
		register(conn)
	case "login":
		login(conn)
	case "add-data":
		addSecureData(conn)
	case "get-data":
		getSecureData(conn)
	case "upload-file":
		uploadFile(conn)
	case "download-file":
		downloadFile(conn)
	default:
		fmt.Println("Неизвестная команда. Доступные команды: register, login, add-data, get-data, upload-file, download-file")
	}
}

// Регистрация нового пользователя
func register(conn *grpc.ClientConn) {
	client := pbAuth.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.Register(ctx, &pbAuth.RegisterRequest{
		Username: *username,
		Password: *password,
	})
	if err != nil {
		log.Fatalf("Ошибка регистрации: %v", err)
	}
	fmt.Println("Успешная регистрация:", resp.Message)
}

// Авторизация пользователя
func login(conn *grpc.ClientConn) {
	client := pbAuth.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.Login(ctx, &pbAuth.LoginRequest{
		Username: *username,
		Password: *password,
	})
	if err != nil {
		log.Fatalf("Ошибка авторизации: %v", err)
	}
	fmt.Println("JWT токен:", resp.Token)
}

// Добавление защищённых данных
func addSecureData(conn *grpc.ClientConn) {
	client := pbSecureData.NewSecureDataServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.AddSecureData(ctx, &pbSecureData.AddSecureDataRequest{
		Token:    *token,
		DataType: *dataType,
		Data:     *data,
		MetaInfo: *metaInfo,
	})
	if err != nil {
		log.Fatalf("Ошибка добавления данных: %v", err)
	}
	fmt.Println("Данные успешно добавлены:", resp.Message)
}

// Получение защищённых данных
func getSecureData(conn *grpc.ClientConn) {
	client := pbSecureData.NewSecureDataServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.GetSecureData(ctx, &pbSecureData.GetSecureDataRequest{
		Token:    *token,
		DataType: *dataType,
	})
	if err != nil {
		log.Fatalf("Ошибка получения данных: %v", err)
	}
	fmt.Println("Полученные данные:")
	for _, item := range resp.Items {
		fmt.Printf("ID: %s, Data: %s, Meta: %s\n", item.Id, item.Data, item.MetaInfo)
	}
}

// Загрузка файла
func uploadFile(conn *grpc.ClientConn) {
	client := pbFile.NewFileServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	stream, err := client.UploadFile(ctx)
	if err != nil {
		log.Fatalf("Ошибка загрузки файла: %v", err)
	}

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("Ошибка открытия файла: %v", err)
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Ошибка чтения файла: %v", err)
		}

		err = stream.Send(&pbFile.UploadFileRequest{
			Token:    *token,
			Filename: *filePath,
			Chunk:    buffer[:n],
		})
		if err != nil {
			log.Fatalf("Ошибка отправки файла: %v", err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Ошибка завершения загрузки файла: %v", err)
	}
	fmt.Println("Файл загружен. ID файла:", resp.FileId)
}

// Скачивание файла
func downloadFile(conn *grpc.ClientConn) {
	client := pbFile.NewFileServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	stream, err := client.DownloadFile(ctx, &pbFile.DownloadFileRequest{
		Token:  *token,
		FileId: *fileID,
	})
	if err != nil {
		log.Fatalf("Ошибка скачивания файла: %v", err)
	}

	file, err := os.Create("downloaded_" + *fileID)
	if err != nil {
		log.Fatalf("Ошибка создания файла: %v", err)
	}
	defer file.Close()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Ошибка получения файла: %v", err)
		}

		_, err = file.Write(resp.Chunk)
		if err != nil {
			log.Fatalf("Ошибка записи в файл: %v", err)
		}
	}

	fmt.Println("Файл успешно скачан:", "downloaded_"+*fileID)
}

// Синхронизация данных
func syncSecureData(conn *grpc.ClientConn) {
	client := pbSecureData.NewSecureDataServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	stream, err := client.SyncSecureData(ctx, &pbSecureData.SyncRequest{
		Token: *token,
	})
	if err != nil {
		log.Fatalf("Ошибка синхронизации данных: %v", err)
	}

	fmt.Println("Синхронизация данных:")
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Ошибка получения данных: %v", err)
		}

		fmt.Printf("ID: %s, Type: %s, Data: %s, Meta: %s, Action: %s\n",
			resp.Id, resp.DataType, resp.Data, resp.MetaInfo, resp.Action)
	}
}
