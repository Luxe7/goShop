package goods

import (
	"context"
	"go.uber.org/zap"
	"goShop_Web/proto"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"goShop_Web/global"
)

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

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
func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func List(ctx *gin.Context) {
	req := &proto.GoodsFilterRequest{}

	pmin, _ := strconv.Atoi(ctx.DefaultQuery("pmin", "0"))
	req.PriceMax = int32(pmin)

	pmax, _ := strconv.Atoi(ctx.DefaultQuery("pmax", "0"))
	req.PriceMax = int32(pmax)

	ih := ctx.DefaultQuery("ih", "0")
	if ih == "1" {
		req.IsHot = true
	}

	in := ctx.DefaultQuery("in", "0")
	if in == "1" {
		req.IsNew = true
	}

	it := ctx.DefaultQuery("it", "0")
	if it == "1" {
		req.IsTab = true
	}

	c, _ := strconv.Atoi(ctx.DefaultQuery("c", "0"))
	req.TopCategory = int32(c)

	pn, _ := strconv.Atoi(ctx.DefaultQuery("pn", "0"))
	req.TopCategory = int32(pn)

	pnum, _ := strconv.Atoi(ctx.DefaultQuery("pnum", "0"))
	req.TopCategory = int32(pnum)

	req.KeyWords = ctx.DefaultQuery("q", "0")

	b, _ := strconv.Atoi(ctx.DefaultQuery("b", "0"))
	req.Brand = int32(b)

	list, err := global.GoodsSrvClient.GoodsList(context.Background(), req)
	if err != nil {
		zap.S().Errorw("[List] 查询 【用户列表】失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	reMap := map[string]interface{}{
		"data": list.Data,
	}
	goodsList := make([]interface{}, 0)
	for _, value := range list.Data {
		goodsList = append(goodsList, map[string]interface{}{
			"id":          value.Id,
			"name":        value.Name,
			"goods_brief": value.GoodsBrief,
			"desc":        value.GoodsDesc,
			"ship_free":   value.ShipFree,
			"images":      value.Images,
			"desc_images": value.DescImages,
			"front_image": value.GoodsFrontImage,
			"shop_price":  value.ShopPrice,
			"category": map[string]interface{}{
				"id":   value.Category.Id,
				"name": value.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   value.Brand.Id,
				"name": value.Brand.Name,
				"logo": value.Brand.Logo,
			},
			"is_hot":  value.IsHot,
			"is_new":  value.IsNew,
			"on_sale": value.OnSale,
		})
	}
	reMap["data"] = goodsList
	ctx.JSON(http.StatusOK, reMap)
}
