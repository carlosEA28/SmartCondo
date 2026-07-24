package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/carlosEA28/smartcondo/internal/models"
	"github.com/google/uuid"
)

// --- fake for PagamentoRepository ---

type fakePagamentoRepository struct {
	findInadimplentesResult []models.Pagamento
	findInadimplentesErr    error
}

func (f *fakePagamentoRepository) FindInadimplentes(_ context.Context) ([]models.Pagamento, error) {
	return f.findInadimplentesResult, f.findInadimplentesErr
}

// --- helpers ---

func newTestMorador(id uuid.UUID, name, email, phone string, apt *models.Apartment) models.User {
	return models.User{
		ID:       id,
		FullName: name,
		Email:    email,
		Phone:    phone,
		Status:   models.UserStatusActive,
		Role:     models.RoleMorador,
		Apartment: apt,
	}
}

func newTestPagamento(id uuid.UUID, moradorID uuid.UUID, valor float64, vencimento time.Time) models.Pagamento {
	return models.Pagamento{
		ID:         id,
		Valor:      valor,
		Vencimento: vencimento,
		Status:     models.PaymentStatusOverdue,
		MoradorID:  moradorID,
	}
}

// --- ListInadimplentes tests ---

func TestInadimplenteServiceListEmpty(t *testing.T) {
	t.Parallel()
	repo := &fakePagamentoRepository{findInadimplentesResult: []models.Pagamento{}}
	svc := NewInadimplenteService(repo)

	result, err := svc.ListInadimplentes(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}

func TestInadimplenteServiceListRepoError(t *testing.T) {
	t.Parallel()
	repo := &fakePagamentoRepository{findInadimplentesErr: errors.New("db down")}
	svc := NewInadimplenteService(repo)

	_, err := svc.ListInadimplentes(context.Background())

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestInadimplenteServiceListSinglePayment(t *testing.T) {
	t.Parallel()
	moradorID := uuid.New()
	aptID := uuid.New()
	morador := newTestMorador(moradorID, "João", "joao@test.com", "11999990000", &models.Apartment{
		ID:     aptID,
		Number: 101,
		Block:  "A",
	})
	payment := newTestPagamento(uuid.New(), moradorID, 250.00, time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC))
	payment.Morador = morador

	repo := &fakePagamentoRepository{findInadimplentesResult: []models.Pagamento{payment}}
	svc := NewInadimplenteService(repo)

	result, err := svc.ListInadimplentes(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 inadimplente, got %d", len(result))
	}

	entry := result[0]
	if entry.Morador.ID != moradorID {
		t.Errorf("morador ID: got %v, want %v", entry.Morador.ID, moradorID)
	}
	if entry.Morador.FullName != "João" {
		t.Errorf("morador name: got %q, want %q", entry.Morador.FullName, "João")
	}
	if entry.TotalOverdue != 250.00 {
		t.Errorf("total overdue: got %v, want 250.00", entry.TotalOverdue)
	}
	if len(entry.Payments) != 1 {
		t.Fatalf("expected 1 payment, got %d", len(entry.Payments))
	}
	if entry.Payments[0].Valor != 250.00 {
		t.Errorf("payment valor: got %v, want 250.00", entry.Payments[0].Valor)
	}
	if entry.Morador.Apartment == nil {
		t.Fatal("expected apartment to be set")
	}
	if entry.Morador.Apartment.Number != 101 {
		t.Errorf("apartment number: got %d, want 101", entry.Morador.Apartment.Number)
	}
}

func TestInadimplenteServiceListMultiplePaymentsSameMorador(t *testing.T) {
	t.Parallel()
	moradorID := uuid.New()
	morador := newTestMorador(moradorID, "Maria", "maria@test.com", "11988887777", nil)

	p1 := newTestPagamento(uuid.New(), moradorID, 100.00, time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC))
	p1.Morador = morador
	p2 := newTestPagamento(uuid.New(), moradorID, 200.00, time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC))
	p2.Morador = morador

	repo := &fakePagamentoRepository{findInadimplentesResult: []models.Pagamento{p1, p2}}
	svc := NewInadimplenteService(repo)

	result, err := svc.ListInadimplentes(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 inadimplente, got %d", len(result))
	}

	entry := result[0]
	if entry.TotalOverdue != 300.00 {
		t.Errorf("total overdue: got %v, want 300.00", entry.TotalOverdue)
	}
	if len(entry.Payments) != 2 {
		t.Fatalf("expected 2 payments, got %d", len(entry.Payments))
	}
	if entry.Morador.Apartment != nil {
		t.Errorf("expected nil apartment, got %+v", entry.Morador.Apartment)
	}
}

func TestInadimplenteServiceListMultipleMoradores(t *testing.T) {
	t.Parallel()
	morador1ID := uuid.New()
	morador2ID := uuid.New()

	morador1 := newTestMorador(morador1ID, "Ana", "ana@test.com", "11911112222", nil)
	morador2 := newTestMorador(morador2ID, "Carlos", "carlos@test.com", "11933334444", nil)

	p1 := newTestPagamento(uuid.New(), morador1ID, 150.00, time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC))
	p1.Morador = morador1
	p2 := newTestPagamento(uuid.New(), morador2ID, 300.00, time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC))
	p2.Morador = morador2

	repo := &fakePagamentoRepository{findInadimplentesResult: []models.Pagamento{p1, p2}}
	svc := NewInadimplenteService(repo)

	result, err := svc.ListInadimplentes(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 inadimplentes, got %d", len(result))
	}
	if result[0].Morador.ID != morador1ID {
		t.Errorf("first inadimplente: got %v, want %v", result[0].Morador.ID, morador1ID)
	}
	if result[1].Morador.ID != morador2ID {
		t.Errorf("second inadimplente: got %v, want %v", result[1].Morador.ID, morador2ID)
	}
	if result[0].TotalOverdue != 150.00 {
		t.Errorf("first total: got %v, want 150.00", result[0].TotalOverdue)
	}
	if result[1].TotalOverdue != 300.00 {
		t.Errorf("second total: got %v, want 300.00", result[1].TotalOverdue)
	}
}

func TestInadimplenteServiceListPaymentWithApartment(t *testing.T) {
	t.Parallel()
	moradorID := uuid.New()
	aptID := uuid.New()
	morador := newTestMorador(moradorID, "Pedro", "pedro@test.com", "11955556666", &models.Apartment{
		ID:     aptID,
		Number: 302,
		Block:  "B",
	})

	payment := newTestPagamento(uuid.New(), moradorID, 500.00, time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC))
	payment.Morador = morador

	repo := &fakePagamentoRepository{findInadimplentesResult: []models.Pagamento{payment}}
	svc := NewInadimplenteService(repo)

	result, err := svc.ListInadimplentes(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Morador.Apartment == nil {
		t.Fatal("expected apartment, got nil")
	}
	if result[0].Morador.Apartment.ID != aptID {
		t.Errorf("apartment ID: got %v, want %v", result[0].Morador.Apartment.ID, aptID)
	}
	if result[0].Morador.Apartment.Number != 302 {
		t.Errorf("apartment number: got %d, want 302", result[0].Morador.Apartment.Number)
	}
	if result[0].Morador.Apartment.Block != "B" {
		t.Errorf("apartment block: got %q, want %q", result[0].Morador.Apartment.Block, "B")
	}
}
