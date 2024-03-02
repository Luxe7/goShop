package initialize

import (
	"github.com/gin-gonic/gin"
	"goShop_Web/middlewares"
	"goShop_Web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	//配置跨域
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("/g/v1")
	router.InitGoodsRouter(ApiGroup)
	return Router
}
