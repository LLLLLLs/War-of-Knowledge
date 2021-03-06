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
	GetId() int
	GetType() string
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
func (m *MiddleCreature) GetId() int {
	return m.ID
}
func (m *MiddleCreature) GetType() string {
	return m.Type
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

func ForceMove(m Middle, room *Room) {
	timer := time.NewTimer(time.Second)
	<-timer.C
	room.Lock.Lock()
	for _, pp := range room.Players {
		for _, h := range pp.Heros {
			distance := GetDistance(m.GetTF(), h.Transform)
			if distance > m.GetSelfRadius()+1.5 {
				continue
			}
			angle := getAngle(*m.GetTF(), *h.Transform)
			z := (m.GetSelfRadius() + 1.5) * math.Sin(angle)
			x := (m.GetSelfRadius() + 1.5) * math.Cos(angle)
			h.Transform.Position[0] = m.GetTF().Position[0] + x
			h.Transform.Position[2] = m.GetTF().Position[2] + z
			for _, aa := range room.User2Agent {
				if aa != nil {
					(*aa).WriteMsg(&msg.UpdatePosition{
						Id:       h.ID,
						RoomId:   room.RoomId,
						TfServer: *h.Transform,
					})
				}
			}
		}
	}
	room.Lock.Unlock()
}

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
		3.0,
		7.0,
		time.Second * 15,
		20.0,
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
	for _, aa := range room.User2Agent {
		if aa == nil {
			continue
		}
		(*aa).WriteMsg(&msg.UpdateMiddleState{hf.ID, hf.HP})
		(*aa).WriteMsg(&msg.Damage{Id: hf.ID, Damage: damage})
		if hf.HP == 0 {
			(*aa).WriteMsg(&msg.DeleteMiddle{hf.ID})
		}
	}
}

func (hf *HealFlower) TakeAction(room *Room) {
	hf.Invincible = true
	timer := time.NewTimer(time.Second)
	<-timer.C
	hf.Invincible = false
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
			for _, pp := range room.Players {
				for _, h := range pp.Heros {
					distance := GetDistance(hf.TF, h.Transform)
					if distance < hf.Radius {
						h.Heal(hf.Heal, *room, "HP")
					}
				}
			}
		case <-quit:
			if room.Closed || !HasMiddle(hf.ID, room) {
				return
			}
			room.DeleteMiddle(hf.ID)
			for _, aa := range room.User2Agent {
				if aa == nil {
					continue
				}
				(*aa).WriteMsg(&msg.DeleteMiddle{
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
		3.0,
		time.Second * 15,
		200,
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
	for _, aa := range room.User2Agent {
		if aa == nil {
			continue
		}
		(*aa).WriteMsg(&msg.UpdateMiddleState{bt.ID, bt.HP})
		(*aa).WriteMsg(&msg.Damage{Id: bt.ID, Damage: damage})
		if bt.HP == 0 {
			(*aa).WriteMsg(&msg.DeleteMiddle{bt.ID})
		}
	}
}

func (bt *BarrierTree) TakeAction(room *Room) {
	bt.Invincible = true
	timer := time.NewTimer(time.Second)
	<-timer.C
	bt.Invincible = false
	timer1 := time.NewTimer(bt.Duration)
	ticker := time.NewTicker(time.Second * 2)
	go func() {
		for {
			select {
			case <-ticker.C:
				if room.Closed {
					return
				}
			case <-timer1.C:
				if !room.Closed || HasMiddle(bt.ID, room) {
					room.DeleteMiddle(bt.ID)
					for _, aa := range room.User2Agent {
						if aa == nil {
							continue
						}
						(*aa).WriteMsg(&msg.DeleteMiddle{bt.ID})
					}
				}
				return
			}
		}
	}()

}

type Resource struct {
	MiddleCreature
	Duration time.Duration
}

func (r *Resource) TakeAction(room *Room) {
	ticker := time.NewTicker(r.Duration)
	ticker2 := time.NewTicker(time.Second * 3)

	for {
		select {
		case <-ticker.C:
			if !room.Closed || HasMiddle(r.ID, room) {
				log.Debug("资源 %d 销毁", r.ID)
				room.DeleteMiddle(r.ID)
				for _, aa := range room.User2Agent {
					if aa == nil {
						continue
					}
					(*aa).WriteMsg(&msg.DeleteMiddle{r.ID,})
				}
			}
			return
		case <-ticker2.C:
			if room.Closed || !HasMiddle(r.ID, room) {
				log.Debug("资源 %d 删除", r.ID)
				return
			}
		}
	}
}

type Gold struct {
	Resource
	Value int
}

func NewGold(id int, tf msg.TFServer) *Gold {
	rand.Seed(time.Now().Unix())
	value := rand.Intn(20) + 30
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
		50,
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
		20,
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
		3.5,
		time.Second * 30,
		time.Second * 3,
		100.0,
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
	for _, aa := range room.User2Agent {
		if aa == nil {
			continue
		}
		(*aa).WriteMsg(&msg.UpdateMiddleState{rt.ID, rt.HP})
		(*aa).WriteMsg(&msg.Damage{Id: rt.ID, Damage: damage})
		if rt.HP == 0 {
			(*aa).WriteMsg(&msg.DeleteMiddle{rt.ID})
		}
	}
}

func (rt *ResourceTree) TakeAction(room *Room) {
	rt.Invincible = true
	timer := time.NewTimer(time.Second)
	<-timer.C
	rt.Invincible = false
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
			if room.Closed {
				log.Debug("房间%d关闭-资源树销毁", room.RoomId)
				return
			}
			radius := rt.Radius
			resourceTF := getRandomTFInCircle(radius, rt.GetSelfRadius(), *rt.GetTF())
			rand.Seed(time.Now().Unix())
			randi := rand.Intn(90)
			if randi < 30 {
				resource := NewGold(room.Count+1, resourceTF)
				room.Count += 1
				log.Debug("资源树生成金币 %d", resource.ID)
				go resource.TakeAction(room)
				room.SetMiddle(resource.ID, resource)
				for _, aa := range room.User2Agent {
					if aa == nil {
						continue
					}
					(*aa).WriteMsg(&msg.CreateMiddle{
						resource.ID,
						*resource.TF,
						resource.Type,
					})
				}
			} else if randi < 60 {
				resource := NewBlood(room.Count+1, resourceTF)
				room.Count += 1
				log.Debug("资源树生成血包 %d", resource.ID)
				go resource.TakeAction(room)
				room.SetMiddle(resource.ID, resource)
				for _, aa := range room.User2Agent {
					if aa == nil {
						continue
					}
					(*aa).WriteMsg(&msg.CreateMiddle{
						resource.ID,
						*resource.TF,
						resource.Type,
					})
				}
			} else {
				resource := NewMana(room.Count+1, resourceTF)
				room.Count += 1
				log.Debug("资源树生成蓝包 %d", resource.ID)
				go resource.TakeAction(room)
				room.SetMiddle(resource.ID, resource)
				for _, aa := range room.User2Agent {
					if aa == nil {
						continue
					}
					(*aa).WriteMsg(&msg.CreateMiddle{
						resource.ID,
						*resource.TF,
						resource.Type,
					})
				}
			}

		case <-quit:
			if !room.Closed || HasMiddle(rt.ID, room) {
				room.DeleteMiddle(rt.ID)
				for _, aa := range room.User2Agent {
					if aa == nil {
						continue
					}
					(*aa).WriteMsg(&msg.DeleteMiddle{
						rt.ID,
					})
				}
			}
			return
		}
	}
}

func getRandomTFInCircle(radius, selfRadius float64, tf msg.TFServer) msg.TFServer {
	time.Sleep(time.Nanosecond * 10)
	rand.Seed(time.Now().UnixNano())
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
		time.Second * 10,
		-10.0,
	}
}

func (eb *ElectricBall) UpdateBallPosition(tf msg.TFServer) {
	eb.TF = &tf
}

func (eb *ElectricBall) TakeAction(room *Room) {
	ticker := time.NewTicker(eb.Duration)
	<-ticker.C
	if !room.Closed || HasMiddle(eb.ID, room) {
		room.DeleteMiddle(eb.ID)
		for _, aa := range room.User2Agent {
			if aa == nil {
				continue
			}
			(*aa).WriteMsg(&msg.DeleteMiddle{
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
		60.0,
		time.Second * 3,
	}
}

func (sb *StraightBall) GetAttack() (float64, bool) {
	return sb.Attack, true
}

func (sb *StraightBall) TakeAction(room *Room) {
	ticker := time.NewTicker(sb.Duration)

	select {
	case <-ticker.C:
		if !room.Closed || HasMiddle(sb.ID, room) {
			room.DeleteMiddle(sb.ID)
			for _, aa := range room.User2Agent {
				if aa == nil {
					continue
				}
				(*aa).WriteMsg(&msg.DeleteMiddle{sb.ID})
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
		200.0,
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
	for _, aa := range room.User2Agent {
		if aa == nil {
			continue
		}
		(*aa).WriteMsg(&msg.UpdateMiddleState{iw.ID, iw.HP})
		(*aa).WriteMsg(&msg.Damage{Id: iw.ID, Damage: damage})
		if iw.HP == 0 {
			(*aa).WriteMsg(&msg.DeleteMiddle{iw.ID})
		}
	}
}

func (iw *IceWall) TakeAction(room *Room) {
	for {
		ticker := time.NewTicker(iw.Duration)
		select {
		case <-ticker.C:
			if !room.Closed || HasMiddle(iw.ID, room) {
				room.DeleteMiddle(iw.ID)
				for _, aa := range room.User2Agent {
					if aa == nil {
						continue
					}
					(*aa).WriteMsg(&msg.DeleteMiddle{
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
		for _, aa := range room.User2Agent {
			if aa == nil {
				continue
			}
			(*aa).WriteMsg(&msg.DeleteMiddle{ID: fb.ID})
			(*aa).WriteMsg(&msg.CreateMiddle{
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
		time.Second * 10,
		time.Second * 5,
		5,
		time.Millisecond * 300,
		60,
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
				if !room.Closed || HasMiddle(fs.ID, room) {
					room.DeleteMiddle(fs.ID)
					log.Debug("delete fire sea")
					for _, aa := range room.User2Agent {
						if aa == nil {
							continue
						}
						(*aa).WriteMsg(&msg.DeleteMiddle{fs.ID})
					}
				}
				return
			}
		}
	}(quit)
}

type Stone struct {
	MiddleCreature
	Duration   time.Duration
	HP         float64
	SelfRadius float64
}

func NewStone(id int, tf msg.TFServer) *Stone {
	return &Stone{
		MiddleCreature: MiddleCreature{
			ID:         id,
			Type:       "100",
			TF:         &tf,
			Invincible: false,
		},
		Duration:   time.Second * 15,
		HP:         200,
		SelfRadius: 3.5,
	}
}

func (s *Stone) GetSelfRadius() float64 {
	return s.SelfRadius
}

func (s *Stone) TakeAction(room *Room) {
	s.Invincible = true
	timer := time.NewTimer(time.Second)
	<-timer.C
	s.Invincible = false
	timer1 := time.NewTimer(s.Duration)
	ticker := time.NewTicker(time.Second * 2)
	go func() {
		for {
			select {
			case <-ticker.C:
				if room.Closed {
					return
				}
			case <-timer1.C:
				if !room.Closed || HasMiddle(s.ID, room) {
					room.DeleteMiddle(s.ID)
					for _, aa := range room.User2Agent {
						if aa == nil {
							continue
						}
						(*aa).WriteMsg(&msg.DeleteMiddle{s.ID})
					}
				}
				return
			}
		}
	}()
}

func (s *Stone) SubHp(damage float64, room Room) {
	s.HP -= damage
	if s.HP <= 0 {
		s.HP = 0
		room.DeleteMiddle(s.ID)
	}
	for _, aa := range room.User2Agent {
		if aa == nil {
			continue
		}
		(*aa).WriteMsg(&msg.UpdateMiddleState{s.ID, s.HP})
		(*aa).WriteMsg(&msg.Damage{Id: s.ID, Damage: damage})
		if s.HP == 0 {
			(*aa).WriteMsg(&msg.DeleteMiddle{s.ID})
		}
	}
}

type Chest struct {
	MiddleCreature
	Duration   time.Duration
	HP         float64
	SelfRadius float64
}

func NewChest(id int, tf msg.TFServer) *Chest {
	return &Chest{
		MiddleCreature: MiddleCreature{
			ID:         id,
			Type:       "200",
			TF:         &tf,
			Invincible: false,
		},
		Duration:   time.Second * 30,
		HP:         200,
		SelfRadius: 3.5,
	}
}

func (c *Chest) GetSelfRadius() float64 {
	return c.SelfRadius
}

func (c *Chest) TakeAction(room *Room) {
	c.Invincible = true
	timer := time.NewTimer(time.Second)
	<-timer.C
	c.Invincible = false
	timer1 := time.NewTimer(c.Duration)
	ticker := time.NewTicker(time.Second * 2)
	go func() {
		for {
			select {
			case <-ticker.C:
				if room.Closed {
					return
				}
			case <-timer1.C:
				if !room.Closed || HasMiddle(c.ID, room) {
					room.DeleteMiddle(c.ID)
					log.Debug("宝箱时间到")
					for _, aa := range room.User2Agent {
						if aa == nil {
							continue
						}
						(*aa).WriteMsg(&msg.DeleteMiddle{c.ID})
					}
				}
				return
			}
		}
	}()
}

func (c *Chest) SubHp(damage float64, room Room) {
	c.HP -= damage
	if c.HP <= 0 {
		c.HP = 0
	}
	for _, aa := range room.User2Agent {
		if aa == nil {
			continue
		}
		(*aa).WriteMsg(&msg.UpdateMiddleState{c.ID, c.HP})
		(*aa).WriteMsg(&msg.Damage{Id: c.ID, Damage: damage})
		if c.HP == 0 {
			(*aa).WriteMsg(&msg.DeleteMiddle{c.ID})
			room.DeleteMiddle(c.ID)
		}
	}
	if c.HP == 0 {
		// 生成奖励:5枚金币
		for i := 0; i < 5; i++ {
			tf := getRandomTFInCircle(5, c.GetSelfRadius(), *c.GetTF())
			gold := NewGold(room.Count, tf)
			room.Count += 1
			room.SetMiddle(gold.ID, gold)
			go gold.TakeAction(&room)
			for _, aa := range room.User2Agent {
				if aa != nil {
					(*aa).WriteMsg(&msg.CreateMiddle{
						ID:   gold.ID,
						TF:   *gold.TF,
						Type: gold.Type,
					})
				}
			}
		}
	}
}

type Wind struct {
	MiddleCreature
	Speed      int
	SelfRadius float64
	Damage     float64
}

func NewWind(id int, tf msg.TFServer) *Wind {
	return &Wind{
		MiddleCreature: MiddleCreature{
			ID:         id,
			TF:         &tf,
			Type:       "300",
			Invincible: true,
		},
		Speed:      5,
		SelfRadius: 2,
		Damage:     5,
	}
}

func (w *Wind) TakeAction(room *Room) {

}

func (w *Wind) MoveTo(room *Room, target msg.TFServer) {
	ticker := time.NewTicker(time.Millisecond * 100)
	for {
		select {
		case <-ticker.C:
			if room.Closed {
				return
			}
			angle := getAngle(*w.TF, target)
			distance := float64(w.Speed) / 10.0
			z := distance * math.Sin(angle)
			x := distance * math.Cos(angle)
			w.TF.Position[0] += x
			w.TF.Position[2] += z
			// 判定英雄是否在范围内
			for _, pp := range room.Players {
				for _, h := range pp.Heros {
					d := GetDistance(w.TF, h.Transform)
					if d < (w.SelfRadius + 2) {
						h.SubHP(w.Damage, *room)
						for _, aa := range room.User2Agent {
							if aa != nil {
								(*aa).WriteMsg(&msg.Damage{
									Id:     h.ID,
									Damage: w.Damage,
								})
							}
						}
					}
				}
			}
			// 判定龙卷风是否走到尽头
			remain := GetDistance(w.TF, &target)
			//log.Debug("距离尽头还有 %f", remain)
			if remain < 2 {
				log.Debug("龙卷风已到尽头")
				room.DeleteMiddle(w.ID)
				for _, aa := range room.User2Agent {
					if aa != nil {
						(*aa).WriteMsg(&msg.DeleteMiddle{
							ID: w.ID,
						})
					}
				}
				return
			}
		}
	}
}
