package global

import (
	"goShop/UserSrv/config"

	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
)

//func init() {
//	c := ServerConfig.MysqlInfo
//	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
//		c.User, c.Password, c.Host, c.Port, c.Name)
//	//dsn := "root:123456@tcp(127.0.0.1:3306)/goshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
//	newLogger := logger.New(
//		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
//		logger.Config{
//			SlowThreshold:             time.Second,   // Slow SQL threshold
//			LogLevel:                  logger.Silent, // Log level Info会打印出所有的SQL语句
//			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
//			ParameterizedQueries:      true,          // Don't include params in the SQL log
//			Colorful:                  false,         // Disable color
//		},
//	)
//	var err error
//	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
//		NamingStrategy: schema.NamingStrategy{
//			SingularTable: true,
//		},
//		Logger: newLogger,
//	})
//	if err != nil {
//		panic(err)
//	}
//
//}
