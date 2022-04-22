// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"chat-app/chat"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./frontend/index.html")
}

// https://github.com/gorilla/websocket/tree/master/examples/chat 옮기기

/*
	기반 기술 검증
	1) client ↔ server websocket 통신 테스트 - OK!
	2) server와 redis 사이의 통신 테스트 -
	3) 다중 server와 redis 사이의 통신 테스트 -
*/
func main() {
	// var ctx = context.Background()

	// var redisClient = redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "qwer1234",
	// })

	// out := redisClient.Set(ctx, "kaye4", "3", 0)
	// fmt.Println("out", out)
	// if out.Err() != nil {
	// 	panic("??")
	// }

	// res := redisClient.Get(ctx, "kaye4")
	// fmt.Println("res", res)
	// if res.Err() != nil {
	// 	panic("???")
	// }

	// go subscribe(ctx, redisClient)

	// time.Sleep(time.Second * 2)

	// for i := 0; i < 5; i++ {
	// 	publish(ctx, redisClient)
	// }

	flag.Parse()
	hub := chat.NewHub()
	go hub.Run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func publish(ctx context.Context, redisClient *redis.Client) {
	if err := redisClient.Publish(ctx, "test", "messgae-test").Err(); err != nil {
		panic(err)
	}
}

func subscribe(ctx context.Context, redisClient *redis.Client) {
	subscriber := redisClient.Subscribe(ctx, "test")

	for {
		msg, err := subscriber.ReceiveMessage(ctx)

		if err != nil {
			panic(err)
		}

		fmt.Println("msg: ", msg)
	}
}
