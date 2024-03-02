package handler

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"goShop/GoodsSrv/global"
	"goShop/GoodsSrv/model"
	"goShop/GoodsSrv/proto"
)

func (s *GoodsServer) GetAllCategorysList(ctx context.Context, req *emptypb.Empty) (*proto.CategoryListResponse, error) {
	resp := &proto.CategoryListResponse{}
	var categorys []model.Category
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)
	b, _ := json.Marshal(&categorys)
	resp.JsonData = string(b)
	//resp的其他字段都是空的
	return resp, nil
}

func (s *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	resp := &proto.SubCategoryListResponse{}
	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	resp.Info = &proto.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		ParentCategory: category.ParentCategoryID,
		Level:          category.Level,
		IsTab:          category.IsTab,
	}
	var categorys []model.Category
	global.DB.Where(&model.Category{Level: 1}).Find(&categorys)
	for _, category := range categorys {
		resp.SubCategorys = append(resp.SubCategorys, &proto.CategoryInfoResponse{
			Id:             category.ID,
			Name:           category.Name,
			ParentCategory: category.ParentCategoryID,
			Level:          category.Level,
			IsTab:          category.IsTab,
		})
	}
	return resp, nil

}
func (s *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	resp := &proto.CategoryInfoResponse{}
	if result := global.DB.First(&model.Category{}, req.Id); result.RowsAffected != 0 {
		return nil, status.Errorf(codes.AlreadyExists, "已存在该商品目录")
	}
	category := model.Category{
		Name:             req.Name,
		ParentCategoryID: req.ParentCategory,
		Level:            req.Level,
		IsTab:            req.IsTab,
	}
	global.DB.Save(&category)
	resp.Id = category.ID
	return resp, nil
}
func (s *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Category{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	return &emptypb.Empty{}, nil
}
func (s *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Category{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	category := model.Category{}
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentCategory != 0 {
		category.ParentCategoryID = req.ParentCategory
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	if req.IsTab {
		category.IsTab = req.IsTab
	}

	global.DB.Save(&category)

	return &emptypb.Empty{}, nil
}
