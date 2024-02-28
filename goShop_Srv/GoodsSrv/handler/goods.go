package handler

type GoodsServer struct {

}

// 商品接口
GoodsList(context.Context, *GoodsFilterRequest) (*GoodsListResponse, error)
// 现在用户提交订单有多个商品，你得批量查询商品的信息吧
BatchGetGoods(context.Context, *BatchGoodsIdInfo) (*GoodsListResponse, error)
CreateGoods(context.Context, *CreateGoodsInfo) (*GoodsInfoResponse, error)
DeleteGoods(context.Context, *DeleteGoodsInfo) (*emptypb.Empty, error)
UpdateGoods(context.Context, *CreateGoodsInfo) (*emptypb.Empty, error)
GetGoodsDetail(context.Context, *GoodInfoRequest) (*GoodsInfoResponse, error)

