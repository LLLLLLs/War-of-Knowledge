package msg

import (
	_ "encoding/json"
)

// req
type Login struct {
	UserName string
	UserPwd  string
}

// resp
type LoginStat struct {
	Status   int    `json:"status"`
	Msg      string `json:"msg"`
	PlayerId int    `json:"playerId,omitempty"`
}

// req
type Match struct {
	PlayerId int
}

// resp
type MatchStat struct {
	Status      int    `json:"status"` // 0 匹配成功 ; 1 匹配中 ; 2 错误
	Msg         string `json:"msg"`
	RoomId      int    `json:"roomId"`
	WhichPlayer int    `json:"whichPlayer"` //0 左;1 右
}

type User struct {
	UserName string `json:"userName"`
}
