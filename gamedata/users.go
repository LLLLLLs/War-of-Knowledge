package gamedata

//type Users struct {
//	Users []User
//}

//type User struct {
//	InBattle bool
//	RoomId   int
//	Login    bool
//	UserName string
//	UserPwd  string
//}

var UsersMap = make(map[string]*UserData)

//func loadUsers() {
//	//Us := Users{}
//	//data, err := ioutil.ReadFile("conf/users.json")
//	//if err != nil {
//	//	log.Fatal("%v", err)
//	//}
//	//err = json.Unmarshal(data, &Us)
//	//if err != nil {
//	//	log.Fatal("%v", err)
//	//}
//	//for _, user := range Us.Users {
//	//	UsersMap[user.UserName] = &User{
//	//		InBattle: false,
//	//		RoomId:   0,
//	//		Login:    false,
//	//		UserName: user.UserName,
//	//		UserPwd:  user.UserPwd,
//	//	}
//	//}
//}
//
//func init() {
//	loadUsers()
//	log.Debug("玩家人数 %d", len(UsersMap))
//}
