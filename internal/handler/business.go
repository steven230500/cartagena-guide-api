package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/steven230500/cartagena-api/internal/domain"
	"github.com/steven230500/cartagena-api/internal/middleware"
	"github.com/steven230500/cartagena-api/internal/service"
)

// BusinessInput es el body que espera el panel admin para crear/editar un comercio.
type BusinessInput struct {
	Name             string   `json:"name" binding:"required"`
	Slug             string   `json:"slug" binding:"required"`
	Type             string   `json:"type" binding:"required"`
	Subtype          string   `json:"subtype"`
	Barrio           string   `json:"barrio" binding:"required"`
	Lat              float64  `json:"lat"`
	Lng              float64  `json:"lng"`
	Image            string   `json:"image"`
	Tags             []string `json:"tags"`
	Description      string   `json:"description"`
	Hours            *string  `json:"hours"`
	PriceHint        *string  `json:"price_hint"`
	PriceTypicalNote *string  `json:"price_typical_note"`
	Phone            *string  `json:"phone"`
	Web              *string  `json:"web"`
	Email            *string  `json:"email"`
	Instagram        *string  `json:"instagram"`
	Verified         bool     `json:"verified"`
}

func (in BusinessInput) toDomain() domain.Business {
	return domain.Business{
		Name: in.Name, Slug: in.Slug, Type: in.Type, Subtype: in.Subtype, Barrio: in.Barrio,
		Lat: in.Lat, Lng: in.Lng, Image: in.Image, Tags: in.Tags, Description: in.Description,
		Hours: in.Hours, PriceHint: in.PriceHint, PriceTypicalNote: in.PriceTypicalNote,
		Phone: in.Phone, Web: in.Web, Email: in.Email, Instagram: in.Instagram, Verified: in.Verified,
	}
}

type BusinessHandler struct {
	svc *service.BusinessService
}

func NewBusinessHandler(svc *service.BusinessService) *BusinessHandler {
	return &BusinessHandler{svc: svc}
}

func (h *BusinessHandler) List(c *gin.Context) {
	filter := domain.BusinessFilter{
		Type:   c.Query("type"),
		Barrio: c.Query("barrio"),
		Q:      c.Query("q"),
	}

	businesses, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, businesses)
}

func (h *BusinessHandler) GetBySlug(c *gin.Context) {
	business, err := h.svc.GetBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "business not found"})
		return
	}

	c.JSON(http.StatusOK, business)
}

func (h *BusinessHandler) Create(c *gin.Context) {
	var in BusinessInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	business, err := h.svc.Create(c.Request.Context(), in.toDomain())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, business)
}

func (h *BusinessHandler) Update(c *gin.Context) {
	var in BusinessInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	business, err := h.svc.Update(c.Request.Context(), c.Param("id"), in.toDomain())
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "business not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, business)
}

func (h *BusinessHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// MyBusinessInput es el subconjunto de campos que puede tocar el dueño de un
// negocio (Fase 5) — nombre/slug/tipo/subtipo/barrio/coords/verified quedan
// admin-only, así que ni siquiera se aceptan en este body.
type MyBusinessInput struct {
	Description      string   `json:"description"`
	Hours            *string  `json:"hours"`
	PriceHint        *string  `json:"price_hint"`
	PriceTypicalNote *string  `json:"price_typical_note"`
	Phone            *string  `json:"phone"`
	Web              *string  `json:"web"`
	Email            *string  `json:"email"`
	Instagram        *string  `json:"instagram"`
	Image            string   `json:"image"`
	Tags             []string `json:"tags"`
}

func (in MyBusinessInput) toPatch() domain.BusinessOwnerPatch {
	return domain.BusinessOwnerPatch{
		Description: in.Description, Hours: in.Hours, PriceHint: in.PriceHint,
		PriceTypicalNote: in.PriceTypicalNote, Phone: in.Phone, Web: in.Web,
		Email: in.Email, Instagram: in.Instagram, Image: in.Image, Tags: in.Tags,
	}
}

func (h *BusinessHandler) ListMine(c *gin.Context) {
	userID := c.GetString(middleware.UserIDContextKey)

	businesses, err := h.svc.ListByOwner(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, businesses)
}

func (h *BusinessHandler) UpdateMine(c *gin.Context) {
	userID := c.GetString(middleware.UserIDContextKey)

	var in MyBusinessInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	business, err := h.svc.UpdateAsOwner(c.Request.Context(), c.Param("id"), userID, in.toPatch())
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "no sos el dueño de este negocio"})
			return
		}
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "business not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, business)
}
