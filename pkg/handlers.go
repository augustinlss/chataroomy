package pkg

import (
	"encoding/json"
	"log"
	"net/http"
)

type RoomCreationRequest struct {
	roomName string
}

func handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	var roomCreationRequest RoomCreationRequest

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	roomID, err := GenerateRandomToken(32)

	if err != nil {
		http.Error(w, "Internal server error. Unable to generate room token...", http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&roomCreationRequest)

	if err != nil {
		http.Error(w, "Internal server error. Unable to read reqeust body...", http.StatusInternalServerError)
	}

	defer r.Body.Close()

	roomName := roomCreationRequest.roomName

	roomsLock.Lock()
	if _, exists := rooms[roomID]; exists {
		log.Printf("Generated duplicate room id. Retrying...")
		roomsLock.Unlock()
		handleCreateRoom(w, r)
		return
	}

	newRoom := NewRoom(roomID, roomName)
	rooms[roomID] = newRoom
	roomsLock.Unlock()

	log.Printf("Successfuly created room %s", roomID)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(roomID))
}

func handleJoinRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	queryParams := r.URL.Query()

	roomID := queryParams.Get("roomID")
	username := queryParams.Get("username")

	roomsLock.RLock()
	room, exists := rooms[roomID]
	roomsLock.RUnlock()

	if !exists {
		log.Printf("Failed to find room with id %s.", roomID)
		http.Error(w, "Room not found...", http.StatusNotFound)
		roomsLock.RUnlock()
		return
	}

	if username == "" {
		log.Printf("No username.")
		http.Error(w, "Invalid request. Username not provided.", http.StatusMethodNotAllowed)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("Problem upgrading http connection for room %s: %v", roomID, err)
		http.Error(w, "Unable to upgrade connection...", http.StatusInternalServerError)
		return
	}

	userID, err := GenerateRandomToken(32)

	if err != nil {
		log.Printf("Error generating user id.")
		http.Error(w, "Internal server error. Error generating user ID.", http.StatusInternalServerError)
		conn.Close()
		return
	}

	client := &Client{
		username: username,
		roomID:   roomID,
		room:     room,
		send:     make(chan []byte, 256),
		userID:   userID,
	}

}
