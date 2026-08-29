package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/steven230500/cartagena-api/internal/domain"
	"github.com/steven230500/cartagena-api/internal/service"
)

type ScheduleInput struct {
	Day   string   `json:"day" binding:"required"`
	Times []string `json:"times"`
}

type ParishInput struct {
	Name         string          `json:"name" binding:"required"`
	Address      string          `json:"address"`
	Neighborhood string          `json:"neighborhood" binding:"required"`
	Phone        *string         `json:"phone"`
	Schedules    []ScheduleInput `json:"schedules"`
}

func (in ParishInput) toDomain() domain.Parish {
	schedules := make([]domain.Schedule, len(in.Schedules))
	for i, s := range in.Schedules {
		schedules[i] = domain.Schedule{Day: s.Day, Times: s.Times}
	}
	return domain.Parish{
		Name: in.Name, Address: in.Address, Neighborhood: in.Neighborhood, Phone: in.Phone,
		Schedules: schedules,
	}
}

type ParishHandler struct {
	svc *service.ParishService
}

func NewParishHandler(svc *service.ParishService) *ParishHandler {
	return &ParishHandler{svc: svc}
}

func (h *ParishHandler) List(c *gin.Context) {
	parishes, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, parishes)
}

func (h *ParishHandler) Create(c *gin.Context) {
	var in ParishInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parish, err := h.svc.Create(c.Request.Context(), in.toDomain())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, parish)
}

func (h *ParishHandler) Update(c *gin.Context) {
	var in ParishInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parish, err := h.svc.Update(c.Request.Context(), c.Param("id"), in.toDomain())
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "parish not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, parish)
}

func (h *ParishHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
