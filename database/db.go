package database

import (
	"a-easy-memo/internal/model"
	"a-easy-memo/zlog"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DB(name, password string) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(mysql:3306)/go_test?charset=utf8mb4&parseTime=True&loc=Local", name, password)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zlog.Fatal("连接失败", zap.Error(err))
	}
	err = db.AutoMigrate(&model.Member{}, &model.Task{})
	if err != nil {
		zlog.Fatal("自动迁移失败", zap.Error(err))
	}
	zlog.Info("自动迁移成功")
	return db
}
