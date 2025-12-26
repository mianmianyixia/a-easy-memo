package routers

import (
	"a-easy-memo/internal/api/serve"
	"a-easy-memo/internal/config"
	"a-easy-memo/internal/dao"
	"a-easy-memo/internal/login"
	"a-easy-memo/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Routers() {
	r := gin.Default()
	db := dao.NewGorm(config.DB)
	res := dao.NewRedis(config.REDIS)
	account := r.Group("/account")
	{
		account.POST("/register", login.Register(db))
		account.POST("/", login.Login(db))
		account.DELETE("/", middleware.Middleware(), login.Del(db))
		account.PUT("/", middleware.Middleware(), login.Change(db))
	}
	task := r.Group("/task")
	{
		task.Use(middleware.Middleware())
		task.PUT("/", serve.Change(db, res, config.CTX))
		task.DELETE("/", serve.Del(db, res, config.CTX))
		task.POST("/", serve.Add(db, res, config.CTX))
		task.GET("/", serve.Find(db, res, config.CTX))
	}
	r.Run(":8080")
}
