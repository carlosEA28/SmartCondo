package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/carlosEA28/smartcondo/internal/apperrors"
	"github.com/carlosEA28/smartcondo/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// --- fake for reservaService ---

type fakeReservaService struct {
	createResult *dto.ReservaResponse
	createErr    error
	listResult   []dto.ReservaResponse
	listErr      error
	cancelResult *dto.ReservaResponse
	cancelErr    error
}

func (f *fakeReservaService) CreateReserva(_ context.Context, _ uuid.UUID, _ dto.CreateReservaRequest) (*dto.ReservaResponse, error) {
	return f.createResult, f.createErr
}

func (f *fakeReservaService) ListMinhasReservas(_ context.Context, _ uuid.UUID) ([]dto.ReservaResponse, error) {
	return f.listResult, f.listErr
}

func (f *fakeReservaService) CancelReserva(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*dto.ReservaResponse, error) {
	return f.cancelResult, f.cancelErr
}

// --- helpers ---

func setupReservaRouter(handler *reservaHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	moradorID := uuid.New()

	router.POST("/morador/reservas", func(c *gin.Context) {
		c.Set("user_id", moradorID)
		c.Next()
	}, handler.create)

	router.GET("/morador/reservas", func(c *gin.Context) {
		c.Set("user_id", moradorID)
		c.Next()
	}, handler.listMinhas)

	router.PATCH("/morador/reservas/:id/cancelar", func(c *gin.Context) {
		c.Set("user_id", moradorID)
		c.Next()
	}, handler.cancelar)

	return router
}

func testReservaResponse() *dto.ReservaResponse {
	return &dto.ReservaResponse{
		ID:         uuid.New(),
		Data:       time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC),
		HoraInicio: "14:00",
		HoraFim:    "16:00",
		Status:     "PENDENTE",
		MoradorID:  uuid.New(),
		AreaComum: dto.AreaComumResponse{
			ID:         uuid.New(),
			Nome:       "Piscina",
			Descricao:  "Climatizada",
			Capacidade: 20,
		},
	}
}

// --- Create tests ---

func TestReservaHandlerCreateSuccess(t *testing.T) {
	t.Parallel()
	svc := &fakeReservaService{createResult: testReservaResponse()}
	handler := newReservaHandler(svc)
	router := setupReservaRouter(handler)

	body := `{"areacomum_id":"` + uuid.New().String() + `","data":"2026-08-15","horaInicio":"14:00","horaFim":"16:00"}`
	req := httptest.NewRequest(http.MethodPost, "/morador/reservas", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusCreated)
	}
}

func TestReservaHandlerCreateInvalidBody(t *testing.T) {
	t.Parallel()
	svc := &fakeReservaService{}
	handler := newReservaHandler(svc)
	router := setupReservaRouter(handler)

	body := `{"invalid": "json"}`
	req := httptest.NewRequest(http.MethodPost, "/morador/reservas", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestReservaHandlerCreateAreaNotFound(t *testing.T) {
	t.Parallel()
	svc := &fakeReservaService{createErr: apperrors.ErrAreaComumNotFound}
	handler := newReservaHandler(svc)
	router := setupReservaRouter(handler)

	body := `{"areacomum_id":"` + uuid.New().String() + `","data":"2026-08-15","horaInicio":"14:00","horaFim":"16:00"}`
	req := httptest.NewRequest(http.MethodPost, "/morador/reservas", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestReservaHandlerCreateInadimplente(t *testing.T) {
	t.Parallel()
	svc := &fakeReservaService{createErr: apperrors.ErrMoradorInadimplente}
	handler := newReservaHandler(svc)
	router := setupReservaRouter(handler)

	body := `{"areacomum_id":"` + uuid.New().String() + `","data":"2026-08-15","horaInicio":"14:00","horaFim":"16:00"}`
	req := httptest.NewRequest(http.MethodPost, "/morador/reservas", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestReservaHandlerCreateConflict(t *testing.T) {
	t.Parallel()
	svc := &fakeReservaService{createErr: apperrors.ErrReservaConflito}
	handler := newReservaHandler(svc)
	router := setupReservaRouter(handler)

	body := `{"areacomum_id":"` + uuid.New().String() + `","data":"2026-08-15","horaInicio":"14:00","horaFim":"16:00"}`
	req := httptest.NewRequest(http.MethodPost, "/morador/reservas", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusConflict)
	}
}

func TestReservaHandlerCreateGenericError(t *testing.T) {
	t.Parallel()
	svc := &fakeReservaService{createErr: errors.New("db failure")}
	handler := newReservaHandler(svc)
	router := setupReservaRouter(handler)

	body := `{"areacomum_id":"` + uuid.New().String() + `","data":"2026-08-15","horaInicio":"14:00","horaFim":"16:00"}`
	req := httptest.NewRequest(http.MethodPost, "/morador/reservas", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// --- ListMinhas tests ---

func TestReservaHandlerListMinhasSuccess(t *testing.T) {
	t.Parallel()
	svc := &fakeReservaService{listResult: []dto.ReservaResponse{*testReservaResponse()}}
	handler := newReservaHandler(svc)
	router := setupReservaRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/morador/reservas", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusOK)
	}

	var result []dto.ReservaResponse
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 reserva, got %d", len(result))
	}
}

func TestReservaHandlerListMinhasServiceError(t *testing.T) {
	t.Parallel()
	svc := &fakeReservaService{listErr: errors.New("db failure")}
	handler := newReservaHandler(svc)
	router := setupReservaRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/morador/reservas", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestReservaHandlerListMinhasEmpty(t *testing.T) {
	t.Parallel()
	svc := &fakeReservaService{listResult: []dto.ReservaResponse{}}
	handler := newReservaHandler(svc)
	router := setupReservaRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/morador/reservas", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusOK)
	}

	var result []dto.ReservaResponse
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}

// --- Cancelar tests ---

func TestReservaHandlerCancelarSuccess(t *testing.T) {
	t.Parallel()
	cancelled := testReservaResponse()
	cancelled.Status = "CANCELADA"
	svc := &fakeReservaService{cancelResult: cancelled}
	handler := newReservaHandler(svc)
	router := setupReservaRouter(handler)

	req := httptest.NewRequest(http.MethodPatch, "/morador/reservas/"+uuid.New().String()+"/cancelar", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusOK)
	}
}

func TestReservaHandlerCancelarInvalidID(t *testing.T) {
	t.Parallel()
	svc := &fakeReservaService{}
	handler := newReservaHandler(svc)
	router := setupReservaRouter(handler)

	req := httptest.NewRequest(http.MethodPatch, "/morador/reservas/not-a-uuid/cancelar", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestReservaHandlerCancelarNotFound(t *testing.T) {
	t.Parallel()
	svc := &fakeReservaService{cancelErr: apperrors.ErrReservaNotFound}
	handler := newReservaHandler(svc)
	router := setupReservaRouter(handler)

	req := httptest.NewRequest(http.MethodPatch, "/morador/reservas/"+uuid.New().String()+"/cancelar", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestReservaHandlerCancelarNotOwner(t *testing.T) {
	t.Parallel()
	svc := &fakeReservaService{cancelErr: apperrors.ErrReservaNotOwner}
	handler := newReservaHandler(svc)
	router := setupReservaRouter(handler)

	req := httptest.NewRequest(http.MethodPatch, "/morador/reservas/"+uuid.New().String()+"/cancelar", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestReservaHandlerCancelarGenericError(t *testing.T) {
	t.Parallel()
	svc := &fakeReservaService{cancelErr: errors.New("db failure")}
	handler := newReservaHandler(svc)
	router := setupReservaRouter(handler)

	req := httptest.NewRequest(http.MethodPatch, "/morador/reservas/"+uuid.New().String()+"/cancelar", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
