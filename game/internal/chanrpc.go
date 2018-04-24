package internal

import (
	"server/gamedata"
	"time"
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/gate"
)

var (
	Users      = make(map[gate.Agent]string)
	Agent2Room = make(map[gate.Agent]int)
)

func init() {
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
	skeleton.RegisterChanRPC("Login", rpcLogin)
	skeleton.RegisterChanRPC("RecoverBattle", rpcRecoverBattle)
}

func rpcNewAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	log.Debug("%v", a.RemoteAddr())
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	if _, ok := Users[a]; ok {
		gamedata.UsersMap[Users[a]].Login = 0
	} else {
		log.Debug("%v 断开连接", a.RemoteAddr())
		return
	}
	userData := gamedata.UsersMap[Users[a]]
	gamedata.Db.Id(userData.Id).Cols("login").Update(gamedata.UsersMap[Users[a]])
	if roomId, ok := Agent2Room[a]; ok {
		if room, ok := GetRoom(roomId); ok {
			if room.InBattle == true {
				room.User2Agent[Users[a]] = nil
				userName := Users[a]
				delete(Users, a)
				timer := time.NewTimer(time.Second * 30)
				<-timer.C
				if gamedata.UsersMap[userName].Login == 1 {
					return
				}
				EndBattle(roomId, a)
				DeleteRoom(roomId, a)
				return
			}
			if room.Mode == Match {
				QuitMatch(a)
			} else {
				ExitRoom(a, room)
			}
		}
	}
}

func rpcLogin(args []interface{}) {
	a := args[0].(gate.Agent)
	name := args[1].(string)
	Users[a] = name
}

func rpcRecoverBattle(args []interface{}) {
	a := args[0].(gate.Agent)
	room, ok := GetRoom(gamedata.UsersMap[Users[a]].RoomId)
	if !ok {
		log.Debug("重连失败,房间已关闭")
		return
	}
	Agent2Room[a] = room.RoomId
	RecoverBattle(a, room)
}
