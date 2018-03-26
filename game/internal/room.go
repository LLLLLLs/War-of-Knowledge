package internal

import (
	"fmt"
	"github.com/name5566/leaf/gate"
)

type Room struct {
	RoomId      int
	RoomPlayers map[int]gate.Agent
	RoomStat    int // 0 : open , 1 : close
}

func NewRoom(roomId int) *Room {
	fmt.Println("newRoom", roomId)
	room := Room{RoomId: roomId, RoomStat: 0}

	return &room
}
