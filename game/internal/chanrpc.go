package internal

import (
	"github.com/name5566/leaf/gate"
	"server/gamedata"
)

var agents = make(map[gate.Agent]struct{})
var Users = make(map[gate.Agent]string)

func init() {
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
	skeleton.RegisterChanRPC("Login", rpcLogin)
}

func rpcNewAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	agents[a] = struct{}{}
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	delete(agents, a)
	gamedata.UsersMap[Users[a]].Login = false
	delete(Users, a)
}

func rpcLogin(args []interface{}) {
	a := args[0].(gate.Agent)
	name := args[1].(string)
	Users[a] = name
	gamedata.UsersMap[name].Login = true
}
