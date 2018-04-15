package internal

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"server/msg"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
)

var (
	ID = 0
)

func init() {
	handler(&msg.Match{}, handleMatch)
	handler(&msg.CreateHero{}, handleCreateHero)
	handler(&msg.UpdatePosition{}, handleUpdatePosition)
	handler(&msg.UseSkill{}, handleUseSkill)
	handler(&msg.GetResource{}, handleGetResource)
	handler(&msg.SkillCrash{}, handleSkillCrash)
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handleMatch(args []interface{}) {
	m := args[0].(*msg.Match)
	a := args[1].(gate.Agent)
	log.Debug("Call Match from %v", a.RemoteAddr())
	room := new(Room)

	if LastRoomId == 0 || GetRoom(LastRoomId).PlayerCount == 2 {
		roomId := LastRoomId + 1
		room = NewRoom(roomId)
		AddRoom(room)
		LastRoomId = roomId
	} else {
		room = GetRoom(LastRoomId)
	}
	fmt.Println("playerId:", m.PlayerId)
	room.RoomPlayers[m.PlayerId] = a
	room.PlayerCount += 1
	fmt.Println("RoomStat:", room)
	if room.PlayerCount == 1 {
		a.WriteMsg(&msg.MatchStat{
			Status: 1,
			Msg:    "匹配中",
		})
	} else {
		rand.Seed(time.Now().Unix())
		r := rand.Intn(100)
		fmt.Println(r)
		for i, aa := range room.RoomPlayers {
			if r >= 50 {
				if i != m.PlayerId {
					player1 := NewPlayer(0)
					player1.Which = 0
					room.Players[aa] = player1
					aa.WriteMsg(&msg.MatchStat{
						Status:      0,
						Msg:         "匹配成功！",
						RoomId:      LastRoomId,
						WhichPlayer: 0,
					})
				} else {
					player2 := NewPlayer(1)
					player2.Which = 1
					room.Players[aa] = player2
					aa.WriteMsg(&msg.MatchStat{
						Status:      0,
						Msg:         "匹配成功！",
						RoomId:      LastRoomId,
						WhichPlayer: 1,
					})
				}
			} else {
				if i != m.PlayerId {
					player1 := NewPlayer(1)
					player1.Which = 1
					room.Players[aa] = player1
					aa.WriteMsg(&msg.MatchStat{
						Status:      0,
						Msg:         "匹配成功！",
						RoomId:      LastRoomId,
						WhichPlayer: 1,
					})
				} else {
					player2 := NewPlayer(0)
					player2.Which = 0
					room.Players[aa] = player2
					aa.WriteMsg(&msg.MatchStat{
						Status:      0,
						Msg:         "匹配成功！",
						RoomId:      LastRoomId,
						WhichPlayer: 0,
					})
				}
			}
		}
		go StartBattle(room)
	}
}

func handleCreateHero(args []interface{}) {
	m := args[0].(*msg.CreateHero)
	a := args[1].(gate.Agent)
	log.Debug("Call CreateHero from %v", a.RemoteAddr())
	room := GetRoom(m.RoomId)
	player := room.Players[a]
	which := int(player.Which)
	hero, err := player.CreateHero(m.HeroType, room.Count+1, which)
	go HealByHot(*room, hero.ID, a)
	if err != nil {

	} else {
		room.Count += 1
		for aa, _ := range room.Players {
			aa.WriteMsg(&msg.CreateHeroInf{
				HeroType:    hero.Type,
				TFServer:    *hero.Transform,
				WhichPlayer: which,
				ID:          room.Count,
				HPMax:       hero.HPMax,
				HP:          hero.HP,
				HPHot:       hero.HPHot,
				MPMax:       hero.MPMax,
				MP:          hero.MP,
				MPHot:       hero.MPHot,
				Speed:       hero.Speed,
				Attack:      hero.Attack,
				Def:         hero.Def,
			})
		}
		a.WriteMsg(&msg.MoneyLeft{
			MoneyLeft: player.Base.Money,
		})
	}
}

func handleUpdatePosition(args []interface{}) {
	m := args[0].(*msg.UpdatePosition)
	a := args[1].(gate.Agent)
	room := GetRoom(m.RoomId)
	p := room.Players[a]
	h, ok := p.GetHeros(m.Id)
	if ok {
		h.UpdatePosition(m.TfServer)
		for aa := range room.Players {
			if aa != a {
				aa.WriteMsg(m)
			}
		}
	} else {
		middle, ok := room.GetMiddle(m.Id)
		if !ok {
			log.Debug("fail to get middle")
			return
		}
		if v, ok := middle.(*ElectricBall); ok {
			v.UpdateBallPosition(m.TfServer)
		}
	}
}

func handleUseSkill(args []interface{}) {
	m := args[0].(*msg.UseSkill)
	a := args[1].(gate.Agent)
	log.Debug("call skill %s from %v", m.SkillID, a.RemoteAddr())
	room := GetRoom(m.RoomId)
	p := room.Players[a]
	h, ok := p.GetHeros(m.Id)
	if !ok {
		log.Debug("can't get hero")
	}
	skill := GetSkill(m.SkillID)
	skill.InitSkill()
	skill.Cast(a, room, h, m.TfServer)
}

func handleGetResource(args []interface{}) {
	m := args[0].(*msg.GetResource)
	a := args[1].(gate.Agent)
	log.Debug("call getResource %d from %v", m.ItemId, a.RemoteAddr())
	room := GetRoom(m.RoomId)
	p := room.Players[a]
	item, ok := room.GetMiddle(m.ItemId)
	if !ok {
		log.Debug("no middle id:%d", m.ItemId)
		return
	}
	switch middle := item.(type) {
	case *Gold:
		p.Base.Money += middle.Value
		a.WriteMsg(&msg.MoneyLeft{p.Base.Money})
		room.DeleteMiddle(m.ItemId)
		for aa := range room.Players {
			aa.WriteMsg(&msg.DeleteMiddle{m.ItemId})
		}
	case *Blood:
		hero, ok := p.GetHeros(m.HeroId)
		if !ok {
			log.Debug("no hero id:%d", m.HeroId)
		}
		hero.AddHP(float64(middle.value), a, *room)
	case *Mana:
		hero, ok := p.GetHeros(m.HeroId)
		if !ok {
			log.Debug("no hero id:%d", m.HeroId)
		}
		hero.AddMP(float64(middle.value), a, *room)
	}
}

func handleSkillCrash(args []interface{}) {
	m := args[0].(*msg.SkillCrash)
	a := args[1].(gate.Agent)
	log.Debug("call skill-crash from %d to %d from %v", m.FromHeroId, m.ToId, a.RemoteAddr())
	room := GetRoom(m.RoomId)
	p := room.Players[a]
	enemy := GetEnemy(a, *room)
	fromHero, _ := p.GetHeros(m.FromHeroId)
	fromItem, ok := room.GetMiddle(m.FromItemId)
	if !ok {
		log.Debug("not has item %d", m.FromItemId)
		return
	}
	var damage float64
	if att, ok := fromItem.GetAttack(); !ok {
		log.Debug("middle %d has no attack", m.FromItemId)
		return
	} else {
		damage = fromHero.Attack + att
	}
	if m.ToId == 0 || m.ToId == 1 {
		if m.ToId != enemy.Which {
			return
		}
		toBase := enemy.Base
		toBase.SubHP(damage, enemy.Which, *room)
	} else {
		toHero, ok := enemy.GetHeros(m.ToId)
		if ok {
			toHero.SubHP(damage, *room)
		} else {
			toMiddle, ok := room.GetMiddle(m.ToId)
			if !ok {
				log.Debug("no object %d", m.ToId)
			}
			if toMiddle.IsInvincible() {
				return
			} else {
				toMiddle.SubHp(damage, *room)
			}
		}
	}
	room.DeleteMiddle(m.FromItemId)
	for aa := range room.Players {
		aa.WriteMsg(&msg.DeleteMiddle{m.FromItemId})
	}
}
