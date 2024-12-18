package domain

import "github.com/gorilla/websocket"

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Player struct {
	ID string `json:"id"`
	Position Position `json:"position"`
	Ws *websocket.Conn
}

