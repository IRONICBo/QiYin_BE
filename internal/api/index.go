package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping
// @Tags Ping
// @Summary Ping
// @Description Test API
// @Produce application/json
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/ping [post].
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
