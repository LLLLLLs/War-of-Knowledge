package internal

import (
	"server/msg"
	"time"
	"github.com/name5566/leaf/log"
	"math/rand"
	"math"
)

type Middle interface {
	IsInvincible() bool
	GetTF() *msg.TFServer
	GetSelfRadius() float64
	SubHp(damage float64, room Room)
	TakeAction(room *Room)
	GetAttack() (float64, bool)
}

func HasMiddle(id int, room *Room) bool {
	_, ok := room.GetMiddle(id)
	return ok
}

type MiddleCreature struct {
	ID         int
	Type       string
	TF         *msg.TFServer
	Invincible bool //判定是否无敌
}

func (m *MiddleCreature) IsInvincible() bool {
	return m.Invincible
}

func (m *MiddleCreature) GetTF() *msg.TFServer {
	return m.TF
}

func (m *MiddleCreature) GetSelfRadius() (r float64) {
	return 0.0
}

func (m *MiddleCreature) SubHp(damage float64, room Room) {}

func (m *MiddleCreature) TakeAction(room *Room) {}

func (m *MiddleCreature) GetAttack() (float64, bool) { return 0, false }

type HealFlower struct {
	MiddleCreature
	SelfRadius float64
	Radius     float64
	Duration   time.Duration
	Heal       float64
	HP         float64
}

func NewFlower(id int, tf msg.TFServer) *HealFlower {
	return &HealFlower{
		MiddleCreature{
			id,
			"001",
			&tf,
			false,
		},
		2.0,
		8.0,
		time.Second * 30,
		30.0,
		100,
	}
}

func (hf *HealFlower) GetSelfRadius() float64 {
	return hf.SelfRadius
}

func (hf *HealFlower) SubHp(damage float64, room Room) {
	hf.HP -= damage
	if hf.HP <= 0 {
		hf.HP = 0
		room.DeleteMiddle(hf.ID)
	}
	for aa := range room.Players {
		aa.WriteMsg(&msg.UpdateMiddleState{hf.ID, hf.HP})
		aa.WriteMsg(&msg.Damage{Id: hf.ID, Damage: damage})
		if hf.HP == 0 {
			aa.WriteMsg(&msg.DeleteMiddle{hf.ID})
		}
	}
}

func (hf *HealFlower) TakeAction(room *Room) {
	quit := make(chan int)
	go func(quit chan int) {
		ticker := time.NewTicker(hf.Duration)
		select {
		case <-ticker.C:
			{
				quit <- 0
				return
			}
		}
	}(quit)
	for {
		ticker := time.NewTicker(time.Second * 2)

		select {
		case <-ticker.C:
			if !HasMiddle(hf.ID, room) {
				return
			}
			for aa, pp := range room.Players {
				for _, h := range pp.Heros {
					distance := GetDistance(hf.TF, h.Transform)
					if distance < hf.Radius {
						h.Heal(hf.Heal, aa, *room, "HP")
					}
				}
			}
		case <-quit:
			if !HasMiddle(hf.ID, room) {
				return
			}
			room.DeleteMiddle(hf.ID)
			for aa := range room.Players {
				aa.WriteMsg(&msg.DeleteMiddle{
					hf.ID,
				})
			}
			log.Debug("delete middle %d", hf.ID)
			return
		}
	}
}

type BarrierTree struct {
	MiddleCreature
	SelfRadius float64
	Duration   time.Duration
	HP         float64
}

func NewBarrierTree(id int, tf msg.TFServer) *BarrierTree {
	return &BarrierTree{
		MiddleCreature{
			id,
			"002",
			&tf,
			false,
		},
		2.0,
		time.Second * 60,
		100,
	}
}

func (bt *BarrierTree) GetSelfRadius() float64 {
	return bt.SelfRadius
}

func (bt *BarrierTree) SubHp(damage float64, room Room) {
	bt.HP -= damage
	if bt.HP <= 0 {
		bt.HP = 0
		room.DeleteMiddle(bt.ID)
	}
	for aa := range room.Players {
		aa.WriteMsg(&msg.UpdateMiddleState{bt.ID, bt.HP})
		aa.WriteMsg(&msg.Damage{Id: bt.ID, Damage: damage})
		if bt.HP == 0 {
			aa.WriteMsg(&msg.DeleteMiddle{bt.ID})
		}
	}
}

func (bt *BarrierTree) TakeAction(room *Room) {
	ticker := time.NewTicker(bt.Duration)

	select {
	case <-ticker.C:
		if HasMiddle(bt.ID, room) {
			room.DeleteMiddle(bt.ID)
			for aa := range room.Players {
				aa.WriteMsg(&msg.DeleteMiddle{bt.ID})
			}
		}
		return
	}
}

type Resource struct {
	MiddleCreature
	Duration time.Duration
}

func (r *Resource) TakeAction(room *Room) {
	ticker := time.NewTicker(r.Duration)

	for {
		select {
		case <-ticker.C:
			if HasMiddle(r.ID, room) {
				room.DeleteMiddle(r.ID)
				for aa := range room.Players {
					aa.WriteMsg(&msg.DeleteMiddle{
						r.ID,
					})
				}
			}
			return
		}
	}
}

type Gold struct {
	Resource
	Value int
}

func NewGold(id int, tf msg.TFServer) *Gold {
	rand.Seed(time.Now().Unix())
	value := rand.Intn(4) + 3
	return &Gold{
		Resource{
			MiddleCreature{
				id,
				"000",
				&tf,
				true,
			},
			time.Second * 30,
		},
		value,
	}
}

type Blood struct {
	Resource
	value int
}

func NewBlood(id int, tf msg.TFServer) *Blood {
	return &Blood{
		Resource{
			MiddleCreature{
				id,
				"010",
				&tf,
				true,
			},
			time.Second * 30,
		},
		15,
	}
}

type Mana struct {
	Resource
	value int
}

func NewMana(id int, tf msg.TFServer) *Mana {
	return &Mana{
		Resource{
			MiddleCreature{
				id,
				"020",
				&tf,
				true,
			},
			time.Second * 30,
		},
		5,
	}
}

type ResourceTree struct {
	MiddleCreature
	Radius     float64
	SelfRadius float64
	Duration   time.Duration
	Interval   time.Duration
	HP         float64
}

func NewResourceTree(id int, tf msg.TFServer) *ResourceTree {
	return &ResourceTree{
		MiddleCreature{
			id,
			"003",
			&tf,
			false,
		},
		10.0,
		2.5,
		time.Second * 20,
		time.Second * 2,
		200.0,
	}
}

func (rt *ResourceTree) GetSelfRadius() float64 {
	return rt.SelfRadius
}

func (rt *ResourceTree) SubHp(damage float64, room Room) {
	rt.HP -= damage
	if rt.HP <= 0 {
		rt.HP = 0
		room.DeleteMiddle(rt.ID)
	}
	for aa := range room.Players {
		aa.WriteMsg(&msg.UpdateMiddleState{rt.ID, rt.HP})
		aa.WriteMsg(&msg.Damage{Id: rt.ID, Damage: damage})
		if rt.HP == 0 {
			aa.WriteMsg(&msg.DeleteMiddle{rt.ID})
		}
	}
}

func (rt *ResourceTree) TakeAction(room *Room) {
	quit := make(chan int, 1)
	go func(chan int) {
		ticker1 := time.NewTicker(rt.Duration)
		select {
		case <-ticker1.C:
			quit <- 1
			return
		}
	}(quit)
	for {
		ticker2 := time.NewTicker(rt.Interval)
		select {
		case <-ticker2.C:
			radius := rt.Radius
			resourceTF := getRandonTFInCircle(radius, rt.GetSelfRadius(), *rt.GetTF())
			rand.Seed(time.Now().Unix())
			randi := rand.Intn(90)
			if randi < 30 {
				resource := NewGold(room.Count+1, resourceTF)
				room.Count += 1
				log.Debug("rt create gold %d", resource.ID)
				go resource.TakeAction(room)
				room.SetMiddle(resource.ID, resource)
				for aa := range room.Players {
					aa.WriteMsg(&msg.CreateMiddle{
						resource.ID,
						*resource.TF,
						resource.Type,
					})
				}
			} else if randi < 60 {
				resource := NewBlood(room.Count+1, resourceTF)
				room.Count += 1
				log.Debug("rt create blood %d", resource.ID)
				go resource.TakeAction(room)
				room.SetMiddle(resource.ID, resource)
				for aa := range room.Players {
					aa.WriteMsg(&msg.CreateMiddle{
						resource.ID,
						*resource.TF,
						resource.Type,
					})
				}
			} else {
				resource := NewMana(room.Count+1, resourceTF)
				room.Count += 1
				log.Debug("rt create mana %d", resource.ID)
				go resource.TakeAction(room)
				room.SetMiddle(resource.ID, resource)
				for aa := range room.Players {
					aa.WriteMsg(&msg.CreateMiddle{
						resource.ID,
						*resource.TF,
						resource.Type,
					})
				}
			}

		case <-quit:
			if HasMiddle(rt.ID, room) {
				room.DeleteMiddle(rt.ID)
				for aa := range room.Players {
					aa.WriteMsg(&msg.DeleteMiddle{
						rt.ID,
					})
				}
			}
			return
		}
	}
}

func getRandonTFInCircle(radius, selfRadius float64, tf msg.TFServer) msg.TFServer {
	rand.Seed(time.Now().Unix())
	angle := float64(rand.Intn(360))
	length := float64(rand.Intn(int(radius-selfRadius))) + selfRadius
	xoffset := math.Sin(angle*math.Pi/180) * length
	zoffset := math.Cos(angle*math.Pi/180) * length
	x := tf.Position[0] + xoffset
	z := tf.Position[2] + zoffset
	return msg.TFServer{
		Position: []float64{x, 0.07, z},
		Rotation: []float64{0, 90, 0},
	}
}

type ElectricBall struct {
	MiddleCreature
	Duration time.Duration
	Attack   float64
}

func NewElectricBall(id int, tf msg.TFServer) *ElectricBall {
	return &ElectricBall{
		MiddleCreature{
			id,
			"004",
			&tf,
			true,
		},
		time.Second * 60,
		10.0,
	}
}

func (eb *ElectricBall) UpdateBallPosition(tf msg.TFServer) {
	eb.TF = &tf
}

func (eb *ElectricBall) TakeAction(room *Room) {
	ticker := time.NewTicker(eb.Duration)
	<-ticker.C
	if HasMiddle(eb.ID, room) {
		room.DeleteMiddle(eb.ID)
		for aa := range room.Players {
			aa.WriteMsg(&msg.DeleteMiddle{
				eb.ID,
			})
		}
	}
}

func (eb *ElectricBall) GetAttack() (float64, bool) {
	return eb.Attack, true
}

type StraightBall struct {
	MiddleCreature
	Attack   float64
	Duration time.Duration
}

func NewStraightBall(id int, tf msg.TFServer) *StraightBall {
	return &StraightBall{
		MiddleCreature{
			id,
			"006",
			&tf,
			true,
		},
		20.0,
		time.Second * 2,
	}
}

func (sb *StraightBall) GetAttack() (float64, bool) {
	return sb.Attack, true
}

func (sb *StraightBall) TakeAction(room *Room) {
	ticker := time.NewTicker(sb.Duration)

	select {
	case <-ticker.C:
		if _, ok := room.GetMiddle(sb.ID); !ok {
			return
		} else {
			room.DeleteMiddle(sb.ID)
			for aa := range room.Players {
				aa.WriteMsg(&msg.DeleteMiddle{sb.ID})
			}
		}
	}
}

type IceWall struct {
	MiddleCreature
	HP         float64
	SelfRadius float64
	Duration   time.Duration
}

func NewIceWall(id int, tf msg.TFServer) *IceWall {
	return &IceWall{
		MiddleCreature{
			id,
			"005",
			&tf,
			false,
		},
		100.0,
		5.0,
		time.Second * 10,
	}
}

func (iw *IceWall) SubHp(damage float64, room Room) {
	iw.HP -= damage
	if iw.HP <= 0 {
		iw.HP = 0
		room.DeleteMiddle(iw.ID)
	}
	for aa := range room.Players {
		aa.WriteMsg(&msg.UpdateMiddleState{iw.ID, iw.HP})
		aa.WriteMsg(&msg.Damage{Id: iw.ID, Damage: damage})
		if iw.HP == 0 {
			aa.WriteMsg(&msg.DeleteMiddle{iw.ID})
		}
	}
}

func (iw *IceWall) TakeAction(room *Room) {
	for {
		ticker := time.NewTicker(iw.Duration)
		select {
		case <-ticker.C:
			if HasMiddle(iw.ID, room) {
				room.DeleteMiddle(iw.ID)
				for aa := range room.Players {
					aa.WriteMsg(&msg.DeleteMiddle{
						iw.ID,
					})
				}
			}
			return
		}
	}
}

type FireBottle struct {
	MiddleCreature
	Target *msg.TFServer
	Hero   *Hero
}

func NewFireBottle(id int, target msg.TFServer, hero *Hero) *FireBottle {
	return &FireBottle{
		MiddleCreature: MiddleCreature{
			ID:         id,
			Type:       "008",
			TF:         hero.Transform,
			Invincible: true,
		},
		Target: &target,
		Hero:   hero,
	}
}

func (fb *FireBottle) UpdateFireBottle(tf msg.TFServer, room *Room) {
	fb.TF = &tf
	distace := GetDistance(fb.Target, fb.TF)
	if distace < 0.5 {
		room.DeleteMiddle(fb.ID)
		firesea := NewFireSea(room.Count+1, *fb.Target)
		for aa := range room.Players {
			aa.WriteMsg(&msg.DeleteMiddle{ID: fb.ID})
			aa.WriteMsg(&msg.CreateMiddle{
				ID:   firesea.ID,
				TF:   *firesea.TF,
				Type: firesea.Type,
			})
		}
		room.Count += 1
		room.SetMiddle(firesea.ID, firesea)
		firesea.TakeAction_(room, fb.Hero)
	}
}

type FireSea struct {
	MiddleCreature
	Duration time.Duration
	BuffTime time.Duration
	Radius   float64
	Interval time.Duration
	Attack   float64
}

func NewFireSea(id int, tf msg.TFServer) *FireSea {
	return &FireSea{
		MiddleCreature{
			id,
			"007",
			&tf,
			true,
		},
		time.Second * 8,
		time.Second * 5,
		3,
		time.Millisecond * 300,
		50,
	}
}

func (fs *FireSea) TakeAction_(room *Room, h *Hero) {
	var enemy *Player
	for _, pp := range room.Players {
		if _, ok := pp.Heros[h.ID]; !ok {
			enemy = pp
		}
	}

	quit := make(chan int, 1)
	go func(q chan int) {
		timer := time.NewTimer(fs.Duration)
		select {
		case <-timer.C:
			quit <- 1
		}
	}(quit)
	go func(quit chan int) {
		ticker := time.NewTicker(fs.Interval)
		for {
			select {
			case <-ticker.C:
				for _, hero := range enemy.Heros {
					distance := GetDistance(hero.Transform, fs.TF)
					if distance > fs.Radius {
						continue
					}
					if hero.Debuff == nil {
						hero.Debuff = &Burn{
							Timer:    *time.NewTimer(fs.BuffTime),
							Ticker:   *time.NewTicker(time.Second * 1),
							Dot:      (h.Attack + fs.Attack) / 10,
							IsEffect: false,
						}
					}
					if hero.Debuff.IsEffect == false {
						hero.Debuff.Timer.Reset(fs.BuffTime)
						hero.Debuff.IsEffect = true
						go hero.Burning(enemy, *room)
					} else {
						hero.Debuff.Timer.Reset(fs.BuffTime)
					}
				}
			case <-quit:
				room.DeleteMiddle(fs.ID)
				log.Debug("delete fire sea")
				for aa := range room.Players {
					aa.WriteMsg(&msg.DeleteMiddle{fs.ID})
				}
				return
			}
		}
	}(quit)
}
