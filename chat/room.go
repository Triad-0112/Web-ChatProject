package main

import (
	"log"
	"net/http"
	"webchat/tracer"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

//room struct
type room struct {
	forward chan *message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	tracer  tracer.Tracer
}

//create allocated room
func newRoom(avatar Avatar) *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  tracer.Off(),
	}
}

//Functionality inside room wheneever the room is still exist or connection not closed
func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			//Check if client join
			r.clients[client] = true
			r.tracer.Trace("New Client joined")
		case client := <-r.leave:
			//Check if client leave
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client Left")
		case msg := <-r.forward:
			r.tracer.Trace("Message received: ", msg.Message)
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

//Room activity, used to send signal whether client join or leave along with client writing function in room
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//upgrade http to web socket connection
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
		return
	}
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
