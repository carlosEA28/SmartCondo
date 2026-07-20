package services

import (
	"context"
	"errors"
	"testing"

	"github.com/carlosEA28/smartcondo/internal/dto"
	"github.com/carlosEA28/smartcondo/internal/models"
	"github.com/carlosEA28/smartcondo/internal/repositories"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepository struct {
	findByIDResult    *models.User
	findByIDErr       error
	findByEmailResult *models.User
	findByEmailErr    error
	listResult        []models.User
	listErr           error
	createdUser       *models.User
	createdApartment  *models.Apartment
	createErr         error
}

func (f *fakeUserRepository) FindByID(context.Context, uuid.UUID) (*models.User, error) {
	return f.findByIDResult, f.findByIDErr
}

func (f *fakeUserRepository) FindByEmail(context.Context, string) (*models.User, error) {
	return f.findByEmailResult, f.findByEmailErr
}

func (f *fakeUserRepository) List(context.Context) ([]models.User, error) {
	return f.listResult, f.listErr
}

func (f *fakeUserRepository) Create(_ context.Context, user *models.User, apartment *models.Apartment) error {
	f.createdUser = user
	f.createdApartment = apartment
	return f.createErr
}

func (f *fakeUserRepository) Update(context.Context, *models.User, *models.Apartment) error {
	return nil
}

func (f *fakeUserRepository) Delete(context.Context, uuid.UUID) error {
	return nil
}

func TestUserServiceGetUserReturnsUser(t *testing.T) {
	id := uuid.New()
	repository := &fakeUserRepository{findByIDResult: &models.User{
		ID:       id,
		FullName: "Maria Silva",
		Email:    "maria@example.com",
		Status:   models.UserStatusActive,
		Role:     models.RoleMorador,
	}}
	service := NewUserService(repository)

	response, err := service.GetUser(context.Background(), id)
	if err != nil {
		t.Fatalf("GetUser() error = %v", err)
	}
	if response.ID != id || response.Email != repository.findByIDResult.Email {
		t.Fatalf("GetUser() response = %#v", response)
	}
}

func TestUserServiceGetUserReturnsNotFound(t *testing.T) {
	repository := &fakeUserRepository{findByIDErr: repositories.ErrUserNotFound}
	service := NewUserService(repository)

	_, err := service.GetUser(context.Background(), uuid.New())
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("GetUser() error = %v, want %v", err, ErrUserNotFound)
	}
}

func TestUserServiceListUsersReturnsUsers(t *testing.T) {
	repository := &fakeUserRepository{listResult: []models.User{
		{ID: uuid.New(), FullName: "Maria Silva", Role: models.RoleMorador},
		{ID: uuid.New(), FullName: "João Souza", Role: models.RolePorteiro},
	}}
	service := NewUserService(repository)

	response, err := service.ListUsers(context.Background())
	if err != nil {
		t.Fatalf("ListUsers() error = %v", err)
	}
	if len(response) != len(repository.listResult) {
		t.Fatalf("ListUsers() length = %d, want %d", len(response), len(repository.listResult))
	}
}

func TestUserServiceCreateUserCreatesResidentAndApartment(t *testing.T) {
	repository := &fakeUserRepository{findByEmailErr: repositories.ErrUserNotFound}
	service := NewUserService(repository)
	input := dto.CreateUserDTO{
		FullName:    "  Maria Silva  ",
		Email:       "  MARIA@EXAMPLE.COM ",
		Password:    "password123",
		Phone:       " 11999999999 ",
		Responsible: true,
		Apartment: &dto.CreateApartmentDTO{
			Number: 101,
			Block:  " A ",
		},
	}

	response, err := service.CreateUser(context.Background(), input)
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	if repository.createdUser == nil || repository.createdApartment == nil {
		t.Fatal("CreateUser() did not persist both user and apartment")
	}
	if repository.createdUser.Email != "maria@example.com" {
		t.Fatalf("created email = %q, want %q", repository.createdUser.Email, "maria@example.com")
	}
	if repository.createdApartment.Block != "A" {
		t.Fatalf("created apartment block = %q, want %q", repository.createdApartment.Block, "A")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(repository.createdUser.Password), []byte(input.Password)); err != nil {
		t.Fatalf("stored password is not a valid hash: %v", err)
	}
	if response.Apartment == nil || response.Apartment.Number != input.Apartment.Number {
		t.Fatalf("response apartment = %#v, want number %d", response.Apartment, input.Apartment.Number)
	}
	if response.Status != string(models.UserStatusActive) {
		t.Fatalf("response status = %q, want %q", response.Status, models.UserStatusActive)
	}
}

func TestUserServiceCreateUserRejectsDuplicateEmail(t *testing.T) {
	repository := &fakeUserRepository{findByEmailResult: &models.User{}}
	service := NewUserService(repository)

	_, err := service.CreateUser(context.Background(), validResidentInput())
	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Fatalf("CreateUser() error = %v, want %v", err, ErrUserAlreadyExists)
	}
	if repository.createdUser != nil {
		t.Fatal("CreateUser() persisted a duplicate user")
	}
}

func TestUserServiceCreateUserRequiresApartmentForResident(t *testing.T) {
	repository := &fakeUserRepository{}
	service := NewUserService(repository)
	input := validResidentInput()
	input.Apartment = nil

	_, err := service.CreateUser(context.Background(), input)
	if !errors.Is(err, ErrApartmentRequired) {
		t.Fatalf("CreateUser() error = %v, want %v", err, ErrApartmentRequired)
	}
}

func validResidentInput() dto.CreateUserDTO {
	return dto.CreateUserDTO{
		FullName: "Maria Silva",
		Email:    "maria@example.com",
		Password: "password123",
		Phone:    "11999999999",
		Apartment: &dto.CreateApartmentDTO{
			Number: 101,
			Block:  "A",
		},
	}
}
