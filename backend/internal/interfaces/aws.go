package interfaces

import (
	"context"
	"mime/multipart"

	"github.com/carlosEA28/smartcondo/internal/dto"
)

type S3Interface interface {
	UploadFile(file *multipart.FileHeader, path string) (string, error)
	DeleteFile(path string) error
	GetFileURL(path string) (string, error)
	GetFile(path string) (*multipart.FileHeader, error)
}

type CognitoInterface interface {
	CreateUser(ctx context.Context, user *dto.CreateUserDTO) (string, error)
	DeleteUser(ctx context.Context, userID string) error
	GetUser(ctx context.Context, userID string) (*dto.RegisterRequest, error)
	GetUserList(ctx context.Context) ([]*dto.RegisterRequest, error)
}
