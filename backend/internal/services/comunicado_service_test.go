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

// --- fakes for ComunicadoRepository ---

type fakeComunicadoRepository struct {
	findByIDResult *models.Comunicado
	findByIDErr    error
	findAllResult  []models.Comunicado
	findAllErr     error
	createErr      error
	deleteErr      error
	createdComunicado *models.Comunicado
}

func (f *fakeComunicadoRepository) FindByID(_ context.Context, id uuid.UUID) (*models.Comunicado, error) {
	if f.findByIDResult != nil {
		return f.findByIDResult, f.findByIDErr
	}
	return nil, f.findByIDErr
}

func (f *fakeComunicadoRepository) FindAll(context.Context) ([]models.Comunicado, error) {
	return f.findAllResult, f.findAllErr
}

func (f *fakeComunicadoRepository) Create(_ context.Context, c *models.Comunicado) error {
	f.createdComunicado = c
	return f.createErr
}

func (f *fakeComunicadoRepository) Delete(context.Context, uuid.UUID, uuid.UUID) error {
	return f.deleteErr
}

// --- PublishComunicado tests ---

func TestComunicadoServicePublishSindicoNotFound(t *testing.T) {
	comunicadoRepo := &fakeComunicadoRepository{}
	userRepo := &fakeUserRepository{findByIDErr: apperrors.ErrUserNotFound}
	service := NewComunicadoService(comunicadoRepo, userRepo)

	_, err := service.PublishComunicado(context.Background(), uuid.New(), &dto.CreateComunicadoDTO{
		Titulo:    "Aviso Teste",
		Descricao: "Descrição do aviso",
	})
	if !errors.Is(err, apperrors.ErrUnauthorizedSindico) {
		t.Fatalf("PublishComunicado() error = %v, want %v", err, apperrors.ErrUnauthorizedSindico)
	}
}

func TestComunicadoServicePublishSindicoGenericError(t *testing.T) {
	comunicadoRepo := &fakeComunicadoRepository{}
	userRepo := &fakeUserRepository{findByIDErr: errors.New("database error")}
	service := NewComunicadoService(comunicadoRepo, userRepo)

	_, err := service.PublishComunicado(context.Background(), uuid.New(), &dto.CreateComunicadoDTO{
		Titulo:    "Aviso Teste",
		Descricao: "Descrição do aviso",
	})
	if err == nil {
		t.Fatal("PublishComunicado() expected error, got nil")
	}
	if errors.Is(err, apperrors.ErrUnauthorizedSindico) {
		t.Fatal("PublishComunicado() should not return ErrUnauthorizedSindico for generic error")
	}
}

func TestComunicadoServicePublishCreateError(t *testing.T) {
	comunicadoRepo := &fakeComunicadoRepository{createErr: errors.New("database write failed")}
	userRepo := &fakeUserRepository{findByIDResult: &models.User{
		ID:       uuid.New(),
		FullName: "Síndico Teste",
	}}
	service := NewComunicadoService(comunicadoRepo, userRepo)

	_, err := service.PublishComunicado(context.Background(), uuid.New(), &dto.CreateComunicadoDTO{
		Titulo:    "Aviso Teste",
		Descricao: "Descrição do aviso",
	})
	if err == nil {
		t.Fatal("PublishComunicado() expected error, got nil")
	}
}

func TestComunicadoServicePublishSuccess(t *testing.T) {
	sindicoID := uuid.New()
	comunicadoID := uuid.New()
	now := time.Now()

	comunicadoRepo := &fakeComunicadoRepository{
		findByIDResult: &models.Comunicado{
			ID:             comunicadoID,
			Titulo:         "Aviso Teste",
			Descricao:      "Descrição do aviso",
			DataPublicacao: now,
			SindicoID:      sindicoID,
			Sindico:        models.User{ID: sindicoID, FullName: "Síndico Geral"},
		},
	}
	userRepo := &fakeUserRepository{findByIDResult: &models.User{
		ID:       sindicoID,
		FullName: "Síndico Geral",
	}}
	service := NewComunicadoService(comunicadoRepo, userRepo)

	response, err := service.PublishComunicado(context.Background(), sindicoID, &dto.CreateComunicadoDTO{
		Titulo:    "Aviso Teste",
		Descricao: "Descrição do aviso",
	})
	if err != nil {
		t.Fatalf("PublishComunicado() error = %v", err)
	}
	if response.Titulo != "Aviso Teste" {
		t.Fatalf("PublishComunicado() Titulo = %q, want %q", response.Titulo, "Aviso Teste")
	}
	if response.SindicoID != sindicoID {
		t.Fatalf("PublishComunicado() SindicoID = %v, want %v", response.SindicoID, sindicoID)
	}
	if response.SindicoNome != "Síndico Geral" {
		t.Fatalf("PublishComunicado() SindicoNome = %q, want %q", response.SindicoNome, "Síndico Geral")
	}
}

func TestComunicadoServicePublishFindByIDAfterCreateError(t *testing.T) {
	sindicoID := uuid.New()
	comunicadoRepo := &fakeComunicadoRepository{
		findByIDErr: errors.New("database error"),
	}
	userRepo := &fakeUserRepository{findByIDResult: &models.User{
		ID:       sindicoID,
		FullName: "Síndico Geral",
	}}
	service := NewComunicadoService(comunicadoRepo, userRepo)

	_, err := service.PublishComunicado(context.Background(), sindicoID, &dto.CreateComunicadoDTO{
		Titulo:    "Aviso Teste",
		Descricao: "Descrição do aviso",
	})
	if err == nil {
		t.Fatal("PublishComunicado() expected error, got nil")
	}
}

// --- ListComunicados tests ---

func TestComunicadoServiceListSuccess(t *testing.T) {
	comunicadoRepo := &fakeComunicadoRepository{
		findAllResult: []models.Comunicado{
			{
				ID:         uuid.New(),
				Titulo:     "Aviso 1",
				Descricao:  "Desc 1",
				SindicoID:  uuid.New(),
				Sindico:    models.User{FullName: "Síndico 1"},
			},
			{
				ID:         uuid.New(),
				Titulo:     "Aviso 2",
				Descricao:  "Desc 2",
				SindicoID:  uuid.New(),
				Sindico:    models.User{FullName: "Síndico 2"},
			},
		},
	}
	userRepo := &fakeUserRepository{}
	service := NewComunicadoService(comunicadoRepo, userRepo)

	response, err := service.ListComunicados(context.Background())
	if err != nil {
		t.Fatalf("ListComunicados() error = %v", err)
	}
	if len(response) != 2 {
		t.Fatalf("ListComunicados() length = %d, want 2", len(response))
	}
}

func TestComunicadoServiceListError(t *testing.T) {
	comunicadoRepo := &fakeComunicadoRepository{findAllErr: errors.New("database error")}
	userRepo := &fakeUserRepository{}
	service := NewComunicadoService(comunicadoRepo, userRepo)

	_, err := service.ListComunicados(context.Background())
	if err == nil {
		t.Fatal("ListComunicados() expected error, got nil")
	}
}

func TestComunicadoServiceListEmpty(t *testing.T) {
	comunicadoRepo := &fakeComunicadoRepository{findAllResult: []models.Comunicado{}}
	userRepo := &fakeUserRepository{}
	service := NewComunicadoService(comunicadoRepo, userRepo)

	response, err := service.ListComunicados(context.Background())
	if err != nil {
		t.Fatalf("ListComunicados() error = %v", err)
	}
	if len(response) != 0 {
		t.Fatalf("ListComunicados() length = %d, want 0", len(response))
	}
}

// --- GetComunicado tests ---

func TestComunicadoServiceGetSuccess(t *testing.T) {
	id := uuid.New()
	comunicadoRepo := &fakeComunicadoRepository{
		findByIDResult: &models.Comunicado{
			ID:        id,
			Titulo:    "Aviso Teste",
			Descricao: "Descrição",
			SindicoID: uuid.New(),
			Sindico:   models.User{FullName: "Síndico"},
		},
	}
	userRepo := &fakeUserRepository{}
	service := NewComunicadoService(comunicadoRepo, userRepo)

	response, err := service.GetComunicado(context.Background(), id)
	if err != nil {
		t.Fatalf("GetComunicado() error = %v", err)
	}
	if response.ID != id {
		t.Fatalf("GetComunicado() ID = %v, want %v", response.ID, id)
	}
}

func TestComunicadoServiceGetNotFound(t *testing.T) {
	comunicadoRepo := &fakeComunicadoRepository{findByIDErr: apperrors.ErrComunicadoNotFound}
	userRepo := &fakeUserRepository{}
	service := NewComunicadoService(comunicadoRepo, userRepo)

	_, err := service.GetComunicado(context.Background(), uuid.New())
	if !errors.Is(err, apperrors.ErrComunicadoNotFound) {
		t.Fatalf("GetComunicado() error = %v, want %v", err, apperrors.ErrComunicadoNotFound)
	}
}

func TestComunicadoServiceGetGenericError(t *testing.T) {
	comunicadoRepo := &fakeComunicadoRepository{findByIDErr: errors.New("database error")}
	userRepo := &fakeUserRepository{}
	service := NewComunicadoService(comunicadoRepo, userRepo)

	_, err := service.GetComunicado(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("GetComunicado() expected error, got nil")
	}
	if errors.Is(err, apperrors.ErrComunicadoNotFound) {
		t.Fatal("GetComunicado() should not return ErrComunicadoNotFound for generic error")
	}
}

// --- DeleteComunicado tests ---

func TestComunicadoServiceDeleteSuccess(t *testing.T) {
	comunicadoRepo := &fakeComunicadoRepository{}
	userRepo := &fakeUserRepository{}
	service := NewComunicadoService(comunicadoRepo, userRepo)

	err := service.DeleteComunicado(context.Background(), uuid.New(), uuid.New())
	if err != nil {
		t.Fatalf("DeleteComunicado() error = %v", err)
	}
}

func TestComunicadoServiceDeleteNotFound(t *testing.T) {
	comunicadoRepo := &fakeComunicadoRepository{deleteErr: apperrors.ErrComunicadoNotFound}
	userRepo := &fakeUserRepository{}
	service := NewComunicadoService(comunicadoRepo, userRepo)

	err := service.DeleteComunicado(context.Background(), uuid.New(), uuid.New())
	if !errors.Is(err, apperrors.ErrComunicadoNotFound) {
		t.Fatalf("DeleteComunicado() error = %v, want %v", err, apperrors.ErrComunicadoNotFound)
	}
}

func TestComunicadoServiceDeleteNotOwner(t *testing.T) {
	comunicadoRepo := &fakeComunicadoRepository{deleteErr: apperrors.ErrComunicadoNotOwner}
	userRepo := &fakeUserRepository{}
	service := NewComunicadoService(comunicadoRepo, userRepo)

	err := service.DeleteComunicado(context.Background(), uuid.New(), uuid.New())
	if !errors.Is(err, apperrors.ErrComunicadoNotOwner) {
		t.Fatalf("DeleteComunicado() error = %v, want %v", err, apperrors.ErrComunicadoNotOwner)
	}
}

func TestComunicadoServiceDeleteGenericError(t *testing.T) {
	comunicadoRepo := &fakeComunicadoRepository{deleteErr: errors.New("database error")}
	userRepo := &fakeUserRepository{}
	service := NewComunicadoService(comunicadoRepo, userRepo)

	err := service.DeleteComunicado(context.Background(), uuid.New(), uuid.New())
	if err == nil {
		t.Fatal("DeleteComunicado() expected error, got nil")
	}
	if errors.Is(err, apperrors.ErrComunicadoNotFound) || errors.Is(err, apperrors.ErrComunicadoNotOwner) {
		t.Fatal("DeleteComunicado() should not return domain error for generic error")
	}
}

// --- comunicadoToResponse tests ---

func TestComunicadoToResponseEmptySindicoName(t *testing.T) {
	c := &models.Comunicado{
		ID:        uuid.New(),
		Titulo:    "Teste",
		Descricao: "Desc",
		SindicoID: uuid.New(),
		Sindico:   models.User{FullName: ""},
	}

	response := comunicadoToResponse(c)
	if response.SindicoNome != "" {
		t.Fatalf("comunicadoToResponse() SindicoNome = %q, want empty", response.SindicoNome)
	}
}

func TestComunicadoToResponseWithSindicoName(t *testing.T) {
	c := &models.Comunicado{
		ID:        uuid.New(),
		Titulo:    "Teste",
		Descricao: "Desc",
		SindicoID: uuid.New(),
		Sindico:   models.User{FullName: "Síndico Geral"},
	}

	response := comunicadoToResponse(c)
	if response.SindicoNome != "Síndico Geral" {
		t.Fatalf("comunicadoToResponse() SindicoNome = %q, want %q", response.SindicoNome, "Síndico Geral")
	}
}
