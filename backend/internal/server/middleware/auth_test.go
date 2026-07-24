package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carlosEA28/smartcondo/internal/apperrors"
	"github.com/carlosEA28/smartcondo/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type fakeUserRepo struct {
	findByIDResult *models.User
	findByIDErr    error
}

func (f *fakeUserRepo) FindByID(_ context.Context, _ uuid.UUID) (*models.User, error) {
	return f.findByIDResult, f.findByIDErr
}

func (f *fakeUserRepo) FindByEmail(_ context.Context, _ string) (*models.User, error) {
	return nil, nil
}

func (f *fakeUserRepo) FindApartmentByNumberAndBlock(_ context.Context, _ int, _ string) (*models.Apartment, error) {
	return nil, nil
}

func (f *fakeUserRepo) List(_ context.Context) ([]models.User, error) {
	return nil, nil
}

func (f *fakeUserRepo) Create(_ context.Context, _ *models.User, _ *models.Apartment) error {
	return nil
}

func (f *fakeUserRepo) Save(_ context.Context, _ *models.User) error {
	return nil
}

func (f *fakeUserRepo) Delete(_ context.Context, _ uuid.UUID) error {
	return nil
}

func setupMiddlewareRouter(userRepo *fakeUserRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequireSindicoRole(userRepo))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	return router
}

func TestRequireSindicoRoleMissingHeader(t *testing.T) {
	userRepo := &fakeUserRepo{}
	router := setupMiddlewareRouter(userRepo)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("RequireSindicoRole() status = %d, want %d", w.Code, http.StatusUnauthorized)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if response["error"] != apperrors.ErrMissingAuthHeader.Error() {
		t.Fatalf("RequireSindicoRole() error = %q, want %q", response["error"], apperrors.ErrMissingAuthHeader.Error())
	}
}

func TestRequireSindicoRoleInvalidUUID(t *testing.T) {
	userRepo := &fakeUserRepo{}
	router := setupMiddlewareRouter(userRepo)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("X-User-ID", "not-a-uuid")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("RequireSindicoRole() status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestRequireSindicoRoleUserNotFound(t *testing.T) {
	userRepo := &fakeUserRepo{findByIDErr: apperrors.ErrUserNotFound}
	router := setupMiddlewareRouter(userRepo)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("X-User-ID", uuid.New().String())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("RequireSindicoRole() status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestRequireSindicoRoleGenericError(t *testing.T) {
	userRepo := &fakeUserRepo{findByIDErr: errors.New("database error")}
	router := setupMiddlewareRouter(userRepo)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("X-User-ID", uuid.New().String())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("RequireSindicoRole() status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestRequireSindicoRoleWrongRole(t *testing.T) {
	userRepo := &fakeUserRepo{
		findByIDResult: &models.User{
			ID:   uuid.New(),
			Role: models.RoleMorador,
		},
	}
	router := setupMiddlewareRouter(userRepo)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("X-User-ID", uuid.New().String())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("RequireSindicoRole() status = %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestRequireSindicoRolePorteiroRole(t *testing.T) {
	userRepo := &fakeUserRepo{
		findByIDResult: &models.User{
			ID:   uuid.New(),
			Role: models.RolePorteiro,
		},
	}
	router := setupMiddlewareRouter(userRepo)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("X-User-ID", uuid.New().String())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("RequireSindicoRole() status = %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestRequireSindicoRoleSuccess(t *testing.T) {
	userID := uuid.New()
	userRepo := &fakeUserRepo{
		findByIDResult: &models.User{
			ID:   userID,
			Role: models.RoleSindico,
		},
	}
	router := setupMiddlewareRouter(userRepo)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("X-User-ID", userID.String())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("RequireSindicoRole() status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRequireSindicoRoleSetsUserID(t *testing.T) {
	userID := uuid.New()
	userRepo := &fakeUserRepo{
		findByIDResult: &models.User{
			ID:   userID,
			Role: models.RoleSindico,
		},
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	var capturedID interface{}

	router.Use(RequireSindicoRole(userRepo))
	router.GET("/check", func(c *gin.Context) {
		capturedID, _ = c.Get("user_id")
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/check", nil)
	req.Header.Set("X-User-ID", userID.String())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("RequireSindicoRole() status = %d, want %d", w.Code, http.StatusOK)
	}
	if capturedID == nil {
		t.Fatal("RequireSindicoRole() did not set user_id in context")
	}
	if id, ok := capturedID.(uuid.UUID); !ok || id != userID {
		t.Fatalf("RequireSindicoRole() user_id = %v, want %v", capturedID, userID)
	}
}
