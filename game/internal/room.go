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
	Users       map[string]*msg.User
	Count       int // 物体数量,包括英雄与中立生物/地形
	RoomId      int
	Players     map[gate.Agent]*Player
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
	for agent, player := range r.Players {
		go player.Base.GetMoneyByTime(agent)
	}
	go HealByHot(r)
	r.RandomResource(time.Second*10, time.Second*5)
}

func StartBattle(room *Room) {
	rand.Seed(time.Now().Unix())
	r := rand.Intn(100)
	flag := 0
	i := 0
	if r >= 50 {
		flag = 1
	}
	for aa := range room.Players {
		which := 0
		if flag == i {
			which = 1
		}
		player := NewPlayer(which)
		room.Players[aa] = player
		aa.WriteMsg(&msg.MatchStat{
			Status:      0,
			Msg:         "开始战斗",
			RoomId:      room.RoomId,
			WhichPlayer: which,
		})
		// 设置玩家信息为战斗中(用于断线重连)
		gamedata.UsersMap[Users[aa]].InBattle = true
		gamedata.UsersMap[Users[aa]].RoomId = LastRoomId
		i += 1
	}
	go room.StartMapEvent()
}

func EndBattle(roomId int, lose gate.Agent) {
	room, ok := GetRoom(roomId)
	if !ok {
		log.Debug("结束战斗时获取房间失败", lose.RemoteAddr())
		return
	}
	room.Closed = true
	for aa, pp := range room.Players {
		pp.Base.Timer.Reset(time.Millisecond)
		aa.WriteMsg(&msg.EndBattle{
			IsWin: aa != lose,
		})
		gamedata.UsersMap[Users[aa]].InBattle = false
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
			for aa := range r.Players {
				aa.WriteMsg(&msg.CreateMiddle{
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
		Count:       2,
		RoomId:      roomId,
		Players:     make(map[gate.Agent]*Player),
		PlayerCount: 0,
		Middle:      make(map[int]Middle),
		Lock:        sync.Mutex{},
		Closed:      false,
		InBattle:    false,
		Mode:        mode,
	}
	AddRoom(&room)
	room.Players[a] = nil
	Agent2Room[a] = room.RoomId
	user := msg.User{
		UserName: Users[a],
		KeyOwner: mode == Spec,
	}
	room.Users[user.UserName] = &user
	room.PlayerCount += 1
	return &room
}

func QuitMatch(a gate.Agent) {
	if roomId, ok := Agent2Room[a]; ok {
		if room, ok := GetRoom(roomId); ok {
			if room.PlayerCount == 1 {
				DeleteRoom(roomId, a)
				log.Debug("退出成功,房间%d删除", roomId)
				return
			}
		}
	}
	log.Debug("获取房间失败")
}

func ExitRoom(a gate.Agent, room *Room) {
	if room.PlayerCount == 1 {
		DeleteRoom(room.RoomId, a)
	} else {
		room.PlayerCount -= 1
		delete(Agent2Room, a)
		delete(room.Players, a)
		if room.Users[Users[a]].KeyOwner == true {
			for aa, user := range room.Users {
				if aa != Users[a] {
					user.KeyOwner = true
				}
			}
		}
		delete(room.Users, Users[a])
		UpdateRoomInfo(room)
	}
}
