package middleware

import (
	"errors"
	"net/http"

	"github.com/carlosEA28/smartcondo/internal/apperrors"
	"github.com/carlosEA28/smartcondo/internal/models"
	"github.com/carlosEA28/smartcondo/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequireSindicoRole(userRepo repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetHeader("X-User-ID")
		if userIDStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrMissingAuthHeader.Error()})
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
			return
		}

		user, err := userRepo.FindByID(c.Request.Context(), userID)
		if err != nil {
			if errors.Is(err, apperrors.ErrUserNotFound) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUnauthorizedSindico.Error()})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "authentication failed"})
			return
		}

		if user.Role != models.RoleSindico {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": apperrors.ErrUnauthorizedSindico.Error()})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
