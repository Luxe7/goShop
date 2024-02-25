package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goShop_Web/api"
	"goShop_Web/middlewares"
)

func InitUserRouter(router *gin.RouterGroup) {
	UserRouter := router.Group("user")
	zap.S().Debug("login the router of user")
	{
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("pwd_login", api.PassWordLogin)
	}
}
