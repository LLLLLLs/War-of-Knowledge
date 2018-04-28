/*
Author  : Leshuo Lian
Time    : 2018\4\28 0028 
*/

package internal

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestFilter(t *testing.T) {
	ast := assert.New(t)
	b := filter("helloWorld12345")
	fmt.Println(b)
	ast.True(b)
	b = filter("你好111")
	ast.False(b)
	b = filter("asdf aaa")
	ast.False(b)
	b = filter("a123456 1")
	ast.False(b)
	b = filter("laksnfoa21354")
	ast.True(b)
	b = filter("asdf-asdf")
	ast.False(b)
}
