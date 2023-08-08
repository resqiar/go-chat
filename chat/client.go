package chat

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

// Client is responsible for encapsulating a single connected user.
type Client struct {
	// hold a websocket connection for current user
	socket *websocket.Conn

	// a channel which used to send a message
	send chan []byte
	// hold a reference to which room the current user is
	room *Room
}

type Msg struct {
	Chat string `json:"chat_message"`
}

func (c *Client) Read() {
	defer c.socket.Close()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			break
		}

		// pass the message to room channel
		c.room.forward <- message
	}
}

func (c *Client) Write() {
	defer c.socket.Close()

	for message := range c.send {
		var m Msg
		if err := json.Unmarshal(message, &m); err != nil {
			return
		}

		html := []byte(fmt.Sprintf(
			"<div id='chat_room' hx-swap-oob='beforeend'><p>%d:%d - %v</p></div>",
			time.Now().Hour(),
			time.Now().Minute(),
			m.Chat))
		if err := c.socket.WriteMessage(websocket.TextMessage, html); err != nil {
			return
		}
	}
}
