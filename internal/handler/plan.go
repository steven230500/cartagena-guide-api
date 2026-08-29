package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/steven230500/cartagena-api/internal/domain"
	"github.com/steven230500/cartagena-api/internal/service"
)

type PlanInput struct {
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description"`
	Type         string  `json:"type" binding:"required"`
	Price        string  `json:"price" binding:"required"`
	PriceAmount  *string `json:"price_amount"`
	Date         string  `json:"date"`
	Time         string  `json:"time"`
	Location     string  `json:"location"`
	Neighborhood string  `json:"neighborhood" binding:"required"`
}

func (in PlanInput) toDomain() domain.Plan {
	return domain.Plan{
		Title: in.Title, Description: in.Description, Type: in.Type, Price: in.Price,
		PriceAmount: in.PriceAmount, Date: in.Date, Time: in.Time, Location: in.Location,
		Neighborhood: in.Neighborhood,
	}
}

type PlanHandler struct {
	svc *service.PlanService
}

func NewPlanHandler(svc *service.PlanService) *PlanHandler {
	return &PlanHandler{svc: svc}
}

func (h *PlanHandler) List(c *gin.Context) {
	plans, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plans)
}

func (h *PlanHandler) Create(c *gin.Context) {
	var in PlanInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := h.svc.Create(c.Request.Context(), in.toDomain())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

func (h *PlanHandler) Update(c *gin.Context) {
	var in PlanInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := h.svc.Update(c.Request.Context(), c.Param("id"), in.toDomain())
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plan)
}

func (h *PlanHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
