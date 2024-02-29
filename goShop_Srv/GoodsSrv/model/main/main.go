package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"goShop/GoodsSrv/model"
)

/*
对用户密码要进行加密，防止彩虹表撞库，我们要对密码进行加盐
*/

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
