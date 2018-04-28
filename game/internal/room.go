package internal

import (
	"github.com/name5566/leaf/gate"
	"time"
	"math/rand"
	"server/msg"
	"github.com/name5566/leaf/log"
	"sync"
	"server/gamedata"
)

const (
	Match = "匹配"
	Spec  = "指定"
)

type Room struct {
	Name        string
	Users       map[string]*msg.User   // userName -> User
	User2Agent  map[string]*gate.Agent // userName -> Agent
	Players     map[string]*Player     // userName -> Player
	Count       int                    // 物体数量,包括英雄与中立生物/地形
	RoomId      int
	PlayerCount int
	Middle      map[int]Middle
	Lock        sync.Mutex
	Closed      bool
	InBattle    bool
	Mode        string
}

func (r *Room) GetMiddle(k int) (Middle, bool) {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	m, ok := r.Middle[k]
	return m, ok
}

func (r *Room) SetMiddle(k int, v Middle) {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	r.Middle[k] = v
}

func (r *Room) DeleteMiddle(k int) {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	delete(r.Middle, k)
}

func (r *Room) StartMapEvent() {
	r.InBattle = true
	for user := range r.User2Agent {
		pp := r.Players[user]
		go pp.Base.GetMoneyByTime(user)
	}
	go HealByHot(r)
	r.RandomResource(time.Second*10, time.Second*5)
}

func (r *Room) SyncItems() {
	itemList := []int{0, 1}
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			if r.Closed {
				return
			}
			for _, pp := range r.Players {
				for _, h := range pp.Heros {
					itemList = append(itemList, h.ID)
				}
			}
			for _, m := range r.Middle {
				itemList = append(itemList, m.GetId())
			}
			for _, aa := range r.User2Agent {
				if aa != nil {
					(*aa).WriteMsg(&msg.SyncItems{
						ItemList: itemList,
					})
				}
			}
		}
	}
}

func StartBattle(room *Room) {
	deleteRoomInfo(room.RoomId) // 开始战斗后从房间列表中删除该房间
	rand.Seed(time.Now().Unix())
	r := rand.Intn(100)
	flag := 0
	i := 0
	if r >= 50 {
		flag = 1
	}
	for user, aa := range room.User2Agent {
		which := 0
		if flag == i {
			which = 1
		}
		i += 1
		player := NewPlayer(which)
		room.Players[user] = player
		(*aa).WriteMsg(&msg.MatchStat{
			Status:      0,
			Msg:         "开始战斗",
			RoomId:      room.RoomId,
			WhichPlayer: which,
		})
	}
	timer := time.NewTimer(time.Second * 3)
	<-timer.C
	for _, aa := range room.User2Agent {
		// 设置玩家信息为战斗中(用于断线重连)
		userData := gamedata.UsersMap[Users[*aa]]
		userData.InBattle = 1
		userData.RoomId = room.RoomId
		cond := gamedata.UserData{
			Id: userData.Id,
		}
		gamedata.Db.Update(userData, cond)
	}
	go room.StartMapEvent()
	go room.SyncItems()
}

func EndBattle(roomId int, lose gate.Agent) {
	log.Debug("战斗结束-正常")
	room, ok := GetRoom(roomId)
	if !ok {
		log.Debug("结束战斗时获取房间失败", lose.RemoteAddr())
		return
	}
	room.Closed = true
	for user, aa := range room.User2Agent {
		userData := gamedata.UsersMap[user]
		condi := gamedata.UserData{
			Id: userData.Id,
		}
		pp := room.Players[user]
		pp.Base.Timer.Reset(time.Millisecond)
		var isWin bool
		if aa == nil || (*aa) == lose {
			isWin = false
			userData.Total += 1
			userData.Defeat += 1
		} else {
			isWin = true
			userData.Total += 1
			userData.Victory += 1
		}
		userData.Rate = int(userData.Victory * 100 / userData.Total)
		userData.InBattle = 0

		effect, err := gamedata.Db.Update(userData, condi)
		if err != nil || int(effect) != 1 {
			log.Debug("更新数据失败")
		}
		cond := gamedata.UserData{
			Id: userData.Id,
		}
		gamedata.Db.Cols("in_battle").Update(userData, cond)
		if aa != nil {
			(*aa).WriteMsg(&msg.EndBattle{
				IsWin: isWin,
			})
		}
	}
}

func RecoverBattle(a gate.Agent, room *Room) {
	userName := Users[a]
	room.User2Agent[userName] = &a
	a.WriteMsg(&msg.MatchStat{
		Status:      0,
		Msg:         "重连成功",
		RoomId:      room.RoomId,
		WhichPlayer: room.Players[userName].Which,
	})
	log.Debug("%s 同步数据", userName)
	room.Count += 1
	timer := time.NewTimer(time.Second * 3)
	<-timer.C
	for user, pp := range room.Players {
		if user == userName {
			a.WriteMsg(&msg.MoneyLeft{
				MoneyLeft: pp.Base.Money,
			})
		}
		for _, hero := range pp.Heros {
			a.WriteMsg(&msg.CreateHeroInf{
				Msg:         "ok",
				HeroType:    hero.Type,
				TFServer:    *hero.Transform,
				WhichPlayer: pp.Which,
				ID:          hero.ID,
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
	}
	for _, middle := range room.Middle {
		a.WriteMsg(&msg.CreateMiddle{
			ID:   middle.GetId(),
			TF:   *middle.GetTF(),
			Type: middle.GetType(),
		})
	}
}

func (r *Room) RandomResource(beforeTime, interval time.Duration) {
	ticker1 := time.NewTicker(beforeTime)
	ticker2 := time.NewTicker(interval)

	<-ticker1.C

	for {
		select {
		case <-ticker2.C:
			if r.Closed {
				log.Debug("房间%d关闭-资源生成关闭", r.RoomId)
				return
			}
			rand.Seed(time.Now().Unix())
			x := float64(rand.Intn(6700))/100.0 + 29.00
			z := float64(rand.Intn(3300))/100.0 + 9.00
			tf := msg.TFServer{
				Position: []float64{x, 0.07, z},
				Rotation: []float64{0.0, 90.0, 0.0},
			}
			gold := NewGold(r.Count+1, tf)
			r.Count += 1
			r.SetMiddle(gold.ID, gold)
			for _, aa := range r.User2Agent {
				if aa == nil {
					continue
				}
				(*aa).WriteMsg(&msg.CreateMiddle{
					ID:   gold.ID,
					TF:   *gold.TF,
					Type: gold.Type,
				})
			}
			go gold.TakeAction(r)
		}
	}

}

func NewRoom(roomId int, name string, mode string, a gate.Agent) *Room {
	log.Debug("新建 %s 房间 %d", mode, roomId)
	room := Room{
		Name:        name,
		Users:       make(map[string]*msg.User),
		User2Agent:  make(map[string]*gate.Agent),
		Players:     make(map[string]*Player),
		Count:       2,
		RoomId:      roomId,
		PlayerCount: 0,
		Middle:      make(map[int]Middle),
		Lock:        sync.Mutex{},
		Closed:      false,
		InBattle:    false,
		Mode:        mode,
	}
	Rooms[room.RoomId] = &room
	userName := Users[a]
	room.Players[userName] = nil
	Agent2Room[a] = room.RoomId
	room.User2Agent[userName ] = &a
	user := msg.User{
		UserName: userName,
		KeyOwner: mode == Spec,
	}
	room.Users[userName] = &user
	room.PlayerCount += 1
	return &room
}

func QuitMatch(a gate.Agent) {
	if roomId, ok := Agent2Room[a]; ok {
		if room, ok := GetRoom(roomId); ok {
			if room.PlayerCount == 1 {
				delete(Agent2Room, a)
				delete(Rooms, room.RoomId)
				log.Debug("退出成功,房间%d删除", roomId)
				return
			}
			log.Debug("战斗已经开始，无法退出")
		}
	}
	log.Debug("获取房间失败")
}

func ExitRoom(a gate.Agent, room *Room) {
	if room.PlayerCount == 1 {
		DeleteRoom(room.RoomId, a, false)
	} else {
		room.PlayerCount -= 1
		delete(Agent2Room, a)
		delete(room.Players, Users[a])
		if room.Users[Users[a]].KeyOwner == true {
			for aa, user := range room.Users {
				if aa != Users[a] {
					user.KeyOwner = true
				}
			}
		}
		delete(room.Users, Users[a])
		delete(room.User2Agent, Users[a])
		UpdateRoomInfo(room)
	}
}
