package internal

import (
	"reflect"
	"server/msg"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"fmt"
	"server/gamedata"
)

func init() {
	handler(&msg.Match{}, handleMatch)
	handler(&msg.QuitMatch{}, handleQuitMatch)
	handler(&msg.GetUserInfo{}, handleGetUserInfo)

	handler(&msg.CreateRoom{}, handleCreateRoom)
	handler(&msg.GetRoomList{}, handleGetRoomList)
	handler(&msg.EnterRoom{}, handleEnterRoom)
	handler(&msg.ExitRoom{}, handleExitRoom)
	handler(&msg.StartBattle{}, handleStartBattle)

	handler(&msg.CreateHero{}, handleCreateHero)
	handler(&msg.UpdatePosition{}, handleUpdatePosition)
	handler(&msg.UseSkill{}, handleUseSkill)
	handler(&msg.GetResource{}, handleGetResource)
	handler(&msg.SkillCrash{}, handleSkillCrash)
	handler(&msg.Surrender{}, handleSurrender)
	handler(&msg.FireBottleCrash{}, handleFireBottleCrash)
	handler(&msg.Upgrade{}, handleUpgrade)

	handler(&msg.Test{}, handleTest)
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handleMatch(args []interface{}) {
	a := args[1].(gate.Agent)
	log.Debug("玩家 %s 匹配中...", Users[a])
	room, ok := GetRoom(LastMatchId)
	if !ok || room.PlayerCount == 2 {
		roomId := LastRoomId + 1
		LastRoomId += 1
		LastMatchId = LastRoomId
		room = NewRoom(roomId, fmt.Sprintf("房间_%d", roomId), Match, a)
	}
	if room.PlayerCount == 1 {
		a.WriteMsg(&msg.MatchStat{
			Status: 1,
			Msg:    "匹配中",
		})
	} else { // 匹配成功
		//UpdateRoomInfo(room)
		StartBattle(room)
	}
}

func handleCreateHero(args []interface{}) {
	m := args[0].(*msg.CreateHero)
	a := args[1].(gate.Agent)
	userName := Users[a]
	log.Debug("玩家 %s 创建英雄 %s", Users[a], m.HeroType)
	room, ok := GetRoom(m.RoomId)
	if !ok {
		log.Debug("获取房间失败")
		return
	}
	player := room.Players[userName]
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
		for _, aa := range room.User2Agent {
			if aa == nil {
				continue
			}
			(*aa).WriteMsg(&msg.CreateHeroInf{
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
	p := room.Players[Users[a]]
	h, ok := p.GetHeros(m.Id)
	if ok {
		h.UpdatePosition(m.TfServer)
		for _, aa := range room.User2Agent {
			if aa == nil {
				continue
			}
			if (*aa) != a {
				(*aa).WriteMsg(m)
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
	p := room.Players[Users[a]]
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
	p := room.Players[Users[a]]
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
		hero.Heal(float64(middle.value), *room, "HP")
	case *Mana:
		hero, ok := p.GetHeros(m.HeroId)
		if !ok {
			log.Debug("no hero id:%d", m.HeroId)
		}
		hero.Heal(float64(middle.value), *room, "MP")
	}
	room.DeleteMiddle(m.ItemId)
	for _, aa := range room.User2Agent {
		if aa == nil {
			continue
		}
		(*aa).WriteMsg(&msg.DeleteMiddle{m.ItemId})
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
	p := room.Players[Users[a]]
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
	for _, aa := range room.User2Agent {
		if aa == nil {
			continue
		}
		(*aa).WriteMsg(&msg.DeleteMiddle{m.FromItemId})
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
		for _, aa := range room.User2Agent {
			if aa == nil {
				continue
			}
			(*aa).WriteMsg(&msg.DeleteMiddle{ID: v.ID})
			(*aa).WriteMsg(&msg.CreateMiddle{
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
	log.Debug("玩家 %s 投降", Users[a])
	EndBattle(m.RoomId, a)
	DeleteRoom(m.RoomId, a)
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
	for user, pp := range room.Players {
		if user == Users[a] {
			pp.Upgrade(room, m.Id, m.TypeOld, m.TypeNew)
		}
	}
}

func handleQuitMatch(args []interface{}) {
	a := args[1].(gate.Agent)
	QuitMatch(a)
}

func handleTest(args []interface{}) {
	m := args[0].(*msg.Test)
	log.Debug("测试数据 %v", m)
	a := args[1].(gate.Agent)
	a.WriteMsg(&msg.RoomInfo{
		Msg:    "ok",
		RoomId: 1,
		Name:   "Test1",
		Users: map[string]*msg.User{
			"User1": {UserName: "User1"},
			"User2": {UserName: "User2"},
		},
	})
	a.WriteMsg(&msg.RoomInfo{
		Msg:    "ok",
		RoomId: 2,
		Name:   "Test2",
		Users: map[string]*msg.User{
			"User2": {UserName: "User2"},
		},
	})
}

func handleGetRoomList(args []interface{}) {
	m := args[0].(*msg.GetRoomList)
	a := args[1].(gate.Agent)
	log.Debug("正在获取第 %d 页...", m.PageNum)
	roomList := []*msg.RoomInfo{}
	pageNum := m.PageNum - 1
	if len(RoomList) > (pageNum+1)*5 {
		roomList = RoomList[pageNum*5 : (pageNum+1)*5]
	} else if len(RoomList) > pageNum*5 {
		roomList = RoomList[pageNum*5:]
	}
	log.Debug("房间列表 %v", roomList)
	a.WriteMsg(&msg.RoomList{
		RoomList: roomList,
	})
}

func handleCreateRoom(args []interface{}) {
	m := args[0].(*msg.CreateRoom)
	a := args[1].(gate.Agent)
	LastRoomId += 1
	roomId := LastRoomId
	room := NewRoom(roomId, m.Name, Spec, a)
	UpdateRoomInfo(room)
}

func handleEnterRoom(args []interface{}) {
	m := args[0].(*msg.EnterRoom)
	a := args[1].(gate.Agent)
	roomId := m.RoomId
	room, ok := GetRoom(roomId)
	if !ok {
		a.WriteMsg(&msg.RoomInfo{
			Msg: "房间已关闭,请刷新列表",
		})
		return
	} else if room.PlayerCount == 2 || room.InBattle == true {
		a.WriteMsg(&msg.RoomInfo{
			Msg: "房间已满",
		})
		return
	}
	Agent2Room[a] = roomId
	userData := gamedata.UsersMap[Users[a]]
	room.Users[Users[a]] = &msg.User{
		UserName: userData.Name,
		Photo:    userData.Photo,
		Total:    userData.Total,
		Victory:  userData.Victory,
		Defeat:   userData.Defeat,
		Rate:     userData.Defeat,
		KeyOwner: false,
	}
	room.PlayerCount += 1
	room.Players[Users[a]] = nil
	room.User2Agent[Users[a]] = &a
	UpdateRoomInfo(room)
}

func handleExitRoom(args []interface{}) {
	a := args[1].(gate.Agent)
	ExitRoom(a, Rooms[Agent2Room[a]])
}

func handleStartBattle(args []interface{}) {
	a := args[1].(gate.Agent)
	room, ok := Rooms[Agent2Room[a]]
	if !ok {
		log.Debug("房间不存在,开启战斗失败")
		return
	}
	if room.PlayerCount != 2 {
		log.Debug("人数过少,开启战斗失败")
		return
	}
	if room.Users[Users[a]].KeyOwner {
		StartBattle(room)
	}
}

func handleGetUserInfo(args []interface{}) {
	a := args[1].(gate.Agent)
	userData := new(gamedata.UserData)
	has, err := gamedata.Db.Where("name=?", Users[a]).Get(userData)
	if err != nil || !has {
		log.Debug("获取角色信息失败")
		return
	}
	a.WriteMsg(&msg.User{
		UserName: userData.Name,
		Photo:    userData.Photo,
		Total:    userData.Total,
		Victory:  userData.Victory,
		Defeat:   userData.Defeat,
		Rate:     userData.Rate,
		KeyOwner: false,
	})
}
