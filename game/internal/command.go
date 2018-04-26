/*
Author  : Leshuo Lian
Time    : 2018\4\24 0024 
*/

package internal

import (
	"fmt"
	"github.com/name5566/leaf/log"
	"strconv"
)

func init() {
	skeleton.RegisterCommand("echo", "echo user inputs", commandEcho)
	skeleton.RegisterCommand("logout", "logout users", commandLogout)
	skeleton.RegisterCommand("close", "close room with RoomId", commandClose)
}

func commandEcho(args []interface{}) interface{} {
	return fmt.Sprintf("%v", args)
}

func commandLogout(args []interface{}) interface{} {
	for _, arg := range args {
		userName := arg.(string)
		log.Debug("%s 登出游戏", userName)
		for aa, user := range Users {
			if user == userName {
				aa.Close()
			}
		}
	}
	return nil
}

func commandClose(args []interface{}) interface{} {
	for _, arg := range args {
		str := arg.(string)
		roomId, err := strconv.Atoi(str)
		if err != nil {
			return fmt.Sprintf("输入错误")
		}
		if room, ok := Rooms[roomId]; ok {
			room.Closed = true
			for _, aa := range room.User2Agent {
				if aa != nil {
					(*aa).Close()
				}
			}
			DeleteRoom(roomId, nil, false)
		}
	}
	return nil
}
