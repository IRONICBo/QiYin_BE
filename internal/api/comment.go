package api

import (
	"time"

	"github.com/IRONICBo/QiYin_BE/internal/common"
	"github.com/IRONICBo/QiYin_BE/internal/common/response"
	"github.com/IRONICBo/QiYin_BE/internal/dal/dao"
	requestparams "github.com/IRONICBo/QiYin_BE/internal/params/request"
	"github.com/IRONICBo/QiYin_BE/internal/service"
	"github.com/IRONICBo/QiYin_BE/internal/utils"
	"github.com/IRONICBo/QiYin_BE/pkg/log"
	"github.com/gin-gonic/gin"
)

// CommentList
// @Tags comment
// @Summary CommentList
// @Description Test API
// @Produce application/json
// @Param videoId query string true "videoId"
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/comment/list [get].
func CommentList(c *gin.Context) {
	videoId := c.Query("videoId")

	svc := service.NewCommentService(c)
	u, err := svc.GetList(videoId)
	if err != nil {
		log.Debug("Get comment list error", err)
		response.FailWithCode(common.ERROR, c)
		return
	}

	response.SuccessWithData(u, c)
}

// CommentDelete
// @Tags comment
// @Summary CommentAction
// @Description Test API
// @Produce application/json
// @Param data body requestparams.CommentDelParams true "CommentDelParams"
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/comment/delete [post].
func CommentDelete(c *gin.Context) {
	var params requestparams.CommentDelParams
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
	svc := service.NewCommentService(c)
	err = svc.CommentDelete(userId, params.CommentId)

	if err != nil {
		log.Debug("Favorite operation error", err)
		response.FailWithCode(common.ERROR, c)
		return
	}

	response.Success(c)
}

// CommentAdd
// @Tags comment
// @Summary CommentAdd
// @Description Test API
// @Produce application/json
// @Param data body requestparams.CommentAddParams true "CommentAddParams"
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/comment/add [post].
func CommentAdd(c *gin.Context) {
	var params requestparams.CommentAddParams
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

	// 数据准备
	commentInfo := dao.Comment{
		VideoId:     params.VideoId,     // 评论视频id传入
		UserId:      userId,             // 评论用户id传入
		CommentText: params.CommentText, // 评论内容传入
		Cancel:      utils.ValidComment, // 评论状态，0，有效
		CreateDate:  time.Now(),         // 评论时间
	}

	svc := service.NewCommentService(c)
	commentData, err := svc.CommentAdd(commentInfo)
	if err != nil {
		log.Debug("Favorite operation error", err)
		response.FailWithCode(common.ERROR, c)
		return
	}

	response.SuccessWithData(commentData, c)
}
