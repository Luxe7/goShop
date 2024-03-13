package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"goShop/GoodsSrv/model"
)

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/goshop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	// 全局模式
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&model.Category{}, &model.Brands{}, &model.GoodsCategoryBrand{}, &model.Banner{}, &model.Goods{})
}

func Mysql2Es() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/goshop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	// 全局模式
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	host := fmt.Sprintf("https://%s:%d", "192.168.171.130", 9200)
	logger := log.New(os.Stdout, "goShop", log.LstdFlags)

	EsClient, err := elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false), elastic.SetTraceLog(logger))
	if err != nil {
		panic(err)
	}

	var goods []model.Goods
	db.Find(&goods)
	for _, good := range goods {
		esModel := model.EsGoods{
			ID:          good.ID,
			CategoryID:  good.CategoryID,
			BrandsID:    good.BrandsID,
			OnSale:      good.OnSale,
			ShipFree:    good.ShipFree,
			IsNew:       good.IsNew,
			IsHot:       good.IsHot,
			Name:        good.Name,
			ClickNum:    good.ClickNum,
			SoldNum:     good.SoldNum,
			FavNum:      good.FavNum,
			MarketPrice: good.MarketPrice,
			GoodsBrief:  good.GoodsBrief,
			ShopPrice:   good.ShopPrice,
		}
		_, err = EsClient.Index().Index(esModel.GetIndexName()).BodyJson(esModel).Id(strconv.Itoa(int(good.ID))).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}

}
