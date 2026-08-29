package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/steven230500/cartagena-api/internal/domain"
	"github.com/steven230500/cartagena-api/internal/middleware"
	"github.com/steven230500/cartagena-api/internal/service"
)

type RouteProgressHandler struct {
	svc *service.RouteProgressService
}

func NewRouteProgressHandler(svc *service.RouteProgressService) *RouteProgressHandler {
	return &RouteProgressHandler{svc: svc}
}

func (h *RouteProgressHandler) Get(c *gin.Context) {
	userID := c.GetString(middleware.UserIDContextKey)

	progress, err := h.svc.Get(c.Request.Context(), userID, c.Param("route_id"))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "sin progreso guardado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, progress)
}

type updateRouteProgressRequest struct {
	CurrentStep int `json:"current_step"`
}

func (h *RouteProgressHandler) Upsert(c *gin.Context) {
	userID := c.GetString(middleware.UserIDContextKey)

	var req updateRouteProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cuerpo inválido"})
		return
	}

	progress, err := h.svc.Upsert(c.Request.Context(), userID, c.Param("route_id"), req.CurrentStep)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no se pudo guardar el progreso"})
		return
	}

	c.JSON(http.StatusOK, progress)
}
