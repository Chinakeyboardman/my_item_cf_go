package redis

import (
	"context"
	"sync"
	"time"

	"github.com/spf13/viper"

	"github.com/go-redis/redis/v8"
)

type connect struct {
	client *redis.Client
}

var once = sync.Once{}

var _connect *connect

func connectRedis() {

	cxt, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	defer cancel()

	// conf := &redis.Options{
	// 	Addr: "127.0.0.1:6379",
	// 	DB:   0,
	// }

	conf := &redis.Options{
		Addr: viper.GetString("redis.addr"),
		DB:   viper.GetInt("redis.db"),
	}

	c := redis.NewClient(conf)

	re := c.Ping(cxt)

	if re.Err() != nil {

		panic(re.Err())

	}

	_connect = &connect{
		client: c,
	}

}

func Client() *redis.Client {

	if _connect == nil {

		once.Do(func() {

			connectRedis()
		})

	}

	return _connect.client

}
