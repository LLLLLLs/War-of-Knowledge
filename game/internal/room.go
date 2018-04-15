package internal

import (
	"fmt"
	"github.com/name5566/leaf/gate"
	"time"
	"math/rand"
	"server/msg"
	"github.com/name5566/leaf/log"
	"sync"
)

type Room struct {
	Count       int // 物体数量,包括英雄与中立生物/地形
	RoomId      int
	RoomPlayers map[int]gate.Agent
	Players     map[gate.Agent]*Player
	PlayerCount int
	Middle      map[int]Middle
	Lock        sync.Mutex
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

func StartBattle(room *Room) {
	agents := make([]gate.Agent, 2)
	players := make([]*Player, 2)
	i := int(0)
	for agent, player := range room.Players {
		agents[i] = agent
		players[i] = player
		i += 1
	}
	//go players[0].Base.GetMoneyByTime(agents[0])
	//go players[1].Base.GetMoneyByTime(agents[1])
	room.RandomResource(time.Second*10, time.Second*5)
}

func (r *Room) RandomResource(beforeTime, interval time.Duration) {
	ticker1 := time.NewTicker(beforeTime)
	ticker2 := time.NewTicker(interval)

	<-ticker1.C

	for {
		select {
		case <-ticker2.C:
			rand.Seed(time.Now().Unix())
			x := float64(rand.Intn(6700))/100.0 + 29.00
			z := float64(rand.Intn(3300))/100.0 + 9.00
			tf := msg.TFServer{
				[]float64{x, 0.07, z},
				[]float64{0.0, 90.0, 0.0},
			}
			gold := NewGold(r.Count+1, tf)
			r.Count += 1
			r.SetMiddle(gold.ID, gold)
			for aa := range r.Players {
				aa.WriteMsg(&msg.CreateMiddle{
					gold.ID,
					*gold.TF,
					gold.Type,
				})
			}
			log.Debug("room: %d create resource %d", r.RoomId, gold.ID)
			go gold.TakeAction(r)
		}
	}

}

func NewRoom(roomId int) *Room {
	fmt.Println("newRoom:", roomId)
	room := Room{
		Count:       2,
		RoomId:      roomId,
		RoomPlayers: make(map[int]gate.Agent),
		Players:     make(map[gate.Agent]*Player),
		PlayerCount: 0,
		Middle:      make(map[int]Middle),
		Lock:        sync.Mutex{},
	}
	return &room
}
