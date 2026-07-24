package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carlosEA28/smartcondo/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// --- fake for areaComumService ---

type fakeAreaComumService struct {
	listResult []dto.AreaComumResponse
	listErr    error
}

func (f *fakeAreaComumService) ListAreas(_ context.Context) ([]dto.AreaComumResponse, error) {
	return f.listResult, f.listErr
}

// --- helper ---

func setupAreaComumRouter(handler *areaComumHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/morador/areas-comuns", handler.listAreas)
	return router
}

// --- tests ---

func TestAreaComumHandlerListSuccess(t *testing.T) {
	t.Parallel()
	svc := &fakeAreaComumService{
		listResult: []dto.AreaComumResponse{
			{ID: uuid.New(), Nome: "Piscina", Descricao: "Climatizada", Capacidade: 20},
			{ID: uuid.New(), Nome: "Academia", Descricao: "Equipamentos novos", Capacidade: 15},
		},
	}
	handler := newAreaComumHandler(svc)
	router := setupAreaComumRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/morador/areas-comuns", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusOK)
	}

	var result []dto.AreaComumResponse
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 areas, got %d", len(result))
	}
	if result[0].Nome != "Piscina" {
		t.Errorf("first name: got %q, want %q", result[0].Nome, "Piscina")
	}
}

func TestAreaComumHandlerListServiceError(t *testing.T) {
	t.Parallel()
	svc := &fakeAreaComumService{listErr: errors.New("db failure")}
	handler := newAreaComumHandler(svc)
	router := setupAreaComumRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/morador/areas-comuns", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestAreaComumHandlerListEmpty(t *testing.T) {
	t.Parallel()
	svc := &fakeAreaComumService{listResult: []dto.AreaComumResponse{}}
	handler := newAreaComumHandler(svc)
	router := setupAreaComumRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/morador/areas-comuns", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusOK)
	}

	var result []dto.AreaComumResponse
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}
