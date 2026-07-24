package services

import (
	"context"
	"errors"
	"testing"

	"github.com/carlosEA28/smartcondo/internal/models"
	"github.com/google/uuid"
)

// --- fake for AreaComumRepository ---

type fakeAreaComumRepository struct {
	findByIDResult *models.AreaComum
	findByIDErr    error
	findAllResult  []models.AreaComum
	findAllErr     error
}

func (f *fakeAreaComumRepository) FindByID(_ context.Context, _ uuid.UUID) (*models.AreaComum, error) {
	return f.findByIDResult, f.findByIDErr
}

func (f *fakeAreaComumRepository) FindAll(_ context.Context) ([]models.AreaComum, error) {
	return f.findAllResult, f.findAllErr
}

// --- ListAreas tests ---

func TestAreaComumServiceListEmpty(t *testing.T) {
	t.Parallel()
	repo := &fakeAreaComumRepository{findAllResult: []models.AreaComum{}}
	svc := NewAreaComumService(repo)

	result, err := svc.ListAreas(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}

func TestAreaComumServiceListSuccess(t *testing.T) {
	t.Parallel()
	repo := &fakeAreaComumRepository{
		findAllResult: []models.AreaComum{
			{ID: uuid.New(), Nome: "Piscina", Descricao: "Piscina climatizada", Capacidade: 20},
			{ID: uuid.New(), Nome: "Salão de Festas", Descricao: "Salão principal", Capacidade: 50},
		},
	}
	svc := NewAreaComumService(repo)

	result, err := svc.ListAreas(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 areas, got %d", len(result))
	}
	if result[0].Nome != "Piscina" {
		t.Errorf("first area name: got %q, want %q", result[0].Nome, "Piscina")
	}
	if result[1].Capacidade != 50 {
		t.Errorf("second area capacity: got %d, want 50", result[1].Capacidade)
	}
}

func TestAreaComumServiceListRepoError(t *testing.T) {
	t.Parallel()
	repo := &fakeAreaComumRepository{findAllErr: errors.New("db down")}
	svc := NewAreaComumService(repo)

	_, err := svc.ListAreas(context.Background())

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
