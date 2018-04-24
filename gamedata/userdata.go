/*
Author  : Leshuo Lian
Time    : 2018\4\23 0023 
*/

package gamedata

type UserData struct {
	Id       int    `xorm:"not null pk autoincr INT(11)"`
	Name     string `xorm:"not null VARCHAR(16)"`
	PwdHash  string `xorm:"not null VARCHAR(63)"`
	Rate     int    `xorm:"not null default 0 INT(11)"`
	Victory  int    `xorm:"not null default 0 INT(11)"`
	Defeat   int    `xorm:"not null default 0 INT(11)"`
	Total    int    `xorm:"not null default 0 INT(11)"`
	Photo    int    `xorm:"not null INT(11)"`
	Login    int    `xorm:"not null default 0 TINYINT(1)"`
	InBattle int    `xorm:"not null default 0 TINYINT(1)"`
	RoomId   int    `xorm:"not null default 0 INT(11)"`
}

func (u *UserData) TableName() string {
	return "user_data"
}
