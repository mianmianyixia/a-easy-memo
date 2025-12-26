package main

import (
	"a-easy-memo/internal/config"
	"a-easy-memo/internal/routers"
	"fmt"

	"log"
)

func main() {
	fmt.Println("请先登录数据库")
	err := config.Command.Execute()
	if err != nil {
		log.Fatal(err)
	}
	routers.Routers()
}
