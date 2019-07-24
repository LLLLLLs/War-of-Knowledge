package main

import (
	"flag"
	"fmt"
	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
	"github.com/name5566/leaf/log"
	"server/conf"
	"server/game"
	"server/gamedata"
	"server/gate"
	"server/login"
)

func main() {
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	lconf.LogFlag = conf.LogFlag
	lconf.ConsolePort = conf.Server.ConsolePort
	lconf.ProfilePath = conf.Server.ProfilePath
	var debug = flag.Bool("sql", true, "")
	flag.Parse()
	gamedata.SetDebug(*debug)
	fmt.Println(*debug)
	initApp()
	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)

}

func initApp() {
	gamedata.InitDB(gamedata.GetDbUri())
	log.Debug("数据清除...")
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
	log.Debug("数据清除完成...")
}
