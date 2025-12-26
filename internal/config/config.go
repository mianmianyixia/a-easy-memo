package config

import (
	"a-easy-memo/database"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type Config struct {
	Username string
	Password string
}

var Account Config
var Command *cobra.Command
var DB *gorm.DB
var REDIS *redis.Client
var CTX context.Context

func init() {
	var rootCmd = &cobra.Command{
		Use:   "login",
		Short: "数据库账号登录",
		Run: func(cmd *cobra.Command, args []string) {
			DB = database.DB(Account.Username, Account.Password)
			REDIS, CTX = database.ConnectRedis()
		},
	}
	rootCmd.Flags().StringVar(&Account.Username, "account", "", "输入账户")
	rootCmd.Flags().StringVar(&Account.Password, "password", "", "输入密码")
	Command = rootCmd
}
