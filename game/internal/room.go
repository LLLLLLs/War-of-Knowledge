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
	Count       int // 物体数量,包括英雄与中立生物/地形
	RoomId      int
	Players     map[gate.Agent]*Player
	PlayerCount int
	Middle      map[int]Middle
	Lock        sync.Mutex
	Timer       *RoomTimer
}

type RoomTimer struct {
	EndTimer *time.Timer
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
	for agent, player := range r.Players {
		go player.Base.GetMoneyByTime(agent)
	}
	r.RandomResource(time.Second*10, time.Second*5)
}

func EndBattle(roomId int, lose gate.Agent) {
	room, ok := GetRoom(roomId)
	if !ok {
		log.Debug("%v get room fail", lose.RemoteAddr())
		return
	}
	room.Timer.EndTimer.Reset(time.Millisecond)
	for aa, pp := range room.Players {
		pp.Base.Timer.Reset(time.Millisecond)
		aa.WriteMsg(&msg.EndBattle{
			IsWin: aa != lose,
		})
		gamedata.UsersMap[Users[aa]].InBattle = false
	}
	if r, ok := Rooms[roomId]; ok {
		r.Timer.EndTimer.Reset(time.Millisecond)
		delete(Room2Agent, roomId)
		delete(Rooms, roomId)
	}
	log.Debug("end battle with wrong roomId")
}

func (r *Room) RandomResource(beforeTime, interval time.Duration) {
	ticker1 := time.NewTicker(beforeTime)
	ticker2 := time.NewTicker(interval)

	<-ticker1.C

	for {
		select {
		case <-r.Timer.EndTimer.C:
			log.Debug("game over! random resource stop")
			return
		case <-ticker2.C:
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
		Players:     make(map[gate.Agent]*Player),
		PlayerCount: 0,
		Middle:      make(map[int]Middle),
		Lock:        sync.Mutex{},
		Timer: &RoomTimer{
			EndTimer: time.NewTimer(time.Minute * 30),
		},
	}
	return &room
}
