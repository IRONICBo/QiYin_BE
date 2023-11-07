package api

import (
	"github.com/IRONICBo/QiYin_BE/internal/common"
	"github.com/IRONICBo/QiYin_BE/internal/common/response"
	"github.com/IRONICBo/QiYin_BE/internal/dal/dao"
	requestparams "github.com/IRONICBo/QiYin_BE/internal/params/request"
	"github.com/IRONICBo/QiYin_BE/internal/service"
	"github.com/IRONICBo/QiYin_BE/pkg/log"
	"github.com/gin-gonic/gin"
)

// UserLogin
// @Tags user
// @Summary UserLogin
// @Description user login
// @Produce application/json
// @Param data body requestparams.UserParams true "UserParams"
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/login [post].
func UserLogin(c *gin.Context) {
	var params requestparams.UserParams
	err := c.ShouldBindJSON(&params)
	if err != nil {
		response.FailWithCode(common.INVALID_PARAMS, c)
		return
	}

	svc := service.NewUserService(c)
	u, err := svc.Login(&params)
	if err != nil {
		log.Debug("Login error", err)
		response.FailWithCode(common.ERROR, c)
		return
	}

	response.SuccessWithData(u, c)
}

// UserRegister
// @Tags user
// @Summary UserRegister
// @Description user register
// @Produce application/json
// @Param data body requestparams.UserParams true "UserParams"
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/register [post].
func UserRegister(c *gin.Context) {
	var params requestparams.UserParams
	err := c.ShouldBindJSON(&params)
	if err != nil {
		response.FailWithCode(common.INVALID_PARAMS, c)
		return
	}

	svc := service.NewUserService(c)
	err = svc.Register(&params)
	if err != nil {
		log.Debug("register error", err)
		response.FailWithCode(common.ERROR, c)
		return
	}

	response.Success(c)
}

// UserInfo
// @Tags user
// @Summary UserInfo
// @Description get userinfo by id
// @Produce application/json
// @Param userId query string true "query user id"
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/userinfo [get].
func UserInfo(c *gin.Context) {
	userId := c.Query("userId")

	svc := service.NewUserService(c)
	u, err := svc.GetUserById(userId)
	if err != nil {
		log.Debug("user doesn't exit", err)
		response.FailWithCode(common.ERROR, c)
		return
	}
	response.SuccessWithData(u, c)
}

// CheckToken
// @Tags user
// @Summary CheckToken
// @Description Check whether the token is valid
// @Produce application/json
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/check [get].
func CheckToken(c *gin.Context) {
	userId := c.GetString("userId")
	if len(userId) == 0 {
		response.FailWithData(false, c)
		return
	}
	//	 简单的根据id查找用户信息 并返回
	user, err := dao.GetTableUserById(userId)
	if err != nil {
		response.FailWithData(false, c)
		return
	}
	response.SuccessWithData(user, c)
}

// SearchUser
// @Tags user
// @Summary SearchUser
// @Description search user by name
// @Produce application/json
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/searchUser [get].
func SearchUser(c *gin.Context) {
	searchValue := c.Query("searchValue")

	svc := service.NewUserService(c)
	u, err := svc.Search(searchValue)
	if err != nil {
		log.Debug("user doesn't exit", err)
		response.FailWithCode(common.ERROR, c)
		return
	}
	response.SuccessWithData(u, c)
}

// SetStyle
// @Tags user
// @Summary SetStyle
// @Description set user style
// @Produce application/json
// @Param data body requestparams.StyleParams true "StyleParams"
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/setStyle [post].
func SetStyle(c *gin.Context) {
	var params requestparams.StyleParams
	err := c.ShouldBindJSON(&params)
	if err != nil {
		response.FailWithCode(common.INVALID_PARAMS, c)
		return
	}

	// 只有登录的用户才可以选择风格
	userId := c.GetString("userId")
	if len(userId) == 0 {
		response.FailWithData(false, c)
		return
	}

	err = dao.SetStyle(userId, params.Style)
	if err != nil {
		log.Debug("user doesn't exit", err)
		response.FailWithCode(common.ERROR, c)
		return
	}
	response.Success(c)
}

// SetUser
// @Tags user
// @Summary SetUser
// @Description update user info
// @Produce application/json
// @Param data body requestparams.UserInfoParams true "UserInfoParams"
// @Success 200 {object}  response.Response{msg=string} "Success"
// @Router /api/v1/setUser [post].
func SetUser(c *gin.Context) {
	var params requestparams.UserInfoParams
	err := c.ShouldBindJSON(&params)
	if err != nil {
		response.FailWithCode(common.INVALID_PARAMS, c)
		return
	}

	// 只有登录的用户才可以修改自身用户信息
	userId := c.GetString("userId")
	if len(userId) == 0 {
		response.FailWithData(false, c)
		return
	}

	svc := service.NewUserService(c)
	err = svc.UpdateUser(userId, &params)
	if err != nil {
		log.Debug("user update failed", err)
		response.FailWithCode(common.ERROR, c)
		return
	}
	response.Success(c)
}
