package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/carlosEA28/smartcondo/internal/apperrors"
	"github.com/carlosEA28/smartcondo/internal/dto"
	"github.com/carlosEA28/smartcondo/internal/models"
	"github.com/carlosEA28/smartcondo/internal/repositories"
	"github.com/carlosEA28/smartcondo/internal/utils"
	"github.com/google/uuid"
)

type CognitoProvider interface {
	CreateUser(ctx context.Context, user *dto.CreateUserDTO) (bool, error)
}

type UserService struct {
	userRepository  repositories.UserRepository
	cognitoProvider CognitoProvider
}

func NewUserService(userRepository repositories.UserRepository, cognitoProvider CognitoProvider) *UserService {
	return &UserService{userRepository: userRepository, cognitoProvider: cognitoProvider}
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*dto.UserResponseDTO, error) {
	user, err := s.userRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	return userToResponse(user), nil
}

func (s *UserService) ListUsers(ctx context.Context) ([]dto.UserResponseDTO, error) {
	users, err := s.userRepository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	response := make([]dto.UserResponseDTO, 0, len(users))
	for index := range users {
		response = append(response, *userToResponse(&users[index]))
	}

	return response, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, input *dto.UpdateUserDTO) (*dto.UserResponseDTO, error) {
	user, err := s.userRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.FullName = strings.TrimSpace(input.FullName)
	user.Phone = strings.TrimSpace(input.Phone)
	if user.FullName == "" || user.Phone == "" {
		return nil, apperrors.ErrInvalidUserData
	}

	if err := s.userRepository.Save(ctx, user); err != nil {
		return nil, err
	}

	return s.GetUser(ctx, id)
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := s.userRepository.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUserNotFound):
			return apperrors.ErrUserNotFound
		case errors.Is(err, apperrors.ErrUserInUse):
			return apperrors.ErrUserInUse
		default:
			return fmt.Errorf("delete user: %w", err)
		}
	}

	return nil
}

func (s *UserService) CreateUser(ctx context.Context, input *dto.CreateUserDTO) (*dto.UserResponseDTO, error) {

	userExists, err := s.userRepository.FindByEmail(ctx, input.Email)
	if err == nil && userExists.Email == input.Email {
		return nil, apperrors.ErrUserAlreadyExists
	}

	s.cognitoProvider.CreateUser(ctx, input)

	passwordHash, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash user password: %w", err)

	}

	if input.Apartment == nil {
		return nil, apperrors.ErrApartmentRequired
	}
	apartment := &models.Apartment{
		ID:     uuid.New(),
		Number: input.Apartment.Number,
		Block:  strings.TrimSpace(input.Apartment.Block),
	}

	validNumber, err := utils.ValidatePhoneNumber(input.Phone)
	if err != nil {
		return nil, err
	}

	fmt.Println("TESTE NUMERO AQUI" + validNumber)

	user := &models.User{
		ID:          uuid.New(),
		FullName:    input.FullName,
		Email:       input.Email,
		Password:    string(passwordHash),
		Phone:       input.Phone,
		Status:      models.UserStatusActive,
		Role:        models.RoleMorador,
		Responsible: input.Responsible,
		Apartment:   apartment,
	}
	user.ApartmentID = &apartment.ID

	if err := s.userRepository.Create(ctx, user, apartment); err != nil {
		if errors.Is(err, apperrors.ErrUserAlreadyExists) {
			return nil, apperrors.ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	return userToResponse(user), nil
}

func userToResponse(user *models.User) *dto.UserResponseDTO {
	response := &dto.UserResponseDTO{
		ID:          user.ID,
		FullName:    user.FullName,
		Email:       user.Email,
		Phone:       user.Phone,
		Status:      string(user.Status),
		Role:        string(user.Role),
		Responsible: user.Responsible,
	}

	if user.Apartment != nil {
		response.Apartment = &dto.ApartmentResponseDTO{
			ID:     user.Apartment.ID,
			Number: user.Apartment.Number,
			Block:  user.Apartment.Block,
		}
	}

	return response

}
