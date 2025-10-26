package handler

import (
	"net/http"

	"github.com/codewithwan/gopilot/internal/domain"
	"github.com/codewithwan/gopilot/internal/service"
	"github.com/gin-gonic/gin"
)

type URLShortenerHandler struct {
	service *service.URLShortenerService
}

func NewURLShortenerHandler(service *service.URLShortenerService) *URLShortenerHandler {
	return &URLShortenerHandler{service: service}
}

// CreateShortURL godoc
// @Summary Create a short URL
// @Description Create a shortened URL with optional custom alias and expiration
// @Tags url-shortener
// @Accept json
// @Produce json
// @Param request body domain.CreateShortURLRequest true "Short URL request"
// @Success 200 {object} domain.ShortURL
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/shorten [post]
func (h *URLShortenerHandler) CreateShortURL(c *gin.Context) {
	var req domain.CreateShortURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortURL, err := h.service.CreateShortURL(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shortURL)
}

// GetShortURL godoc
// @Summary Get short URL details
// @Description Get details and statistics of a short URL
// @Tags url-shortener
// @Produce json
// @Param code path string true "Short URL code"
// @Success 200 {object} domain.ShortURL
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/shorten/{code} [get]
func (h *URLShortenerHandler) GetShortURL(c *gin.Context) {
	code := c.Param("code")

	shortURL, err := h.service.GetShortURL(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shortURL)
}

// RedirectShortURL godoc
// @Summary Redirect to original URL
// @Description Redirect to the original URL and record click statistics
// @Tags url-shortener
// @Param code path string true "Short URL code"
// @Success 302 "Redirect to original URL"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /s/{code} [get]
func (h *URLShortenerHandler) RedirectShortURL(c *gin.Context) {
	code := c.Param("code")

	shortURL, err := h.service.GetShortURL(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}

	// Record click
	referrer := c.Request.Referer()
	userAgent := c.Request.UserAgent()
	ipAddress := c.ClientIP()

	if err := h.service.RecordClick(c.Request.Context(), shortURL, referrer, userAgent, ipAddress); err != nil {
		// Log error but don't fail the redirect
		c.Error(err)
	}

	c.Redirect(http.StatusFound, shortURL.OriginalURL)
}
