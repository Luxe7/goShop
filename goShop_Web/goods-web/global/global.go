package global

import (
	ut "github.com/go-playground/universal-translator"
	"goShop_Web/config"
	"goShop_Web/proto"
)

var (
	ServerConfig   *config.ServerConfig = &config.ServerConfig{}
	Trans          ut.Translator
	GoodsSrvClient proto.GoodsClient
	NacosConfig    *config.NacosConfig = &config.NacosConfig{}
)
