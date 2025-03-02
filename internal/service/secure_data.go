package service

import (
	"context"
	"errors"
	"log"

	pb "github.com/evildead81/gophkeeper/api/secure_data"
	"github.com/evildead81/gophkeeper/internal/repository"
	"github.com/evildead81/gophkeeper/pkg/session"
)

type SecureDataService struct {
	repo     repository.SecureDataRepositoryInterface
	userRepo repository.UserRepositoryInterface
	session  *session.SessionManager
	pb.UnimplementedSecureDataServiceServer
}

// NewSecureDataService создаёт новый сервис
func NewSecureDataService(repo repository.SecureDataRepositoryInterface, userRepo repository.UserRepositoryInterface, session *session.SessionManager) *SecureDataService {
	return &SecureDataService{
		repo:     repo,
		userRepo: userRepo,
		session:  session,
	}
}

// AddSecureData добавляет зашифрованные данные пользователя
func (s *SecureDataService) AddSecureData(ctx context.Context, req *pb.AddSecureDataRequest) (*pb.AddSecureDataResponse, error) {
	userID, err := s.session.GetUserID(req.Token)
	if err != nil {
		return nil, errors.New("пользователь не авторизован")
	}

	err = s.repo.SaveSecureData(ctx, userID, req.DataType, req.Data, req.MetaInfo)
	if err != nil {
		return nil, err
	}

	log.Printf("Добавлены данные типа %s для пользователя %s", req.DataType, userID)
	return &pb.AddSecureDataResponse{Message: "Данные сохранены"}, nil
}

// GetSecureData получает зашифрованные данные пользователя
func (s *SecureDataService) GetSecureData(ctx context.Context, req *pb.GetSecureDataRequest) (*pb.GetSecureDataResponse, error) {
	// Проверяем, авторизован ли пользователь
	userID, err := s.session.GetUserID(req.Token)
	if err != nil {
		return nil, errors.New("пользователь не авторизован")
	}

	data, err := s.repo.GetSecureData(ctx, userID, req.DataType)
	if err != nil {
		return nil, err
	}

	var response []*pb.SecureData
	for _, item := range data {
		response = append(response, &pb.SecureData{
			Id:       item.ID,
			Data:     item.Data,
			MetaInfo: item.MetaInfo,
		})
	}

	log.Printf("Пользователь %s запросил данные типа %s", userID, req.DataType)
	return &pb.GetSecureDataResponse{Items: response}, nil
}

// DeleteSecureData удаляет данные пользователя
func (s *SecureDataService) DeleteSecureData(ctx context.Context, req *pb.DeleteSecureDataRequest) (*pb.DeleteSecureDataResponse, error) {
	userID, err := s.session.GetUserID(req.Token)
	if err != nil {
		return nil, errors.New("пользователь не авторизован")
	}

	err = s.repo.DeleteSecureData(ctx, userID, req.Id)
	if err != nil {
		return nil, err
	}

	log.Printf("Удалены данные с ID %s у пользователя %s", req.Id, userID)
	return &pb.DeleteSecureDataResponse{Message: "Данные удалены"}, nil
}

// SyncSecureData синхронизирует данные пользователя
func (s *SecureDataService) SyncSecureData(req *pb.SyncRequest, stream pb.SecureDataService_SyncSecureDataServer) error {
	userID, err := s.session.GetUserID(req.Token)
	if err != nil {
		return errors.New("пользователь не авторизован")
	}

	updates, err := s.repo.GetLastUpdatedData(stream.Context(), userID)
	if err != nil {
		return err
	}

	for _, update := range updates {
		err := stream.Send(&pb.SyncResponse{
			Id:       update.ID,
			DataType: update.DataType,
			Data:     update.Data,
			MetaInfo: update.MetaInfo,
			Action:   update.Action,
		})
		if err != nil {
			return err
		}
	}

	log.Printf("Синхронизированы данные пользователя %s", userID)
	return nil
}
