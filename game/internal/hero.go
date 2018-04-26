package internal

import (
	"server/msg"
	"errors"
	"github.com/name5566/leaf/gate"
	"time"
	"github.com/name5566/leaf/log"
	"fmt"
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

func (p *Player) CreateHero(heroType string, id, which int, tf *msg.TFServer) (*Hero, error) {
	born := msg.TFServer{}
	if which == 0 {
		born = BornPosition1
	} else {
		born = BornPosition2
	}
	switch heroType {
	case "000", "100":
		if p.Base.Money < 400 {
			return nil, errors.New("金钱不足")
		}
		p.Base.Money -= 400
		hero := &Hero{
			ID:        id,
			Type:      heroType,
			Transform: &born,
			// Skills
			HPMax:  300.0,
			HP:     300.0,
			HPHot:  2.0,
			MPMax:  150.0,
			MP:     150.0,
			MPHot:  4.0,
			Speed:  10.0,
			Attack: 10.0,
			Def:    5.0,

			Debuff: nil,
		}
		return hero, nil

	case "001", "101":
		if p.Base.Money < 700 {
			return nil, errors.New("金钱不足")
		}
		p.Base.Money -= 700
		hero := &Hero{
			ID:        id,
			Type:      heroType,
			Transform: &born,
			// Skills
			HPMax:  300.0,
			HP:     300.0,
			HPHot:  2.0,
			MPMax:  200.0,
			MP:     200.0,
			MPHot:  5.0,
			Speed:  10.0,
			Attack: 20.0,
			Def:    20.0,

			Debuff: nil,
		}
		return hero, nil
	case "002", "102":
		if p.Base.Money < 700 {
			return nil, errors.New("金钱不足")
		}
		p.Base.Money -= 700
		hero := &Hero{
			ID:        id,
			Type:      heroType,
			Transform: &born,
			// Skills
			HPMax:  600.0,
			HP:     600.0,
			HPHot:  5.0,
			MPMax:  200.0,
			MP:     200.0,
			MPHot:  5.0,
			Speed:  12.0,
			Attack: 40.0,
			Def:    8.0,

			Debuff: nil,
		}
		return hero, nil

	case "010", "110":
		if p.Base.Money < 400 {
			return nil, errors.New("金钱不足")
		}
		p.Base.Money -= 400
		hero := &Hero{
			ID:        id,
			Type:      heroType,
			Transform: &born,
			// Skills
			HPMax:  200.0,
			HP:     200.0,
			HPHot:  2.0,
			MPMax:  100.0,
			MP:     100.0,
			MPHot:  4.0,
			Speed:  8.0,
			Attack: 20.0,
			Def:    2.0,

			Debuff: nil,
		}
		return hero, nil

	case "011", "111":
		if p.Base.Money < 900 {
			return nil, errors.New("金钱不足")
		}
		p.Base.Money -= 900
		hero := &Hero{
			ID:        id,
			Type:      heroType,
			Transform: &born,
			// Skills
			HPMax:  400.0,
			HP:     400.0,
			HPHot:  5.0,
			MPMax:  150.0,
			MP:     150.0,
			MPHot:  5.0,
			Speed:  12.0,
			Attack: 40.0,
			Def:    10.0,

			Debuff: nil,
		}
		return hero, nil

	case "020", "120":
		if p.Base.Money < 300 {
			return nil, errors.New("金钱不足")
		}
		p.Base.Money -= 300
		hero := &Hero{
			ID:        id,
			Type:      heroType,
			Transform: &born,
			// Skills
			HPMax:  200.0,
			HP:     200.0,
			HPHot:  2.0,
			MPMax:  100.0,
			MP:     100.0,
			MPHot:  2.0,
			Speed:  8.0,
			Attack: 10.0,
			Def:    2.0,

			Debuff: nil,
		}
		return hero, nil

	case "021", "121":
		if p.Base.Money < 600 {
			return nil, errors.New("金钱不足")
		}
		p.Base.Money -= 600
		hero := &Hero{
			ID:        id,
			Type:      heroType,
			Transform: &born,
			// Skills
			HPMax:  400.0,
			HP:     400.0,
			HPHot:  2.0,
			MPMax:  150.0,
			MP:     150.0,
			MPHot:  3.0,
			Speed:  8.0,
			Attack: 10.0,
			Def:    4.0,

			Debuff: nil,
		}
		return hero, nil

	default:
		return nil, errors.New("类型错误")
	}
}

func (p *Player) Upgrade(room *Room, id int, old, new string) {
	var a gate.Agent
	for user, pp := range room.Players {
		if pp == p {
			a = (*room.User2Agent[user])
		}
	}
	if len(old) != 3 || len(new) != 3 || (old[0:2] != new[0:2]) || (new[2] <= old[2]) {
		log.Debug("类型错误,进阶失败,%s->%s", old, new)
		a.WriteMsg(&msg.CreateHeroInf{Msg: "类型错误,进阶失败"})
		return
	}
	oldh, ok := p.GetHeros(id)
	if !ok {
		log.Debug("英雄死亡,进阶失败")
		a.WriteMsg(&msg.CreateHeroInf{Msg: "英雄死亡,进阶失败"})
		return
	}
	tf := oldh.Transform
	distance := GetDistance(p.Base.TF, tf)
	if distance > 20 {
		log.Debug("距离基地太远,进阶失败")
		a.WriteMsg(&msg.CreateHeroInf{Msg: "距离基地太远,进阶失败"})
		return
	}
	newh, err := p.CreateHero(new, room.Count+1, p.Which, tf)
	if err != nil {
		log.Debug("进阶失败,%s", err.Error())
		a.WriteMsg(&msg.CreateHeroInf{Msg: fmt.Sprintf("进阶失败,%s", err.Error())})
		return
	}
	room.Count += 1
	p.DeleteHero(id)
	p.SetHeros(newh.ID, newh)
	for _, aa := range room.User2Agent {
		if aa == nil {
			continue
		}
		(*aa).WriteMsg(&msg.DeleteHero{ID: id})
		(*aa).WriteMsg(&msg.CreateHeroInf{
			Msg:         "ok",
			HeroType:    newh.Type,
			TFServer:    *newh.Transform,
			WhichPlayer: p.Which,
			ID:          newh.ID,
			HPMax:       newh.HPMax,
			HP:          newh.HP,
			HPHot:       newh.HPHot,
			MPMax:       newh.MPMax,
			MP:          newh.MP,
			MPHot:       newh.MPHot,
			Speed:       newh.Speed,
			Attack:      newh.Attack,
			Def:         newh.Def,
		})
	}
	a.WriteMsg(&msg.MoneyLeft{MoneyLeft: p.Base.Money})
}

func HealByHot(room *Room) {
	ticker := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-ticker.C:
			if room.Closed {
				log.Debug("房间%d关闭-自动回血关闭", room.RoomId)
				return
			}
			for _, pp := range room.Players {
				for _, h := range pp.Heros {
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
					for _, aa := range room.User2Agent {
						if aa == nil {
							continue
						}
						(*aa).WriteMsg(&msg.UpdateHeroState{
							h.ID,
							h.HP,
							h.MP,
						})
					}
				}
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
	for user, pp := range room.Players {
		aa := room.User2Agent[user]
		var flag = (aa != nil)

		if h.HP == 0 {
			if _, ok := pp.GetHeros(h.ID); ok {
				pp.DeleteHero(h.ID)
			}
			if flag {
				(*aa).WriteMsg(&msg.DeleteHero{h.ID})
			}
		}
		if !flag {
			continue
		}
		(*aa).WriteMsg(&msg.UpdateHeroState{
			Id: h.ID,
			Hp: h.HP,
			Mp: h.MP,
		})
		(*aa).WriteMsg(&msg.Damage{
			Id:     h.ID,
			Damage: damage,
		})
	}
}

func (h *Hero) Heal(heal float64, room Room, t string) {
	if t == "HP" {
		h.HP += heal
		if h.HP > h.HPMax {
			h.HP = h.HPMax
		}
	} else {
		h.MP += heal
		if h.MP > h.MPMax {
			h.MP = h.MPMax
		}
	}
	for _, aa := range room.User2Agent {
		if aa == nil {
			continue
		}
		(*aa).WriteMsg(&msg.UpdateHeroState{
			Id: h.ID,
			Hp: h.HP,
			Mp: h.MP,
		})
	}
}
