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
	go r.RandomResource(time.Second*10, time.Second*5)
	go r.MapEvent()
}

func (r *Room) MapEvent() {
	eventTime := time.NewTicker(time.Second * 20)
	checkRoom := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-checkRoom.C:
			if r.Closed {
				return
			}
		case <-eventTime.C:
			go RandomMapEvent(r)
		}
	}
}

func RandomMapEvent(room *Room) {
	rand.Seed(time.Now().Unix())
	randn := rand.Intn(3)
	randn = 1

	switch randn {
	case 0:
		log.Debug("房间 %d 生成石阵", room.RoomId)
		go StoneEvent(room)
	case 1:
		log.Debug("房间 %d 生成宝箱", room.RoomId)
		go BonusEvent(room)
	case 2:
		log.Debug("房间 %d 生成龙卷风", room.RoomId)
		go WindEvent(room)
	default:
		log.Debug("房间 %d 错误事件id", room.RoomId)
	}
}

func StoneEvent(room *Room) {
	ids := []int{}
	stones := []*Stone{}
	room.Count += 5
	for i := 0; i < 5; i++ {
		ids = append(ids, room.Count-i)
	}
	for i := 0; i < 5; i++ {
		tf := msg.TFServer{
			Position: []float64{64.5, 0, float64(36 - 5*i)},
		}
		stones[i] = NewStone(ids[4-i], tf)
		room.SetMiddle(stones[i].ID, stones[i])
		go stones[i].TakeAction(room)
		go ForceMove(stones[i], room)
	}
	for _, aa := range room.User2Agent {
		if aa != nil {
			(*aa).WriteMsg(&msg.MapEvent{
				Msg:  "生成石阵",
				Type: 0,
				TFServer: msg.TFServer{
					Position: []float64{64.5, 0, 36},
				},
				ID: stones[0].ID,
			})
		}
	}
}

func BonusEvent(room *Room) {
	tf := getRandTF()
	chest := NewChest(room.Count+1, tf)
	room.Count += 1
	room.SetMiddle(chest.ID, chest)
	go chest.TakeAction(room)
	go ForceMove(chest, room)
	for _, aa := range room.User2Agent {
		if aa != nil {
			(*aa).WriteMsg(&msg.MapEvent{
				Msg:      "生成宝箱",
				Type:     1,
				TFServer: tf,
				ID:       chest.ID,
			})
		}
	}
}

func WindEvent(room *Room) {
	windPosition := [6][]float64{
		{32, 0, 41}, // 左上
		{97, 0, 11}, // 右下
		{97, 0, 41}, // 右上
		{31, 0, 11}, // 左下
		{42, 0, 26}, // 左
		{87, 0, 26}, // 右
	}
	rand.Seed(time.Now().Unix())
	num := rand.Intn(2)
	var (
		wind1 *Wind
		wind2 *Wind
		tf1   msg.TFServer
		tg1   msg.TFServer
		tf2   msg.TFServer
		tg2   msg.TFServer
	)
	if num == 0 {
		tf1 = msg.TFServer{
			Position: windPosition[0],
		}
		tf2 = msg.TFServer{
			Position: windPosition[2],
		}
		tg1 = msg.TFServer{
			Position: windPosition[1],
		}
		tg2 = msg.TFServer{
			Position: windPosition[3],
		}
		wind1 = NewWind(room.Count+1, tf1)
		room.Count += 1
		wind2 = NewWind(room.Count+1, tf2)
		room.Count += 1
	} else {
		tf1 = msg.TFServer{
			Position: []float64{42, 0, 26},
		}
		tg1 = msg.TFServer{
			Position: []float64{87, 0, 26},
		}
		tf2 = msg.TFServer{
			Position: []float64{87, 0, 26},
		}
		tg2 = msg.TFServer{
			Position: []float64{42, 0, 26},
		}
		wind1 = NewWind(room.Count+1, tf1)
		room.Count += 1
		wind2 = NewWind(room.Count+1, tf2)
		room.Count += 1
	}
	room.SetMiddle(wind1.ID, wind1)
	room.SetMiddle(wind2.ID, wind2)
	go wind1.MoveTo(room, tg1)
	go wind2.MoveTo(room, tg2)
	for _, aa := range room.User2Agent {
		if aa != nil {
			(*aa).WriteMsg(&msg.MapEvent{
				Msg:  "生成龙卷风",
				Type: 2,
				Num:  num,
				ID:   wind1.ID,
			})
		}
	}
}

func (r *Room) SyncItems() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			itemList := []int{0, 1}
			if r.Closed {
				return
			}
			r.Lock.Lock()
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
			r.Lock.Unlock()
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
	room.Lock.Lock()
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
	room.Lock.Unlock()
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
	room.PlayerCount += 1
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
	log.Debug("%s 同步数据成功", userName)
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
			tf := getRandTF()
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

func getRandTF() msg.TFServer {
	rand.Seed(time.Now().Unix())
	x := float64(rand.Intn(6700))/100.0 + 29.00
	z := float64(rand.Intn(3300))/100.0 + 9.00
	tf := msg.TFServer{
		Position: []float64{x, 0.07, z},
		Rotation: []float64{0.0, 90.0, 0.0},
	}
	return tf
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
	userData := gamedata.UsersMap[userName]
	room.Players[userName] = nil
	Agent2Room[a] = room.RoomId
	room.User2Agent[userName ] = &a
	user := msg.User{
		UserName: userName,
		Photo:    userData.Photo,
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
		for index, roomInfo := range RoomList {
			if roomInfo.RoomId == room.RoomId {
				RoomList = append(RoomList[:index], RoomList[index+1:]...)
			}
		}
		delete(Rooms, room.RoomId)
		delete(Agent2Room, a)
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
