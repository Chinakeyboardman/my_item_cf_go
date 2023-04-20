package main

import (
	"fmt"
	"my_item_cf_go/kernel"
	"my_item_cf_go/plugin/myorm"
	"my_item_cf_go/routes"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {

	// 加载viper配置中心
	initConfig()

	// 加载gorm
	_, err := myorm.Connect()
	if err != nil {
		fmt.Print("Failed to connect to the database: " + err.Error())
		// panic("Failed to connect to the database: " + err.Error())
	}

	r := gin.Default()

	//加载全局变量
	kernel.Load()

	routes.Load(r)

	r.Run()
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Failed to read configuration file: " + err.Error())
	}
}
