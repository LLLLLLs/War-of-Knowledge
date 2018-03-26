package msg

import (
	"github.com/name5566/leaf/network/json"
)

var Processor = json.NewProcessor()

func init() {
	Processor.Register(&Hello{})
	Processor.Register(&Calculate{})
	Processor.Register(&Login{})
	Processor.Register(&Match{})
}

type Hello struct {
	Name string
}

type Calculate struct {
	X int
	Y int
}
