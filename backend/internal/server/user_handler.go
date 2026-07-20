package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/carlosEA28/smartcondo/internal/dto"
	"github.com/carlosEA28/smartcondo/internal/services"
	"github.com/gin-gonic/gin"
)

type userCreator interface {
	CreateUser(ctx context.Context, input dto.CreateUserDTO) (*dto.UserResponseDTO, error)
}

type userHandler struct {
	service userCreator
}

func newUserHandler(service userCreator) *userHandler {
	return &userHandler{service: service}
}

func (h *userHandler) create(c *gin.Context) {
	var input dto.CreateUserDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrApartmentRequired),
			errors.Is(err, services.ErrApartmentNotAllowed),
			errors.Is(err, services.ErrResponsibleNotAllowed):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}
