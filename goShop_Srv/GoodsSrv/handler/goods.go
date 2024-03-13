package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"goShop/GoodsSrv/global"
	"goShop/GoodsSrv/model"
	"goShop/GoodsSrv/proto"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

func ModelToResponse(goods model.Goods) proto.GoodsInfoResponse {
	return proto.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brands.ID,
			Name: goods.Brands.Name,
			Logo: goods.Brands.Logo,
		},
	}
}
func (s *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	resp := &proto.GoodsListResponse{}

	q := elastic.NewBoolQuery()

	localDB := global.DB.Model(model.Goods{})
	if req.KeyWords != "" {
		//localDB = localDB.Where("name LIKE ?", "%"+req.KeyWords+"%")
		q = q.Must(elastic.NewMultiMatchQuery(req.KeyWords, "name", "goods_brief"))
	}
	if req.IsNew {
		//localDB = localDB.Where(model.Goods{IsNew: true})
		q = q.Filter(elastic.NewTermQuery("is_new", req.IsNew))
	}
	if req.IsHot {
		//localDB = localDB.Where(model.Goods{IsHot: true})
		q = q.Filter(elastic.NewTermQuery("is_hot", req.IsHot))
	}
	if req.PriceMin > 0 {
		//localDB = localDB.Where("shop_price >= ?", req.PriceMin)
		q = q.Must(elastic.NewRangeQuery("shop_price").Gte(req.PriceMin))
	}
	if req.PriceMax > 0 {
		//localDB = localDB.Where("shop_price <= ?", req.PriceMax)
		q = q.Must(elastic.NewRangeQuery("shop_price").Lte(req.PriceMax))
	}

	if req.Brand > 0 {
		//localDB = localDB.Where("brand_id=?", req.Brand)
		q = q.Filter(elastic.NewTermQuery("brand_id", req.Brand))
	}
	//通过不同目录去查询商品
	//通过category去查询商品
	var categoryIds []interface{}
	var subQuery string
	if req.TopCategory > 0 {
		var category model.Category
		if result := localDB.Preload("Category").Preload("Brands").First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}

		if category.Level == 1 {
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category WHERE parent_category_id=%d)", req.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("select id from category WHERE parent_category_id=%d", req.TopCategory)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("select id from category WHERE id=%d", req.TopCategory)
		}
		type result struct {
			id int32
		}
		var results []result

		global.DB.Model(model.Category{}).Raw(subQuery).Scan(&results)
		for _, r := range results {
			categoryIds = append(categoryIds, r.id)
		}
		//生成terms查询
		q = q.Filter(elastic.NewTermsQuery("category_id", categoryIds...))
		//localDB = localDB.Where(fmt.Sprintf("category_id in (%s)", subQuery))
	}

	if req.Pages == 0 {
		req.Pages = 1
	}
	switch {
	case req.PagePerNums > 100:
		req.PagePerNums = 100
	case req.PagePerNums <= 0:
		req.PagePerNums = 10
	}
	searchResult, err := global.EsClient.Search().Index(model.EsGoods{}.GetIndexName()).Query(q).From(int(req.Pages)).Size(int(req.PagePerNums)).Do(context.Background())
	if err != nil {
		return nil, err
	}
	resp.Total = int32(searchResult.Hits.TotalHits.Value)
	goodsIds := make([]int32, 0)
	for _, value := range searchResult.Hits.Hits {
		goods := model.EsGoods{}
		_ = json.Unmarshal(value.Source, &goods)
		goodsIds = append(goodsIds, goods.ID)
	}

	//var count int64
	//localDB.Count(&count)
	//resp.Total = int32(count)

	//查询id在某个数组中的值
	var goods []model.Goods
	result := localDB.Preload("Category").Preload("Brands").Find(&goods, goodsIds)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, good := range goods {
		goodsInfo := ModelToResponse(good)
		resp.Data = append(resp.Data, &goodsInfo)
	}
	return resp, nil
}

func (s *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	resp := &proto.GoodsListResponse{}
	var goods []model.Goods
	result := global.DB.Find(&goods, req.Ids)
	for _, good := range goods {
		goodInfo := ModelToResponse(good)
		resp.Data = append(resp.Data, &goodInfo)
	}
	resp.Total = int32(result.RowsAffected)
	return resp, nil
}

func (s *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}
	//先检查redis中是否有这个token
	//防止同一个token的数据重复插入到数据库中，如果redis中没有这个token则放入redis
	//这里没有看到图片文件是如何上传， 在微服务中 普通的文件上传已经不再使用
	goods := model.Goods{
		Brands:          brand,
		BrandsID:        brand.ID,
		Category:        category,
		CategoryID:      category.ID,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.GoodsBrief,
		ShipFree:        req.ShipFree,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		OnSale:          req.OnSale,
	}

	//srv之间互相调用了
	tx := global.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &proto.GoodsInfoResponse{
		Id: goods.ID,
	}, nil
}
func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Goods{BaseModel: model.BaseModel{ID: req.Id}}, req.Id); result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	var goods model.Goods

	if result := global.DB.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	goods.Brands = brand
	goods.BrandsID = brand.ID
	goods.Category = category
	goods.CategoryID = category.ID
	goods.Name = req.Name
	goods.GoodsSn = req.GoodsSn
	goods.MarketPrice = req.MarketPrice
	goods.ShopPrice = req.ShopPrice
	goods.GoodsBrief = req.GoodsBrief
	goods.ShipFree = req.ShipFree
	goods.Images = req.Images
	goods.DescImages = req.DescImages
	goods.GoodsFrontImage = req.GoodsFrontImage
	goods.IsNew = req.IsNew
	goods.IsHot = req.IsHot
	goods.OnSale = req.OnSale

	tx := global.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}
func (s *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var goods model.Goods

	if result := global.DB.Preload("Category").Preload("Brands").First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	goodsInfoResponse := ModelToResponse(goods)
	return &goodsInfoResponse, nil
}
