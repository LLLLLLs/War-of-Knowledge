package internal

import (
	"github.com/name5566/leaf/gate"
	"server/gamedata"
)

var (
	Users      = make(map[gate.Agent]string)
	Room2Agent = make(map[int][]gate.Agent)
	Agent2Room = make(map[gate.Agent]int)
)

func init() {
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
	skeleton.RegisterChanRPC("Login", rpcLogin)
	skeleton.RegisterChanRPC("RecoverBattle", rpcRecoverBattle)
}

func rpcNewAgent(args []interface{}) {}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	if _, ok := Users[a]; ok {
		gamedata.UsersMap[Users[a]].Login = false
	}
	if roomId, ok := Agent2Room[a]; ok {
		if _, ok := GetRoom(roomId); ok {
			//在room里增加 断线计时器 和 重连计时器
			//TODO
			EndBattle(roomId, a)
			delete(Agent2Room, a)
		}
		delete(Agent2Room, a)
	}
	delete(Users, a)
}

func rpcLogin(args []interface{}) {
	a := args[0].(gate.Agent)
	name := args[1].(string)
	Users[a] = name
	gamedata.UsersMap[name].Login = true
}

func rpcRecoverBattle(args []interface{}) {
	// TODO
	//a := args[0].(gate.Agent)
	//roomId := gamedata.UsersMap[Users[a]].RoomId
}
