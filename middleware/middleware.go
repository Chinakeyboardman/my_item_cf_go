package middleware

//middleware/middleware.go

import (
	"fmt"
	"my_item_cf_go/context"
)

func M1(c *context.Context) {

	fmt.Println("我是1")

}

func M2(c *context.Context) {

	fmt.Println("我是2")

}

func M3(c *context.Context) {

	fmt.Println("我是3")

}
