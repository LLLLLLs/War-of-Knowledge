package gate

import (
	"server/game"
	"server/login"
	"server/msg"
)

func init() {
	msg.Processor.SetRouter(&msg.Login{}, login.ChanRPC)
	msg.Processor.SetRouter(&msg.Match{}, game.ChanRPC)
}
