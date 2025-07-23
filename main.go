package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // return true for now, maybe ill do cors later.
	},
}

var rooms = make(map[string]*Room)
var roomsLock sync.RWMutex

type Room struct {
	clients    map[*Client]bool
	roomID     string
	roomName   string
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// so this function will be run as a goroutine,
// allowing us to spin up multiple rooms concurrently.
func (r *Room) run() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
			log.Printf("Client registered in room %s. Total clients %d", r.roomID, len(r.clients))

		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)

				log.Printf("Client unregistered successfully from room %s", r.roomID)

				if len(r.clients) == 0 {
					log.Printf("Room %s is empty. Closing room...", r.roomID)

					close(r.broadcast)

					roomsLock.Lock()
					delete(rooms, r.roomID)
					roomsLock.Unlock()
					return // end the goroutine for this room
				}
			}

		case message := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(r.clients, client)
					log.Printf("Failed to send message to client, unregistering client...")
				}
			}
		}

	}
}

func handleCreateRoom(w http.ResponseWriter) (bool, error) {
	roomID, err := GenerateRandomToken(32)

	if err != nil {
		return false, err
	}

	roomsLock.Lock()
	rooms[roomID] = &Room{clients: make(map[*Client]bool), roomID: "123", roomName: "name"}
	roomsLock.Unlock()

	w.Write([]byte(roomID))
	return true, nil
}

func GenerateRandomToken(byteLength int) (string, error) {
	b := make([]byte, byteLength)
	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
