package pkg

import (
	"log"
	"sync"
)

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

func NewRoom(roomID string, roomName string) *Room {
	room := &Room{
		clients:    make(map[*Client]bool),
		roomID:     roomID,
		roomName:   roomName,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	go room.run()
	return room
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
