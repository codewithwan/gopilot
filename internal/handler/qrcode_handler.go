package handler

import (
	"net/http"

	"github.com/codewithwan/gopilot/internal/domain"
	"github.com/codewithwan/gopilot/internal/service"
	"github.com/gin-gonic/gin"
)

type QRCodeHandler struct {
	service *service.QRCodeService
}

func NewQRCodeHandler(service *service.QRCodeService) *QRCodeHandler {
	return &QRCodeHandler{service: service}
}

// GenerateQR godoc
// @Summary Generate QR code
// @Description Generate a QR code from text
// @Tags qr-code
// @Accept json
// @Produce json
// @Param request body domain.GenerateQRRequest true "QR code request"
// @Success 200 {object} domain.QRCode
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/qr [post]
func (h *QRCodeHandler) GenerateQR(c *gin.Context) {
	var req domain.GenerateQRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	qr, err := h.service.GenerateQR(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     qr.ID,
		"text":   qr.Text,
		"format": qr.Format,
		"size":   qr.Size,
	})
}

// GetQRCode godoc
// @Summary Get QR code
// @Description Get a previously generated QR code
// @Tags qr-code
// @Produce png
// @Param id path string true "QR code ID"
// @Success 200 {file} image/png
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/qr/{id} [get]
func (h *QRCodeHandler) GetQRCode(c *gin.Context) {
	id := c.Param("id")

	qr, err := h.service.GetQRCode(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "image/png", qr.ImageData)
}
