package internal

import (
	"fmt"
	"github.com/name5566/leaf/gate"
	"time"
	"math/rand"
	"server/msg"
	"github.com/name5566/leaf/log"
	"sync"
	"server/gamedata"
)

type Room struct {
	Name        string
	Count       int // 物体数量,包括英雄与中立生物/地形
	RoomId      int
	Players     map[gate.Agent]*Player
	PlayerCount int
	Middle      map[int]Middle
	Lock        sync.Mutex
	Closed      bool
	Matching    bool
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

func (r *Room) StartBattle() {
	r.Matching = false
	for agent, player := range r.Players {
		go player.Base.GetMoneyByTime(agent)
	}
	go HealByHot(r)
	r.RandomResource(time.Second*10, time.Second*5)
}

func EndBattle(roomId int, lose gate.Agent) {
	room, ok := GetRoom(roomId)
	if !ok {
		log.Debug("结束战斗时获取房间失败", lose.RemoteAddr())
		return
	}
	room.Closed = true
	delete(Room2Agent, roomId)
	delete(Rooms, roomId)
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
			//log.Debug("room: %d create resource %d", r.RoomId, gold.ID)
			go gold.TakeAction(r)
		}
	}

}

func NewRoom(roomId int, name string) *Room {
	fmt.Println("newRoom:", roomId)
	room := Room{
		Name:        name,
		Count:       2,
		RoomId:      roomId,
		Players:     make(map[gate.Agent]*Player),
		PlayerCount: 0,
		Middle:      make(map[int]Middle),
		Lock:        sync.Mutex{},
		Closed:      false,
		Matching:    true,
	}
	return &room
}

func QuitMatch(a gate.Agent) {
	if roomId, ok := Agent2Room[a]; ok {
		if room, ok := GetRoom(roomId); ok {
			if room.PlayerCount == 1 {
				room.PlayerCount -= 1
				delete(Agent2Room, a)
				DeleteRoom(roomId)
				delete(Room2Agent, roomId)
				log.Debug("退出成功,房间%d删除", roomId)
				return
			} else {
				if room.Matching == false {
					log.Debug("正在战斗中,无法退出匹配")
					return
				} else {
					room.PlayerCount -= 1
					delete(room.Players, a)
					for aa := range room.Players {
						aa.WriteMsg(&msg.RoomInfo{
							RoomId: roomId,
							Name:   room.Name,
							Users: []msg.User{
								{UserName: Users[a]},
							},
						})
					}
				}
			}
		}
	}
	log.Debug("获取房间失败")
}
