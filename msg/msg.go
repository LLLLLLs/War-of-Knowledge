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
}

// type Resp struct {
// 	StateCode int    `json:stateCode`
// 	Msg       string `json:msg`
// }
