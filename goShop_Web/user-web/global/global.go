package global

import (
	ut "github.com/go-playground/universal-translator"
	"goShop_Web/config"
	"goShop_Web/proto"
)

var (
	ServerConfig  *config.ServerConfig = &config.ServerConfig{}
	Trans         ut.Translator
	UserSrvClient proto.UserClient
	NacosConfig   *config.NacosConfig = &config.NacosConfig{}
)
