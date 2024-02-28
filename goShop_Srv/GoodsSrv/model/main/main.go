package main

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"goShop/GoodsSrv/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

/*
对用户密码要进行加密，防止彩虹表撞库，我们要对密码进行加盐
*/

func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}

func main() {

	// Using custom options
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode("generic password", options)
	genPassword := fmt.Sprintf("pbkdf2-sha512$%s$%s", salt, encodedPwd)

	passwordInfo := strings.Split(genPassword, "$")
	check := password.Verify("generic password", passwordInfo[1], passwordInfo[2], options)
	fmt.Println(check) // true

	dsn := "root:123456@tcp(127.0.0.1:3306)/goshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
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
	db.AutoMigrate(&model.User{})
	for i := 0; i < 10; i++ {
		user := model.User{
			Nickname: fmt.Sprintf("bobby%d", i),
			Mobil:    fmt.Sprintf("1878222222%d", i),
			Password: "newPassword",
		}
		db.Save(&user)
	}
}
