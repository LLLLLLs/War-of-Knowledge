package internal

import "time"

type Burn struct {
	Timer    time.Timer
	Ticker   time.Ticker
	Dot      float64
	IsEffect bool
}
