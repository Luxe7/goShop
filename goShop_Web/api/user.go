package api

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"goShop_Web/proto"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	//将grpc的状态码转换为http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "Internal Errors",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "Invalid Argument",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "Other Errors",
				})
			}
			return
		}
	}
}
func GetUserList(ctx *gin.Context) {
	zap.S().Debug("获取用户列表页")
	ip := "127.0.0.1"
	port := 50051
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] connect gRpc failed", "msg", err.Error())
	}
	userSrvClient := proto.NewUserClient(userConn)
	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		PN:    0,
		PSize: 0,
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] require UserList failed")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	//result := make(map[string]interface{})
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		data := make(map[string]interface{})
		data["id"] = value.Id
		data["name"] = value.NickName
		data["birthday"] = value.Birthday
		data["gender"] = value.Gender
		data["mobile"] = value.Mobile
		result = append(result, data)
	}
	ctx.JSON(http.StatusOK, result)
}
