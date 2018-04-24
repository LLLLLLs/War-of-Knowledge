/*
Author  : Leshuo Lian
Time    : 2018\4\23 0023 
*/

package gamedata

import (
	"testing"
)

func TestDB(t *testing.T) {
	//cipher := MD5("123456")
	//fmt.Println(cipher)
	//userData := UserData{
	//	Name:    "123456",
	//	PwdHash: cipher,
	//	Photo:   5,
	//}
	//userData1 := UserData{
	//	Name:    "654321",
	//	PwdHash: cipher,
	//	Photo:   5,
	//}
	////Db.Delete(userData)
	//effect, err := Db.Insert(userData)
	//Db.Insert(userData1)
	//ast.Nil(err)
	//ast.Equal(int64(1), effect)
	userData := new(UserData)
	Db.Where("name=?", "654321").Get(userData)

}
