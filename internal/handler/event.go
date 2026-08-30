package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/steven230500/cartagena-api/internal/domain"
	"github.com/steven230500/cartagena-api/internal/service"
)

// EventInput es el body que espera el panel admin. Fechas como "2006-01-02".
type EventInput struct {
	Title         string   `json:"title" binding:"required"`
	TitleEn       *string  `json:"title_en"`
	Slug          string   `json:"slug" binding:"required"`
	StartDate     string   `json:"start_date" binding:"required"`
	EndDate       string   `json:"end_date"`
	Type          string   `json:"type" binding:"required"`
	Venue         string   `json:"venue"`
	Lat           float64  `json:"lat"`
	Lng           float64  `json:"lng"`
	Image         string   `json:"image"`
	Tags          []string `json:"tags"`
	Description   string   `json:"description"`
	DescriptionEn *string  `json:"description_en"`
	Content       *string  `json:"content"`
}

func validDateFormat(s string) bool {
	if s == "" {
		return true
	}
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}

func (in EventInput) toDomain() domain.Event {
	var endDate *string
	if in.EndDate != "" {
		endDate = &in.EndDate
	}
	return domain.Event{
		Title: in.Title, TitleEn: in.TitleEn, Slug: in.Slug, StartDate: in.StartDate, EndDate: endDate,
		Type: in.Type, Venue: in.Venue, Lat: in.Lat, Lng: in.Lng, Image: in.Image,
		Tags: in.Tags, Description: in.Description, DescriptionEn: in.DescriptionEn, Content: in.Content,
	}
}

type EventHandler struct {
	svc *service.EventService
}

func NewEventHandler(svc *service.EventService) *EventHandler {
	return &EventHandler{svc: svc}
}

func (h *EventHandler) List(c *gin.Context) {
	events, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

func (h *EventHandler) GetBySlug(c *gin.Context) {
	event, err := h.svc.GetBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}
	c.JSON(http.StatusOK, event)
}

func (h *EventHandler) Create(c *gin.Context) {
	var in EventInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !validDateFormat(in.StartDate) || !validDateFormat(in.EndDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fecha inválida, usar YYYY-MM-DD"})
		return
	}

	event, err := h.svc.Create(c.Request.Context(), in.toDomain())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

func (h *EventHandler) Update(c *gin.Context) {
	var in EventInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !validDateFormat(in.StartDate) || !validDateFormat(in.EndDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fecha inválida, usar YYYY-MM-DD"})
		return
	}

	event, err := h.svc.Update(c.Request.Context(), c.Param("id"), in.toDomain())
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, event)
}

func (h *EventHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
