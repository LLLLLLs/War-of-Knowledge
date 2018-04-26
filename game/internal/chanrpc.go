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
	log.Debug("%v 连接", a.RemoteAddr())
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	delete(Agent2Room, a)
	if _, ok := Users[a]; !ok {
		log.Debug("未登录...连接断开")
		return
	}
	userData := gamedata.UsersMap[Users[a]]
	userData.Login = 0
	cond := gamedata.UserData{
		Id: userData.Id,
	}
	gamedata.Db.Id(userData.Id).Cols("login").Update(userData, cond)
	roomId := userData.RoomId
	if room, ok := GetRoom(roomId); ok {
		go func() {
			room.PlayerCount -= 1
			if room.InBattle == true {
				if room.PlayerCount == 0 {
					log.Debug("双方退出,游戏结束")
					DeleteRoom(roomId, a, false)
					return
				}
				room.User2Agent[Users[a]] = nil
				delete(Users, a)
				timer := time.NewTimer(time.Second * 30)
				<-timer.C
				userData.Refresh()
				if userData.Login == 1 {
					log.Debug("%s 重连成功...", userData.Name)
					return
				}
				if _, ok := GetRoom(roomId); ok {
					EndBattle(roomId, a)
					DeleteRoom(roomId, a, false)
				}
				return
			}
			if room.Mode == Match {
				QuitMatch(a)
			} else {
				ExitRoom(a, room)
			}
		}()
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
