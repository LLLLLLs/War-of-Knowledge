package internal

import (
	"math/rand"
	"reflect"
	"time"

	"server/msg"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"server/gamedata"
	"fmt"
)

func init() {
	handler(&msg.Match{}, handleMatch)
	handler(&msg.QuitMatch{}, handleQuitMatch)
	handler(&msg.CreateHero{}, handleCreateHero)
	handler(&msg.UpdatePosition{}, handleUpdatePosition)
	handler(&msg.UseSkill{}, handleUseSkill)
	handler(&msg.GetResource{}, handleGetResource)
	handler(&msg.SkillCrash{}, handleSkillCrash)
	handler(&msg.Surrender{}, handleSurrender)
	handler(&msg.FireBottleCrash{}, handleFireBottleCrash)
	handler(&msg.Upgrade{}, handleUpgrade)

}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handleMatch(args []interface{}) {
	a := args[1].(gate.Agent)
	log.Debug("Call Match from %v", a.RemoteAddr())
	room := new(Room)
	lastRoom, ok := GetRoom(LastRoomId)
	if !ok || lastRoom.PlayerCount == 2 {
		roomId := LastRoomId
		if ok {
			roomId = LastRoomId + 1
			LastRoomId += 1
		}
		room = NewRoom(roomId, fmt.Sprintf("房间_%d", roomId))
		AddRoom(room)
	} else {
		room, ok = GetRoom(LastRoomId)
		if !ok {
			log.Debug("%v get room fail", a.RemoteAddr())
			return
		}
	}
	Room2Agent[room.RoomId] = append(Room2Agent[room.RoomId], a)
	Agent2Room[a] = room.RoomId
	room.PlayerCount += 1
	if room.PlayerCount == 1 {
		a.WriteMsg(&msg.MatchStat{
			Status: 1,
			Msg:    "匹配中",
		})
	} else { // 匹配成功
		rand.Seed(time.Now().Unix())
		r := rand.Intn(100)
		flag := 0
		if r >= 50 {
			flag = 1
		}
		for i, aa := range Room2Agent[room.RoomId] {
			which := 0
			if flag == i {
				which = 1
			}
			player := NewPlayer(which)
			room.Players[aa] = player
			aa.WriteMsg(&msg.MatchStat{
				Status:      0,
				Msg:         "匹配成功",
				RoomId:      LastRoomId,
				WhichPlayer: which,
			})
			// 设置玩家信息为战斗中(用于断线重连)
			gamedata.UsersMap[Users[aa]].InBattle = true
			gamedata.UsersMap[Users[aa]].RoomId = LastRoomId
		}
		go room.StartBattle()
	}
}

func handleCreateHero(args []interface{}) {
	m := args[0].(*msg.CreateHero)
	a := args[1].(gate.Agent)
	log.Debug("call createHero from %v", a.RemoteAddr())
	room, ok := GetRoom(m.RoomId)
	if !ok {
		log.Debug("%v get room fail", a.RemoteAddr())
		return
	}
	player := room.Players[a]
	which := int(player.Which)
	hero, err := player.CreateHero(m.HeroType, room.Count+1, which, nil)
	if err != nil {
		log.Debug("创建英雄失败,%s", err.Error())
		a.WriteMsg(&msg.CreateHeroInf{
			Msg: err.Error(),
		})
		return
	} else {
		log.Debug("创建英雄成功,%s", m.HeroType)
		player.SetHeros(hero.ID, hero)
		room.Count += 1
		for aa := range room.Players {
			aa.WriteMsg(&msg.CreateHeroInf{
				Msg:         "ok",
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
	room, ok := GetRoom(m.RoomId)
	if !ok {
		log.Debug("UpdatePosition: %v 房间不存在", a.RemoteAddr())
		return
	}
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
			log.Debug("UpdatePosition:获取中立生物失败")
			return
		}
		switch v := middle.(type) {
		case *ElectricBall:
			v.UpdateBallPosition(m.TfServer)
		case *FireBottle:
			v.UpdateFireBottle(m.TfServer, room)
		}
	}
}

func handleUseSkill(args []interface{}) {
	m := args[0].(*msg.UseSkill)
	a := args[1].(gate.Agent)
	log.Debug("call skill %s from %v", m.SkillID, a.RemoteAddr())
	room, ok := GetRoom(m.RoomId)
	if !ok {
		log.Debug("%v get room fail", a.RemoteAddr())
		return
	}
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
	room, ok := GetRoom(m.RoomId)
	if !ok {
		log.Debug("%v get room fail", a.RemoteAddr())
		return
	}
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
	case *Blood:
		hero, ok := p.GetHeros(m.HeroId)
		if !ok {
			log.Debug("no hero id:%d", m.HeroId)
		}
		hero.Heal(float64(middle.value), a, *room, "HP")
	case *Mana:
		hero, ok := p.GetHeros(m.HeroId)
		if !ok {
			log.Debug("no hero id:%d", m.HeroId)
		}
		hero.Heal(float64(middle.value), a, *room, "MP")
	}
	room.DeleteMiddle(m.ItemId)
	for aa := range room.Players {
		aa.WriteMsg(&msg.DeleteMiddle{m.ItemId})
	}
}

func handleSkillCrash(args []interface{}) {
	m := args[0].(*msg.SkillCrash)
	a := args[1].(gate.Agent)
	log.Debug("call skill-crash from %d to %d from %v", m.FromHeroId, m.ToId, a.RemoteAddr())
	room, ok := GetRoom(m.RoomId)
	if !ok {
		log.Debug("%v get room fail", a.RemoteAddr())
		return
	}
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

func handleFireBottleCrash(args []interface{}) {
	m := args[0].(*msg.FireBottleCrash)
	room, ok := GetRoom(m.RoomId)
	if !ok {
		log.Debug("wrong room id")
		return
	}
	log.Debug("fire bottle burst")
	fb, ok := room.GetMiddle(m.ItemId)
	if !ok {
		log.Debug("no fire bottle id %d", m.ItemId)
		return
	}
	room.DeleteMiddle(m.ItemId)
	if v, ok := fb.(*FireBottle); ok {
		firesea := NewFireSea(room.Count+1, *v.TF)
		for aa := range room.Players {
			aa.WriteMsg(&msg.DeleteMiddle{ID: v.ID})
			aa.WriteMsg(&msg.CreateMiddle{
				ID:   firesea.ID,
				TF:   *firesea.TF,
				Type: firesea.Type,
			})
		}
		room.Count += 1
		room.SetMiddle(firesea.ID, firesea)
		firesea.TakeAction_(room, v.Hero)
	}
}

func handleSurrender(args []interface{}) {
	m := args[0].(*msg.Surrender)
	a := args[1].(gate.Agent)
	log.Debug("call surrender from %v", a.RemoteAddr())
	EndBattle(m.RoomId, a)
}

func handleUpgrade(args []interface{}) {
	m := args[0].(*msg.Upgrade)
	a := args[1].(gate.Agent)
	log.Debug("英雄 %d 升级", m.Id)
	room, ok := GetRoom(m.RoomId)
	if !ok {
		log.Debug("升级时获取房间失败")
		return
	}
	for aa, pp := range room.Players {
		if aa == a {
			pp.Upgrade(room, m.Id, m.TypeOld, m.TypeNew)
		}
	}
}

func handleQuitMatch(args []interface{}) {
	a := args[1].(gate.Agent)
	QuitMatch(a)
}
