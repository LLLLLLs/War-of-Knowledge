package gamedata

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"fmt"
)

var (
	Db    *xorm.Engine
	debug = true
)

const (
	DB_HOST     = "118.89.46.254"
	DB_PORT     = 3306
	DB_USER     = "root"
	DB_PWD      = "longjing"
	DB_DATABASE = "lls_gradu_project"
)

func SetDebug(b bool) {
	debug = b
}

func InitDB(dburi string) {
	engine, err := xorm.NewEngine("mysql", dburi)
	if err != nil {
		panic(err)
	}
	if debug {
		engine.ShowSQL(false)
	}
	Db = engine
}

func init() {
	dburi := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		DB_USER,
		DB_PWD,
		DB_HOST,
		DB_PORT,
		DB_DATABASE,
	)
	InitDB(dburi)
}
