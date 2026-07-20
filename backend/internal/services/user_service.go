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
	ErrUserNotFound          = errors.New("user not found")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrApartmentRequired     = errors.New("apartment is required for residents")
	ErrApartmentNotAllowed   = errors.New("apartment is only allowed for residents")
	ErrResponsibleNotAllowed = errors.New("only residents can be responsible for an apartment")
	ErrInvalidUserData       = errors.New("invalid user data")
	ErrUserInUse             = errors.New("user has related records")
)

type UserService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*dto.UserResponseDTO, error) {
	user, err := s.userRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, ErrUserNotFound
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

func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, input dto.UpdateUserDTO) (*dto.UserResponseDTO, error) {
	user, err := s.userRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("find user to update: %w", err)
	}

	if input.FullName != nil {
		fullName := strings.TrimSpace(*input.FullName)
		if fullName == "" {
			return nil, ErrInvalidUserData
		}
		user.FullName = fullName
	}
	if input.Phone != nil {
		phone := strings.TrimSpace(*input.Phone)
		if phone == "" {
			return nil, ErrInvalidUserData
		}
		user.Phone = phone
	}
	if input.Status != nil {
		status := models.UserStatus(*input.Status)
		if status != models.UserStatusActive && status != models.UserStatusInactive && status != models.UserStatusBlocked {
			return nil, ErrInvalidUserData
		}
		user.Status = status
	}
	if input.Email != nil {
		email := strings.ToLower(strings.TrimSpace(*input.Email))
		if email == "" {
			return nil, ErrInvalidUserData
		}
		if email != user.Email {
			existing, findErr := s.userRepository.FindByEmail(ctx, email)
			switch {
			case findErr == nil && existing.ID != user.ID:
				return nil, ErrUserAlreadyExists
			case findErr != nil && !errors.Is(findErr, repositories.ErrUserNotFound):
				return nil, fmt.Errorf("check existing user: %w", findErr)
			}
			user.Email = email
		}
	}

	if input.Role != nil {
		user.Role = models.Role(*input.Role)
	}
	var apartmentToSave *models.Apartment
	switch user.Role {
	case models.RoleMorador:
		if user.Apartment == nil {
			if input.Apartment == nil || input.Apartment.Number == nil || input.Apartment.Block == nil {
				return nil, ErrApartmentRequired
			}
			block := strings.TrimSpace(*input.Apartment.Block)
			if block == "" {
				return nil, ErrInvalidUserData
			}
			user.Apartment = &models.Apartment{
				ID:     uuid.New(),
				Number: *input.Apartment.Number,
				Block:  block,
			}
			user.ApartmentID = &user.Apartment.ID
			apartmentToSave = user.Apartment
		} else if input.Apartment != nil {
			if input.Apartment.Number != nil {
				user.Apartment.Number = *input.Apartment.Number
			}
			if input.Apartment.Block != nil {
				block := strings.TrimSpace(*input.Apartment.Block)
				if block == "" {
					return nil, ErrInvalidUserData
				}
				user.Apartment.Block = block
			}
			apartmentToSave = user.Apartment
		}
		if input.Responsible != nil {
			user.Responsible = *input.Responsible
		}
	case models.RolePorteiro, models.RoleSindico:
		if input.Apartment != nil {
			return nil, ErrApartmentNotAllowed
		}
		if input.Responsible != nil && *input.Responsible {
			return nil, ErrResponsibleNotAllowed
		}
		user.ApartmentID = nil
		user.Apartment = nil
		user.Responsible = false
	default:
		return nil, ErrInvalidUserData
	}

	if err := s.userRepository.Update(ctx, user, apartmentToSave); err != nil {
		switch {
		case errors.Is(err, repositories.ErrUserNotFound):
			return nil, ErrUserNotFound
		case errors.Is(err, repositories.ErrUserAlreadyExists):
			return nil, ErrUserAlreadyExists
		default:
			return nil, fmt.Errorf("update user: %w", err)
		}
	}

	return userToResponse(user), nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := s.userRepository.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, repositories.ErrUserNotFound):
			return ErrUserNotFound
		case errors.Is(err, repositories.ErrUserInUse):
			return ErrUserInUse
		default:
			return fmt.Errorf("delete user: %w", err)
		}
	}

	return nil
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
