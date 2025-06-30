package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := "root:password@tcp(localhost:3306)/call_center_db?charset=utf8mb4&parseTime=True&loc=Asia%2FHo_Chi_Minh"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to db: " + err.Error())
	}
	return db
}