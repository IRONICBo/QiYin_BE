package api

import (
	"github.com/IRONICBo/QiYin_BE/internal/common"
	"github.com/IRONICBo/QiYin_BE/internal/common/response"
	"github.com/IRONICBo/QiYin_BE/internal/middleware/jwt"
	requestparams "github.com/IRONICBo/QiYin_BE/internal/params/request"
	"github.com/IRONICBo/QiYin_BE/internal/service"
	"github.com/IRONICBo/QiYin_BE/pkg/log"
	"github.com/gin-gonic/gin"
)

// Search
// @Tags video
// @Summary Search
// @Description search videos by text
// @Produce application/json
// @Param searchValue query string true "searchValue"
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/video/search [get].
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

// GetVideos
// @Tags video
// @Summary GetVideos
// @Description get videos by userId
// @Produce application/json
// @Param searchValue query string true "searchValue"
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/video/list [get].
func GetVideos(c *gin.Context) {
	userId := c.Query("userId")
	//是否登录
	curId := jwt.GetUserId(c)

	svc := service.NewVideoService(c)
	u, err := svc.GetVideoByUserId(userId, curId)

	if err != nil {
		log.Debug("video doesn't exit", err)
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
// @Param searchValue query string true "searchValue"
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/video/hots [get].
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

// UploadVideo
// @Tags video
// @Summary UploadVideo
// @Description hot list
// @Produce application/json
// @Param data body requestparams.VideoUpdateParams true "VideoUpdateParams"
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/video/upload [post].
func UploadVideo(c *gin.Context) {
	var params requestparams.VideoUpdateParams
	err := c.ShouldBindJSON(&params)
	if err != nil {
		response.FailWithCode(common.INVALID_PARAMS, c)
		return
	}

	userId := c.GetString("userId")
	if len(userId) == 0 {
		response.FailWithCode(common.INVALID_PARAMS, c)
		return
	}
	svc := service.NewVideoService(c)
	err = svc.UploadVideo(userId, &params)

	if err != nil {
		log.Debug("upload operation error", err)
		response.FailWithCode(common.ERROR, c)
		return
	}

	response.Success(c)
}

// GetHistory
// @Tags video
// @Summary GetHistory
// @Description get history list
// @Produce application/json
// @Param data body requestparams.VideoUpdateParams true "VideoUpdateParams"
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/video/getHistory [get].
func GetHistory(c *gin.Context) {
	userId := c.GetString("userId")
	if len(userId) == 0 {
		response.FailWithCode(common.INVALID_PARAMS, c)
		return
	}

	svc := service.NewVideoService(c)
	videos, err := svc.GetHisVideos(userId)

	if err != nil {
		log.Debug("upload operation error", err)
		response.FailWithCode(common.ERROR, c)
		return
	}

	response.SuccessWithData(videos, c)
}

// SaveVideoHis
// @Tags video
// @Summary SaveVideoHis
// @Description video history
// @Produce application/json
// @Param data body requestparams.VideoHisParams true "VideoHisParams"
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/video/save [post].
func SaveVideoHis(c *gin.Context) {
	//是否登录
	curId := jwt.GetUserId(c)
	if curId == "" {
		response.FailWithCode(common.ERROR, c)
		return
	}

	var params requestparams.VideoHisParams
	err := c.ShouldBindJSON(&params)
	if err != nil {
		response.FailWithCode(common.INVALID_PARAMS, c)
		return
	}

	svc := service.NewVideoService(c)
	err = svc.SaveVideoHis(curId, &params)

	if err != nil {
		log.Debug("save operation error", err)
		response.FailWithCode(common.ERROR, c)
		return
	}

	response.Success(c)
}
