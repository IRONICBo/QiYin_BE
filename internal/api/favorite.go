package api

import (
	"github.com/IRONICBo/QiYin_BE/internal/common"
	"github.com/IRONICBo/QiYin_BE/internal/common/response"
	requestparams "github.com/IRONICBo/QiYin_BE/internal/params/request"
	"github.com/IRONICBo/QiYin_BE/internal/service"
	"github.com/IRONICBo/QiYin_BE/pkg/log"
	"github.com/gin-gonic/gin"
	"strconv"
)

// FavoriteAction
// @Tags favorite
// @Summary FavoriteAction
// @Description Test API
// @Produce application/json
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
	err = svc.FavouriteAction(userId, strconv.FormatInt(params.VideoId, 10), params.ActionType)

	if err != nil {
		log.Debug("Favorite operation error", err)
		response.FailWithCode(common.ERROR, c)
		return
	}

	response.Success(c)
}

//// GetFavouriteList 获取点赞列表;
//func GetFavouriteList(c *gin.Context) {
//	strUserId := c.Query("user_id")
//	strCurId := c.GetString("userId")
//	userId, _ := strconv.ParseInt(strUserId, 10, 64)
//	curId, _ := strconv.ParseInt(strCurId, 10, 64)
//	like := GetVideo()
//	videos, err := like.GetFavouriteList(userId, curId)
//	if err == nil {
//		log.Printf("方法like.GetFavouriteList(userid) 成功")
//		c.JSON(http.StatusOK, GetFavouriteListResponse{
//			StatusCode: 0,
//			StatusMsg:  "get favouriteList success",
//			VideoList:  videos,
//		})
//	} else {
//		log.Printf("方法like.GetFavouriteList(userid) 失败：%v", err)
//		c.JSON(http.StatusOK, GetFavouriteListResponse{
//			StatusCode: 1,
//			StatusMsg:  "get favouriteList fail ",
//		})
//	}
//}
