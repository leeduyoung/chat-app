// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chat

import (
	"context"
	"encoding/json"
	"fmt"

	"chat-app/chat/event"
	"chat-app/chat/redis"
)

type BroadcastMessage struct {
	targetID []string
	message  []byte
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// 입장
	enter chan *Client

	// 퇴장
	leave chan *Client

	// 채팅 메시지
	broadcastMessage chan BroadcastMessage

	// 연결된 유저
	userIDToClient map[string]*Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),

		enter:            make(chan *Client),
		leave:            make(chan *Client),
		broadcastMessage: make(chan BroadcastMessage),
		userIDToClient:   make(map[string]*Client),
	}
}

/*
	1. Reids subscribe -> 유저가 보낸 메시지를 받아온다.
	2. 메시지 안에 roomID를 사용해서 redis에서 user ID목록 조회
	3. userID 목록을 순회하면서 client를 찾고 메시지 send
*/
func (h *Hub) Subscribe() {
	ctx := context.Background()
	subscriber := redisClient.Subscribe(redis.MsgToSub{
		Channels: []string{chatChannel},
	})

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Println("[subscribe] msg: ", msg)
		data := event.ChatEvent{}
		json.Unmarshal([]byte(msg.Payload), &data)

		users := []string{}
		response, err := redisClient.Get(data.RoomID)
		if err != nil {
			fmt.Println("redist client get error: ", err)
		}

		json.Unmarshal([]byte(response), &users)

		// 중복 유저 제거
		m := make(map[string]bool)
		targetIDs := []string{}
		for _, userID := range users {
			if _, ok := m[userID]; !ok {
				m[userID] = true
				targetIDs = append(targetIDs, userID)
			}
		}

		h.broadcastMessage <- BroadcastMessage{
			targetID: targetIDs,
			message:  []byte(msg.Payload),
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		case client := <-h.enter:
			h.userIDToClient[client.userID] = client
		case client := <-h.leave:
			if _, ok := h.userIDToClient[client.userID]; ok {
				delete(h.userIDToClient, client.userID)
				close(client.send)
			}
		case bMessage := <-h.broadcastMessage:
			for _, userID := range bMessage.targetID {
				if val, ok := h.userIDToClient[userID]; ok {
					val.send <- bMessage.message
				}
			}
		}
	}
}
