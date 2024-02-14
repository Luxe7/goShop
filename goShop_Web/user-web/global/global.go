package global

import (
	ut "github.com/go-playground/universal-translator"
	"goShop_Web/config"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	Trans        ut.Translator
)
