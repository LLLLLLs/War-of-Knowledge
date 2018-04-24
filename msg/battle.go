package msg

import (
	_ "encoding/json"
)

type CreateHero struct {
	RoomId   int    `json:"roomId"`
	HeroType string `json:"heroType"`
}

type CreateHeroInf struct {
	Msg         string   `json:"msg"`
	HeroType    string   `json:"heroType"`
	TFServer    TFServer `json:"tfServer"`
	WhichPlayer int      `json:"whichPlayer"`
	ID          int      `json:"id"`
	HPMax       float64  `json:"HPmax"`
	HP          float64  `json:"HP"`
	HPHot       float64  `json:"HPHot"`
	MPMax       float64  `json:"MPmax"`
	MP          float64  `json:"MP"`
	MPHot       float64  `json:"MPHot"`
	Speed       float64  `json:"speed"`
	Attack      float64  `json:"attack"`
	Def         float64  `json:"armor"`
}

type MoneyLeft struct {
	MoneyLeft int `json:"moneyLeft"`
}

type UpdatePosition struct {
	Id       int      `json:"id"`
	RoomId   int      `json:"roomId"`
	AnimNum  int      `json:"animNum"`
	TfServer TFServer `json:"tfServer"`
}

type MoveTo struct {
	Id       int      `json:"id"`
	TFNow    TFServer `json:"tfNow"`
	TFTarget TFServer `json:"tfServer"`
}

type UseSkill struct {
	RoomId   int
	Id       int
	SkillID  string
	TfServer TFServer
}

type UseSkillInf struct {
	Id        int      `json:"id"`
	MiddleId  int      `json:"middleId"`
	SkillID   string   `json:"skillID"`
	IsSuccess bool     `json:"isSuccess"`
	TfServer  TFServer `json:"tfServer"`
}

type UpdateHeroState struct {
	Id int     `json:"id"`
	Hp float64 `json:"hp"`
	Mp float64 `json:"mp"`
}

type UpdateBaseState struct {
	Which int     `json:"which"`
	Hp    float64 `json:"hp"`
}

type UpdateMiddleState struct {
	Id int     `json:"id"`
	Hp float64 `json:"hp"`
}

type DeleteHero struct {
	ID int `json:"id"`
}

type CreateMiddle struct {
	ID   int      `json:"id"`
	TF   TFServer `json:"tfServer"`
	Type string   `json:"type"`
}

type DeleteMiddle struct {
	ID int `json:"id"`
}

type GetResource struct {
	RoomId int `json:"roomID"`
	HeroId int `json:"heroID"`
	ItemId int `json:"itemID"`
}

type SkillCrash struct {
	RoomId     int
	FromHeroId int
	FromItemId int
	ToId       int
}

type Surrender struct {
	RoomId int `json:"roomId"`
}

type EndBattle struct {
	IsWin bool `json:"isWin"`
}

type Damage struct {
	Id     int     `json:"id"`
	Damage float64 `json:"damage"`
}

type FireBottleCrash struct {
	RoomId int
	ItemId int
}

type Upgrade struct {
	RoomId  int    `json:"roomId"`
	Id      int    `json:"id"`
	TypeOld string `json:"typeOld"`
	TypeNew string `json:"typeNew"`
}

type CreateRoom struct {
	Name string
}

type GetRoomList struct {
	PageNum int
}

type RoomInfo struct {
	Msg    string           `json:"msg"`
	RoomId int              `json:"roomId"`
	Name   string           `json:"name"`
	Users  map[string]*User `json:"users"`
}

type RoomList struct {
	RoomList []*RoomInfo `json:"roomList"`
}

type EnterRoom struct {
	RoomId int `json:"roomId"`
}

type QuitMatch struct {
}

type ExitRoom struct {
}

type HeartBeat struct {
}

type Test struct {
}

type StartBattle struct {
}
