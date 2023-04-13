package routes

//routes/routes.go

import (
	"my_item_cf_go/controller"
	"my_item_cf_go/kernel"

	"github.com/gin-gonic/gin"
	"my_item_cf_go/middleware"
)

func config(router group) {

	//router.Group("/api", func(api group) {
	//
	//	api.Group("/user", func(user group) {
	//
	//		user.Registered(GET, "/info", controller.Index)
	//		user.Registered(GET, "/order", controller.Index)
	//		user.Registered(GET, "/money", controller.Index)
	//
	//	})
	//
	//})

	router.Registered(GET, "/", controller.Index)
	router.Registered(GET, "/index2", controller.Index2)
	router.Registered(GET, "/index3", controller.Index3)
	router.Registered(GET, "/index4", controller.Index4)

	// 控制器
	router.Group("/api", func(api group) {

		api.Group("/item_cf", func(item_cf group) {

			item_cf.Registered(GET, "/testItemCf", controller.TestItemCf, middleware.JWTAuth)

		}, middleware.M2)

	}, middleware.M1)

}

func Load(r *gin.Engine) {

	router := newRouter(r)

	router.Group("", func(g group) {

		config(g)

	}, kernel.Middleware...) //加载全局中间件

}
