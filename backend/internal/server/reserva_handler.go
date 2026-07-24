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

type reservaService interface {
	CreateReserva(ctx context.Context, moradorID uuid.UUID, req dto.CreateReservaRequest) (*dto.ReservaResponse, error)
	ListMinhasReservas(ctx context.Context, moradorID uuid.UUID) ([]dto.ReservaResponse, error)
	CancelReserva(ctx context.Context, reservaID uuid.UUID, moradorID uuid.UUID) (*dto.ReservaResponse, error)
}

type reservaHandler struct {
	service reservaService
}

func newReservaHandler(service reservaService) *reservaHandler {
	return &reservaHandler{service: service}
}

func (h *reservaHandler) create(c *gin.Context) {
	moradorID := c.MustGet("user_id").(uuid.UUID)

	var req dto.CreateReservaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": apperrors.ErrInvalidReservaData.Error()})
		return
	}

	reserva, err := h.service.CreateReserva(c.Request.Context(), moradorID, req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrAreaComumNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, apperrors.ErrMoradorInadimplente):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, apperrors.ErrReservaConflito):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, apperrors.ErrInvalidReservaData):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create reservation"})
		}
		return
	}

	c.JSON(http.StatusCreated, reserva)
}

func (h *reservaHandler) listMinhas(c *gin.Context) {
	moradorID := c.MustGet("user_id").(uuid.UUID)

	reservas, err := h.service.ListMinhasReservas(c.Request.Context(), moradorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list reservations"})
		return
	}

	c.JSON(http.StatusOK, reservas)
}

func (h *reservaHandler) cancelar(c *gin.Context) {
	moradorID := c.MustGet("user_id").(uuid.UUID)

	reservaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reservation id"})
		return
	}

	reserva, err := h.service.CancelReserva(c.Request.Context(), reservaID, moradorID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrReservaNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, apperrors.ErrReservaNotOwner):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel reservation"})
		}
		return
	}

	c.JSON(http.StatusOK, reserva)
}
