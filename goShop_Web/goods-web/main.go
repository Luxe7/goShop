package main

import (
	"fmt"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"

	"goShop_Web/global"
	"goShop_Web/initialize"
	"goShop_Web/utils/register/consul"
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
	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err := registerClient.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败:", err.Error())
	}
	zap.S().Debug("Run the Web server，Port：%d", global.ServerConfig.Port) //它返回的就是一个全局的sugar logger，和上边的代码是等价的

	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panicf("启动失败", err.Error())
	}
}
