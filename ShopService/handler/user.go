package handler

import (
	"context"
	"goShop/ShopService/global"
	"goShop/ShopService/model"
	"goShop/ShopService/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserServer struct {
}

func ModelToResponse(user model.User) proto.UserInfoResponse {
	//在grpc的message中字段有默认值，不能随便nil进去，否则会报错
	//需要分清哪些有默认值
	UserInfoResponse := proto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		NickName: user.Nickname,
		Gender:   user.Gender,
		Role:     user.Role,
	}
	if user.Birthday != nil {
		UserInfoResponse.Birthday = uint64(user.Birthday.Unix())
	}
	return UserInfoResponse
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (s UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)

	global.DB.Scopes(Paginate(int(req.PN), int(req.PSize))).Find(&users)
	//scopes用于在查询数据库时应用一些预定义或自定义的条件或操作。
	//find是立即执行方法，表示根据前面的条件，查询数据库并将结果赋值给 users 这个切片变量

	for _, user := range users {
		userInfoRsp := ModelToResponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}
	return rsp, nil

}
func (s UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobil: req.Mobile}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}
func (s UserServer) GetUserByID(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil

}
