package initialize

import (
	"github.com/gin-gonic/gin"

	"goShop_Web/goods-web/middlewares"
	"goShop_Web/goods-web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	//配置跨域
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("/g/v1")
	router.InitGoodsRouter(ApiGroup)
	router.InitCategoryRouter(ApiGroup)
	router.InitBrandRouter(ApiGroup)
	router.InitBannerRouter(ApiGroup)
	return Router
}
