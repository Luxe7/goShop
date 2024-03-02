package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goShop_Web/api/goods"
)

func InitGoodsRouter(router *gin.RouterGroup) {
	GoodsRouter := router.Group("goods")
	zap.S().Debug("login the router of goods")
	{
		GoodsRouter.GET("", goods.List)
	}
}
