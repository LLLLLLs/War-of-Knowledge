package internal

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"server/msg"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
)

type Users struct {
	Users []User
}

type User struct {
	UserName string
	UserPwd  string
}

var UsersMap = make(map[string]string) //map[name] = pwd
var PlayerId = 1

func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func loadUsers() {
	Us := Users{}
	data, err := ioutil.ReadFile("conf/users.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &Us)
	if err != nil {
		log.Fatal("%v", err)
	}
	for _, user := range Us.Users {
		UsersMap[user.UserName] = user.UserPwd
	}
}

func init() {
	loadUsers()
	for k, v := range UsersMap {
		log.Debug("%v %v", k, v)
	}
	handleMsg(&msg.Login{}, handleAuth)
}

func handleAuth(args []interface{}) {
	m := args[0].(*msg.Login)
	a := args[1].(gate.Agent)
	pwd, ok := UsersMap[m.UserName]
	log.Debug("Call Login from %v", a.RemoteAddr())
	if ok && pwd == m.UserPwd {
		a.WriteMsg(&msg.LoginStat{
			Status:   0,
			Msg:      "Success",
			PlayerId: PlayerId,
		})
		PlayerId += 1
	} else {
		//TODO LoginFail
		a.WriteMsg(&msg.LoginStat{
			Status:   1,
			Msg:      "Account or password is wrong",
			PlayerId: PlayerId,
		})
	}
}
