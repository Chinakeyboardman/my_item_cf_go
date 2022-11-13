package controller

//controller/Index.go

import (
	"my_item_cf_go/component/lock"
	"my_item_cf_go/context"
	"my_item_cf_go/response"
	"strconv"
	"time"
)

func Index(context *context.Context) *response.Response {

	l := lock.NewLock("test", 1*time.Second)

	defer l.Release()

	if l.Get() {

		time.Sleep(4 * time.Second)

		return response.Resp().String("拿锁成功")
	}

	return response.Resp().String("拿锁失败")
}

func Block(context *context.Context) *response.Response {

	l := lock.NewLock("test", 10*time.Second)

	defer l.Release()

	if l.Block(5 * time.Second) {

		return response.Resp().String("拿锁成功")

	}

	return response.Resp().String("拿锁失败")

}

func Index2(context *context.Context) *response.Response {

	//msg, _ := context.Session().Get("msg")

	//fmt.Println(limiter.GlobalLimiters)

	return response.Resp().String("nice")
}

func Index3(context *context.Context) *response.Response {

	context.Session().Remove("msg")

	return response.Resp().String("")
}

func Index4(context *context.Context) *response.Response {

	session := context.Session()

	for i := 0; i < 100; i++ {

		go func(index int) {

			session.Set("msg"+strconv.Itoa(index), index)

		}(i)
	}

	return response.Resp().String("")
}
