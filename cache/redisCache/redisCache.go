package redisCache

import (
	"context"
	"my_item_cf_go/redis"
	"time"
)

type redisCache struct {
}

func NewRedisCache() *redisCache {

	return &redisCache{}
}

func (r redisCache) Put(key string, value string, ttl time.Duration) error {
	//TODO implement me

	_, err := redis.Client().Set(context.TODO(), key, value, ttl).Result()

	return err
}

func (r redisCache) Get(key string) (string, error) {
	//TODO implement me

	return redis.Client().Get(context.TODO(), key).Result()

}
