package main

import (
	"log"
	"net/http"
	"strconv"
	"testProject/learning/example/websocket/go-pb-custom/trace"

	"github.com/gorilla/websocket"
)

type room struct {

	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clients.
	forward chan []byte

	// join is a channel for clients wishing to join the room.
	join chan *client

	// leave is a channel for clients wishing to leave the room.
	leave chan *client

	// clients holds all current clients in this room.
	clients map[string]*client
	// clients map[*client]bool

	// tracer will receive trace information of activity
	// in the room.
	tracer trace.Tracer
}

// newRoom makes a new room that is ready to
// go.
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		// clients: make(map[*client]bool),
		clients: make(map[string]*client),
		tracer:  trace.Off(),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			//
			if _, exists := r.clients[client.ID]; exists {
				log.Println("Something wrong, client already here!!!")
				close(client.send)
				continue
			}
			// joining
			r.clients[client.ID] = client
			// r.clients[client] = true

			r.tracer.Trace("New client joined", client.ID)
		case client := <-r.leave:
			// leaving
			// delete(r.clients, client)
			delete(r.clients, client.ID)
			close(client.send)
			r.tracer.Trace("Client left", client.ID)
		case msg := <-r.forward:
			r.tracer.Trace("Message received: ", string(msg))
			// forward message to all clients
			for id, client := range r.clients {
				client.send <- msg
				r.tracer.Trace(" -- sent to client", id)
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
	WriteBufferSize: socketBufferSize,
}

// when client.js call connect -> that request will be handle by this func
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// upgrade http request to ws 101
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	// gen id
	// TODO: should be username
	// TODO: support case one userA open 2 tab or open in 2 device
	// -> then that create 2 request connect to server
	// -> ???
	clientID := len(r.clients) + 1

	// new client
	client := &client{
		ID:     strconv.Itoa(clientID),
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}

	// join the room
	r.join <- client
	defer func() { r.leave <- client }()

	go client.write()

	client.read()
}
