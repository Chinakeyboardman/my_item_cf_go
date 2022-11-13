package kernel

//middleware/session/session.go

import (
	"my_item_cf_go/context"
	"my_item_cf_go/middleware/exception"
	"my_item_cf_go/middleware/session"
)

// Middleware 全局中间件
var Middleware []context.HandlerFunc

func Load() {

	Middleware = []context.HandlerFunc{
		exception.Exception,
		session.Session,
	}

}
