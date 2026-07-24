package server

import (
	"context"
	"net/http"

	"github.com/carlosEA28/smartcondo/internal/dto"
	"github.com/gin-gonic/gin"
)

type inadimplenteService interface {
	ListInadimplentes(ctx context.Context) ([]dto.InadimplenteResponseDTO, error)
}

type inadimplenteHandler struct {
	service inadimplenteService
}

func newInadimplenteHandler(service inadimplenteService) *inadimplenteHandler {
	return &inadimplenteHandler{service: service}
}

func (h *inadimplenteHandler) list(c *gin.Context) {
	inadimplentes, err := h.service.ListInadimplentes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list inadimplentes"})
		return
	}

	c.JSON(http.StatusOK, inadimplentes)
}
