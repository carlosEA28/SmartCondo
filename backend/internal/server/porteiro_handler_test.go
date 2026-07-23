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

type fakePorteiroService struct {
	searchResult  []dto.VisitorResponseDTO
	searchErr     error
	releaseResult *dto.VisitorResponseDTO
	releaseErr    error
}

func (f *fakePorteiroService) SearchVisitors(_ context.Context, _ *dto.VisitorFilterDTO) ([]dto.VisitorResponseDTO, error) {
	return f.searchResult, f.searchErr
}

func (f *fakePorteiroService) ReleaseVisitor(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ *uuid.UUID) (*dto.VisitorResponseDTO, error) {
	return f.releaseResult, f.releaseErr
}

func setupPorteiroRouter(handler *porteiroHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/porteiros/visitantes", handler.search)
	router.PATCH("/porteiros/visitantes/:id/liberar", handler.release)
	return router
}

// --- search tests ---

func TestPorteiroHandlerSearchSuccess(t *testing.T) {
	service := &fakePorteiroService{
		searchResult: []dto.VisitorResponseDTO{
			{ID: uuid.New(), Name: "João Silva", CPF: "12345678901"},
		},
	}
	handler := newPorteiroHandler(service)
	router := setupPorteiroRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/porteiros/visitantes?nome=João", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Search() status = %d, want %d", w.Code, http.StatusOK)
	}

	var response []dto.VisitorResponseDTO
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Search() failed to unmarshal response: %v", err)
	}
	if len(response) != 1 {
		t.Fatalf("Search() length = %d, want 1", len(response))
	}
}

func TestPorteiroHandlerSearchEmptyFilter(t *testing.T) {
	service := &fakePorteiroService{
		searchErr: apperrors.ErrFilterRequired,
	}
	handler := newPorteiroHandler(service)
	router := setupPorteiroRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/porteiros/visitantes", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Search() status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestPorteiroHandlerSearchServiceError(t *testing.T) {
	service := &fakePorteiroService{
		searchErr: errors.New("database error"),
	}
	handler := newPorteiroHandler(service)
	router := setupPorteiroRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/porteiros/visitantes?nome=João", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Search() status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestPorteiroHandlerSearchByCPF(t *testing.T) {
	service := &fakePorteiroService{
		searchResult: []dto.VisitorResponseDTO{
			{ID: uuid.New(), Name: "Maria Santos", CPF: "98765432100"},
		},
	}
	handler := newPorteiroHandler(service)
	router := setupPorteiroRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/porteiros/visitantes?cpf=98765432100", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Search() status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestPorteiroHandlerSearchByTelefone(t *testing.T) {
	service := &fakePorteiroService{
		searchResult: []dto.VisitorResponseDTO{
			{ID: uuid.New(), Name: "Pedro Lima", CPF: "11122233344"},
		},
	}
	handler := newPorteiroHandler(service)
	router := setupPorteiroRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/porteiros/visitantes?telefone=777", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Search() status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestPorteiroHandlerSearchByLiberado(t *testing.T) {
	service := &fakePorteiroService{
		searchResult: []dto.VisitorResponseDTO{
			{ID: uuid.New(), Name: "Ana Costa", CPF: "55566677788", Liberado: true},
		},
	}
	handler := newPorteiroHandler(service)
	router := setupPorteiroRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/porteiros/visitantes?liberado=true", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Search() status = %d, want %d", w.Code, http.StatusOK)
	}
}

// --- release tests ---

func TestPorteiroHandlerReleaseSuccess(t *testing.T) {
	service := &fakePorteiroService{
		releaseResult: &dto.VisitorResponseDTO{
			ID:       uuid.New(),
			Name:     "João Silva",
			CPF:      "12345678901",
			Liberado: true,
		},
	}
	handler := newPorteiroHandler(service)
	router := setupPorteiroRouter(handler)

	visitorID := uuid.New()
	body := `{"porteiro_id": "` + uuid.New().String() + `"}`

	req := httptest.NewRequest(http.MethodPatch, "/porteiros/visitantes/"+visitorID.String()+"/liberar", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Release() status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestPorteiroHandlerReleaseInvalidID(t *testing.T) {
	service := &fakePorteiroService{}
	handler := newPorteiroHandler(service)
	router := setupPorteiroRouter(handler)

	body := `{"porteiro_id": "` + uuid.New().String() + `"}`

	req := httptest.NewRequest(http.MethodPatch, "/porteiros/visitantes/invalid-id/liberar", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Release() status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestPorteiroHandlerReleaseInvalidBody(t *testing.T) {
	service := &fakePorteiroService{}
	handler := newPorteiroHandler(service)
	router := setupPorteiroRouter(handler)

	visitorID := uuid.New()
	body := `{"invalid": "json"}`

	req := httptest.NewRequest(http.MethodPatch, "/porteiros/visitantes/"+visitorID.String()+"/liberar", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Release() status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestPorteiroHandlerReleaseVisitorNotFound(t *testing.T) {
	service := &fakePorteiroService{
		releaseErr: apperrors.ErrVisitorNotFound,
	}
	handler := newPorteiroHandler(service)
	router := setupPorteiroRouter(handler)

	visitorID := uuid.New()
	body := `{"porteiro_id": "` + uuid.New().String() + `"}`

	req := httptest.NewRequest(http.MethodPatch, "/porteiros/visitantes/"+visitorID.String()+"/liberar", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Release() status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestPorteiroHandlerReleasePorteiroNotFound(t *testing.T) {
	service := &fakePorteiroService{
		releaseErr: apperrors.ErrPorteiroNotFound,
	}
	handler := newPorteiroHandler(service)
	router := setupPorteiroRouter(handler)

	visitorID := uuid.New()
	body := `{"porteiro_id": "` + uuid.New().String() + `"}`

	req := httptest.NewRequest(http.MethodPatch, "/porteiros/visitantes/"+visitorID.String()+"/liberar", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Release() status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestPorteiroHandlerReleaseGenericError(t *testing.T) {
	service := &fakePorteiroService{
		releaseErr: errors.New("database error"),
	}
	handler := newPorteiroHandler(service)
	router := setupPorteiroRouter(handler)

	visitorID := uuid.New()
	body := `{"porteiro_id": "` + uuid.New().String() + `"}`

	req := httptest.NewRequest(http.MethodPatch, "/porteiros/visitantes/"+visitorID.String()+"/liberar", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Release() status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestPorteiroHandlerReleaseWithMorador(t *testing.T) {
	service := &fakePorteiroService{
		releaseResult: &dto.VisitorResponseDTO{
			ID:       uuid.New(),
			Name:     "João Silva",
			CPF:      "12345678901",
			Liberado: true,
		},
	}
	handler := newPorteiroHandler(service)
	router := setupPorteiroRouter(handler)

	visitorID := uuid.New()
	moradorID := uuid.New()
	body := `{"porteiro_id": "` + uuid.New().String() + `", "morador_id": "` + moradorID.String() + `"}`

	req := httptest.NewRequest(http.MethodPatch, "/porteiros/visitantes/"+visitorID.String()+"/liberar", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Release() status = %d, want %d", w.Code, http.StatusOK)
	}
}
