package config

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
type JwtConfig struct {
	SigningKey string `mapstructure:"key"`
}
type AliSmsConfig struct {
	ApiKey     string `mapstructure:"key" json:"key"`
	ApiSecrect string `mapstructure:"secrect" json:"secrect"`
}

type RedisConfig struct {
	Host   string `mapstructure:"host" json:"host"`
	Port   int    `mapstructure:"port" json:"port"`
	Expire int    `mapstructure:"expire" json:"expire"`
}
type ServerConfig struct {
	Name          string        `mapstructure:"name"`
	Port          int           `mapstructure:"port"`
	UserSrvConfig UserSrvConfig `mapstructure:"user_srv"`
	JWTInfo       JwtConfig     `mapstructure:"jwt"`
	AliSmsInfo    AliSmsConfig  `mapstructure:"sms" json:"sms"`
	RedisInfo     RedisConfig   `mapstructure:"redis"`
}
