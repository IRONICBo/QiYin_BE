package router

import (
	"github.com/IRONICBo/QiYin_BE/internal/api"
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
	}

	return r
}
