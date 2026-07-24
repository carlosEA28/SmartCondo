package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/carlosEA28/smartcondo/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// --- fake for inadimplenteService ---

type fakeInadimplenteService struct {
	listResult []dto.InadimplenteResponseDTO
	listErr    error
}

func (f *fakeInadimplenteService) ListInadimplentes(_ context.Context) ([]dto.InadimplenteResponseDTO, error) {
	return f.listResult, f.listErr
}

// --- helper ---

func setupInadimplenteRouter(handler *inadimplenteHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/sindico/inadimplentes", handler.list)
	return router
}

func TestInadimplenteHandlerListSuccess(t *testing.T) {
	t.Parallel()
	moradorID := uuid.New()
	svc := &fakeInadimplenteService{
		listResult: []dto.InadimplenteResponseDTO{
			{
				Morador: dto.UserResponseDTO{
					ID:       moradorID,
					FullName: "João Silva",
					Email:    "joao@test.com",
					Phone:    "11999990000",
					Status:   "ATIVO",
					Role:     "MORADOR",
				},
				TotalOverdue: 500.00,
				Payments: []dto.PagamentoResumoDTO{
					{
						ID:         uuid.New(),
						Valor:      500.00,
						Vencimento: mustParseTime("2026-06-01"),
						Status:     "ATRASADO",
					},
				},
			},
		},
	}
	handler := newInadimplenteHandler(svc)
	router := setupInadimplenteRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/sindico/inadimplentes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusOK)
	}

	var result []dto.InadimplenteResponseDTO
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 inadimplente, got %d", len(result))
	}
	if result[0].Morador.FullName != "João Silva" {
		t.Errorf("name: got %q, want %q", result[0].Morador.FullName, "João Silva")
	}
	if result[0].TotalOverdue != 500.00 {
		t.Errorf("total: got %v, want 500.00", result[0].TotalOverdue)
	}
}

func TestInadimplenteHandlerListServiceError(t *testing.T) {
	t.Parallel()
	svc := &fakeInadimplenteService{listErr: errors.New("db failure")}
	handler := newInadimplenteHandler(svc)
	router := setupInadimplenteRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/sindico/inadimplentes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestInadimplenteHandlerListEmpty(t *testing.T) {
	t.Parallel()
	svc := &fakeInadimplenteService{listResult: []dto.InadimplenteResponseDTO{}}
	handler := newInadimplenteHandler(svc)
	router := setupInadimplenteRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/sindico/inadimplentes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusOK)
	}

	var result []dto.InadimplenteResponseDTO
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}

func mustParseTime(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}
