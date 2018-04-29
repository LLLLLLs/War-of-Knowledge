package internal

import (
	"reflect"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"server/msg"
	"server/gamedata"
	"server/game"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"regexp"
)

var PlayerId = 1

func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handleMsg(&msg.Login{}, handleAuth)
	handleMsg(&msg.Register{}, handleRegister)
}

func handleAuth(args []interface{}) {
	m := args[0].(*msg.Login)
	a := args[1].(gate.Agent)
	log.Debug("玩家 %s 登录", m.UserName)
	userData := &gamedata.UserData{
		Name: m.UserName,
	}
	has, err := gamedata.Db.Get(userData)
	if err != nil || !has {
		log.Debug("账号不存在")
		a.WriteMsg(&msg.LoginStat{
			Status:   1,
			Msg:      "账号不存在",
			PlayerId: PlayerId,
		})
		return
	}
	//cipher := gamedata.MD5(m.UserPwd)
	if m.UserPwd != userData.PwdHash {
		log.Debug("账号密码不匹配")
		a.WriteMsg(&msg.LoginStat{
			Status:   1,
			Msg:      "账号密码不匹配",
			PlayerId: PlayerId,
		})
	} else if userData.Login == 1 {
		log.Debug("账号已登录")
		a.WriteMsg(&msg.LoginStat{
			Status:   1,
			Msg:      "账号已登录",
			PlayerId: PlayerId,
		})
	} else {
		PlayerId += 1
		userData.Login = 1
		gamedata.UsersMap[userData.Name] = userData
		log.Debug("玩家 %s 登录成功", userData.Name)
		a.WriteMsg(&msg.LoginStat{
			Status:   0,
			Msg:      "ok",
			PlayerId: PlayerId,
		})
		a.WriteMsg(&msg.User{
			UserName: userData.Name,
			Photo:    userData.Photo,
			Total:    userData.Total,
			Victory:  userData.Victory,
			Defeat:   userData.Defeat,
			Rate:     userData.Rate,
			KeyOwner: false,
		})
		cond := gamedata.UserData{
			Id: userData.Id,
		}
		gamedata.Db.Update(userData, cond)
		game.ChanRPC.Go("Login", a, userData.Name)
		if userData.InBattle == 1 {
			log.Debug("%s 重连中...", userData.Name)
			game.ChanRPC.Go("RecoverBattle", a)
		}
	}
}

func handleRegister(args []interface{}) {
	m := args[0].(*msg.Register)
	a := args[1].(gate.Agent)
	userData := &gamedata.UserData{
		Name: m.Name,
	}
	has, err := gamedata.Db.Get(userData)
	if err != nil {
		log.Debug("数据库错误")
		a.WriteMsg(&msg.RegisterInfo{
			Msg: "服务器繁忙",
		})
		return
	} else if has {
		log.Debug("用户名已存在")
		a.WriteMsg(&msg.RegisterInfo{
			Msg: "用户名已存在",
		})
		return
	}
	//err = validParam(m.Name, 6, 10)
	//if err != nil {
	//	log.Debug("名称%s", err.Error())
	//	a.WriteMsg(&msg.RegisterInfo{
	//		Msg: fmt.Sprintf("名称%s", err.Error()),
	//	})
	//	return
	//}
	//err = validParam(m.Password, 6, 12)
	//if err != nil {
	//	log.Debug("密码%s", err.Error())
	//	a.WriteMsg(&msg.RegisterInfo{
	//		Msg: fmt.Sprintf("密码%s", err.Error()),
	//	})
	//	return
	//}
	// 开始注册账号
	rand.Seed(time.Now().Unix())
	photo := rand.Intn(10)
	//pwdHash := gamedata.MD5(m.Password)
	userData.Name = m.Name
	userData.PwdHash = m.Password
	userData.Photo = photo
	effect, err := gamedata.Db.Insert(userData)
	if err != nil || int(effect) != 1 {
		log.Debug("数据库错误")
		a.WriteMsg(&msg.RegisterInfo{
			Msg: "服务器繁忙",
		})
		return
	}
	log.Debug("注册成功 %s", userData.Name)
	a.WriteMsg(&msg.RegisterInfo{
		Msg: "ok",
	})
}

func validParam(word string, min, max int) (error) {
	if !filter(word) {
		return errors.New("仅支持字母和数字")
	} else if !validLen(word, min, max) {
		return errors.New(fmt.Sprintf("长度仅限%d~%d位", min, max))
	}
	return nil
}

func validLen(word string, min, max int) bool {
	return len(word) >= min && len(word) <= max
}

func filter(word string) bool {
	reg := regexp.MustCompile(`[^\w]`)
	if len(reg.FindAllString(word, -1)) != 0 {
		return false
	}
	return true
}
