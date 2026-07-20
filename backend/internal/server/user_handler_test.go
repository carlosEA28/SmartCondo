package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carlosEA28/smartcondo/internal/apperrors"
	"github.com/carlosEA28/smartcondo/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type fakeUserService struct {
	createUserResult *dto.UserResponseDTO
	createUserErr    error
	getUserResult    *dto.UserResponseDTO
	getUserErr       error
	listUsersResult  []dto.UserResponseDTO
	listUsersErr     error
	updateUserResult *dto.UserResponseDTO
	updateUserErr    error
	deleteUserErr    error
}

func (f *fakeUserService) CreateUser(_ context.Context, _ *dto.CreateUserDTO) (*dto.UserResponseDTO, error) {
	return f.createUserResult, f.createUserErr
}

func (f *fakeUserService) GetUser(_ context.Context, _ uuid.UUID) (*dto.UserResponseDTO, error) {
	return f.getUserResult, f.getUserErr
}

func (f *fakeUserService) ListUsers(_ context.Context) ([]dto.UserResponseDTO, error) {
	return f.listUsersResult, f.listUsersErr
}

func (f *fakeUserService) UpdateUser(_ context.Context, _ uuid.UUID, _ *dto.UpdateUserDTO) (*dto.UserResponseDTO, error) {
	return f.updateUserResult, f.updateUserErr
}

func (f *fakeUserService) DeleteUser(_ context.Context, _ uuid.UUID) error {
	return f.deleteUserErr
}

func setupRouter(handler *userHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/users", handler.create)
	router.GET("/users", handler.list)
	router.GET("/users/:id", handler.getByID)
	router.PUT("/users/:id", handler.update)
	router.DELETE("/users/:id", handler.delete)
	return router
}

func TestUserHandlerCreateSuccess(t *testing.T) {
	service := &fakeUserService{
		createUserResult: &dto.UserResponseDTO{
			ID:       uuid.New(),
			FullName: "Maria Silva",
			Email:    "maria@example.com",
		},
	}
	handler := newUserHandler(service)
	router := setupRouter(handler)

	body := `{
		"full_name": "Maria Silva",
		"email": "maria@example.com",
		"password": "password123",
		"phone": "11999999999",
		"apartment": {
			"number": 101,
			"block": "A"
		}
	}`

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Create() status = %d, want %d", w.Code, http.StatusCreated)
	}
}

func TestUserHandlerCreateInvalidBody(t *testing.T) {
	service := &fakeUserService{}
	handler := newUserHandler(service)
	router := setupRouter(handler)

	body := `{"invalid": "json"}`

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Create() status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestUserHandlerCreateDuplicateEmail(t *testing.T) {
	service := &fakeUserService{
		createUserErr: apperrors.ErrUserAlreadyExists,
	}
	handler := newUserHandler(service)
	router := setupRouter(handler)

	body := `{
		"full_name": "Maria Silva",
		"email": "maria@example.com",
		"password": "password123",
		"phone": "11999999999",
		"apartment": {
			"number": 101,
			"block": "A"
		}
	}`

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("Create() status = %d, want %d", w.Code, http.StatusConflict)
	}
}

func TestUserHandlerCreateApartmentRequired(t *testing.T) {
	service := &fakeUserService{
		createUserErr: apperrors.ErrApartmentRequired,
	}
	handler := newUserHandler(service)
	router := setupRouter(handler)

	body := `{
		"full_name": "Maria Silva",
		"email": "maria@example.com",
		"password": "password123",
		"phone": "11999999999",
		"apartment": {
			"number": 101,
			"block": "A"
		}
	}`

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("Create() status = %d, want %d", w.Code, http.StatusUnprocessableEntity)
	}
}

func TestUserHandlerGetByIDSuccess(t *testing.T) {
	id := uuid.New()
	service := &fakeUserService{
		getUserResult: &dto.UserResponseDTO{
			ID:       id,
			FullName: "Maria Silva",
			Email:    "maria@example.com",
		},
	}
	handler := newUserHandler(service)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/users/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GetByID() status = %d, want %d", w.Code, http.StatusOK)
	}

	var response dto.UserResponseDTO
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("GetByID() failed to unmarshal response: %v", err)
	}
	if response.ID != id {
		t.Fatalf("GetByID() ID = %v, want %v", response.ID, id)
	}
}

func TestUserHandlerGetByIDInvalidID(t *testing.T) {
	service := &fakeUserService{}
	handler := newUserHandler(service)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/users/invalid-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("GetByID() status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestUserHandlerGetByIDNotFound(t *testing.T) {
	id := uuid.New()
	service := &fakeUserService{
		getUserErr: apperrors.ErrUserNotFound,
	}
	handler := newUserHandler(service)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/users/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("GetByID() status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestUserHandlerGetByIDInternalError(t *testing.T) {
	id := uuid.New()
	service := &fakeUserService{
		getUserErr: errors.New("database error"),
	}
	handler := newUserHandler(service)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/users/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("GetByID() status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestUserHandlerListSuccess(t *testing.T) {
	service := &fakeUserService{
		listUsersResult: []dto.UserResponseDTO{
			{ID: uuid.New(), FullName: "Maria Silva"},
			{ID: uuid.New(), FullName: "João Souza"},
		},
	}
	handler := newUserHandler(service)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("List() status = %d, want %d", w.Code, http.StatusOK)
	}

	var response []dto.UserResponseDTO
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("List() failed to unmarshal response: %v", err)
	}
	if len(response) != 2 {
		t.Fatalf("List() length = %d, want 2", len(response))
	}
}

func TestUserHandlerListInternalError(t *testing.T) {
	service := &fakeUserService{
		listUsersErr: errors.New("database error"),
	}
	handler := newUserHandler(service)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("List() status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestUserHandlerUpdateSuccess(t *testing.T) {
	id := uuid.New()
	service := &fakeUserService{
		updateUserResult: &dto.UserResponseDTO{
			ID:       id,
			FullName: "Maria Santos",
			Phone:    "11888888888",
		},
	}
	handler := newUserHandler(service)
	router := setupRouter(handler)

	body := `{
		"full_name": "Maria Santos",
		"phone": "11888888888"
	}`

	req := httptest.NewRequest(http.MethodPut, "/users/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Update() status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestUserHandlerUpdateInvalidID(t *testing.T) {
	service := &fakeUserService{}
	handler := newUserHandler(service)
	router := setupRouter(handler)

	body := `{"full_name": "Maria Santos", "phone": "11888888888"}`

	req := httptest.NewRequest(http.MethodPut, "/users/invalid-id", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Update() status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestUserHandlerUpdateNotFound(t *testing.T) {
	id := uuid.New()
	service := &fakeUserService{
		updateUserErr: apperrors.ErrUserNotFound,
	}
	handler := newUserHandler(service)
	router := setupRouter(handler)

	body := `{"full_name": "Maria Santos", "phone": "11888888888"}`

	req := httptest.NewRequest(http.MethodPut, "/users/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Update() status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestUserHandlerUpdateInvalidData(t *testing.T) {
	id := uuid.New()
	service := &fakeUserService{
		updateUserErr: apperrors.ErrInvalidUserData,
	}
	handler := newUserHandler(service)
	router := setupRouter(handler)

	body := `{"full_name": "Maria Santos", "phone": "11888888888"}`

	req := httptest.NewRequest(http.MethodPut, "/users/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("Update() status = %d, want %d", w.Code, http.StatusUnprocessableEntity)
	}
}

func TestUserHandlerDeleteSuccess(t *testing.T) {
	id := uuid.New()
	service := &fakeUserService{}
	handler := newUserHandler(service)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodDelete, "/users/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("Delete() status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestUserHandlerDeleteInvalidID(t *testing.T) {
	service := &fakeUserService{}
	handler := newUserHandler(service)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodDelete, "/users/invalid-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Delete() status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestUserHandlerDeleteNotFound(t *testing.T) {
	id := uuid.New()
	service := &fakeUserService{
		deleteUserErr: apperrors.ErrUserNotFound,
	}
	handler := newUserHandler(service)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodDelete, "/users/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Delete() status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestUserHandlerDeleteInUse(t *testing.T) {
	id := uuid.New()
	service := &fakeUserService{
		deleteUserErr: apperrors.ErrUserInUse,
	}
	handler := newUserHandler(service)
	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodDelete, "/users/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("Delete() status = %d, want %d", w.Code, http.StatusConflict)
	}
}
