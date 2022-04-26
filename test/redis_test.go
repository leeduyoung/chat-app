package test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"chat-app/chat/event"
	"chat-app/chat/redis"
)

var (
	redisClient redis.IRedisClientService
	ctx         context.Context
)

func init() {
	ctx = context.Background()
	redisClient = redis.New(ctx)
}

func TestRedisSetGet(t *testing.T) {
	const (
		roomID = "1"
		userID = "kaye"
	)
	users := []string{"kaye", "jade"}

	// Set
	err := redisClient.Set(roomID, users)
	if err != nil {
		t.Error(err)
		return
	}

	// Get
	val, err := redisClient.Get(roomID)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("[Success] TestRedis - val: ", val)
}

func TestRedisGet(t *testing.T) {
	const notExistRoomID = "273ghfjis8"

	t.Run("[TestRedisGet] Not found get key", func(t *testing.T) {
		_, err := redisClient.Get(notExistRoomID)
		if err != nil {
			t.Log(err)
			return
		}

		t.Error(err)
	})
}

func TestRedisPubSub(t *testing.T) {
	redisPubSubChannel := "test-channel"

	go t.Run("subscribe", func(t *testing.T) {
		msgToSub := redis.MsgToSub{
			Channels: []string{redisPubSubChannel},
		}

		subscriber := redisClient.Subscribe(msgToSub)

		for {
			msg, err := subscriber.ReceiveMessage(ctx)
			if err != nil {
				panic(err)
			}

			fmt.Println("msg: ", msg)
		}
	})

	time.Sleep(time.Second)

	t.Run("publish", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			msg := event.ChatEvent{
				UserID:      fmt.Sprintf("user[%d]", i+1),
				RoomID:      fmt.Sprintf("room[%d]", i+1),
				Nickname:    fmt.Sprintf("nickname[%d]", i+1),
				Message:     "test",
				MessageType: event.MessageTypeEnter,
			}
			bytes, _ := json.Marshal(msg)

			msgToPub := redis.MsgToPub{
				Channel: redisPubSubChannel,
				Message: string(bytes),
			}
			redisClient.Publish(msgToPub)
		}
	})
}
