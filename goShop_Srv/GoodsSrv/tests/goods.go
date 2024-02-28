package main

import (
	"context"
	"goShop/GoodsSrv/proto"
	"google.golang.org/grpc"
)

var goodsClient proto.GoodsClient
var conn *grpc.ClientConn

func Init() {

	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	goodsClient = proto.NewGoodsClient(conn)
}
func TestGetUserList() {
	rsp, err := goodsClient.BrandList(context.Background(), &proto.BrandFilterRequest{})

	if err != nil {
		panic(err)
	}
	for _, brand := range rsp.Data {
		println(brand.Name)
	}
}

func main() {
	Init()
	TestGetUserList()
	//TestCreateUser()
	conn.Close()
}
