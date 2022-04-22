package chat

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type MyRedisClient struct {
	redisClient *redis.Client
}

type YourRedisClient interface {
	Publish(ctx context.Context, redisClient *redis.Client)
	Subscribe(ctx context.Context, redisClient *redis.Client)
}

type MsgToPub struct {
	Channel string
	Message string
}

type MsgToSub struct {
	Channels []string
}

func NewMyRedisClient() *MyRedisClient {
	return &MyRedisClient{
		redisClient: redis.NewClient(
			&redis.Options{
				Addr:     "localhost:6379",
				Password: "qwer1234",
			},
		),
	}
}

func (mrc *MyRedisClient) Publish(ctx context.Context, msgToPub MsgToPub) {
	if err := mrc.redisClient.Publish(ctx, msgToPub.Channel, msgToPub.Message).Err(); err != nil {
		panic(err)
	}
}

func (mrc *MyRedisClient) Subscribe(ctx context.Context, msgToSub MsgToSub) {
	subscriber := mrc.redisClient.Subscribe(ctx, msgToSub.Channels...)

	for {
		msg, err := subscriber.ReceiveMessage(ctx)

		if err != nil {
			panic(err)
		}

		fmt.Println("msg: ", msg)
	}
}
