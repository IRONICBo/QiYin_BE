package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// GetPfopCallback
// @Tags QiNiu
// @Summary GetPfopCallback
// @Description Get QiNiu Pfop callback result
// @Produce application/json
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/qiniu/pfop/callback [get].
func GetPfopCallback(c *gin.Context) {
	// Print the callback result.
	fmt.Println("GetPfopCallback", "GetPfopCallback", c.Request.Body)
}
