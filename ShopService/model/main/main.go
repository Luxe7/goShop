package main

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"io"
	"strings"
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

	//dsn := "root:123456@tcp(127.0.0.1:3306)/goshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	//	logger.Config{
	//		SlowThreshold:             time.Second,   // Slow SQL threshold
	//		LogLevel:                  logger.Silent, // Log level Info会打印出所有的SQL语句
	//		IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
	//		ParameterizedQueries:      true,          // Don't include params in the SQL log
	//		Colorful:                  false,         // Disable color
	//	},
	//)
	//
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	//	NamingStrategy: schema.NamingStrategy{
	//		SingularTable: true,
	//	},
	//	Logger: newLogger,
	//})
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = db.AutoMigrate(&model.User{})
	//if err != nil {
	//	panic(err)
	//}
	// Using the default options

	// Using custom options
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode("generic password", options)
	genPassword := fmt.Sprintf("pbkdf2-sha512$%s$%s", salt, encodedPwd)
	passwordInfo := strings.Split(genPassword, "$")
	check := password.Verify("generic password", passwordInfo[1], passwordInfo[2], options)
	fmt.Println(check) // true

}
