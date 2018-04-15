package internal

import (
	"server/msg"
	"math"
	"github.com/name5566/leaf/gate"
)

func GetDistance(tf1, tf2 *msg.TFServer) float64 {
	sposition := tf1.Position
	eposition := tf2.Position
	d := float64(0.0)
	for i := 0; i < 3; i++ {
		d += math.Pow((sposition[i] - eposition[i]), 2)
	}
	distance := math.Sqrt(d)
	return distance
}

func GetDamage(hatk, def float64, satk float64) float64 {
	return hatk + satk - def
}

func GetEnemyAgent(a gate.Agent, room Room) gate.Agent {
	for aa := range room.Players {
		if aa != a {
			return aa
		}
	}
	return nil
}

func GetEnemy(a gate.Agent, room Room) *Player {
	for aa, p := range room.Players {
		if aa != a {
			return p
		}
	}
	return nil
}
