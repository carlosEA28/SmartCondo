package services

import (
	"context"
	"errors"
	"testing"

	"github.com/carlosEA28/smartcondo/internal/apperrors"
	"github.com/carlosEA28/smartcondo/internal/dto"
	"github.com/carlosEA28/smartcondo/internal/models"
	"github.com/google/uuid"
)

// --- fakes for VisitRepository ---

type fakeVisitRepository struct {
	createErr error
}

func (f *fakeVisitRepository) Create(context.Context, *models.Visit) error {
	return f.createErr
}

// --- fakes for VisitorRepository (reused from visitor_service_test if same package) ---

type fakeVisitorRepo struct {
	findByIDResult    *models.Visitor
	findByIDErr       error
	findByCPFResult   *models.Visitor
	findByCPFErr      error
	listResult        []models.Visitor
	listErr           error
	createdVisitor    *models.Visitor
	createErr         error
	deleteErr         error
	searchResult      []models.Visitor
	searchErr         error
	updateLiberadoErr error
}

func (f *fakeVisitorRepo) FindByID(context.Context, uuid.UUID) (*models.Visitor, error) {
	return f.findByIDResult, f.findByIDErr
}

func (f *fakeVisitorRepo) FindByCPF(context.Context, string) (*models.Visitor, error) {
	return f.findByCPFResult, f.findByCPFErr
}

func (f *fakeVisitorRepo) List(context.Context) ([]models.Visitor, error) {
	return f.listResult, f.listErr
}

func (f *fakeVisitorRepo) Create(_ context.Context, v *models.Visitor) error {
	f.createdVisitor = v
	return f.createErr
}

func (f *fakeVisitorRepo) Delete(context.Context, uuid.UUID) error {
	return f.deleteErr
}

func (f *fakeVisitorRepo) Search(_ context.Context, _ *dto.VisitorFilterDTO) ([]models.Visitor, error) {
	return f.searchResult, f.searchErr
}

func (f *fakeVisitorRepo) UpdateLiberado(context.Context, uuid.UUID) error {
	return f.updateLiberadoErr
}

// --- fakes for UserRepository ---

type fakeUserRepoForPorteiro struct {
	findByIDResult *models.User
	findByIDErr    error
}

func (f *fakeUserRepoForPorteiro) FindByID(context.Context, uuid.UUID) (*models.User, error) {
	return f.findByIDResult, f.findByIDErr
}

func (f *fakeUserRepoForPorteiro) FindByEmail(context.Context, string) (*models.User, error) {
	return nil, nil
}

func (f *fakeUserRepoForPorteiro) FindApartmentByNumberAndBlock(context.Context, int, string) (*models.Apartment, error) {
	return nil, nil
}

func (f *fakeUserRepoForPorteiro) List(context.Context) ([]models.User, error) {
	return nil, nil
}

func (f *fakeUserRepoForPorteiro) Create(context.Context, *models.User, *models.Apartment) error {
	return nil
}

func (f *fakeUserRepoForPorteiro) Save(context.Context, *models.User) error {
	return nil
}

func (f *fakeUserRepoForPorteiro) Delete(context.Context, uuid.UUID) error {
	return nil
}

// --- SearchVisitors tests ---

func TestPorteiroServiceSearchVisitorsEmptyFilterReturnsError(t *testing.T) {
	visitorRepo := &fakeVisitorRepo{}
	visitRepo := &fakeVisitRepository{}
	userRepo := &fakeUserRepoForPorteiro{}
	service := NewPorteiroService(nil, visitorRepo, visitRepo, userRepo)

	filter := &dto.VisitorFilterDTO{}
	_, err := service.SearchVisitors(context.Background(), filter)
	if !errors.Is(err, apperrors.ErrFilterRequired) {
		t.Fatalf("SearchVisitors() error = %v, want %v", err, apperrors.ErrFilterRequired)
	}
}

func TestPorteiroServiceSearchVisitorsByNome(t *testing.T) {
	visitorRepo := &fakeVisitorRepo{
		searchResult: []models.Visitor{
			{ID: uuid.New(), Name: "João Silva", CPF: "12345678901", Phone: "11999999999"},
		},
	}
	visitRepo := &fakeVisitRepository{}
	userRepo := &fakeUserRepoForPorteiro{}
	service := NewPorteiroService(nil, visitorRepo, visitRepo, userRepo)

	filter := &dto.VisitorFilterDTO{Nome: "João"}
	response, err := service.SearchVisitors(context.Background(), filter)
	if err != nil {
		t.Fatalf("SearchVisitors() error = %v", err)
	}
	if len(response) != 1 {
		t.Fatalf("SearchVisitors() length = %d, want 1", len(response))
	}
	if response[0].Name != "João Silva" {
		t.Fatalf("SearchVisitors() Name = %q, want %q", response[0].Name, "João Silva")
	}
}

func TestPorteiroServiceSearchVisitorsByCPF(t *testing.T) {
	visitorRepo := &fakeVisitorRepo{
		searchResult: []models.Visitor{
			{ID: uuid.New(), Name: "Maria Santos", CPF: "98765432100", Phone: "11888888888"},
		},
	}
	visitRepo := &fakeVisitRepository{}
	userRepo := &fakeUserRepoForPorteiro{}
	service := NewPorteiroService(nil, visitorRepo, visitRepo, userRepo)

	filter := &dto.VisitorFilterDTO{CPF: "98765432100"}
	response, err := service.SearchVisitors(context.Background(), filter)
	if err != nil {
		t.Fatalf("SearchVisitors() error = %v", err)
	}
	if len(response) != 1 {
		t.Fatalf("SearchVisitors() length = %d, want 1", len(response))
	}
}

func TestPorteiroServiceSearchVisitorsByTelefone(t *testing.T) {
	visitorRepo := &fakeVisitorRepo{
		searchResult: []models.Visitor{
			{ID: uuid.New(), Name: "Pedro Lima", CPF: "11122233344", Phone: "11777777777"},
		},
	}
	visitRepo := &fakeVisitRepository{}
	userRepo := &fakeUserRepoForPorteiro{}
	service := NewPorteiroService(nil, visitorRepo, visitRepo, userRepo)

	filter := &dto.VisitorFilterDTO{Telefone: "777"}
	response, err := service.SearchVisitors(context.Background(), filter)
	if err != nil {
		t.Fatalf("SearchVisitors() error = %v", err)
	}
	if len(response) != 1 {
		t.Fatalf("SearchVisitors() length = %d, want 1", len(response))
	}
}

func TestPorteiroServiceSearchVisitorsByLiberado(t *testing.T) {
	liberado := true
	visitorRepo := &fakeVisitorRepo{
		searchResult: []models.Visitor{
			{ID: uuid.New(), Name: "Ana Costa", CPF: "55566677788", Phone: "11666666666", Liberado: true},
		},
	}
	visitRepo := &fakeVisitRepository{}
	userRepo := &fakeUserRepoForPorteiro{}
	service := NewPorteiroService(nil, visitorRepo, visitRepo, userRepo)

	filter := &dto.VisitorFilterDTO{Liberado: &liberado}
	response, err := service.SearchVisitors(context.Background(), filter)
	if err != nil {
		t.Fatalf("SearchVisitors() error = %v", err)
	}
	if len(response) != 1 {
		t.Fatalf("SearchVisitors() length = %d, want 1", len(response))
	}
	if !response[0].Liberado {
		t.Fatal("SearchVisitors() Liberado = false, want true")
	}
}

func TestPorteiroServiceSearchVisitorsRepoError(t *testing.T) {
	visitorRepo := &fakeVisitorRepo{searchErr: errors.New("database error")}
	visitRepo := &fakeVisitRepository{}
	userRepo := &fakeUserRepoForPorteiro{}
	service := NewPorteiroService(nil, visitorRepo, visitRepo, userRepo)

	filter := &dto.VisitorFilterDTO{Nome: "João"}
	_, err := service.SearchVisitors(context.Background(), filter)
	if err == nil {
		t.Fatal("SearchVisitors() expected error, got nil")
	}
}

func TestPorteiroServiceSearchVisitorsEmptyResult(t *testing.T) {
	visitorRepo := &fakeVisitorRepo{searchResult: []models.Visitor{}}
	visitRepo := &fakeVisitRepository{}
	userRepo := &fakeUserRepoForPorteiro{}
	service := NewPorteiroService(nil, visitorRepo, visitRepo, userRepo)

	filter := &dto.VisitorFilterDTO{Nome: "NotExist"}
	response, err := service.SearchVisitors(context.Background(), filter)
	if err != nil {
		t.Fatalf("SearchVisitors() error = %v", err)
	}
	if len(response) != 0 {
		t.Fatalf("SearchVisitors() length = %d, want 0", len(response))
	}
}

// --- ReleaseVisitor tests ---

func TestPorteiroServiceReleaseVisitorNotFound(t *testing.T) {
	visitorRepo := &fakeVisitorRepo{findByIDErr: apperrors.ErrVisitorNotFound}
	visitRepo := &fakeVisitRepository{}
	userRepo := &fakeUserRepoForPorteiro{}
	service := NewPorteiroService(nil, visitorRepo, visitRepo, userRepo)

	_, err := service.ReleaseVisitor(context.Background(), uuid.New(), uuid.New(), nil)
	if !errors.Is(err, apperrors.ErrVisitorNotFound) {
		t.Fatalf("ReleaseVisitor() error = %v, want %v", err, apperrors.ErrVisitorNotFound)
	}
}

func TestPorteiroServiceReleaseVisitorFindByIDGenericError(t *testing.T) {
	visitorRepo := &fakeVisitorRepo{findByIDErr: errors.New("database error")}
	visitRepo := &fakeVisitRepository{}
	userRepo := &fakeUserRepoForPorteiro{}
	service := NewPorteiroService(nil, visitorRepo, visitRepo, userRepo)

	_, err := service.ReleaseVisitor(context.Background(), uuid.New(), uuid.New(), nil)
	if err == nil {
		t.Fatal("ReleaseVisitor() expected error, got nil")
	}
	if errors.Is(err, apperrors.ErrVisitorNotFound) {
		t.Fatal("ReleaseVisitor() should not return ErrVisitorNotFound for generic error")
	}
}

func TestPorteiroServiceReleasePorteiroNotFound(t *testing.T) {
	visitorRepo := &fakeVisitorRepo{
		findByIDResult: &models.Visitor{
			ID:    uuid.New(),
			Name:  "João Silva",
			CPF:   "12345678901",
			Phone: "11999999999",
		},
	}
	visitRepo := &fakeVisitRepository{}
	userRepo := &fakeUserRepoForPorteiro{findByIDErr: apperrors.ErrUserNotFound}
	service := NewPorteiroService(nil, visitorRepo, visitRepo, userRepo)

	_, err := service.ReleaseVisitor(context.Background(), uuid.New(), uuid.New(), nil)
	if !errors.Is(err, apperrors.ErrPorteiroNotFound) {
		t.Fatalf("ReleaseVisitor() error = %v, want %v", err, apperrors.ErrPorteiroNotFound)
	}
}

func TestPorteiroServiceReleasePorteiroGenericError(t *testing.T) {
	visitorRepo := &fakeVisitorRepo{
		findByIDResult: &models.Visitor{
			ID:    uuid.New(),
			Name:  "João Silva",
			CPF:   "12345678901",
			Phone: "11999999999",
		},
	}
	visitRepo := &fakeVisitRepository{}
	userRepo := &fakeUserRepoForPorteiro{findByIDErr: errors.New("database error")}
	service := NewPorteiroService(nil, visitorRepo, visitRepo, userRepo)

	_, err := service.ReleaseVisitor(context.Background(), uuid.New(), uuid.New(), nil)
	if err == nil {
		t.Fatal("ReleaseVisitor() expected error, got nil")
	}
	if errors.Is(err, apperrors.ErrPorteiroNotFound) {
		t.Fatal("ReleaseVisitor() should not return ErrPorteiroNotFound for generic error")
	}
}

// Note: ReleaseVisitor tests that reach the db.Transaction call require an integration test.
// The pre-transaction logic (visitor/porteiro validation) is tested above.
// The transaction itself (UpdateLiberado + Create Visit) should be tested with a real DB.
