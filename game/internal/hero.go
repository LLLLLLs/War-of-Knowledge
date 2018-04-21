package internal

import (
	"server/msg"
	"errors"
	"github.com/name5566/leaf/gate"
	"time"
	"github.com/name5566/leaf/log"
)

//type:"xyz",x代表左右阵营,y代表物理、化学、生物(0,1,2),z代表初始及各种分支(0,1,2,3)
type Hero struct {
	ID        int
	Type      string
	Transform *msg.TFServer

	HPMax  float64
	HP     float64
	HPHot  float64
	MPMax  float64
	MP     float64
	MPHot  float64
	Speed  float64
	Attack float64
	Def    float64

	Debuff *Burn
}

func (h *Hero) Burning(p *Player, room Room) {
	quit := make(chan int, 1)
	go func(quit chan int) {
		for {
			select {
			case <-h.Debuff.Timer.C:
				log.Debug("debuff clear")
				h.Debuff.IsEffect = false
				quit <- 1
				return
			}
		}
	}(quit)
	go func(quit chan int) {
		for {
			select {
			case <-h.Debuff.Ticker.C:
				if _, ok := p.GetHeros(h.ID); !ok {
					quit <- 1
				}
				h.SubHP(h.Debuff.Dot, room)
			case <-quit:
				log.Debug("debuff quit")
				return
			}
		}
	}(quit)
}

func (p *Player) CreateHero(heroType string, id, which int) (*Hero, error) {
	born := msg.TFServer{}
	if which == 0 {
		born = BornPosition1
	} else {
		born = BornPosition2
	}
	switch heroType {
	case "000":
		if p.Base.Money < 400 {
			return nil, errors.New("not enough money")
		}
		p.Base.Money -= 400
		hero := &Hero{
			ID:        id,
			Type:      "000",
			Transform: &born,
			// Skills
			HPMax:  400.0,
			HP:     400.0,
			HPHot:  10.0,
			MPMax:  100.0,
			MP:     100.0,
			MPHot:  15.0,
			Speed:  10.0,
			Attack: 10.0,
			Def:    10.0,

			Debuff: nil,
		}
		p.SetHeros(id, hero)
		return hero, nil
	default:
		return nil, errors.New("type error")
	}
}

func HealByHot(room Room, id int, a gate.Agent) {
	ticker := time.NewTicker(time.Second * 3)
	p := room.Players[a]
	for {
		select {
		case <-ticker.C:
			h, ok := p.GetHeros(id)
			if !ok {
				return
			}
			if h.HP == h.MPMax && h.MP == h.MPMax {
				continue
			}
			h.HP += h.HPHot
			if h.HP > h.HPMax {
				h.HP = h.HPMax
			}
			h.MP += h.MPHot
			if h.MP > h.MPMax {
				h.MP = h.MPMax
			}
			for aa := range room.Players {
				aa.WriteMsg(&msg.UpdateHeroState{
					id,
					h.Type,
					h.HP,
					h.MP,
				})
			}
		}
	}
}

func (h *Hero) UpdatePosition(t msg.TFServer) {
	h.Transform = &t
}

func (h *Hero) SubHP(damage float64, room Room) {
	h.HP -= damage
	if h.HP <= 0 {
		h.HP = 0
	}
	for aa, pp := range room.Players {
		aa.WriteMsg(&msg.UpdateHeroState{
			Id:   h.ID,
			Type: h.Type,
			Hp:   h.HP,
			Mp:   h.MP,
		})
		aa.WriteMsg(&msg.Damage{
			Id:     h.ID,
			Damage: damage,
		})
		if h.HP == 0 {
			if _, ok := pp.GetHeros(h.ID); ok {
				pp.DeleteHero(h.ID)
			}
			aa.WriteMsg(&msg.DeleteHero{h.ID})
		}
	}
}

func (h *Hero) Heal(heal float64, a gate.Agent, room Room, t string) {
	if t == "HP" {
		h.HP += heal
		if h.HP > h.HPMax {
			h.HP = h.HPMax
		}
	} else if t == "MP" {
		h.MP += heal
		if h.MP > h.MPMax {
			h.MP = h.MPMax
		}
	}
	for aa := range room.Players {
		aa.WriteMsg(&msg.UpdateHeroState{
			Id:   h.ID,
			Type: h.Type,
			Hp:   h.HP,
			Mp:   h.MP,
		})
	}
}
