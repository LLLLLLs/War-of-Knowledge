package internal

import (
	"time"

	"server/msg"

	"sync"
	"github.com/name5566/leaf/log"
)

var (
	BornPosition1 = msg.TFServer{
		Position: []float64{37.0, 0.07, 26.0},
		Rotation: []float64{0.0, 90.0, 0.0},
	}
	BornPosition2 = msg.TFServer{
		Position: []float64{88.0, 0.07, 26.0},
		Rotation: []float64{0.0, -90.0, 0.0},
	}
)

type Player struct {
	Which int
	Base  *Base
	Heros map[int]*Hero
	Lock  sync.Mutex
}

func (p *Player) GetHeros(k int) (*Hero, bool) {
	p.Lock.Lock()
	defer p.Lock.Unlock()
	h, ok := p.Heros[k]
	return h, ok
}

func (p *Player) SetHeros(k int, v *Hero) {
	p.Lock.Lock()
	defer p.Lock.Unlock()
	p.Heros[k] = v
}

func (p *Player) DeleteHero(k int) {
	p.Lock.Lock()
	defer p.Lock.Unlock()
	delete(p.Heros, k)
}

type Base struct {
	ID     int
	Money  int
	Hp     float64
	TF     *msg.TFServer
	Radius float64
	Timer  *time.Timer
}

func (b *Base) GetMoneyByTime(user string) {
	ticker := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-b.Timer.C:
			log.Debug("game over base stop get money")
			return
		case <-ticker.C:
			b.Money += 10
			for a, u := range Users {
				if u == user {
					a.WriteMsg(&msg.MoneyLeft{
						MoneyLeft: b.Money,
					})
				}
			}
		}
	}
}

func (b *Base) SubHP(damage float64, which int, room Room) {
	b.Hp -= damage
	if b.Hp <= 0 {
		b.Hp = 0
	}
	for user, aa := range room.User2Agent {
		if aa == nil {
			continue
		}
		(*aa).WriteMsg(&msg.UpdateBaseState{
			Which: which,
			Hp:    b.Hp,
		})
		(*aa).WriteMsg(&msg.Damage{
			Id:     which,
			Damage: damage,
		})
		pp := room.Players[user]
		if b.Hp == 0 && which == pp.Which {
			EndBattle(room.RoomId, *aa)
		}
	}
}

func NewPlayer(which int) *Player {
	tf := new(msg.TFServer)
	if which == 0 {
		tf.Position = []float64{15.0, 0.0, 26.0}
		tf.Rotation = []float64{0.0, 0.0, 0.0}
	} else {
		tf.Position = []float64{110.0, 0.0, 26.0}
		tf.Rotation = []float64{0.0, 0.0, 0.0}
	}
	p := Player{
		Which: which,
		Base: &Base{ID: which,
			Money: 1000,
			Hp: 1000.0,
			TF: tf,
			Radius: 5.0,
			Timer: time.NewTimer(time.Minute * 30),
		},
		Heros: make(map[int]*Hero),
		Lock:  sync.Mutex{},
	}
	return &p
}
