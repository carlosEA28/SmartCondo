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

type fakeComunicadoService struct {
	publishResult *dto.ComunicadoResponseDTO
	publishErr    error
	listResult    []dto.ComunicadoResponseDTO
	listErr       error
	getResult     *dto.ComunicadoResponseDTO
	getErr        error
	deleteErr     error
}

func (f *fakeComunicadoService) PublishComunicado(_ context.Context, _ uuid.UUID, _ *dto.CreateComunicadoDTO) (*dto.ComunicadoResponseDTO, error) {
	return f.publishResult, f.publishErr
}

func (f *fakeComunicadoService) ListComunicados(_ context.Context) ([]dto.ComunicadoResponseDTO, error) {
	return f.listResult, f.listErr
}

func (f *fakeComunicadoService) GetComunicado(_ context.Context, _ uuid.UUID) (*dto.ComunicadoResponseDTO, error) {
	return f.getResult, f.getErr
}

func (f *fakeComunicadoService) DeleteComunicado(_ context.Context, _ uuid.UUID, _ uuid.UUID) error {
	return f.deleteErr
}

func setupComunicadoRouter(handler *comunicadoHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/sindico/comunicados", func(c *gin.Context) {
		c.Set("user_id", uuid.New())
		c.Next()
	}, handler.create)
	router.GET("/sindico/comunicados", handler.list)
	router.GET("/sindico/comunicados/:id", handler.getByID)
	router.DELETE("/sindico/comunicados/:id", func(c *gin.Context) {
		c.Set("user_id", uuid.New())
		c.Next()
	}, handler.delete)
	return router
}

// --- create tests ---

func TestComunicadoHandlerCreateSuccess(t *testing.T) {
	service := &fakeComunicadoService{
		publishResult: &dto.ComunicadoResponseDTO{
			ID:        uuid.New(),
			Titulo:    "Aviso Teste",
			Descricao: "Descrição",
		},
	}
	handler := newComunicadoHandler(service)
	router := setupComunicadoRouter(handler)

	body := `{"titulo": "Aviso Teste", "descricao": "Descrição do aviso"}`

	req := httptest.NewRequest(http.MethodPost, "/sindico/comunicados", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Create() status = %d, want %d", w.Code, http.StatusCreated)
	}
}

func TestComunicadoHandlerCreateInvalidBody(t *testing.T) {
	service := &fakeComunicadoService{}
	handler := newComunicadoHandler(service)
	router := setupComunicadoRouter(handler)

	body := `{"invalid": "json"}`

	req := httptest.NewRequest(http.MethodPost, "/sindico/comunicados", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Create() status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestComunicadoHandlerCreateServiceError(t *testing.T) {
	service := &fakeComunicadoService{
		publishErr: errors.New("service error"),
	}
	handler := newComunicadoHandler(service)
	router := setupComunicadoRouter(handler)

	body := `{"titulo": "Aviso Teste", "descricao": "Descrição"}`

	req := httptest.NewRequest(http.MethodPost, "/sindico/comunicados", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Create() status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// --- list tests ---

func TestComunicadoHandlerListSuccess(t *testing.T) {
	service := &fakeComunicadoService{
		listResult: []dto.ComunicadoResponseDTO{
			{ID: uuid.New(), Titulo: "Aviso 1"},
			{ID: uuid.New(), Titulo: "Aviso 2"},
		},
	}
	handler := newComunicadoHandler(service)
	router := setupComunicadoRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/sindico/comunicados", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("List() status = %d, want %d", w.Code, http.StatusOK)
	}

	var response []dto.ComunicadoResponseDTO
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("List() failed to unmarshal response: %v", err)
	}
	if len(response) != 2 {
		t.Fatalf("List() length = %d, want 2", len(response))
	}
}

func TestComunicadoHandlerListError(t *testing.T) {
	service := &fakeComunicadoService{
		listErr: errors.New("database error"),
	}
	handler := newComunicadoHandler(service)
	router := setupComunicadoRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/sindico/comunicados", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("List() status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestComunicadoHandlerListEmpty(t *testing.T) {
	service := &fakeComunicadoService{
		listResult: []dto.ComunicadoResponseDTO{},
	}
	handler := newComunicadoHandler(service)
	router := setupComunicadoRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/sindico/comunicados", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("List() status = %d, want %d", w.Code, http.StatusOK)
	}
}

// --- getByID tests ---

func TestComunicadoHandlerGetByIDSuccess(t *testing.T) {
	id := uuid.New()
	service := &fakeComunicadoService{
		getResult: &dto.ComunicadoResponseDTO{
			ID:        id,
			Titulo:    "Aviso Teste",
			Descricao: "Descrição",
		},
	}
	handler := newComunicadoHandler(service)
	router := setupComunicadoRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/sindico/comunicados/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GetByID() status = %d, want %d", w.Code, http.StatusOK)
	}

	var response dto.ComunicadoResponseDTO
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("GetByID() failed to unmarshal response: %v", err)
	}
	if response.ID != id {
		t.Fatalf("GetByID() ID = %v, want %v", response.ID, id)
	}
}

func TestComunicadoHandlerGetByIDInvalidID(t *testing.T) {
	service := &fakeComunicadoService{}
	handler := newComunicadoHandler(service)
	router := setupComunicadoRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/sindico/comunicados/invalid-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("GetByID() status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestComunicadoHandlerGetByIDNotFound(t *testing.T) {
	id := uuid.New()
	service := &fakeComunicadoService{
		getErr: apperrors.ErrComunicadoNotFound,
	}
	handler := newComunicadoHandler(service)
	router := setupComunicadoRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/sindico/comunicados/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("GetByID() status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestComunicadoHandlerGetByIDGenericError(t *testing.T) {
	id := uuid.New()
	service := &fakeComunicadoService{
		getErr: errors.New("database error"),
	}
	handler := newComunicadoHandler(service)
	router := setupComunicadoRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/sindico/comunicados/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("GetByID() status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// --- delete tests ---

func TestComunicadoHandlerDeleteSuccess(t *testing.T) {
	service := &fakeComunicadoService{}
	handler := newComunicadoHandler(service)
	router := setupComunicadoRouter(handler)

	id := uuid.New()
	req := httptest.NewRequest(http.MethodDelete, "/sindico/comunicados/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("Delete() status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestComunicadoHandlerDeleteInvalidID(t *testing.T) {
	service := &fakeComunicadoService{}
	handler := newComunicadoHandler(service)
	router := setupComunicadoRouter(handler)

	req := httptest.NewRequest(http.MethodDelete, "/sindico/comunicados/invalid-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Delete() status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestComunicadoHandlerDeleteNotFound(t *testing.T) {
	id := uuid.New()
	service := &fakeComunicadoService{
		deleteErr: apperrors.ErrComunicadoNotFound,
	}
	handler := newComunicadoHandler(service)
	router := setupComunicadoRouter(handler)

	req := httptest.NewRequest(http.MethodDelete, "/sindico/comunicados/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Delete() status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestComunicadoHandlerDeleteNotOwner(t *testing.T) {
	id := uuid.New()
	service := &fakeComunicadoService{
		deleteErr: apperrors.ErrComunicadoNotOwner,
	}
	handler := newComunicadoHandler(service)
	router := setupComunicadoRouter(handler)

	req := httptest.NewRequest(http.MethodDelete, "/sindico/comunicados/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("Delete() status = %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestComunicadoHandlerDeleteGenericError(t *testing.T) {
	id := uuid.New()
	service := &fakeComunicadoService{
		deleteErr: errors.New("database error"),
	}
	handler := newComunicadoHandler(service)
	router := setupComunicadoRouter(handler)

	req := httptest.NewRequest(http.MethodDelete, "/sindico/comunicados/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Delete() status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
