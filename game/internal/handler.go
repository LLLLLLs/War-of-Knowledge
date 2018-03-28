package internal

import (
	"fmt"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"reflect"
	"server/msg"
)

var (
	ID    = 0
	users = make(map[int]gate.Agent)
)

func init() {
	handler(&msg.Match{}, handleMatch)

}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handleMatch(args []interface{}) {
	m := args[0].(*msg.Match)
	a := args[1].(gate.Agent)
	log.Debug("Call Match from %v", a.RemoteAddr())
	room := new(Room)

	if LastRoomId == 0 || GetRoom(LastRoomId).PlayerCount == 2 {
		roomId := LastRoomId + 1
		room = NewRoom(roomId)
		AddRoom(room)
		LastRoomId = roomId
	} else {
		room = GetRoom(LastRoomId)
	}
	fmt.Println("playerId:", m.PlayerId)
	room.RoomPlayers[m.PlayerId] = a
	room.PlayerCount += 1
	fmt.Println("RoomStat:", room)
	if room.PlayerCount == 1 {
		a.WriteMsg(&msg.MatchStat{
			Status: 1,
			Msg:    "匹配中",
		})
	} else {
		for _, aa := range room.RoomPlayers {
			aa.WriteMsg(&msg.MatchStat{
				Status: 0,
				Msg:    "匹配成功！",
			})
		}
	}
}
