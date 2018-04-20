package msg

import (
	_ "encoding/json"
)

type CreateHero struct {
	RoomId   int    `json:"roomId"`
	HeroType string `json:"heroType"`
}

type CreateHeroInf struct {
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
	Id   int     `json:"id"`
	Type string  `json:"type"`
	Hp   float64 `json:"hp"`
	Mp   float64 `json:"mp"`
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
