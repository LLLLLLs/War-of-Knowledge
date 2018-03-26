package msg

type Login struct {
	UserName string
	UserPwd  string
}

type LoginStat struct {
	Status int
	Msg    string
	ID     int
}

type Match struct {
	PlayerId int
}
