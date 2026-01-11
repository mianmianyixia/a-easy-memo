package routers

import (
	"a-easy-memo/internal/api/api"
	"a-easy-memo/internal/config"
	"a-easy-memo/internal/dao"
	"a-easy-memo/internal/login"
	"a-easy-memo/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Routers() {
	r := gin.Default()
	db := dao.NewGorm(config.DB)
	res := dao.NewRedis(config.REDIS, config.CTX)
	r.Use(middleware.ZapLog())
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
		task.PUT("/", api.Change(db, res))
		task.DELETE("/", api.Delete(db, res))
		task.POST("/", api.Change(db, res))
		task.GET("/", api.Find(db, res))
		task.POST("/save", api.Save(db, res))
	}
	r.Run(":8080")
}
