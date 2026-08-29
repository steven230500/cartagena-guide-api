package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/steven230500/cartagena-api/internal/domain"
	"github.com/steven230500/cartagena-api/internal/middleware"
	"github.com/steven230500/cartagena-api/internal/service"
)

type ClaimInput struct {
	BusinessID string `json:"business_id" binding:"required"`
	Message    string `json:"message"`
}

type BusinessClaimHandler struct {
	svc *service.BusinessClaimService
}

func NewBusinessClaimHandler(svc *service.BusinessClaimService) *BusinessClaimHandler {
	return &BusinessClaimHandler{svc: svc}
}

func (h *BusinessClaimHandler) Create(c *gin.Context) {
	userID := c.GetString(middleware.UserIDContextKey)

	var in ClaimInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claim, err := h.svc.Create(c.Request.Context(), in.BusinessID, userID, in.Message)
	if err != nil {
		if errors.Is(err, service.ErrAlreadyOwned) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "este negocio ya tiene dueño"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "no se pudo crear la solicitud"})
		return
	}

	c.JSON(http.StatusCreated, claim)
}

func (h *BusinessClaimHandler) ListMine(c *gin.Context) {
	userID := c.GetString(middleware.UserIDContextKey)

	claims, err := h.svc.ListMine(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, claims)
}

func (h *BusinessClaimHandler) ListPending(c *gin.Context) {
	claims, err := h.svc.ListPending(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, claims)
}

func (h *BusinessClaimHandler) Approve(c *gin.Context) {
	claim, err := h.svc.Approve(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "claim not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, claim)
}

func (h *BusinessClaimHandler) Reject(c *gin.Context) {
	claim, err := h.svc.Reject(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "claim not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, claim)
}
