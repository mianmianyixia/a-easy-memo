package database

import (
	"a-easy-memo/internal/model"
	"fmt"

	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DB(name, password string) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/go_test?charset=utf8mb4&parseTime=True&loc=Local", name, password)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&model.Member{}, &model.Task{})
	if err != nil {
		log.Fatal("自动迁移失败: ", err)
		return nil
	}
	log.Println("自动迁移成功")
	return db
}
