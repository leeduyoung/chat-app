package redis

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

const (
	host           = "localhost"
	port           = "6379"
	password       = "qwer1234"
	expirationTime = 0
)

type redisClientService struct {
	ctx         context.Context
	redisClient *redis.Client
}

type IRedisClientService interface {
	Set(key string, value interface{}) error
	Get(key string) (string, error)
	Publish(msgToPub MsgToPub)
	Subscribe(msgToSub MsgToSub) *redis.PubSub
}

type MsgToPub struct {
	Channel string
	Message string
}

type MsgToSub struct {
	Channels []string
}

func New(ctx context.Context) IRedisClientService {
	return &redisClientService{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     host + ":" + port,
			Password: password,
		}),
		ctx: ctx,
	}
}

func (mrc *redisClientService) Set(key string, value interface{}) error {
	bytes, _ := json.Marshal(value)
	return mrc.redisClient.Set(mrc.ctx, key, string(bytes), expirationTime).Err()
}

func (mrc *redisClientService) Get(key string) (string, error) {
	return mrc.redisClient.Get(mrc.ctx, key).Result()
}

func (mrc *redisClientService) Publish(msgToPub MsgToPub) {
	if err := mrc.redisClient.Publish(mrc.ctx, msgToPub.Channel, msgToPub.Message).Err(); err != nil {
		panic(err)
	}
}

func (mrc *redisClientService) Subscribe(msgToSub MsgToSub) *redis.PubSub {
	return mrc.redisClient.Subscribe(mrc.ctx, msgToSub.Channels...)
}
