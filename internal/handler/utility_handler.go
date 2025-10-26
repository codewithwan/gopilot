package handler

import (
	"net/http"

	"github.com/codewithwan/gopilot/internal/domain"
	"github.com/codewithwan/gopilot/internal/service"
	"github.com/gin-gonic/gin"
)

type UtilityHandler struct {
	hashService      *service.HashService
	converterService *service.ConverterService
	generatorService *service.GeneratorService
	cryptoService    *service.CryptoService
}

func NewUtilityHandler() *UtilityHandler {
	return &UtilityHandler{
		hashService:      service.NewHashService(),
		converterService: service.NewConverterService(),
		generatorService: service.NewGeneratorService(),
		cryptoService:    service.NewCryptoService(),
	}
}

// Hash godoc
// @Summary Hash text
// @Description Hash text using specified algorithm
// @Tags hash-encode
// @Accept json
// @Produce json
// @Param request body domain.HashRequest true "Hash request"
// @Success 200 {object} domain.HashResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/hash [post]
func (h *UtilityHandler) Hash(c *gin.Context) {
	var req domain.HashRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.hashService.Hash(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Encode godoc
// @Summary Encode/decode text
// @Description Encode or decode text using specified operation
// @Tags hash-encode
// @Accept json
// @Produce json
// @Param request body domain.EncodeRequest true "Encode request"
// @Success 200 {object} domain.EncodeResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/encode [post]
func (h *UtilityHandler) Encode(c *gin.Context) {
	var req domain.EncodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.hashService.Encode(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GeneratePassword godoc
// @Summary Generate password
// @Description Generate a random secure password
// @Tags hash-encode
// @Accept json
// @Produce json
// @Param request body domain.GeneratePasswordRequest true "Password request"
// @Success 200 {object} domain.GeneratePasswordResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/generate/password [post]
func (h *UtilityHandler) GeneratePassword(c *gin.Context) {
	var req domain.GeneratePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.hashService.GeneratePassword(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ConvertBase godoc
// @Summary Convert number base
// @Description Convert numbers between different bases
// @Tags converter
// @Accept json
// @Produce json
// @Param request body domain.ConvertBaseRequest true "Base conversion request"
// @Success 200 {object} domain.ConvertBaseResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/convert/base [post]
func (h *UtilityHandler) ConvertBase(c *gin.Context) {
	var req domain.ConvertBaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.converterService.ConvertBase(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ConvertColor godoc
// @Summary Convert color format
// @Description Convert colors between different formats
// @Tags converter
// @Accept json
// @Produce json
// @Param request body domain.ConvertColorRequest true "Color conversion request"
// @Success 200 {object} domain.ConvertColorResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/convert/color [post]
func (h *UtilityHandler) ConvertColor(c *gin.Context) {
	var req domain.ConvertColorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.converterService.ConvertColor(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ConvertTime godoc
// @Summary Convert time format
// @Description Convert time between different formats
// @Tags converter
// @Accept json
// @Produce json
// @Param request body domain.ConvertTimeRequest true "Time conversion request"
// @Success 200 {object} domain.ConvertTimeResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/convert/time [post]
func (h *UtilityHandler) ConvertTime(c *gin.Context) {
	var req domain.ConvertTimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.converterService.ConvertTime(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// FormatJSON godoc
// @Summary Format JSON
// @Description Format or minify JSON
// @Tags formatter
// @Accept json
// @Produce json
// @Param request body domain.FormatJSONRequest true "JSON format request"
// @Success 200 {object} domain.FormatJSONResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/format/json [post]
func (h *UtilityHandler) FormatJSON(c *gin.Context) {
	var req domain.FormatJSONRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.converterService.FormatJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ConvertYAML godoc
// @Summary Convert YAML/JSON
// @Description Convert between YAML and JSON formats
// @Tags formatter
// @Accept json
// @Produce json
// @Param request body domain.ConvertYAMLRequest true "YAML conversion request"
// @Success 200 {object} domain.ConvertYAMLResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/format/yaml [post]
func (h *UtilityHandler) ConvertYAML(c *gin.Context) {
	var req domain.ConvertYAMLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.converterService.ConvertYAML(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GenerateUUID godoc
// @Summary Generate UUID
// @Description Generate UUIDs
// @Tags generator
// @Accept json
// @Produce json
// @Param request body domain.GenerateUUIDRequest true "UUID request"
// @Success 200 {object} domain.GenerateUUIDResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/generate/uuid [post]
func (h *UtilityHandler) GenerateUUID(c *gin.Context) {
	var req domain.GenerateUUIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.generatorService.GenerateUUID(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GenerateToken godoc
// @Summary Generate token
// @Description Generate random tokens
// @Tags generator
// @Accept json
// @Produce json
// @Param request body domain.GenerateTokenRequest true "Token request"
// @Success 200 {object} domain.GenerateTokenResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/generate/token [post]
func (h *UtilityHandler) GenerateToken(c *gin.Context) {
	var req domain.GenerateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.generatorService.GenerateToken(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GenerateLorem godoc
// @Summary Generate lorem ipsum
// @Description Generate lorem ipsum text
// @Tags generator
// @Accept json
// @Produce json
// @Param request body domain.GenerateLoremRequest true "Lorem request"
// @Success 200 {object} domain.GenerateLoremResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/generate/lorem [post]
func (h *UtilityHandler) GenerateLorem(c *gin.Context) {
	var req domain.GenerateLoremRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.generatorService.GenerateLorem(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GenerateFakeUser godoc
// @Summary Generate fake user data
// @Description Generate fake user profiles
// @Tags generator
// @Accept json
// @Produce json
// @Param request body domain.GenerateFakeUserRequest true "Fake user request"
// @Success 200 {object} domain.GenerateFakeUserResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/generate/user [post]
func (h *UtilityHandler) GenerateFakeUser(c *gin.Context) {
	var req domain.GenerateFakeUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.generatorService.GenerateFakeUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GenerateRandomNumber godoc
// @Summary Generate random numbers
// @Description Generate random numbers
// @Tags generator
// @Accept json
// @Produce json
// @Param request body domain.GenerateRandomNumberRequest true "Random number request"
// @Success 200 {object} domain.GenerateRandomNumberResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/generate/number [post]
func (h *UtilityHandler) GenerateRandomNumber(c *gin.Context) {
	var req domain.GenerateRandomNumberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.generatorService.GenerateRandomNumber(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// AESOperation godoc
// @Summary AES encrypt/decrypt
// @Description Perform AES encryption or decryption
// @Tags crypto
// @Accept json
// @Produce json
// @Param request body domain.AESRequest true "AES request"
// @Success 200 {object} domain.AESResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/crypto/aes [post]
func (h *UtilityHandler) AESOperation(c *gin.Context) {
	var req domain.AESRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.cryptoService.AESOperation(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GenerateRSAKeypair godoc
// @Summary Generate RSA keypair
// @Description Generate an RSA public/private keypair
// @Tags crypto
// @Produce json
// @Success 200 {object} domain.RSAKeypairResponse
// @Failure 500 {object} map[string]string
// @Router /v1/crypto/rsa/keygen [post]
func (h *UtilityHandler) GenerateRSAKeypair(c *gin.Context) {
	result, err := h.cryptoService.GenerateRSAKeypair()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// RSAOperation godoc
// @Summary RSA encrypt/decrypt
// @Description Perform RSA encryption or decryption
// @Tags crypto
// @Accept json
// @Produce json
// @Param request body domain.RSARequest true "RSA request"
// @Success 200 {object} domain.RSAResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/crypto/rsa [post]
func (h *UtilityHandler) RSAOperation(c *gin.Context) {
	var req domain.RSARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.cryptoService.RSAOperation(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// HMACOperation godoc
// @Summary HMAC sign/verify
// @Description Sign or verify HMAC signatures
// @Tags crypto
// @Accept json
// @Produce json
// @Param request body domain.HMACRequest true "HMAC request"
// @Success 200 {object} domain.HMACResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/crypto/hmac [post]
func (h *UtilityHandler) HMACOperation(c *gin.Context) {
	var req domain.HMACRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.cryptoService.HMACOperation(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
