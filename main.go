package main

import (
	"my_item_cf_go/kernel"
	"my_item_cf_go/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	//加载全局变量
	kernel.Load()

	routes.Load(r)

	r.Run()
}
