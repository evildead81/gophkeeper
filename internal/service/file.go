package service

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"

	pb "github.com/evildead81/gophkeeper/api/file"
	"github.com/evildead81/gophkeeper/internal/repository"
	"github.com/evildead81/gophkeeper/pkg/session"
)

type FileService struct {
	fileRepo    *repository.FileRepository
	userRepo    *repository.UserRepository
	session     *session.SessionManager
	storagePath string
	pb.UnimplementedFileServiceServer
}

// NewFileService создаёт новый сервис
func NewFileService(fileRepo *repository.FileRepository, userRepo *repository.UserRepository, session *session.SessionManager) *FileService {
	return &FileService{
		fileRepo:    fileRepo,
		userRepo:    userRepo,
		session:     session,
		storagePath: "storage",
	}
}

// UploadFile загружает файл в систему
func (s *FileService) UploadFile(stream pb.FileService_UploadFileServer) error {
	var fileID, userID, filename string
	var file *os.File

	for {
		req, err := stream.Recv()

		userID, err = s.session.GetUserID(req.Token)
		if err != nil {
			return errors.New("пользователь не авторизован")
		}

		if err == io.EOF {
			s.fileRepo.SaveFileInfo(stream.Context(), userID, fileID, filename)
			log.Printf("Файл %s загружен пользователем %s", filename, userID)
			return stream.SendAndClose(&pb.UploadFileResponse{FileId: fileID})
		}
		if err != nil {
			return err
		}

		if file == nil {

			filename = req.Filename
			fileID = "file_" + filename
			filePath := filepath.Join(s.storagePath, fileID)

			file, err = os.Create(filePath)
			if err != nil {
				return err
			}
			defer file.Close()
		}

		_, err = file.Write(req.Chunk)
		if err != nil {
			return err
		}
	}
}

// DownloadFile скачивает файл
func (s *FileService) DownloadFile(req *pb.DownloadFileRequest, stream pb.FileService_DownloadFileServer) error {
	userID, err := s.session.GetUserID(req.Token)
	if err != nil {
		return errors.New("пользователь не авторизован")
	}

	dbUserID, filename, err := s.fileRepo.GetFileInfo(stream.Context(), req.FileId)
	if err != nil {
		return err
	}

	if dbUserID != userID {
		return errors.New("нет доступа к файлу")
	}

	filePath := filepath.Join(s.storagePath, req.FileId)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := stream.Send(&pb.DownloadFileResponse{Chunk: buffer[:n]}); err != nil {
			return err
		}
	}

	log.Printf("Файл %s скачан пользователем %s", filename, userID)
	return nil
}
