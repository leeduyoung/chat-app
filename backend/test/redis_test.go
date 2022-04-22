package test

import (
	"testing"
)

type food struct {
	Name     string
	Calories float64
}

func TestRedisPubSub(t *testing.T) {
	// t.Run("RedisPubSub", func(t *testing.T) {
	// 	// subscribe
	// 	_, err = pubsub.NewSubscriber("food", eat)
	// 	if err != nil {
	// 		log.Println("NewSubscriber() error", err)
	// 	}

	// 	// publish
	// 	pub = pubsub.Service.Publish("food", food{"Pizza", 50.1})
	// 	if err = pub.Err(); err != nil {
	// 		log.Print("PublishString() error", err)
	// 	}
	// })
}
