package api

import (
	"github.com/IRONICBo/QiYin_BE/internal/common"
	"github.com/IRONICBo/QiYin_BE/internal/common/response"
	"github.com/IRONICBo/QiYin_BE/internal/middleware/jwt"
	"github.com/IRONICBo/QiYin_BE/internal/service"
	"github.com/IRONICBo/QiYin_BE/pkg/log"
	"github.com/gin-gonic/gin"
)

// Search
// @Tags video
// @Summary Search
// @Description search videos by text
// @Produce application/json
// @Param searchValue query string true "query text"
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/search [get].
func Search(c *gin.Context) {
	searchValue := c.Query("searchValue")

	userId := jwt.GetUserId(c)

	svc := service.NewVideoService(c)
	u, err := svc.Search(searchValue, userId)
	if err != nil {
		log.Debug("user doesn't exit", err)
		response.FailWithCode(common.ERROR, c)
		return
	}
	response.SuccessWithData(u, c)
}

// GetHots
// @Tags video
// @Summary GetHots
// @Description hot list
// @Produce application/json
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/hots [get].
func GetHots(c *gin.Context) {
	svc := service.NewVideoService(c)
	u, err := svc.GetHots()
	if err != nil {
		log.Debug("user doesn't exit", err)
		response.FailWithCode(common.ERROR, c)
		return
	}
	response.SuccessWithData(u, c)
}
