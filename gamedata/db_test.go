/*
Author  : Leshuo Lian
Time    : 2018\4\23 0023 
*/

package gamedata

import (
	"testing"
	"fmt"
)

func TestDB(t *testing.T) {
	//ast := assert.New(t)
	cipher := MD5("123456789546456413213")
	fmt.Println(cipher)
	userData := UserData{
		Name: "1181947970",
		//PwdHash: cipher,
		Photo: 5,
	}
	Db.Delete(userData)
	//effect, err := Db.Insert(userData)
	//ast.Nil(err)
	//ast.Equal(int64(1), effect)
}
