package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"goShop_Web/forms"
	"goShop_Web/global"
	"goShop_Web/global/response"
	"goShop_Web/middlewares"
	"goShop_Web/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"strings"
	"time"

	"goShop_Web/proto"
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
func GetUserList(ctx *gin.Context) {
	zap.S().Debug("获取用户列表页")
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvConfig.Host, global.ServerConfig.UserSrvConfig.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] connect gRpc failed", "msg", err.Error())
	}
	userSrvClient := proto.NewUserClient(userConn)

	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "0")
	pSizeInt, _ := strconv.Atoi(pSize)

	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		PN:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] require UserList failed")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	//result := make(map[string]interface{})
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		//data := make(map[string]interface{})

		user := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			Birthday: time.Unix(int64(value.Birthday), 0),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}
		//data["id"] = value.Id
		//data["name"] = value.NickName
		//data["birthday"] = value.Birthday
		//data["gender"] = value.Gender
		//data["mobile"] = value.Mobile
		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)
}
func PassWordLogin(c *gin.Context) {
	//表单验证
	zap.S().Infof("PassWordLogin")
	passwordLoginForm := forms.PassWordLoginForm{}
	if err := c.ShouldBind(&passwordLoginForm); err != nil {
		zap.S().Debug(err.Error())
		HandleValidatorError(c, err)
		return
	}
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvConfig.Host, global.ServerConfig.UserSrvConfig.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] connect gRpc failed", "msg", err.Error())
	}
	//生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(userConn)
	//登录的逻辑
	if resp, err := userSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{Mobile: passwordLoginForm.Mobile}); err != nil {
		zap.S().Errorw(err.Error())
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
					"msg": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{
					"msg": "登录失败",
				})
			}
		}
	} else {
		//只是查询了用户没有检查密码
		if passResp, passerr := userSrvClient.CheckPassWord(context.Background(), &proto.PassWordCheckInfo{
			Password:          passwordLoginForm.PassWord,
			EncryptedPassword: resp.PassWord,
		}); passerr != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"msg": "登录失败",
			})
		} else {
			if passResp.Success {
				//生成token
				j := middlewares.NewJWT()
				claim := models.CustomClaims{
					ID:          uint(resp.Id),
					NickName:    resp.NickName,
					AuthorityId: uint(resp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),               //生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
						Issuer:    "luxe7",
					},
				}
				token, err := j.CreateToken(claim)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "create token failed",
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"id":         resp.Id,
					"nick_name":  resp.NickName,
					"token":      token,
					"expires_at": (time.Now().Unix() + 60*60*24*30) * 1000,
				})
			} else {
				c.JSON(http.StatusBadRequest, map[string]string{
					"msg": "密码错误",
				})
			}
		}
	}

}
