package api

import (
	"strconv"

	"github.com/IRONICBo/QiYin_BE/internal/common"
	"github.com/IRONICBo/QiYin_BE/internal/common/response"
	"github.com/IRONICBo/QiYin_BE/internal/middleware/jwt"
	requestparams "github.com/IRONICBo/QiYin_BE/internal/params/request"
	"github.com/IRONICBo/QiYin_BE/internal/service"
	"github.com/IRONICBo/QiYin_BE/pkg/log"
	"github.com/gin-gonic/gin"
)

// FavoriteAction
// @Tags favorite
// @Summary FavoriteAction
// @Description like or dislike
// @Produce application/json
// @Param data body requestparams.FavoriteParams true "FavoriteParams"
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/favorite/action [post].
func FavoriteAction(c *gin.Context) {
	var params requestparams.FavoriteParams
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
	svc := service.NewFavoriteService(c)
	err = svc.FavoriteAction(userId, strconv.FormatInt(params.VideoId, 10), params.ActionType)

	if err != nil {
		log.Debug("Favorite operation error", err)
		response.FailWithCode(common.ERROR, c)
		return
	}

	response.Success(c)
}

// GetFavoriteList
// @Tags favorite
// @Summary GetFavoriteList
// @Description get favorite video list
// @Produce application/json
// @Param userId query string true "query user id"
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/favorite/list [get].
func GetFavoriteList(c *gin.Context) {
	// 要查看的用户的id
	strUserId := c.Query("userId")
	// 自己的id  如果登录之后不是""  这是因为如果登录之后需要获取到curid 是否点赞过该视频，没有的话默认为false
	curUserId := jwt.GetUserId(c)
	svc := service.NewFavoriteService(c)
	u, err := svc.FavoriteList(strUserId, curUserId)
	if err != nil {
		log.Debug("Favorite count error", err)
		response.FailWithCode(common.ERROR, c)
		return
	}

	response.SuccessWithData(u, c)
}
