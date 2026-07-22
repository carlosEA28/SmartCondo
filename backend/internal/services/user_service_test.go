package services

import (
	"context"
	"errors"
	"testing"

	"github.com/carlosEA28/smartcondo/internal/apperrors"
	"github.com/carlosEA28/smartcondo/internal/dto"
	"github.com/carlosEA28/smartcondo/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepository struct {
	findByIDResult      *models.User
	findByIDErr         error
	findByEmailResult   *models.User
	findByEmailErr      error
	findApartmentResult *models.Apartment
	findApartmentErr    error
	listResult          []models.User
	listErr             error
	createdUser         *models.User
	createdApartment    *models.Apartment
	createErr           error
	saveErr             error
	deleteErr           error
}

func (f *fakeUserRepository) FindByID(context.Context, uuid.UUID) (*models.User, error) {
	return f.findByIDResult, f.findByIDErr
}

func (f *fakeUserRepository) FindByEmail(context.Context, string) (*models.User, error) {
	return f.findByEmailResult, f.findByEmailErr
}

func (f *fakeUserRepository) FindApartmentByNumberAndBlock(context.Context, int, string) (*models.Apartment, error) {
	return f.findApartmentResult, f.findApartmentErr
}

func (f *fakeUserRepository) List(context.Context) ([]models.User, error) {
	return f.listResult, f.listErr
}

func (f *fakeUserRepository) Create(_ context.Context, user *models.User, apartment *models.Apartment) error {
	f.createdUser = user
	f.createdApartment = apartment
	return f.createErr
}

func (f *fakeUserRepository) Save(context.Context, *models.User) error {
	return f.saveErr
}

func (f *fakeUserRepository) Delete(context.Context, uuid.UUID) error {
	return f.deleteErr
}

type fakeCognitoProvider struct {
	createUserResult bool
	createUserErr    error
	deleteUserErr    error
}

func (f *fakeCognitoProvider) CreateUser(context.Context, *dto.CreateUserDTO) (bool, error) {
	return f.createUserResult, f.createUserErr
}

func (f *fakeCognitoProvider) DeleteUser(context.Context, string) error {
	return f.deleteUserErr
}

// --- GetUser tests ---

func TestUserServiceGetUserReturnsUser(t *testing.T) {
	id := uuid.New()
	repository := &fakeUserRepository{findByIDResult: &models.User{
		ID:       id,
		FullName: "Maria Silva",
		Email:    "maria@example.com",
		Status:   models.UserStatusActive,
		Role:     models.RoleMorador,
	}}
	cognito := &fakeCognitoProvider{}
	service := NewUserService(repository, cognito)

	response, err := service.GetUser(context.Background(), id)
	if err != nil {
		t.Fatalf("GetUser() error = %v", err)
	}
	if response.ID != id || response.Email != repository.findByIDResult.Email {
		t.Fatalf("GetUser() response = %#v", response)
	}
}

func TestUserServiceGetUserReturnsNotFound(t *testing.T) {
	repository := &fakeUserRepository{findByIDErr: apperrors.ErrUserNotFound}
	cognito := &fakeCognitoProvider{}
	service := NewUserService(repository, cognito)

	_, err := service.GetUser(context.Background(), uuid.New())
	if !errors.Is(err, apperrors.ErrUserNotFound) {
		t.Fatalf("GetUser() error = %v, want %v", err, apperrors.ErrUserNotFound)
	}
}

func TestUserServiceGetUserReturnsGenericError(t *testing.T) {
	repository := &fakeUserRepository{findByIDErr: errors.New("connection timeout")}
	cognito := &fakeCognitoProvider{}
	service := NewUserService(repository, cognito)

	_, err := service.GetUser(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("GetUser() expected error, got nil")
	}
	if errors.Is(err, apperrors.ErrUserNotFound) {
		t.Fatal("GetUser() should not return ErrUserNotFound for generic error")
	}
}

// --- ListUsers tests ---

func TestUserServiceListUsersReturnsUsers(t *testing.T) {
	repository := &fakeUserRepository{listResult: []models.User{
		{ID: uuid.New(), FullName: "Maria Silva", Role: models.RoleMorador},
		{ID: uuid.New(), FullName: "João Souza", Role: models.RolePorteiro},
	}}
	cognito := &fakeCognitoProvider{}
	service := NewUserService(repository, cognito)

	response, err := service.ListUsers(context.Background())
	if err != nil {
		t.Fatalf("ListUsers() error = %v", err)
	}
	if len(response) != len(repository.listResult) {
		t.Fatalf("ListUsers() length = %d, want %d", len(response), len(repository.listResult))
	}
}

func TestUserServiceListUsersReturnsError(t *testing.T) {
	repository := &fakeUserRepository{listErr: errors.New("database connection failed")}
	cognito := &fakeCognitoProvider{}
	service := NewUserService(repository, cognito)

	_, err := service.ListUsers(context.Background())
	if err == nil {
		t.Fatal("ListUsers() expected error, got nil")
	}
}

// --- UpdateUser tests ---

func TestUserServiceUpdateUserReturnsNotFound(t *testing.T) {
	repository := &fakeUserRepository{findByIDErr: apperrors.ErrUserNotFound}
	cognito := &fakeCognitoProvider{}
	service := NewUserService(repository, cognito)

	_, err := service.UpdateUser(context.Background(), uuid.New(), &dto.UpdateUserDTO{
		FullName: "Maria Silva",
		Phone:    "11999999999",
	})
	if !errors.Is(err, apperrors.ErrUserNotFound) {
		t.Fatalf("UpdateUser() error = %v, want %v", err, apperrors.ErrUserNotFound)
	}
}

func TestUserServiceUpdateUserRejectsEmptyFullName(t *testing.T) {
	repository := &fakeUserRepository{findByIDResult: &models.User{
		ID:       uuid.New(),
		FullName: "Maria Silva",
		Phone:    "11999999999",
	}}
	cognito := &fakeCognitoProvider{}
	service := NewUserService(repository, cognito)

	_, err := service.UpdateUser(context.Background(), uuid.New(), &dto.UpdateUserDTO{
		FullName: "  ",
		Phone:    "11999999999",
	})
	if !errors.Is(err, apperrors.ErrInvalidUserData) {
		t.Fatalf("UpdateUser() error = %v, want %v", err, apperrors.ErrInvalidUserData)
	}
}

func TestUserServiceUpdateUserRejectsEmptyPhone(t *testing.T) {
	repository := &fakeUserRepository{findByIDResult: &models.User{
		ID:       uuid.New(),
		FullName: "Maria Silva",
		Phone:    "11999999999",
	}}
	cognito := &fakeCognitoProvider{}
	service := NewUserService(repository, cognito)

	_, err := service.UpdateUser(context.Background(), uuid.New(), &dto.UpdateUserDTO{
		FullName: "Maria Silva",
		Phone:    "  ",
	})
	if !errors.Is(err, apperrors.ErrInvalidUserData) {
		t.Fatalf("UpdateUser() error = %v, want %v", err, apperrors.ErrInvalidUserData)
	}
}

func TestUserServiceUpdateUserSaveFails(t *testing.T) {
	repository := &fakeUserRepository{
		findByIDResult: &models.User{
			ID:       uuid.New(),
			FullName: "Maria Silva",
			Phone:    "11999999999",
		},
		saveErr: errors.New("database write failed"),
	}
	cognito := &fakeCognitoProvider{}
	service := NewUserService(repository, cognito)

	_, err := service.UpdateUser(context.Background(), uuid.New(), &dto.UpdateUserDTO{
		FullName: "Maria Santos",
		Phone:    "11888888888",
	})
	if err == nil {
		t.Fatal("UpdateUser() expected error, got nil")
	}
}

func TestUserServiceUpdateUserSuccess(t *testing.T) {
	id := uuid.New()
	repository := &fakeUserRepository{findByIDResult: &models.User{
		ID:       id,
		FullName: "Maria Silva",
		Phone:    "11999999999",
	}}
	cognito := &fakeCognitoProvider{}
	service := NewUserService(repository, cognito)

	response, err := service.UpdateUser(context.Background(), id, &dto.UpdateUserDTO{
		FullName: "Maria Santos",
		Phone:    "11888888888",
	})
	if err != nil {
		t.Fatalf("UpdateUser() error = %v", err)
	}
	if response.FullName != "Maria Santos" {
		t.Fatalf("UpdateUser() FullName = %q, want %q", response.FullName, "Maria Santos")
	}
}

// --- DeleteUser tests ---

func TestUserServiceDeleteUserReturnsNotFound(t *testing.T) {
	repository := &fakeUserRepository{findByIDErr: apperrors.ErrUserNotFound}
	cognito := &fakeCognitoProvider{}
	service := NewUserService(repository, cognito)

	err := service.DeleteUser(context.Background(), uuid.New())
	if !errors.Is(err, apperrors.ErrUserNotFound) {
		t.Fatalf("DeleteUser() error = %v, want %v", err, apperrors.ErrUserNotFound)
	}
}

func TestUserServiceDeleteUserFindByIDGenericError(t *testing.T) {
	repository := &fakeUserRepository{findByIDErr: errors.New("database connection lost")}
	cognito := &fakeCognitoProvider{}
	service := NewUserService(repository, cognito)

	err := service.DeleteUser(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("DeleteUser() expected error, got nil")
	}
	if errors.Is(err, apperrors.ErrUserNotFound) {
		t.Fatal("DeleteUser() should not return ErrUserNotFound for generic error")
	}
}

func TestUserServiceDeleteUserReturnsInUse(t *testing.T) {
	repository := &fakeUserRepository{
		findByIDResult: &models.User{
			ID:       uuid.New(),
			FullName: "Maria Silva",
			Email:    "maria@example.com",
		},
		deleteErr: apperrors.ErrUserInUse,
	}
	cognito := &fakeCognitoProvider{}
	service := NewUserService(repository, cognito)

	err := service.DeleteUser(context.Background(), uuid.New())
	if !errors.Is(err, apperrors.ErrUserInUse) {
		t.Fatalf("DeleteUser() error = %v, want %v", err, apperrors.ErrUserInUse)
	}
}

func TestUserServiceDeleteUserDeleteGenericError(t *testing.T) {
	repository := &fakeUserRepository{
		findByIDResult: &models.User{
			ID:       uuid.New(),
			FullName: "Maria Silva",
			Email:    "maria@example.com",
		},
		deleteErr: errors.New("database error"),
	}
	cognito := &fakeCognitoProvider{}
	service := NewUserService(repository, cognito)

	err := service.DeleteUser(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("DeleteUser() expected error, got nil")
	}
	if errors.Is(err, apperrors.ErrUserNotFound) {
		t.Fatal("DeleteUser() should not return ErrUserNotFound for generic error")
	}
}

func TestUserServiceDeleteUserSuccess(t *testing.T) {
	repository := &fakeUserRepository{
		findByIDResult: &models.User{
			ID:       uuid.New(),
			FullName: "Maria Silva",
			Email:    "maria@example.com",
		},
	}
	cognito := &fakeCognitoProvider{}
	service := NewUserService(repository, cognito)

	err := service.DeleteUser(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("DeleteUser() error = %v, want nil", err)
	}
}

// --- CreateUser tests ---

func TestUserServiceCreateUserCreatesResidentAndApartment(t *testing.T) {
	repository := &fakeUserRepository{findByEmailErr: apperrors.ErrUserNotFound}
	cognito := &fakeCognitoProvider{createUserResult: true}
	service := NewUserService(repository, cognito)
	input := &dto.CreateUserDTO{
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
	if repository.createdUser.Email != "  MARIA@EXAMPLE.COM " {
		t.Fatalf("created email = %q, want %q", repository.createdUser.Email, "  MARIA@EXAMPLE.COM ")
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
	repository := &fakeUserRepository{findByEmailResult: &models.User{Email: "maria@example.com"}}
	cognito := &fakeCognitoProvider{}
	service := NewUserService(repository, cognito)

	_, err := service.CreateUser(context.Background(), validResidentInput())
	if !errors.Is(err, apperrors.ErrUserAlreadyExists) {
		t.Fatalf("CreateUser() error = %v, want %v", err, apperrors.ErrUserAlreadyExists)
	}
	if repository.createdUser != nil {
		t.Fatal("CreateUser() persisted a duplicate user")
	}
}

func TestUserServiceCreateUserFindByEmailGenericError(t *testing.T) {
	repository := &fakeUserRepository{findByEmailErr: errors.New("database connection failed")}
	cognito := &fakeCognitoProvider{}
	service := NewUserService(repository, cognito)

	_, err := service.CreateUser(context.Background(), validResidentInput())
	if err == nil {
		t.Fatal("CreateUser() expected error, got nil")
	}
}

func TestUserServiceCreateUserRequiresApartmentForResident(t *testing.T) {
	repository := &fakeUserRepository{findByEmailErr: apperrors.ErrUserNotFound}
	cognito := &fakeCognitoProvider{createUserResult: true}
	service := NewUserService(repository, cognito)
	input := validResidentInput()
	input.Apartment = nil

	_, err := service.CreateUser(context.Background(), input)
	if !errors.Is(err, apperrors.ErrApartmentRequired) {
		t.Fatalf("CreateUser() error = %v, want %v", err, apperrors.ErrApartmentRequired)
	}
}

func TestUserServiceCreateUserApartmentAlreadyExists(t *testing.T) {
	repository := &fakeUserRepository{
		findByEmailErr:      apperrors.ErrUserNotFound,
		findApartmentResult: &models.Apartment{Number: 101, Block: "A"},
	}
	cognito := &fakeCognitoProvider{}
	service := NewUserService(repository, cognito)

	_, err := service.CreateUser(context.Background(), validResidentInput())
	if !errors.Is(err, apperrors.ErrApartmentAlreadyExists) {
		t.Fatalf("CreateUser() error = %v, want %v", err, apperrors.ErrApartmentAlreadyExists)
	}
}

func TestUserServiceCreateUserFindApartmentGenericError(t *testing.T) {
	repository := &fakeUserRepository{
		findByEmailErr:   apperrors.ErrUserNotFound,
		findApartmentErr: errors.New("database error"),
	}
	cognito := &fakeCognitoProvider{}
	service := NewUserService(repository, cognito)

	_, err := service.CreateUser(context.Background(), validResidentInput())
	if err == nil {
		t.Fatal("CreateUser() expected error, got nil")
	}
}

func TestUserServiceCreateUserCreateError(t *testing.T) {
	repository := &fakeUserRepository{
		findByEmailErr: apperrors.ErrUserNotFound,
		createErr:      errors.New("database write failed"),
	}
	cognito := &fakeCognitoProvider{}
	service := NewUserService(repository, cognito)

	_, err := service.CreateUser(context.Background(), validResidentInput())
	if err == nil {
		t.Fatal("CreateUser() expected error, got nil")
	}
}

func TestUserServiceCreateUserCreateDuplicateKeyError(t *testing.T) {
	repository := &fakeUserRepository{
		findByEmailErr: apperrors.ErrUserNotFound,
		createErr:      apperrors.ErrUserAlreadyExists,
	}
	cognito := &fakeCognitoProvider{}
	service := NewUserService(repository, cognito)

	_, err := service.CreateUser(context.Background(), validResidentInput())
	if !errors.Is(err, apperrors.ErrUserAlreadyExists) {
		t.Fatalf("CreateUser() error = %v, want %v", err, apperrors.ErrUserAlreadyExists)
	}
}

func TestUserServiceCreateUserCallsCognito(t *testing.T) {
	repository := &fakeUserRepository{findByEmailErr: apperrors.ErrUserNotFound}
	cognito := &fakeCognitoProvider{createUserResult: true}
	service := NewUserService(repository, cognito)
	input := validResidentInput()

	_, err := service.CreateUser(context.Background(), input)
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	if cognito.createUserResult != true {
		t.Fatal("CreateUser() did not call Cognito provider")
	}
}

func TestUserServiceCreateUserCognitoErrorDoesNotFail(t *testing.T) {
	repository := &fakeUserRepository{findByEmailErr: apperrors.ErrUserNotFound}
	cognito := &fakeCognitoProvider{
		createUserResult: false,
		createUserErr:    errors.New("cognito service unavailable"),
	}
	service := NewUserService(repository, cognito)

	response, err := service.CreateUser(context.Background(), validResidentInput())
	if err != nil {
		t.Fatalf("CreateUser() should not fail when Cognito fails, error = %v", err)
	}
	if response == nil {
		t.Fatal("CreateUser() expected response, got nil")
	}
}

func validResidentInput() *dto.CreateUserDTO {
	return &dto.CreateUserDTO{
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
