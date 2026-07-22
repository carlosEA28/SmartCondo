package services

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/carlosEA28/smartcondo/internal/apperrors"
	"github.com/carlosEA28/smartcondo/internal/dto"
	"github.com/carlosEA28/smartcondo/internal/models"
	"github.com/carlosEA28/smartcondo/internal/repositories"
	"github.com/carlosEA28/smartcondo/internal/utils"
	"github.com/google/uuid"
)

type S3Uploader interface {
	UploadFile(file *multipart.FileHeader, path string) (string, error)
	DeleteFile(path string) error
}

type VisitorService struct {
	visitorRepo repositories.VisitorRepository
	s3Uploader  S3Uploader
}

func NewVisitorService(visitorRepo repositories.VisitorRepository, s3Uploader S3Uploader) *VisitorService {
	return &VisitorService{visitorRepo: visitorRepo, s3Uploader: s3Uploader}
}

func (s *VisitorService) CreateVisitor(ctx context.Context, input *dto.CreateVisitorDTO) (*dto.VisitorResponseDTO, error) {
	existing, err := s.visitorRepo.FindByCPF(ctx, input.CPF)
	if err != nil && !errors.Is(err, apperrors.ErrVisitorNotFound) {
		return nil, fmt.Errorf("check visitor cpf: %w", err)
	}
	if existing != nil {
		return nil, apperrors.ErrVisitorAlreadyExists
	}

	phone, err := utils.ValidatePhoneNumber(input.Phone)
	if err != nil {
		return nil, err
	}

	visitor := &models.Visitor{
		ID:    uuid.New(),
		Name:  strings.TrimSpace(input.Name),
		CPF:   input.CPF,
		Phone: phone,
	}

	if input.Photo != nil {
		ext := filepath.Ext(input.Photo.Filename)
		key := fmt.Sprintf("visitors/%s/photo%s", visitor.ID, ext)
		visitor.Photo, err = s.s3Uploader.UploadFile(input.Photo, key)
		if err != nil {
			return nil, fmt.Errorf("upload visitor photo: %w", err)
		}
	}

	if err := s.visitorRepo.Create(ctx, visitor); err != nil {
		return nil, err
	}

	return visitorToResponse(visitor), nil
}

func (s *VisitorService) GetVisitor(ctx context.Context, id uuid.UUID) (*dto.VisitorResponseDTO, error) {
	visitor, err := s.visitorRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrVisitorNotFound) {
			return nil, apperrors.ErrVisitorNotFound
		}
		return nil, fmt.Errorf("get visitor: %w", err)
	}

	return visitorToResponse(visitor), nil
}

func (s *VisitorService) ListVisitors(ctx context.Context) ([]dto.VisitorResponseDTO, error) {
	visitors, err := s.visitorRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list visitors: %w", err)
	}

	response := make([]dto.VisitorResponseDTO, 0, len(visitors))
	for i := range visitors {
		response = append(response, *visitorToResponse(&visitors[i]))
	}

	return response, nil
}

func (s *VisitorService) DeleteVisitor(ctx context.Context, id uuid.UUID) error {
	visitor, err := s.visitorRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrVisitorNotFound) {
			return apperrors.ErrVisitorNotFound
		}
		return fmt.Errorf("find visitor: %w", err)
	}

	if visitor.Photo != "" {
		if err := s.s3Uploader.DeleteFile(visitor.Photo); err != nil {
			return fmt.Errorf("delete visitor photo: %w", err)
		}
	}

	if err := s.visitorRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete visitor: %w", err)
	}

	return nil
}

func visitorToResponse(visitor *models.Visitor) *dto.VisitorResponseDTO {
	return &dto.VisitorResponseDTO{
		ID:    visitor.ID,
		Name:  visitor.Name,
		CPF:   visitor.CPF,
		Phone: visitor.Phone,
		Photo: visitor.Photo,
	}
}
