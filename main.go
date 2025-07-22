package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // return true for now, maybe ill do cors later.
	},
}

type Room struct {
	roomID   string
	roomName string
}

func handleCreateRoom(w http.ResponseWriter, r http.Request) {
	//roomID :=

}

func GenerateRandomToken(byteLength int) (string, error) {
	b := make([]byte, byteLength)
	_, err := rand.Read(b)

	if err != nil {
		fmt.Errorf("failed to read random bytes: %w", err)
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
