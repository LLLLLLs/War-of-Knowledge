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
	handler(&msg.Hello{}, handleHello)
	handler(&msg.Calculate{}, handleCalculate)
	handler(&msg.Match{}, handleMatch)

}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handleMatch(args []interface{}) {
	m := args[0].(*msg.Match)
	a := args[1].(gate.Agent)
	room := new(Room)

	if lastRoomId == 0 || GetRoom(lastRoomId).RoomStat == 1 {
		roomId := lastRoomId + 1
		room = NewRoom(roomId)
		lastRoomId = roomId
	} else {
		room = GetRoom(lastRoomId)
	}
	room.RoomPlayers[m.PlayerId] = a
	if len(GetRoom(lastRoomId).RoomPlayers) == 2 {
		room.RoomStat = 1 // Close the room
	}
}

func handleHello(args []interface{}) {
	m := args[0].(*msg.Hello)
	a := args[1].(gate.Agent)
	users[ID] = a
	ID += 1

	log.Debug("hello %v", m.Name)

	a.WriteMsg(&msg.Hello{
		Name: "client",
	})
}

func handleCalculate(args []interface{}) {
	m := args[0].(*msg.Calculate)
	a := args[1].(gate.Agent)
	for i := 0; i < ID; i++ {
		fmt.Println("ID", i, users[ID] == a)
		log.Debug("ID %v", i)
		users[ID].WriteMsg(&msg.Hello{
			Name: "client",
		})
	}

	log.Debug("hello %v", m.X+m.Y)

	// a.WriteMsg(&msg.Hello{
	// 	Name: "client111",
	// })
}
