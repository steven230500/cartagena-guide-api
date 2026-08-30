package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/steven230500/cartagena-api/internal/domain"
	"github.com/steven230500/cartagena-api/internal/service"
)

type StepInput struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	AudioURL    *string  `json:"audio_url"`
	Duration    *string  `json:"duration"`
	Lat         *float64 `json:"lat"`
	Lng         *float64 `json:"lng"`
	Image       *string  `json:"image"`
}

type RouteInput struct {
	Slug          string      `json:"slug" binding:"required"`
	Title         string      `json:"title" binding:"required"`
	TitleEn       *string     `json:"title_en"`
	Description   string      `json:"description"`
	DescriptionEn *string     `json:"description_en"`
	Duration      string      `json:"duration"`
	Distance      string      `json:"distance"`
	Difficulty    string      `json:"difficulty"`
	Category      string      `json:"category" binding:"required"`
	Image         string      `json:"image"`
	Highlights    []string    `json:"highlights"`
	AudioGuide    bool        `json:"audio_guide"`
	Offline       bool        `json:"offline"`
	Price         string      `json:"price"`
	Steps         []StepInput `json:"steps"`
}

func (in RouteInput) toDomain() domain.Route {
	steps := make([]domain.RouteStep, len(in.Steps))
	for i, s := range in.Steps {
		steps[i] = domain.RouteStep{
			Title: s.Title, Description: s.Description, AudioURL: s.AudioURL,
			Duration: s.Duration, Lat: s.Lat, Lng: s.Lng, Image: s.Image,
		}
	}
	return domain.Route{
		Slug: in.Slug, Title: in.Title, TitleEn: in.TitleEn, Description: in.Description, DescriptionEn: in.DescriptionEn, Duration: in.Duration,
		Distance: in.Distance, Difficulty: in.Difficulty, Category: in.Category, Image: in.Image,
		Highlights: in.Highlights, AudioGuide: in.AudioGuide, Offline: in.Offline, Price: in.Price,
		Steps: steps,
	}
}

type RouteHandler struct {
	svc *service.RouteService
}

func NewRouteHandler(svc *service.RouteService) *RouteHandler {
	return &RouteHandler{svc: svc}
}

func (h *RouteHandler) List(c *gin.Context) {
	routes, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, routes)
}

func (h *RouteHandler) GetBySlug(c *gin.Context) {
	route, err := h.svc.GetBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
		return
	}
	c.JSON(http.StatusOK, route)
}

func (h *RouteHandler) Create(c *gin.Context) {
	var in RouteInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	route, err := h.svc.Create(c.Request.Context(), in.toDomain())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, route)
}

func (h *RouteHandler) Update(c *gin.Context) {
	var in RouteInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	route, err := h.svc.Update(c.Request.Context(), c.Param("id"), in.toDomain())
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, route)
}

func (h *RouteHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
