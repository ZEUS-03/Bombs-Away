package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"

	"ws/internal/utils"
	"ws/internal/domain"
	"ws/internal/game"
)

var rooms = make(map[string]*domain.Room)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial HTTP connection to a WebSocket
	ws, err := utils.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}
	// fmt.Println("Response generated %s/n", r)
	defer ws.Close()

	fmt.Println("Client connected")

	roomId := r.URL.Query().Get("room_id");
	if roomId == "" {
		roomId = "default"
	}

	// Creating new room
	room := game.CreateRoom(rooms, roomId)

	// Generating random player ID
	playerID := uuid.New().String()

	// Creating new player in the room
	game.AddPlayerToRoom(room, roomId, playerID, ws)

	locations := game.GetPlayerLocations(rooms, roomId);

	var newLocations = make(map[string]interface{})

	newLocations["type"] = "all_players_position"
	newLocations["data"] = locations	

	game.BroadcastToPlayers(room.Players, newLocations)

	fmt.Printf("Player %s connected to room %s\n", playerID, roomId)

	for _, player := range rooms[roomId].Players {
		player.Ws.WriteJSON(map[string]interface{}{
				"type": "player_connected",
				"player_id": playerID,
		})
	}	

	// Handle communication
	for {
		// Read message from the client
		
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			delete(rooms[roomId].Players, playerID)
			
			break
		}
		var wsMsg map[string]interface{}
		err = json.Unmarshal(msg, &wsMsg)
		if err != nil {
			fmt.Println("Error unmarshalling message:", err)
			continue
		}

		fmt.Printf("Received: %v\n", wsMsg["x"])
		if wsMsg["type"] == "update_position" {
			x := wsMsg["x"].(float64)
			y := wsMsg["y"].(float64)
			rooms[roomId].Players[playerID].Position = domain.Position{X: x, Y: y}

			for _, p := range rooms[roomId].Players {
        err := p.Ws.WriteJSON(map[string]interface{}{
            "type":    "update_single_player_position",
            "player_id": playerID,
            "x":       x,
            "y":       y,
        })
        if err != nil {
            fmt.Println("Error broadcasting position:", err)
        }
    	}
		}
		fmt.Printf("rooms %v", rooms[roomId].Players); 
		for _, player := range rooms[roomId].Players {
			ws.WriteJSON(player)
		}
	}
	
}