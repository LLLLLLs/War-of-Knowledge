package internal

import (
	"github.com/name5566/leaf/gate"
	"server/msg"
	"github.com/name5566/leaf/log"
	"math"
	"time"
)

var SkillMap = make(map[string]Skill)

type Skill interface {
	InitSkill()
	Cast(a gate.Agent, room *Room, h *Hero, tf msg.TFServer)
}

func init() {
	SkillMap["0000"] = new(Skill0000)
	SkillMap["0010"] = new(Skill0010)
	SkillMap["0011"] = new(Skill0011)
	SkillMap["0020"] = new(Skill0020)
	SkillMap["0100"] = new(Skill0100)
	SkillMap["0110"] = new(Skill0110)
	SkillMap["0111"] = new(Skill0111)
	SkillMap["0200"] = new(Skill0200)
	SkillMap["0210"] = new(Skill0210)
	SkillMap["0211"] = new(Skill0211)
}

func GetSkill(sid string) Skill {
	return SkillMap[sid]
}

type Skill0000 struct {
	Cost   float64
	Attack float64
	Radius float64
}

func (s *Skill0000) InitSkill() {
	s.Attack = 30.0
	s.Cost = 30.0
	s.Radius = 8.0
}

// 圆形AOE
func (s *Skill0000) Cast(a gate.Agent, room *Room, h *Hero, tf msg.TFServer) {
	if h.MP < s.Cost {
		a.WriteMsg(&msg.UseSkillInf{
			h.ID,
			0,
			"0000",
			false,
			tf,
		})
		log.Debug("fail")
		return
	}
	log.Debug("success")
	h.MP -= s.Cost
	for aa := range room.Players {
		aa.WriteMsg(&msg.UseSkillInf{
			h.ID,
			0,
			"0000",
			true,
			tf,
		})
		aa.WriteMsg(&msg.UpdateHeroState{
			h.ID,
			h.HP,
			h.MP,
		})
	}
	enemy := GetEnemy(a, *room)
	//英雄攻击判定
	for _, hero := range enemy.Heros {
		distance := GetDistance(h.Transform, hero.Transform)
		if distance < s.Radius {
			damage := GetDamage(h.Attack, hero.Def, s.Attack)
			hero.SubHP(damage, *room)
		}
	}
	//基地攻击判定
	base := enemy.Base
	distance := GetDistance(h.Transform, base.TF)
	if distance-base.Radius < s.Radius {
		damage := GetDamage(h.Attack, 0.0, s.Attack)
		base.SubHP(damage, enemy.Which, *room)
	}
	//中立物体判定
	for _, middle := range room.Middle {
		if middle.IsInvincible() == true {
			continue
		} else {
			distance = GetDistance(h.Transform, middle.GetTF())
			if distance-middle.GetSelfRadius() < s.Radius {
				damage := GetDamage(h.Attack, 0.0, s.Attack)
				middle.SubHp(damage, *room)
			}
		}
	}
}

type Skill0010 struct {
	Cost   float64
	Radius float64
	Attack float64
}

func (s *Skill0010) InitSkill() {
	s.Cost = 30
	s.Radius = 20
	s.Attack = 30
}

// 扇形冰
func (s *Skill0010) Cast(a gate.Agent, room *Room, h *Hero, tf msg.TFServer) {
	if h.MP < s.Cost {
		a.WriteMsg(&msg.UseSkillInf{
			h.ID,
			0,
			"0010",
			false,
			tf,
		})
		log.Debug("fail")
		return
	}
	log.Debug("success")
	h.MP -= s.Cost
	for aa := range room.Players {
		aa.WriteMsg(&msg.UseSkillInf{
			h.ID,
			0,
			"0010",
			true,
			tf,
		})
		aa.WriteMsg(&msg.UpdateHeroState{
			h.ID,
			h.HP,
			h.MP,
		})
	}
	skillAngle := getAngle(*h.Transform, tf)
	enemy := GetEnemy(a, *room)
	//英雄攻击判定
	for _, hero := range enemy.Heros {
		distance := GetDistance(h.Transform, hero.Transform)
		if distance > s.Radius {
			continue
		}
		heroAngle := getAngle(*h.Transform, *hero.Transform)
		if math.Abs(skillAngle-heroAngle) < math.Pi/6 {
			damage := GetDamage(h.Attack, hero.Def, s.Attack)
			hero.SubHP(damage, *room)
		}
	}
	//基地攻击判定
	//base := enemy.Base
	//distance := GetDistance(h.Transform, base.TF)
	//if distance-base.Radius < s.Radius {
	//	damage := GetDamage(h.Attack, 0.0, s.Attack)
	//	base.SubHP(damage, enemy.Which, *room)
	//}
	//中立物体判定
	for _, middle := range room.Middle {
		if middle.IsInvincible() == true {
			continue
		} else {
			distance := GetDistance(h.Transform, middle.GetTF())
			if distance-middle.GetSelfRadius() < s.Radius {
				middleAngle := getAngle(*h.Transform, *middle.GetTF())
				if math.Abs(skillAngle-middleAngle) < math.Pi/6 {
					damage := GetDamage(h.Attack, 0.0, s.Attack)
					middle.SubHp(damage, *room)
				}
			}
		}
	}
}

func getAngle(tf1, tf2 msg.TFServer) float64 {
	p1, p2 := tf1.Position, tf2.Position
	a := p2[2] - p1[2]
	b := p2[0] - p1[0]
	c := GetDistance(&tf1, &tf2)
	sinA := a / c
	A := math.Asin(sinA)
	if sinA > 0 {
		if b > 0 {
			return A
		} else {
			return math.Pi - A
		}
	} else {
		if b > 0 {
			return A
		} else {
			return -math.Pi - A
		}
	}

}

type Skill0011 struct {
	Cost         float64
	CastDistance float64
}

func (s *Skill0011) InitSkill() {
	s.Cost = 30.0
	s.CastDistance = 15.0
}

// 一圈地形冰柱
func (s *Skill0011) Cast(a gate.Agent, room *Room, h *Hero, tf msg.TFServer) {
	distance := GetDistance(h.Transform, &tf)
	if h.MP < s.Cost || distance > s.CastDistance {
		a.WriteMsg(&msg.UseSkillInf{
			h.ID,
			0,
			"0011",
			false,
			tf,
		})
		log.Debug("fail")
		return
	}
	log.Debug("success")
	h.MP -= s.Cost
	iw := NewIceWall(room.Count+1, tf)
	room.Count += 1
	room.SetMiddle(iw.ID, iw)
	go iw.TakeAction(room)
	for aa := range room.Players {
		aa.WriteMsg(&msg.UpdateHeroState{
			h.ID,
			h.HP,
			h.MP,
		})
		aa.WriteMsg(&msg.UseSkillInf{
			h.ID,
			iw.ID,
			"0011",
			true,
			tf,
		})
		aa.WriteMsg(&msg.CreateMiddle{
			iw.ID,
			*iw.TF,
			iw.Type,
		})
	}
}

type Skill0020 struct {
	Cost         float64
	CastDistance float64
}

func (s *Skill0020) InitSkill() {
	s.Cost = 20
	s.CastDistance = 15
}

// 火焰瓶
func (s *Skill0020) Cast(a gate.Agent, room *Room, h *Hero, tf msg.TFServer) {
	distance := GetDistance(h.Transform, &tf)
	if h.MP < s.Cost || distance > s.CastDistance {
		a.WriteMsg(&msg.UseSkillInf{
			h.ID,
			0,
			"0020",
			false,
			tf,
		})
		log.Debug("fail")
		return
	}
	log.Debug("success")
	h.MP -= s.Cost
	fb := NewFireBottle(room.Count+1, tf, h)
	room.Count += 1
	room.SetMiddle(fb.ID, fb)
	for aa := range room.Players {
		aa.WriteMsg(&msg.UpdateHeroState{
			h.ID,
			h.HP,
			h.MP,
		})
		aa.WriteMsg(&msg.UseSkillInf{
			h.ID,
			fb.ID,
			"0020",
			true,
			tf,
		})
		aa.WriteMsg(&msg.CreateMiddle{
			fb.ID,
			*fb.TF,
			fb.Type,
		})
	}
}

type Skill0100 struct {
	Cost float64
}

func (s *Skill0100) InitSkill() {
	s.Cost = 20.0
}

// 直射电球
func (s *Skill0100) Cast(a gate.Agent, room *Room, h *Hero, tf msg.TFServer) {
	if h.MP < s.Cost {
		a.WriteMsg(&msg.UseSkillInf{
			h.ID,
			0,
			"0100",
			false,
			tf,
		})
		log.Debug("cast skill 0100 fail")
		return
	}
	log.Debug("success")
	h.MP -= s.Cost
	sb := NewStraightBall(room.Count+1, *h.Transform)
	room.Count += 1
	room.SetMiddle(sb.ID, sb)
	go sb.TakeAction(room)
	for aa := range room.Players {
		aa.WriteMsg(&msg.UpdateHeroState{
			h.ID,
			h.HP,
			h.MP,
		})
		aa.WriteMsg(&msg.UseSkillInf{
			h.ID,
			sb.ID,
			"0100",
			true,
			tf,
		})
		aa.WriteMsg(&msg.CreateMiddle{
			sb.ID,
			*sb.TF,
			sb.Type,
		})
	}
}

type Skill0110 struct {
	Cost   float64
	Radius float64
}

func (s *Skill0110) InitSkill() {
	s.Cost = 30.0
	s.Radius = 15.0
}

// 周身4电球
func (s *Skill0110) Cast(a gate.Agent, room *Room, h *Hero, tf msg.TFServer) {
	if h.MP < s.Cost {
		a.WriteMsg(&msg.UseSkillInf{
			h.ID,
			0,
			"0110",
			false,
			tf,
		})
		log.Debug("cast skill 0110 fail")
		return
	}
	log.Debug("success")
	h.MP -= s.Cost
	balls := s.generateElectricBall(room, h, tf, 4)
	for aa := range room.Players {
		aa.WriteMsg(&msg.UpdateHeroState{
			h.ID,
			h.HP,
			h.MP,
		})
		aa.WriteMsg(&msg.UseSkillInf{
			h.ID,
			balls[0].ID,
			"0110",
			true,
			tf,
		})
		for _, ball := range balls {
			aa.WriteMsg(&msg.CreateMiddle{
				ball.ID,
				*ball.TF,
				ball.Type,
			})
		}
	}
}

func (s *Skill0110) generateElectricBall(room *Room, h *Hero, tf msg.TFServer, n int) []*ElectricBall {
	balls := make([]*ElectricBall, n)
	for i := range balls {
		balls[i] = NewElectricBall(room.Count+1, tf)
		room.Count += 1
		room.SetMiddle(balls[i].ID, balls[i])
		go balls[i].TakeAction(room)
	}
	return balls
}

type Skill0111 struct {
	Cost   float64
	Attack float64
}

func (s *Skill0111) InitSkill() {
	s.Cost = 50.0
	s.Attack = 50.0
}

// 对敌方所有物体造成伤害
func (s *Skill0111) Cast(a gate.Agent, room *Room, h *Hero, tf msg.TFServer) {
	if h.MP < s.Cost {
		a.WriteMsg(&msg.UseSkillInf{
			h.ID,
			0,
			"0111",
			false,
			tf,
		})
		log.Debug("技能0111释放失败")
		return
	}
	log.Debug("技能0111释放成功")
	h.MP -= s.Cost
	for aa := range room.Players {
		aa.WriteMsg(&msg.UseSkillInf{
			h.ID,
			0,
			"0111",
			true,
			tf,
		})
		aa.WriteMsg(&msg.UpdateHeroState{
			h.ID,
			h.HP,
			h.MP,
		})
	}
	go func(room *Room) {
		enemy := GetEnemy(a, *room)
		timer := time.NewTimer(time.Second * 4)
		<-timer.C
		for _, hero := range enemy.Heros {
			damage := GetDamage(h.Attack, hero.Def, s.Attack)
			for _, player := range room.Players {
				if _, ok := player.GetHeros(h.ID); ok {
					hero.SubHP(damage, *room)
				}
			}
		}
	}(room)
}

type Skill0200 struct {
	Cost         float64
	CastDistance float64
}

func (s *Skill0200) InitSkill() {
	s.Cost = 30.0
	s.CastDistance = 15.0
}

// 种障碍树
func (s *Skill0200) Cast(a gate.Agent, room *Room, h *Hero, tf msg.TFServer) {
	distance := GetDistance(h.Transform, &tf)
	if h.MP < s.Cost || distance > s.CastDistance {
		a.WriteMsg(&msg.UseSkillInf{
			h.ID,
			0,
			"0200",
			false,
			tf,
		})
		log.Debug("fail")
		return
	}
	log.Debug("success")
	h.MP -= s.Cost
	room.Count += 1
	bt := NewBarrierTree(room.Count, tf)
	room.SetMiddle(bt.ID, bt)
	for aa := range room.Players {
		aa.WriteMsg(&msg.UpdateHeroState{
			h.ID,
			h.HP,
			h.MP,
		})
		aa.WriteMsg(&msg.UseSkillInf{
			h.ID,
			bt.ID,
			"0200",
			true,
			tf,
		})
		aa.WriteMsg(&msg.CreateMiddle{
			bt.ID,
			*bt.TF,
			bt.Type,
		})
	}
}

type Skill0210 struct {
	Cost         float64
	CastDistance float64
}

func (s *Skill0210) InitSkill() {
	s.Cost = 30.0
	s.CastDistance = 20.0
}

// 种子周围回血
func (s *Skill0210) Cast(a gate.Agent, room *Room, h *Hero, tf msg.TFServer) {
	distance := GetDistance(h.Transform, &tf)
	if h.MP < s.Cost || distance > s.CastDistance {
		a.WriteMsg(&msg.UseSkillInf{
			h.ID,
			0,
			"0210",
			false,
			tf,
		})
		log.Debug("fail")
		return
	}
	log.Debug("success")
	h.MP -= s.Cost
	room.Count += 1
	hf := NewFlower(room.Count, tf)
	room.SetMiddle(hf.ID, hf)
	for aa := range room.Players {
		aa.WriteMsg(&msg.UpdateHeroState{
			h.ID,
			h.HP,
			h.MP,
		})
		aa.WriteMsg(&msg.UseSkillInf{
			h.ID,
			hf.ID,
			"0210",
			true,
			tf,
		})
		aa.WriteMsg(&msg.CreateMiddle{
			hf.ID,
			*hf.TF,
			hf.Type,
		})
	}
	go hf.TakeAction(room)
}

type Skill0211 struct {
	Cost         float64
	CastDistance float64
}

func (s *Skill0211) InitSkill() {
	s.Cost = 30.0
	s.CastDistance = 20.0
}

// 资源树
func (s *Skill0211) Cast(a gate.Agent, room *Room, h *Hero, tf msg.TFServer) {
	distance := GetDistance(h.Transform, &tf)
	if h.MP < s.Cost || distance > s.CastDistance {
		a.WriteMsg(&msg.UseSkillInf{
			h.ID,
			0,
			"0211",
			false,
			tf,
		})
		log.Debug("fail")
		return
	}
	log.Debug("success")
	h.MP -= s.Cost
	room.Count += 1
	rt := NewResourceTree(room.Count, tf)
	go rt.TakeAction(room)
	room.SetMiddle(rt.ID, rt)
	for aa := range room.Players {
		aa.WriteMsg(&msg.UpdateHeroState{
			h.ID,
			h.HP,
			h.MP,
		})
		aa.WriteMsg(&msg.UseSkillInf{
			h.ID,
			rt.ID,
			"0211",
			true,
			tf,
		})
		aa.WriteMsg(&msg.CreateMiddle{
			rt.ID,
			*rt.TF,
			rt.Type,
		})
	}
}
