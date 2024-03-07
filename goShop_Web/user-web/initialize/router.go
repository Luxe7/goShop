package initialize

import (
	"github.com/gin-gonic/gin"

	"goShop_Web/user-web/middlewares"
	"goShop_Web/user-web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	//配置跨域
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("/u/v1")
	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)
	return Router
}
