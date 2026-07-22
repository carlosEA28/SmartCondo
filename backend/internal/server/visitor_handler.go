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

type visitorService interface {
	CreateVisitor(ctx context.Context, input *dto.CreateVisitorDTO) (*dto.VisitorResponseDTO, error)
	GetVisitor(ctx context.Context, id uuid.UUID) (*dto.VisitorResponseDTO, error)
	ListVisitors(ctx context.Context) ([]dto.VisitorResponseDTO, error)
}

type visitorHandler struct {
	service visitorService
}

func newVisitorHandler(service visitorService) *visitorHandler {
	return &visitorHandler{service: service}
}

func (h *visitorHandler) create(c *gin.Context) {
	var input dto.CreateVisitorDTO
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	file, err := c.FormFile("photo")
	if err == nil {
		input.Photo = file
	}

	visitor, err := h.service.CreateVisitor(c.Request.Context(), &input)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrVisitorAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create visitor"})
		}
		return
	}

	c.JSON(http.StatusCreated, visitor)
}

func (h *visitorHandler) getByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid visitor id"})
		return
	}

	visitor, err := h.service.GetVisitor(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrVisitorNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get visitor"})
		return
	}

	c.JSON(http.StatusOK, visitor)
}

func (h *visitorHandler) list(c *gin.Context) {
	visitors, err := h.service.ListVisitors(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list visitors"})
		return
	}

	c.JSON(http.StatusOK, visitors)
}
