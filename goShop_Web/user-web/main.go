package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"goShop_Web/global"
	myvalidator "goShop_Web/validator"

	"go.uber.org/zap"
	"goShop_Web/initialize"
)

func main() {
	//初始化logger
	initialize.InitLogger()
	//初始化配置文件
	initialize.InitConfig()
	//初始化routers
	Router := initialize.Routers()
	//初始化翻译
	_ = initialize.InitTrans("zh")
	//初始化用户服务连接
	initialize.InitSrvConn()
	////获取可用端口，在开发的时候不需要，部署的时候需要在可用端口下监听
	//port, err := utils.GetFreePort()
	//if err == nil {
	//	global.ServerConfig.Port = port
	//}
	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0}手机号格式不正确!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	//logger, _ := zap.NewProduction()
	//defer logger.Sync() //这个应该是刷入缓存的
	//sugar := logger.Sugar()

	zap.S().Debug("Run the Web server，Port：%d", global.ServerConfig.Port) //它返回的就是一个全局的sugar logger，和上边的代码是等价的

	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panicf("启动失败", err.Error())
	}
}
