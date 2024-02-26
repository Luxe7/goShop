package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"goShop_Web/global"
	"goShop_Web/proto"
	"google.golang.org/grpc"
)

func InitSrvConn() {
	//从注册中心获取用户服务的信息
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulConfig.Host, global.ServerConfig.ConsulConfig.Port)
	userSrvHost := ""
	userSrvPort := 0

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service=="%s"`, global.ServerConfig.UserSrvConfig.Name))
	if err != nil {
		panic(err)
	}
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}
	if userSrvHost == "" {
		zap.S().Fatalf("[InitSrvConn] failed ")
	}
	//grpc拨号
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] connect gRpc failed", "msg", err.Error())
	}
	//后续的用户服务下线或端口或ip修改 - 负载均衡来做
	//事先建立好的连接，一直可用，不用进行多次TCP的连接
	//一个连接，多个groutine共用，性能问题？-连接池
	global.UserSrvClient = proto.NewUserClient(userConn)
}
