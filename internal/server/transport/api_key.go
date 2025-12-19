package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"fandom/notifications/internal/models"
	"fandom/notifications/internal/service"
)

type APIKeyHandler struct {
	service *service.APIKeyService
}

func NewAPIKeyHandler(apiKeyService *service.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{service: apiKeyService}
}

// CreateAPIKey godoc
// @Summary      Generate a new API key
// @Description  Creates a new API key with the provided name. Requires master API key as query parameter 'api'.
// @Tags         api-keys
// @Accept       json
// @Produce      json
// @Param        api      query     string  true  "Master API Key"
// @Param        request  body      models.CreateAPIKeyRequest  true  "API Key Request"
// @Success      201     {object}  models.CreateAPIKeyResponse
// @Failure      400     {object}  map[string]string
// @Failure      403     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /api-keys [post]
func (h *APIKeyHandler) CreateAPIKey(c *gin.Context) {
	var req models.CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request. 'name' field is required.",
		})
		return
	}

	response, err := h.service.CreateAPIKey(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create API key",
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

