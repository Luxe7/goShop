package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goShop_Web/api"
)

func InitUserRouter(router *gin.RouterGroup) {
	UserRouter := router.Group("user")
	zap.S().Debug("login the router of user")
	{
		UserRouter.GET("list", api.GetUserList)
		UserRouter.POST("pwd_login", api.PassWordLogin)
	}
}
