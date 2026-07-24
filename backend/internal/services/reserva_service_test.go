package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/carlosEA28/smartcondo/internal/apperrors"
	"github.com/carlosEA28/smartcondo/internal/dto"
	"github.com/carlosEA28/smartcondo/internal/models"
	"github.com/google/uuid"
)

// --- fake for ReservaRepository ---

type fakeReservaRepository struct {
	findByIDResult       *models.Reserva
	findByIDErr          error
	findByMoradorResult  []models.Reserva
	findByMoradorErr     error
	findConflictingResult bool
	findConflictingErr   error
	createErr            error
	updateStatusErr      error
}

func (f *fakeReservaRepository) FindByID(_ context.Context, _ uuid.UUID) (*models.Reserva, error) {
	return f.findByIDResult, f.findByIDErr
}

func (f *fakeReservaRepository) FindByMorador(_ context.Context, _ uuid.UUID) ([]models.Reserva, error) {
	return f.findByMoradorResult, f.findByMoradorErr
}

func (f *fakeReservaRepository) FindConflicting(_ context.Context, _ uuid.UUID, _ string, _, _ string, _ *uuid.UUID) (bool, error) {
	return f.findConflictingResult, f.findConflictingErr
}

func (f *fakeReservaRepository) Create(_ context.Context, _ *models.Reserva) error {
	return f.createErr
}

func (f *fakeReservaRepository) UpdateStatus(_ context.Context, _ uuid.UUID, _ models.ReservaStatus) error {
	return f.updateStatusErr
}

// --- helpers ---

func testAreaComum() models.AreaComum {
	return models.AreaComum{
		ID:         uuid.New(),
		Nome:       "Piscina",
		Descricao:  "Piscina climatizada",
		Capacidade: 20,
	}
}

func testReserva(area models.AreaComum, moradorID uuid.UUID) models.Reserva {
	return models.Reserva{
		ID:          uuid.New(),
		Data:        time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC),
		HoraInicio:  "14:00",
		HoraFim:     "16:00",
		Status:      models.ReservaStatusPendente,
		MoradorID:   moradorID,
		AreaComumID: area.ID,
		AreaComum:   area,
	}
}

func testReservaService(areaRepo *fakeAreaComumRepository, reservaRepo *fakeReservaRepository, pagamentoRepo *fakePagamentoRepository) *ReservaService {
	return NewReservaService(reservaRepo, areaRepo, pagamentoRepo)
}

func testPagamentoRepo(overdue bool) *fakePagamentoRepository {
	return &fakePagamentoRepository{hasOverdueResult: overdue}
}

// --- CreateReserva tests ---

func TestReservaServiceCreateAreaNotFound(t *testing.T) {
	t.Parallel()
	areaRepo := &fakeAreaComumRepository{findByIDErr: apperrors.ErrAreaComumNotFound}
	reservaRepo := &fakeReservaRepository{}
	pagamentoRepo := testPagamentoRepo(false)
	svc := testReservaService(areaRepo, reservaRepo, pagamentoRepo)

	req := dto.CreateReservaRequest{
		AreaComumID: uuid.New(),
		Data:        "2026-08-15",
		HoraInicio:  "14:00",
		HoraFim:     "16:00",
	}

	_, err := svc.CreateReserva(context.Background(), uuid.New(), req)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, apperrors.ErrAreaComumNotFound) {
		t.Errorf("expected ErrAreaComumNotFound, got %v", err)
	}
}

func TestReservaServiceCreateAreaRepoError(t *testing.T) {
	t.Parallel()
	areaRepo := &fakeAreaComumRepository{findByIDErr: errors.New("db down")}
	reservaRepo := &fakeReservaRepository{}
	pagamentoRepo := testPagamentoRepo(false)
	svc := testReservaService(areaRepo, reservaRepo, pagamentoRepo)

	req := dto.CreateReservaRequest{
		AreaComumID: uuid.New(),
		Data:        "2026-08-15",
		HoraInicio:  "14:00",
		HoraFim:     "16:00",
	}

	_, err := svc.CreateReserva(context.Background(), uuid.New(), req)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestReservaServiceCreateInadimplente(t *testing.T) {
	t.Parallel()
	area := testAreaComum()
	areaRepo := &fakeAreaComumRepository{findByIDResult: &area}
	reservaRepo := &fakeReservaRepository{}
	pagamentoRepo := testPagamentoRepo(true) // has overdue
	svc := testReservaService(areaRepo, reservaRepo, pagamentoRepo)

	req := dto.CreateReservaRequest{
		AreaComumID: area.ID,
		Data:        "2026-08-15",
		HoraInicio:  "14:00",
		HoraFim:     "16:00",
	}

	_, err := svc.CreateReserva(context.Background(), uuid.New(), req)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, apperrors.ErrMoradorInadimplente) {
		t.Errorf("expected ErrMoradorInadimplente, got %v", err)
	}
}

func TestReservaServiceCreatePagamentoRepoError(t *testing.T) {
	t.Parallel()
	area := testAreaComum()
	areaRepo := &fakeAreaComumRepository{findByIDResult: &area}
	reservaRepo := &fakeReservaRepository{}
	pagamentoRepo := &fakePagamentoRepository{hasOverdueErr: errors.New("db down")}
	svc := testReservaService(areaRepo, reservaRepo, pagamentoRepo)

	req := dto.CreateReservaRequest{
		AreaComumID: area.ID,
		Data:        "2026-08-15",
		HoraInicio:  "14:00",
		HoraFim:     "16:00",
	}

	_, err := svc.CreateReserva(context.Background(), uuid.New(), req)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestReservaServiceCreateConflict(t *testing.T) {
	t.Parallel()
	area := testAreaComum()
	areaRepo := &fakeAreaComumRepository{findByIDResult: &area}
	reservaRepo := &fakeReservaRepository{findConflictingResult: true}
	pagamentoRepo := testPagamentoRepo(false)
	svc := testReservaService(areaRepo, reservaRepo, pagamentoRepo)

	req := dto.CreateReservaRequest{
		AreaComumID: area.ID,
		Data:        "2026-08-15",
		HoraInicio:  "14:00",
		HoraFim:     "16:00",
	}

	_, err := svc.CreateReserva(context.Background(), uuid.New(), req)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, apperrors.ErrReservaConflito) {
		t.Errorf("expected ErrReservaConflito, got %v", err)
	}
}

func TestReservaServiceCreateSuccess(t *testing.T) {
	t.Parallel()
	area := testAreaComum()
	moradorID := uuid.New()
	areaRepo := &fakeAreaComumRepository{findByIDResult: &area}
	reservaRepo := &fakeReservaRepository{}
	pagamentoRepo := testPagamentoRepo(false)
	svc := testReservaService(areaRepo, reservaRepo, pagamentoRepo)

	req := dto.CreateReservaRequest{
		AreaComumID: area.ID,
		Data:        "2026-08-15",
		HoraInicio:  "14:00",
		HoraFim:     "16:00",
	}

	result, err := svc.CreateReserva(context.Background(), moradorID, req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.Status != "PENDENTE" {
		t.Errorf("status: got %q, want %q", result.Status, "PENDENTE")
	}
	if result.MoradorID != moradorID {
		t.Errorf("moradorID: got %v, want %v", result.MoradorID, moradorID)
	}
	if result.AreaComum.Nome != "Piscina" {
		t.Errorf("area name: got %q, want %q", result.AreaComum.Nome, "Piscina")
	}
}

func TestReservaServiceCreateRepoError(t *testing.T) {
	t.Parallel()
	area := testAreaComum()
	areaRepo := &fakeAreaComumRepository{findByIDResult: &area}
	reservaRepo := &fakeReservaRepository{createErr: errors.New("db failure")}
	pagamentoRepo := testPagamentoRepo(false)
	svc := testReservaService(areaRepo, reservaRepo, pagamentoRepo)

	req := dto.CreateReservaRequest{
		AreaComumID: area.ID,
		Data:        "2026-08-15",
		HoraInicio:  "14:00",
		HoraFim:     "16:00",
	}

	_, err := svc.CreateReserva(context.Background(), uuid.New(), req)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- ListMinhasReservas tests ---

func TestReservaServiceListMinhasEmpty(t *testing.T) {
	t.Parallel()
	reservaRepo := &fakeReservaRepository{findByMoradorResult: []models.Reserva{}}
	areaRepo := &fakeAreaComumRepository{}
	pagamentoRepo := testPagamentoRepo(false)
	svc := testReservaService(areaRepo, reservaRepo, pagamentoRepo)

	result, err := svc.ListMinhasReservas(context.Background(), uuid.New())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}

func TestReservaServiceListMinhasSuccess(t *testing.T) {
	t.Parallel()
	moradorID := uuid.New()
	area := testAreaComum()
	r := testReserva(area, moradorID)
	reservaRepo := &fakeReservaRepository{findByMoradorResult: []models.Reserva{r}}
	areaRepo := &fakeAreaComumRepository{}
	pagamentoRepo := testPagamentoRepo(false)
	svc := testReservaService(areaRepo, reservaRepo, pagamentoRepo)

	result, err := svc.ListMinhasReservas(context.Background(), moradorID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 reserva, got %d", len(result))
	}
	if result[0].AreaComum.Nome != "Piscina" {
		t.Errorf("area name: got %q, want %q", result[0].AreaComum.Nome, "Piscina")
	}
}

func TestReservaServiceListMinhasRepoError(t *testing.T) {
	t.Parallel()
	reservaRepo := &fakeReservaRepository{findByMoradorErr: errors.New("db down")}
	areaRepo := &fakeAreaComumRepository{}
	pagamentoRepo := testPagamentoRepo(false)
	svc := testReservaService(areaRepo, reservaRepo, pagamentoRepo)

	_, err := svc.ListMinhasReservas(context.Background(), uuid.New())

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- CancelReserva tests ---

func TestReservaServiceCancelNotFound(t *testing.T) {
	t.Parallel()
	reservaRepo := &fakeReservaRepository{findByIDErr: apperrors.ErrReservaNotFound}
	areaRepo := &fakeAreaComumRepository{}
	pagamentoRepo := testPagamentoRepo(false)
	svc := testReservaService(areaRepo, reservaRepo, pagamentoRepo)

	_, err := svc.CancelReserva(context.Background(), uuid.New(), uuid.New())

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, apperrors.ErrReservaNotFound) {
		t.Errorf("expected ErrReservaNotFound, got %v", err)
	}
}

func TestReservaServiceCancelNotOwner(t *testing.T) {
	t.Parallel()
	area := testAreaComum()
	ownerID := uuid.New()
	otherID := uuid.New()
	r := testReserva(area, ownerID)
	reservaRepo := &fakeReservaRepository{findByIDResult: &r}
	areaRepo := &fakeAreaComumRepository{}
	pagamentoRepo := testPagamentoRepo(false)
	svc := testReservaService(areaRepo, reservaRepo, pagamentoRepo)

	_, err := svc.CancelReserva(context.Background(), r.ID, otherID)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, apperrors.ErrReservaNotOwner) {
		t.Errorf("expected ErrReservaNotOwner, got %v", err)
	}
}

func TestReservaServiceCancelSuccess(t *testing.T) {
	t.Parallel()
	area := testAreaComum()
	moradorID := uuid.New()
	r := testReserva(area, moradorID)
	reservaRepo := &fakeReservaRepository{findByIDResult: &r}
	areaRepo := &fakeAreaComumRepository{}
	pagamentoRepo := testPagamentoRepo(false)
	svc := testReservaService(areaRepo, reservaRepo, pagamentoRepo)

	result, err := svc.CancelReserva(context.Background(), r.ID, moradorID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != "CANCELADA" {
		t.Errorf("status: got %q, want %q", result.Status, "CANCELADA")
	}
}

func TestReservaServiceCancelRepoError(t *testing.T) {
	t.Parallel()
	area := testAreaComum()
	moradorID := uuid.New()
	r := testReserva(area, moradorID)
	reservaRepo := &fakeReservaRepository{
		findByIDResult: &r,
		updateStatusErr: errors.New("db failure"),
	}
	areaRepo := &fakeAreaComumRepository{}
	pagamentoRepo := testPagamentoRepo(false)
	svc := testReservaService(areaRepo, reservaRepo, pagamentoRepo)

	_, err := svc.CancelReserva(context.Background(), r.ID, moradorID)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
