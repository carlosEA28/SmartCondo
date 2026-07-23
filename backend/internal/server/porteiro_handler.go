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

type porteiroService interface {
	SearchVisitors(ctx context.Context, filter *dto.VisitorFilterDTO) ([]dto.VisitorResponseDTO, error)
	ReleaseVisitor(ctx context.Context, visitorID uuid.UUID, porteiroID uuid.UUID, moradorID *uuid.UUID) (*dto.VisitorResponseDTO, error)
}

type porteiroHandler struct {
	service porteiroService
}

func newPorteiroHandler(service porteiroService) *porteiroHandler {
	return &porteiroHandler{service: service}
}

func (h *porteiroHandler) search(c *gin.Context) {
	var filter dto.VisitorFilterDTO
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	visitors, err := h.service.SearchVisitors(c.Request.Context(), &filter)
	if err != nil {
		if errors.Is(err, apperrors.ErrFilterRequired) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search visitors"})
		return
	}

	c.JSON(http.StatusOK, visitors)
}

func (h *porteiroHandler) release(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid visitor id"})
		return
	}

	var req dto.ReleaseRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	visitor, err := h.service.ReleaseVisitor(c.Request.Context(), id, req.PorteiroID, req.MoradorID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrVisitorNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, apperrors.ErrPorteiroNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to release visitor"})
		}
		return
	}

	c.JSON(http.StatusOK, visitor)
}
