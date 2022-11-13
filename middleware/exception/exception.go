package exception

import (
	"fmt"
	"my_item_cf_go/context"
	"runtime/debug"
)

func Exception(c *context.Context) {

	defer func() {
		if r := recover(); r != nil {

			msg := fmt.Sprint(r) + "\n" + string(debug.Stack())

			c.String(500, msg)

			c.Abort()
		}

	}()
	c.Next()
}
