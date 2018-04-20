package msg

import (
	_ "encoding/json"
	"github.com/name5566/leaf/network/json"
)

var Processor = json.NewProcessor()

func init() {
	Processor.Register(&Login{})
	Processor.Register(&LoginStat{})
	Processor.Register(&Match{})
	Processor.Register(&MatchStat{})

	Processor.Register(&MoneyLeft{})
	Processor.Register(&UpdateBaseState{})

	Processor.Register(&UseSkill{})
	Processor.Register(&UseSkillInf{})
	Processor.Register(&SkillCrash{})

	Processor.Register(&CreateHero{})
	Processor.Register(&CreateHeroInf{})
	Processor.Register(&UpdatePosition{})
	Processor.Register(&UpdateHeroState{})
	Processor.Register(&DeleteHero{})
	Processor.Register(&GetResource{})

	Processor.Register(&CreateMiddle{})
	Processor.Register(&UpdateMiddleState{})
	Processor.Register(&DeleteMiddle{})

	Processor.Register(&Surrender{})
	Processor.Register(&EndBattle{})
}
