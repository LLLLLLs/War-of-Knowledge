package gamedata

import (
	"io/ioutil"
	"encoding/json"
	"github.com/name5566/leaf/log"
)

type Users struct {
	Users []User
}

type User struct {
	Login    bool
	UserName string
	UserPwd  string
}

var UsersMap = make(map[string]*User)

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
		UsersMap[user.UserName] = &User{false, user.UserName, user.UserPwd}
	}
}

func init() {
	loadUsers()
	for k, v := range UsersMap {
		log.Debug("%v %v", k, v)
	}
}
