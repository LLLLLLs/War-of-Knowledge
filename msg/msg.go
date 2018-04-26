package msg

import (
	_ "encoding/json"
	"github.com/name5566/leaf/network/json"
)

var Processor = json.NewProcessor()

func init() {
	Processor.Register(&Register{})
	Processor.Register(&RegisterInfo{})
	Processor.Register(&Login{})
	Processor.Register(&LoginStat{})
	Processor.Register(&User{})
	Processor.Register(&Match{})
	Processor.Register(&MatchStat{})
	Processor.Register(&QuitMatch{})
	Processor.Register(&ChangeImage{})
	Processor.Register(&ChangeImageInf{})

	Processor.Register(&GetUserInfo{})

	Processor.Register(&CreateRoom{})
	Processor.Register(&GetRoomList{})
	Processor.Register(&RoomList{})
	Processor.Register(&RoomInfo{})
	Processor.Register(&EnterRoom{})
	Processor.Register(&ExitRoom{})
	Processor.Register(&StartBattle{})

	Processor.Register(&MoneyLeft{})
	Processor.Register(&UpdateBaseState{})

	Processor.Register(&UseSkill{})
	Processor.Register(&UseSkillInf{})
	Processor.Register(&SkillCrash{})
	Processor.Register(&FireBottleCrash{})
	Processor.Register(&Damage{})

	Processor.Register(&CreateHero{})
	Processor.Register(&CreateHeroInf{})
	Processor.Register(&UpdatePosition{})
	Processor.Register(&MoveTo{})

	Processor.Register(&UpdateHeroState{})
	Processor.Register(&DeleteHero{})
	Processor.Register(&GetResource{})
	Processor.Register(&Upgrade{})

	Processor.Register(&CreateMiddle{})
	Processor.Register(&UpdateMiddleState{})
	Processor.Register(&DeleteMiddle{})

	Processor.Register(&Surrender{})
	Processor.Register(&EndBattle{})

	Processor.Register(&HeartBeat{})
	Processor.Register(&Test{})

	Processor.Register(&SyncItems{})
}
