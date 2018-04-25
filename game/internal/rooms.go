package internal

import (
	"server/msg"
	"github.com/name5566/leaf/gate"
)

var (
	Rooms       = make(map[int]*Room)
	LastRoomId  = int(0)
	LastMatchId = int(0)
	RoomList    = []*msg.RoomInfo{}
)

func UpdateRoomInfo(room *Room) {
	flag := 0
	for _, roomInfo := range RoomList {
		if roomInfo.RoomId == room.RoomId {
			roomInfo.Name = room.Name
			roomInfo.Users = room.Users
			break
		}
		flag += 1
	}
	if flag == len(RoomList) {
		RoomList = append(RoomList, &msg.RoomInfo{
			RoomId: room.RoomId,
			Name:   room.Name,
			Users:  room.Users,
		})
	}
	for _, aa := range room.User2Agent {
		if aa == nil {
			continue
		}
		(*aa).WriteMsg(&msg.RoomInfo{
			Msg:    "ok",
			RoomId: room.RoomId,
			Name:   room.Name,
			Users:  room.Users,
		})
	}
}

func deleteRoomInfo(roomId int) {
	for i, roomInfo := range RoomList {
		if roomInfo.RoomId == roomId {
			RoomList = append(RoomList[:i], RoomList[i+1:]...)
			return
		}
	}
}

func DeleteRoom(roomId int, a gate.Agent) {
	room := Rooms[roomId]
	if room.Mode == Spec {
		deleteRoomInfo(roomId)
	}
	if room.InBattle == true && room.Closed == true {
		for _, aa := range room.User2Agent {
			delete(Agent2Room, *aa)
		}
		delete(Rooms, roomId)
		return
	}
	delete(Agent2Room, a)
	delete(Rooms, roomId)
}

func GetRoom(roomId int) (*Room, bool) {
	room, ok := Rooms[roomId]
	return room, ok
}
