package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        uint32    `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

type User struct {
	BaseModel
	Mobil    string     `gorm:"index:index_mobil;unique;varchar(11);not null"`
	Password string     `gorm:"type:varchar(100);not null"`
	Nickname string     `gorm:"type:varchar(20)"`
	Birthday *time.Time `gorm:"type:datetime"`
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6)"`
	Role     uint32     `gorm:"column:role;default:1;type:int comment'1表示普通用户，2表示管理员'"`
}
