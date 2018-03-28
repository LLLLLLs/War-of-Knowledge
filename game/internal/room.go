package internal

import (
	"fmt"
	"github.com/name5566/leaf/gate"
)

type Room struct {
	RoomId      int
	RoomPlayers map[int]gate.Agent
	PlayerCount int
}

func NewRoom(roomId int) *Room {
	fmt.Println("newRoom:", roomId)
	room := Room{RoomId: roomId, RoomPlayers: make(map[int]gate.Agent), PlayerCount: 0}

	return &room
}
