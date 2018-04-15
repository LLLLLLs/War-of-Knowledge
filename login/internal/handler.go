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
		log.Debug("account not exist")
		a.WriteMsg(&msg.LoginStat{
			Status:   1,
			Msg:      "account not exist",
			PlayerId: PlayerId,
		})
		return
	} else if user.Login == true {
		a.WriteMsg(&msg.LoginStat{
			Status:   1,
			Msg:      "account alredy login",
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
		user.Login = true
		log.Debug("user %s login success", m.UserName)
	} else {
		//TODO LoginFail
		a.WriteMsg(&msg.LoginStat{
			Status:   1,
			Msg:      "account or password is wrong",
			PlayerId: PlayerId,
		})
	}
}
