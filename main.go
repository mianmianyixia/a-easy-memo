package main

import (
	"a-easy-memo/internal/config"
	"a-easy-memo/internal/routers"
	"a-easy-memo/zlog"

	"go.uber.org/zap"
)

func main() {
	zlog.Info("登录数据库")
	err := config.Command.Execute()
	if err != nil {
		zlog.Warn("登录失败", zap.Error(err))
		return
	}
	routers.Routers()
}
