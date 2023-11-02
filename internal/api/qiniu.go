package api

import (
	"github.com/IRONICBo/QiYin_BE/internal/common/response"
	"github.com/IRONICBo/QiYin_BE/internal/service"
	"github.com/gin-gonic/gin"
)

// GetUploadToken
// @Tags QiNiu
// @Summary UserLogin
// @Description Get QiNiu upload token
// @Produce application/json
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/qiniu/token [get].
func GetUploadToken(c *gin.Context) {
	svc := service.NewQiNiuService(c)
	uploadToken := svc.GetUploadToken()

	response.SuccessWithData(uploadToken, c)
}
