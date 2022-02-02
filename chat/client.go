package main

import (
	"time"

	"github.com/gorilla/websocket"
)

//client struct
type client struct {
	socket   *websocket.Conn
	send     chan *message
	room     *room
	userData map[string]interface{}
}

//Client channel to read, wait for msg coming.
func (c *client) read() {
	defer c.socket.Close()
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		if avatarUrl, ok := c.userData["avatar_url"]; ok {
			msg.AvatarURL = avatarUrl.(string)
		}
		c.room.forward <- msg
	}
}

//Function to write msg using channel send and sent over to the target
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}
