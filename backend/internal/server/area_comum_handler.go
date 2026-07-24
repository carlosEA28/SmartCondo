package server

import (
	"context"
	"net/http"

	"github.com/carlosEA28/smartcondo/internal/dto"
	"github.com/gin-gonic/gin"
)

type areaComumService interface {
	ListAreas(ctx context.Context) ([]dto.AreaComumResponse, error)
}

type areaComumHandler struct {
	service areaComumService
}

func newAreaComumHandler(service areaComumService) *areaComumHandler {
	return &areaComumHandler{service: service}
}

func (h *areaComumHandler) listAreas(c *gin.Context) {
	areas, err := h.service.ListAreas(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list areas"})
		return
	}

	c.JSON(http.StatusOK, areas)
}
