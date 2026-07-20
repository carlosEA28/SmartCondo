package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/carlosEA28/smartcondo/internal/dto"
	"github.com/carlosEA28/smartcondo/internal/models"
	"github.com/carlosEA28/smartcondo/internal/repositories"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrApartmentRequired     = errors.New("apartment is required for residents")
	ErrApartmentNotAllowed   = errors.New("apartment is only allowed for residents")
	ErrResponsibleNotAllowed = errors.New("only residents can be responsible for an apartment")
)

type UserService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) CreateUser(ctx context.Context, input dto.CreateUserDTO) (*dto.UserResponseDTO, error) {
	input.FullName = strings.TrimSpace(input.FullName)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.Phone = strings.TrimSpace(input.Phone)

	role := models.Role(input.Role)
	var apartment *models.Apartment
	switch role {
	case models.RoleMorador:
		if input.Apartment == nil {
			return nil, ErrApartmentRequired
		}
		apartment = &models.Apartment{
			ID:     uuid.New(),
			Number: input.Apartment.Number,
			Block:  strings.TrimSpace(input.Apartment.Block),
		}
	case models.RolePorteiro, models.RoleSindico:
		if input.Apartment != nil {
			return nil, ErrApartmentNotAllowed
		}
		if input.Responsible {
			return nil, ErrResponsibleNotAllowed
		}
	default:
		return nil, fmt.Errorf("invalid user role: %s", input.Role)
	}

	_, err := s.userRepository.FindByEmail(ctx, input.Email)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}
	if !errors.Is(err, repositories.ErrUserNotFound) {
		return nil, fmt.Errorf("check existing user: %w", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash user password: %w", err)
	}

	user := &models.User{
		ID:          uuid.New(),
		FullName:    input.FullName,
		Email:       input.Email,
		Password:    string(passwordHash),
		Phone:       input.Phone,
		Status:      models.UserStatusActive,
		Role:        role,
		Responsible: input.Responsible,
		Apartment:   apartment,
	}
	if apartment != nil {
		user.ApartmentID = &apartment.ID
	}

	// TODO: provision the user in Amazon Cognito here before persisting it.
	// Cognito failures must abort this operation so no local orphan is created.

	if err := s.userRepository.Create(ctx, user, apartment); err != nil {
		if errors.Is(err, repositories.ErrUserAlreadyExists) {
			return nil, ErrUserAlreadyExists
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
