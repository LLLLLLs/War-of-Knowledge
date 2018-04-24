/*
Author  : Leshuo Lian
Time    : 2018\4\24 0024 
*/

package internal

import (
	"fmt"
)

func init() {
	skeleton.RegisterCommand("echo", "echo user inputs", commandEcho)
}

func commandEcho(args []interface{}) interface{} {
	return fmt.Sprintf("%v", args)
}
