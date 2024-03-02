package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"goShop_Web/user-web/global"
	"goShop_Web/user-web/initialize"
	"goShop_Web/user-web/utils/register/consul"
	myvalidator "goShop_Web/user-web/validator"
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

	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err := registerClient.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败:", err.Error())
	}
	zap.S().Debug("Run the Web server，Port：%d", global.ServerConfig.Port)

	go func() {
		if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panicf("启动失败", err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := registerClient.DeRegister(serviceId); err != nil {
		zap.S().Info("注销失败", err)
	} else {
		zap.S().Info("注销成功")
	}

}
