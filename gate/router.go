package gate

import (
	"server/game"
	"server/login"
	"server/msg"
)

func init() {
	msg.Processor.SetRouter(&msg.Login{}, login.ChanRPC)
	msg.Processor.SetRouter(&msg.Match{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.CreateHero{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.UpdatePosition{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.UseSkill{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.GetResource{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.SkillCrash{}, game.ChanRPC)
}
