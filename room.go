package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Room struct {
	// a channel that holds incoming messages that
	// should be forwarded to other clients with the same room.
	forward chan []byte

	// channel for joining the room
	join chan *Client

	// channel for leaving the room
	leave chan *Client

	// holds all current connected client in the room
	clients map[*Client]bool
}

func NewRoom() *Room {
	return &Room{
		forward: make(chan []byte),
		join:    make(chan *Client),
		leave:   make(chan *Client),
		clients: make(map[*Client]bool),
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true

		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
		case message := <-r.forward:
			// forward the message to all clients in the room
			for client := range r.clients {
				client.send <- message
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: messageBufferSize,
}

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("ERROR UPGRADING", err)
		return
	}

	client := &Client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize), // buffered channel
		room:   r,
	}

	r.join <- client

	// exit and cleanup  client when function exits.
	// it does not matter is it normal or error exit.
	defer func() { r.leave <- client }()

	go client.Write()
	client.Read()
}
