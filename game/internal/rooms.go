package internal

import (
	"server/msg"
	"github.com/name5566/leaf/gate"
	"server/gamedata"
	"time"
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

func DeleteRoom(roomId int, a gate.Agent, surrender bool) {
	room := Rooms[roomId]
	if room.Mode == Spec {
		deleteRoomInfo(roomId)
	}
	// 正常退出 删除双方信息
	if surrender {
		for _, aa := range room.User2Agent {
			delete(Agent2Room, *aa)
		}
		delete(Rooms, roomId)
		return
	}
	// 非正常关闭(双方都掉线)
	room.Closed = true
	for user := range room.Users {
		room.Players[user].Base.Timer.Reset(time.Millisecond)
		userData := gamedata.UsersMap[user]
		cond := gamedata.UserData{
			Id: userData.Id,
		}
		userData.InBattle = 0
		gamedata.Db.Cols("in_battle").Update(userData, cond)
	}
	delete(Agent2Room, a)
	delete(Rooms, roomId)
}

func GetRoom(roomId int) (*Room, bool) {
	room, ok := Rooms[roomId]
	return room, ok
}
