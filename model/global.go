package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"leadboard/config"
)

var DB *gorm.DB

func BuildConnection(cfg config.Config) {
	dsn := cfg.DbUserName + ":" + cfg.DbPassword + "@(" + cfg.DbIP +
		")/" + cfg.DbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	MigrateModels()
}

func MigrateModels() {
	err := DB.AutoMigrate(&User{}, &Submission{})
	if err != nil {
		panic(err)
	}
}
