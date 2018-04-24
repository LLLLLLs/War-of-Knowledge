package internal

import (
	"github.com/name5566/leaf/module"
	"server/base"
	"server/gamedata"
	"github.com/name5566/leaf/log"
	"fmt"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
}

func (m *Module) OnDestroy() {
	fmt.Println("111")
	db := gamedata.Db
	users := make([]gamedata.UserData, 0)
	db.Where("login=? or in_battle=?", 1, 1).Find(&users)
	for _, user := range users {
		user.Login = 0
		user.InBattle = 0
		condition := gamedata.UserData{
			Id: user.Id,
		}
		effect, err := db.Cols("login", "in_battle").Update(user, condition)
		if err != nil {
			log.Debug("获取数据库失败")
			return
		}
		if int(effect) != 1 {
			log.Debug("数据库更新失败")
		}
	}
	fmt.Println("222")
}
