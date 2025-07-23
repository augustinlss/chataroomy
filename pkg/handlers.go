package pkg

import (
	"encoding/json"
	"io"
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

	err = json.NewDecoder(r.Body).Decode(roomCreationRequest)

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

	log.Printf("Successfuly created room %d", roomID)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(roomID))
}
