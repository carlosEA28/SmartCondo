package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/carlosEA28/smartcondo/internal/apperrors"
	"github.com/carlosEA28/smartcondo/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type comunicadoService interface {
	PublishComunicado(ctx context.Context, sindicoID uuid.UUID, input *dto.CreateComunicadoDTO) (*dto.ComunicadoResponseDTO, error)
	ListComunicados(ctx context.Context) ([]dto.ComunicadoResponseDTO, error)
	GetComunicado(ctx context.Context, id uuid.UUID) (*dto.ComunicadoResponseDTO, error)
	DeleteComunicado(ctx context.Context, id uuid.UUID, sindicoID uuid.UUID) error
}

type comunicadoHandler struct {
	service comunicadoService
}

func newComunicadoHandler(service comunicadoService) *comunicadoHandler {
	return &comunicadoHandler{service: service}
}

func (h *comunicadoHandler) create(c *gin.Context) {
	var input dto.CreateComunicadoDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	sindicoID := c.MustGet("user_id").(uuid.UUID)

	comunicado, err := h.service.PublishComunicado(c.Request.Context(), sindicoID, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish comunicado"})
		return
	}

	c.JSON(http.StatusCreated, comunicado)
}

func (h *comunicadoHandler) list(c *gin.Context) {
	comunicados, err := h.service.ListComunicados(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list comunicados"})
		return
	}

	c.JSON(http.StatusOK, comunicados)
}

func (h *comunicadoHandler) getByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comunicado id"})
		return
	}

	comunicado, err := h.service.GetComunicado(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrComunicadoNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get comunicado"})
		return
	}

	c.JSON(http.StatusOK, comunicado)
}

func (h *comunicadoHandler) delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comunicado id"})
		return
	}

	sindicoID := c.MustGet("user_id").(uuid.UUID)

	if err := h.service.DeleteComunicado(c.Request.Context(), id, sindicoID); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrComunicadoNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, apperrors.ErrComunicadoNotOwner):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete comunicado"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
