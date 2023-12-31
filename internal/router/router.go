package router

import (
	"github.com/IRONICBo/QiYin_BE/internal/api"
	"github.com/IRONICBo/QiYin_BE/internal/middleware/jwt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// swagger docs.
	_ "github.com/IRONICBo/QiYin_BE/docs"

	"github.com/gin-gonic/gin"

	"github.com/IRONICBo/QiYin_BE/internal/config"
	urltrie "github.com/IRONICBo/QiYin_BE/internal/middleware/hooks/url_trie"
)

// InitRouter init router.
//
//nolint:funlen
func InitRouter() *gin.Engine {
	if config.GetString("app.debug") == "true" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	// Set MaxMultipartMemory
	r.MaxMultipartMemory = int64(config.Config.Server.MaxFileSize) << 20

	// Enable Hooks
	r.Use(urltrie.RunHook())

	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiv1 := r.Group("/api/v1")
	{
		apiv1.GET("/ping")
		apiv1.POST("/login", api.UserLogin)
		apiv1.POST("/register", api.UserRegister)
		apiv1.GET("/userinfo", api.UserInfo)
		apiv1.GET("/check", jwt.Auth(), api.CheckToken) // 根据token 查看是否登录
		apiv1.GET("/searchUser", api.SearchUser)
		apiv1.POST("/setStyle", jwt.Auth(), api.SetStyle)
		apiv1.POST("/setUser", jwt.Auth(), api.SetUser)

		// 点赞
		favorite := apiv1.Group("/favorite/")
		{
			favorite.POST("/action", jwt.Auth(), api.FavoriteAction)
			favorite.GET("/list", api.GetFavoriteList)
		}

		// 收藏
		collection := apiv1.Group("/collection/")
		{
			collection.POST("/action", jwt.Auth(), api.CollectionAction)
			collection.GET("/list", api.GetCollectionList)
		}

		// 评论
		comment := apiv1.Group("/comment/")
		{
			// 评论列表
			comment.GET("/list", api.CommentList)
			comment.POST("/delete", jwt.Auth(), api.CommentDelete)
			comment.POST("/add", jwt.Auth(), api.CommentAdd)
		}

		video := apiv1.Group("/video/")
		{
			video.GET("/search", api.Search)
			video.GET("/searchTag", api.SearchTag)
			video.GET("/hots", api.GetHots)
			video.GET("/list", api.GetVideos)
			video.GET("/lists", api.GetVideosList)
			video.POST("/upload", jwt.Auth(), api.UploadVideo)
			video.POST("/save", api.SaveVideoHis)
			video.GET("/getHistory", jwt.Auth(), api.GetHistory)
			video.GET("/getVideo", api.GetVideo)
		}

		// QiNiu
		qiniu := apiv1.Group("/qiniu")
		{
			qiniu.POST("/token", api.GetUploadToken)
			qiniu.POST("/pfop/callback", api.GetPfopCallback)
			qiniu.GET("/proxy", api.GetImageByProxy)
		}
	}

	return r
}
