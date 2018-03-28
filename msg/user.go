package msg

import (
	"enconding/json"
)

// req
type Login struct {
	UserName string
	UserPwd  string
}

// resp
type LoginStat struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	ID     int    `json:"id,omitempty"`
}

//req
type Match struct {
	PlayerId int
}
