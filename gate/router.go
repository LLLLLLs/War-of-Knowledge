package gate

import (
	"server/game"
	"server/login"
	"server/msg"
)

func init() {
	msg.Processor.SetRouter(&msg.Register{}, login.ChanRPC)
	msg.Processor.SetRouter(&msg.Login{}, login.ChanRPC)
	msg.Processor.SetRouter(&msg.Match{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.QuitMatch{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.GetUserInfo{}, game.ChanRPC)

	msg.Processor.SetRouter(&msg.CreateHero{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.UpdatePosition{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.MoveTo{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.UseSkill{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.GetResource{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.SkillCrash{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.Surrender{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.FireBottleCrash{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.Upgrade{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.HeartBeat{}, game.ChanRPC)

	msg.Processor.SetRouter(&msg.CreateRoom{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.GetRoomList{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.EnterRoom{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.ExitRoom{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.StartBattle{}, game.ChanRPC)

	msg.Processor.SetRouter(&msg.Test{}, game.ChanRPC)
}
