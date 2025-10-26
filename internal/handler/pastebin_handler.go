package handler

import (
	"net/http"

	"github.com/codewithwan/gopilot/internal/domain"
	"github.com/codewithwan/gopilot/internal/service"
	"github.com/gin-gonic/gin"
)

type PastebinHandler struct {
	service *service.PastebinService
}

func NewPastebinHandler(service *service.PastebinService) *PastebinHandler {
	return &PastebinHandler{service: service}
}

// CreatePaste godoc
// @Summary Create a paste
// @Description Create a new paste/snippet
// @Tags pastebin
// @Accept json
// @Produce json
// @Param request body domain.CreatePasteRequest true "Paste request"
// @Success 200 {object} domain.Paste
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/paste [post]
func (h *PastebinHandler) CreatePaste(c *gin.Context) {
	var req domain.CreatePasteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paste, err := h.service.CreatePaste(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, paste)
}

// GetPaste godoc
// @Summary Get paste content
// @Description Get the content of a paste by ID
// @Tags pastebin
// @Produce json
// @Param id path string true "Paste ID"
// @Success 200 {object} domain.Paste
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /p/{id} [get]
func (h *PastebinHandler) GetPaste(c *gin.Context) {
	id := c.Param("id")

	paste, err := h.service.GetPaste(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, paste)
}

// DeletePaste godoc
// @Summary Delete a paste
// @Description Delete a paste by ID
// @Tags pastebin
// @Param id path string true "Paste ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/paste/{id} [delete]
func (h *PastebinHandler) DeletePaste(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeletePaste(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListRecentPastes godoc
// @Summary List recent pastes
// @Description Get a list of recent public pastes
// @Tags pastebin
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Success 200 {array} domain.Paste
// @Failure 500 {object} map[string]string
// @Router /v1/paste/recent [get]
func (h *PastebinHandler) ListRecentPastes(c *gin.Context) {
	limit := 20
	if l, ok := c.GetQuery("limit"); ok {
		if parsedLimit, err := parseIntQueryParam(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	pastes, err := h.service.ListRecentPastes(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pastes)
}
