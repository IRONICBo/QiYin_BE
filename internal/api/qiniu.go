package api

import (
	"encoding/base64"
	"io"
	"net/http"

	"github.com/IRONICBo/QiYin_BE/internal/common"
	"github.com/IRONICBo/QiYin_BE/internal/common/response"
	requestparams "github.com/IRONICBo/QiYin_BE/internal/params/request"
	"github.com/IRONICBo/QiYin_BE/internal/service"
	"github.com/gin-gonic/gin"
)

// GetUploadToken
// @Tags QiNiu
// @Summary UserLogin
// @Description Get QiNiu upload token
// @Produce application/json
// @Param data body requestparams.QiNiuTokenParams true "QiNiuTokenParams"
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/qiniu/token [post].
func GetUploadToken(c *gin.Context) {
	var params requestparams.QiNiuTokenParams
	err := c.ShouldBindJSON(&params)
	if err != nil {
		response.FailWithCode(common.INVALID_PARAMS, c)

		return
	}
	svc := service.NewQiNiuService(c)
	uploadToken := svc.GetUploadToken(params.Ticket)

	response.SuccessWithData(uploadToken, c)
}

// GetImageByProxy
// @Tags QiNiu
// @Summary GetImageByProxy
// @Description Get QiNiu image by proxy
// @Produce application/json
// @Param url query string true "url"
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/qiniu/proxy [get].
func GetImageByProxy(c *gin.Context) {
	imageURL := c.Query("url")

	resp, err := http.Get(imageURL)
	if err != nil {
		response.FailWithCode(common.INVALID_PARAMS, c)

		return
	}
	defer resp.Body.Close()

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		response.FailWithCode(common.ERROR, c)

		return
	}

	base64Str := base64.StdEncoding.EncodeToString(imageData)

	response.SuccessWithData(base64Str, c)
}
