package internal

import (
	"reflect"

	"server/game"
	"server/msg"
	"server/gamedata"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
)

var PlayerId = 1

func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handleMsg(&msg.Login{}, handleAuth)
}

func handleAuth(args []interface{}) {
	m := args[0].(*msg.Login)
	a := args[1].(gate.Agent)
	log.Debug("call login from %v", a.RemoteAddr())
	user, ok := gamedata.UsersMap[m.UserName]
	if !ok {
		log.Debug("账号不存在")
		a.WriteMsg(&msg.LoginStat{
			Status:   1,
			Msg:      "账号不存在",
			PlayerId: PlayerId,
		})
		return
	} else if user.Login == true {
		a.WriteMsg(&msg.LoginStat{
			Status:   1,
			Msg:      "账号已登陆",
			PlayerId: PlayerId,
		})
		return
	}
	if m.UserPwd == user.UserPwd {
		game.ChanRPC.Go("Login", a, m.UserName)
		a.WriteMsg(&msg.LoginStat{
			Status:   0,
			Msg:      "Success",
			PlayerId: PlayerId,
		})
		PlayerId += 1
		log.Debug("玩家 %s 登陆成功", m.UserName)
		// 如果玩家还在战斗中则恢复战斗状态(断线重连)
		if gamedata.UsersMap[m.UserName].InBattle == true {
			game.ChanRPC.Go("RecoverBattle", a)
		}
	} else {
		a.WriteMsg(&msg.LoginStat{
			Status:   1,
			Msg:      "账号密码不匹配",
			PlayerId: PlayerId,
		})
	}
}
