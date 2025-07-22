package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // return true for now, maybe ill do cors later.
	},
}

var rooms = make(map[string]*Room)
var roomsLock sync.Mutex

type Room struct {
	clients  map[*Client]bool
	roomID   string
	roomName string
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

func handleCreateRoom(w http.ResponseWriter) (bool, error) {
	roomID, err := GenerateRandomToken(32)

	if err != nil {
		return false, err
	}

	roomsLock.Lock()
	rooms[roomID] = &Room{clients: make(map[*Client]bool)}
	roomsLock.Unlock()

	w.Write([]byte(roomID))
	return true, nil
}

func GenerateRandomToken(byteLength int) (string, error) {
	b := make([]byte, byteLength)
	_, err := rand.Read(b)

	if err != nil {
		log.Printf("failed to read random bytes: %w", err)
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
