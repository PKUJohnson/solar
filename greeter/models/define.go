package models

import (
	std "github.com/PKUJohnson/solar/std"
	"github.com/PKUJohnson/solar/std/toolkit"
	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
	dw *std.DBExtension
)

func InitModel(config std.ConfigMysql) {
	db = toolkit.CreateDB(config)
	dw = std.NewDBWrapper(db)
	std.LogInfoLn("start init mysql model")
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
