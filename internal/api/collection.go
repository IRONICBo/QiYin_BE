package api

import (
	"github.com/IRONICBo/QiYin_BE/internal/common"
	"github.com/IRONICBo/QiYin_BE/internal/common/response"
	"github.com/IRONICBo/QiYin_BE/internal/middleware/jwt"
	requestparams "github.com/IRONICBo/QiYin_BE/internal/params/request"
	"github.com/IRONICBo/QiYin_BE/internal/service"
	"github.com/IRONICBo/QiYin_BE/pkg/log"
	"github.com/gin-gonic/gin"
	"strconv"
)

// CollectionAction
// @Tags favorite
// @Summary CollectionAction
// @Description collect or cancel
// @Param data body requestparams.CollectionParams true "CollectionParams"
// @Produce application/json
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/favorite/action [post].
func CollectionAction(c *gin.Context) {
	var params requestparams.CollectionParams
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
	svc := service.NewCollectionService(c)
	err = svc.CollectionAction(userId, strconv.FormatInt(params.VideoId, 10), params.ActionType)

	if err != nil {
		log.Debug("Collection operation error", err)
		response.FailWithCode(common.ERROR, c)
		return
	}

	response.Success(c)
}

// GetCollectionList
// @Tags favorite
// @Summary GetCollectionList
// @Description get collection video list
// @Param userId query string true "query user id"
// @Produce application/json
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/favorite/list [get].
func GetCollectionList(c *gin.Context) {
	// 要查看的用户的id
	strUserId := c.Query("userId")
	// 自己的id  如果登录之后不是""  这是因为如果登录之后需要获取到curid 是否点赞过该视频，没有的话默认为false
	curUserId := jwt.GetUserId(c)
	svc := service.NewCollectionService(c)
	u, err := svc.CollectionList(strUserId, curUserId)
	if err != nil {
		log.Debug("Collection count error", err)
		response.FailWithCode(common.ERROR, c)
		return
	}

	response.SuccessWithData(u, c)
}
