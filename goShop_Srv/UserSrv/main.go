package main

import (
	"flag"
	"fmt"
	"goShop/UserSrv/utils/register/consul"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"goShop/UserSrv/global"
	"goShop/UserSrv/handler"
	"goShop/UserSrv/initialize"
	"goShop/UserSrv/proto"
	"goShop/UserSrv/utils"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	IP := flag.String("ip", "10.66.58.204", "IP地址")
	Port := flag.Int("port", 0, "端口号")

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	zap.S().Info(global.ServerConfig)

	flag.Parse()
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}
	zap.S().Info("Ip:", *IP)
	zap.S().Info("Port:", *Port)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	//服务注册
	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err = registerClient.Register(global.ServerConfig.Host, *Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败:", err.Error())
	}
	zap.S().Debugf("启动服务器, 端口： %d", *Port)

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = registerClient.DeRegister(serviceId); err != nil {
		zap.S().Info("注销失败:", err.Error())
	} else {
		zap.S().Info("注销成功:")
	}
}
