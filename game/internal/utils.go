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
		if i == 1 {
			continue
		}
		d += math.Pow((sposition[i] - eposition[i]), 2)
	}
	distance := math.Sqrt(d)
	return distance
}

func GetDamage(hatk, def float64, satk float64) float64 {
	return (hatk + satk) * (1 / (1 + def*0.05))
}

func GetEnemy(a gate.Agent, room Room) *Player {
	for user, aa := range room.User2Agent {
		if aa == nil {
			continue
		}
		if (*aa) != a {
			return room.Players[user]
		}
	}
	return nil
}

func getAngle(tf1, tf2 msg.TFServer) float64 {
	p1, p2 := tf1.Position, tf2.Position
	a := p2[2] - p1[2]
	b := p2[0] - p1[0]
	c := GetDistance(&tf1, &tf2)
	sinA := a / c
	A := math.Asin(sinA)
	if sinA > 0 {
		if b > 0 {
			return A
		} else {
			return math.Pi - A
		}
	} else {
		if b > 0 {
			return A
		} else {
			return -math.Pi - A
		}
	}

}
