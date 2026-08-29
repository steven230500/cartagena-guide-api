package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/steven230500/cartagena-api/internal/domain"
	"github.com/steven230500/cartagena-api/internal/middleware"
	"github.com/steven230500/cartagena-api/internal/service"
)

type AchievementInput struct {
	Code         string `json:"code" binding:"required"`
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description"`
	Icon         string `json:"icon" binding:"required"`
	CriteriaType string `json:"criteria_type" binding:"required"`
	Threshold    int    `json:"threshold" binding:"required"`
}

func (in AchievementInput) toDomain() domain.Achievement {
	return domain.Achievement{
		Code: in.Code, Title: in.Title, Description: in.Description,
		Icon: in.Icon, CriteriaType: in.CriteriaType, Threshold: in.Threshold,
	}
}

type AchievementHandler struct {
	svc *service.AchievementService
}

func NewAchievementHandler(svc *service.AchievementService) *AchievementHandler {
	return &AchievementHandler{svc: svc}
}

func (h *AchievementHandler) List(c *gin.Context) {
	achievements, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, achievements)
}

func (h *AchievementHandler) Create(c *gin.Context) {
	var in AchievementInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	achievement, err := h.svc.Create(c.Request.Context(), in.toDomain())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, achievement)
}

func (h *AchievementHandler) Update(c *gin.Context) {
	var in AchievementInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	achievement, err := h.svc.Update(c.Request.Context(), c.Param("id"), in.toDomain())
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "achievement not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, achievement)
}

func (h *AchievementHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AchievementHandler) Progress(c *gin.Context) {
	userID := c.GetString(middleware.UserIDContextKey)

	achievements, stats, err := h.svc.Progress(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"achievements": achievements, "stats": stats})
}
