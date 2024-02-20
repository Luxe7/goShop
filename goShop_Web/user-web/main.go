package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"goShop_Web/global"

	"go.uber.org/zap"
	"goShop_Web/initialize"
	myvalidator "goShop_Web/validator"
)

func main() {
	//初始化logger
	initialize.InitLogger()
	//初始化routers
	Router := initialize.Routers()
	initialize.InitConfig()
	_ = initialize.InitTrans("zh")
	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("mobile", myvalidator.ValidateMobile)
	}

	//logger, _ := zap.NewProduction()
	//defer logger.Sync() //这个应该是刷入缓存的
	//sugar := logger.Sugar()

	zap.S().Debug("Run the Web server，Port：%d", global.ServerConfig.Port) //它返回的就是一个全局的sugar logger，和上边的代码是等价的

	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panicf("启动失败", err.Error())
	}
}
