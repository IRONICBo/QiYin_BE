package jwt

import (
	"github.com/IRONICBo/QiYin_BE/internal/common"
	"github.com/IRONICBo/QiYin_BE/internal/common/response"
	"github.com/IRONICBo/QiYin_BE/internal/utils"
	"github.com/gin-gonic/gin"
)

// Auth 鉴权中间件
// 若用户携带的token正确,解析token,将userId放入上下文context中并放行;否则,返回错误信息
func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		//auth := c.PostForm("token")
		token := context.GetHeader("token")
		if len(token) == 0 {
			token = context.PostForm("token")
			if len(token) == 0 {
				context.Abort()
				response.FailWithCode(common.UNAUTHORIZED, context)
				return
			}
		}
		parseToken, err := utils.ParseJwtToken(token)
		if err != nil {
			context.Abort()
			response.FailWithCode(common.ERROR, context)
			return
		}
		context.Set("userId", parseToken.UserUUID)
		context.Next()
	}
}
