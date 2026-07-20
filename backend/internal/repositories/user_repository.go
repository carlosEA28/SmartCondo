package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/carlosEA28/smartcondo/internal/apperrors"
	"github.com/carlosEA28/smartcondo/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	List(ctx context.Context) ([]models.User, error)
	Create(ctx context.Context, user *models.User, apartment *models.Apartment) error
	Save(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Apartment").First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	return &user, nil
}

func (r *GormUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	return &user, nil
}

func (r *GormUserRepository) List(ctx context.Context) ([]models.User, error) {
	users := make([]models.User, 0)
	if err := r.db.WithContext(ctx).Preload("Apartment").Order("nome ASC").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	return users, nil
}

func (r *GormUserRepository) Save(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", user.ID).Updates(map[string]any{
		"nome":     user.FullName,
		"telefone": user.Phone,
	})
	if result.Error != nil {
		return fmt.Errorf("save user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}

func (r *GormUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrForeignKeyViolated) {
			return apperrors.ErrUserInUse
		}
		return fmt.Errorf("delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}

func (r *GormUserRepository) Create(ctx context.Context, user *models.User, apartment *models.Apartment) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if apartment != nil {
			if err := tx.Create(apartment).Error; err != nil {
				return fmt.Errorf("create apartment: %w", err)
			}
			user.ApartmentID = &apartment.ID
			user.Apartment = apartment
		}

		if err := tx.Omit("Apartment").Create(user).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return apperrors.ErrUserAlreadyExists
			}
			return fmt.Errorf("create user: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("create user transaction: %w", err)
	}

	return nil
}
