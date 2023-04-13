package controller

//controller/itemCf.go

import (
	"my_item_cf_go/component/lock"
	"my_item_cf_go/context"
	"my_item_cf_go/response"
	"time"
)

func TestItemCf(context *context.Context) *response.Response {

	l := lock.NewLock("test", 1*time.Second)

	defer l.Release()

	if l.Block(4 * time.Second) {

		data := map[string]interface{}{
			"msg":    "拿锁成功",
			"status": 200,
		}

		return response.Resp().Json(data)
	}

	return response.Resp().String("拿锁失败")
}
