package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/steven230500/cartagena-api/internal/middleware"
	"github.com/steven230500/cartagena-api/internal/service"
)

type FavoriteHandler struct {
	svc *service.FavoriteService
}

func NewFavoriteHandler(svc *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{svc: svc}
}

func (h *FavoriteHandler) List(c *gin.Context) {
	userID := c.GetString(middleware.UserIDContextKey)

	favorites, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, favorites)
}

func (h *FavoriteHandler) Add(c *gin.Context) {
	userID := c.GetString(middleware.UserIDContextKey)

	if err := h.svc.Add(c.Request.Context(), userID, c.Param("business_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no se pudo agregar a favoritos"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *FavoriteHandler) Remove(c *gin.Context) {
	userID := c.GetString(middleware.UserIDContextKey)

	if err := h.svc.Remove(c.Request.Context(), userID, c.Param("business_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no se pudo quitar de favoritos"})
		return
	}

	c.Status(http.StatusNoContent)
}
