package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HelloResponse struct {
	Message string `json:"message" example:"Hello, World!"`
}

// Hello godoc
// @Summary      Hello endpoint
// @Description  Returns a greeting message. Requires API key as query parameter 'api'.
// @Tags         hello
// @Produce      json
// @Param        api   query     string  true  "API Key"
// @Success      200   {object}  HelloResponse
// @Failure      403   {object}  map[string]string
// @Router       /hello [get]
func hello(c *gin.Context) {
	c.JSON(http.StatusOK, HelloResponse{Message: "Hello, World!"})
}
