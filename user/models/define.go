package models

import (
	"time"
	"github.com/jinzhu/gorm"
	std "github.com/PKUJohnson/solar/std"
	"github.com/PKUJohnson/solar/std/toolkit"
)

var (
	db *gorm.DB
	dw *std.DBExtension
)

func InitModel(config std.ConfigMysql) {
	db = toolkit.CreateDB(config)
	dw = std.NewDBWrapper(db)
	std.LogInfoLn("start init mysql model")
	db.AutoMigrate(&User{})
	std.LogInfoLn("end init mysql model")
}

func DB() *std.DBExtension {
	return dw
}

func Session() *gorm.DB {
	return db.Begin()
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

type User struct {
	UserId               int64        `gorm:"column:user_id; type:bigint(20); primary_key" json:"user_id"`
	AppId                string       `gorm:"column:app_id; type:varchar(32)" json:"app_id"`
	SessionKey           string       `gorm:"column:session_key; type:varchar(255)" json:"session_key"`
	OpenId               string       `gorm:"column:open_id; type:varchar(255)" json:"open_id"`
	UnionId              string       `gorm:"column:union_id; type:varchar(255)" json:"union_id"`
	UserName             string       `gorm:"column:user_name; type:varchar(255)" json:"user_name"`
	UserDisplayName      string       `gorm:"column:user_display_name; type:varchar(255)" json:"user_display_name"`
	Gender               string       `gorm:"column:gender; type:varchar(255)" json:"gender"`
	Mobile               string       `gorm:"column:mobile; type:varchar(255)" json:"mobile"`
	City                 string       `gorm:"column:city; type:varchar(255)" json:"city"`
	Province             string       `gorm:"column:province; type:varchar(255)" json:"province"`
	Country              string       `gorm:"column:country; type:varchar(255)" json:"country"`
	Language             string       `gorm:"column:language; type:varchar(255)" json:"language"`
	Email                string       `gorm:"column:email; type:varchar(255)" json:"email"`
	Signature            string       `gorm:"column:signature; type:varchar(1024)" json:"signature"`
	Expiration           int64        `gorm:"column:expiration; type:int(11)" json:"expiration"`
	FUnionId             string       `gorm:"column:f_union_id; type:varchar(64)" json:"f_union_id"`
	Avatar               string       `gorm:"column:avatar; type:varchar(512)" json:"avatar"`
	CreatedAt            time.Time    `gorm:"column:created_at; type:timestamp(6)" json:"created_at"`
	UpdatedAt            time.Time    `gorm:"column:updated_at; type:timestamp(6)" json:"updated_at"`
}

func (User) TableName() string {
	return "solar.solar_users"
}

type FUser struct {
	FUnionId             string       `gorm:"column:f_union_id; type:varchar(64)" json:"f_union_id"`
	Mobile               string       `gorm:"column:mobile; type:varchar(255)" json:"mobile"`
	Password             string       `gorm:"column:password; type:varchar(255)" json:"password"`
}

func (FUser) TableName() string {
	return "solar.f_users"
}
