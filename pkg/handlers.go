package pkg

import "net/http"

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
