package handler

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/emptypb"

	"goShop/GoodsSrv/global"
	"goShop/GoodsSrv/model"
	"goShop/GoodsSrv/proto"
)

func (s *GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	resp := &proto.BrandListResponse{}

	var brands []model.Brands
	result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)
	if result.Error != nil {
		zap.S().Error("[BrandList]", result.Error)
		return nil, status.Error(codes.Internal, "DB find fail")
	}
	var total int64
	global.DB.Model(&model.Brands{}).Count(&total)

	var brandResponses []*proto.BrandInfoResponse
	for _, brand := range brands {
		brandResponses = append(brandResponses, &proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}
	resp.Total = int32(total)
	resp.Data = brandResponses
	return resp, nil

}
func (s *GoodsServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	resp := &proto.BrandInfoResponse{}
	if result := global.DB.First(&model.Brands{}, req.Id); result.RowsAffected != 0 {
		return nil, status.Error(codes.InvalidArgument, "brand is exist")
	}
	brand := model.Brands{
		Name: req.Name,
		Logo: req.Logo,
	}
	global.DB.Save(brand)
	resp.Id = brand.ID
	return resp, nil

}
func (s *GoodsServer) DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Brands{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "brand is not exist")
	}
	return nil, nil
}
func (s *GoodsServer) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	brand := model.Brands{}
	if result := global.DB.First(&brand, req.Id); result.RowsAffected != 0 {
		return nil, status.Error(codes.InvalidArgument, "brand is exist")
	}
	if req.Name != "" {
		brand.Name = req.Name
	}
	if req.Logo != "" {
		brand.Logo = req.Logo
	}
	global.DB.Save(&brand)
	return nil, nil
}
