package models

import (
	"github.com/astaxie/beego"
	// "github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var (
	alias      = beego.AppConfig.String("db::alias")
	driver     = beego.AppConfig.String("db::driver")
	username   = beego.AppConfig.String("db::username")
	password   = beego.AppConfig.String("db::password")
	database   = beego.AppConfig.String("db::database")
	charset    = beego.AppConfig.String("db::charset")
	maxIdle, _ = beego.AppConfig.Int("db::maxIdle")

	showSql, _   = beego.AppConfig.Bool("db:xorm::showSql")
	showDebug, _ = beego.AppConfig.Bool("db:xorm::showDebug")
	showErr, _   = beego.AppConfig.Bool("db:xorm::showErr")
	showWarn, _  = beego.AppConfig.Bool("db:xorm::showWarn")
	showInfo, _  = beego.AppConfig.Bool("db:xorm::showInfo")

	engine *xorm.Engine
)

func init() {
	// orm.RegisterDataBase(alias, driver, username+":@/"+database+"?charset="+charset, maxIdle)
	var err error
	engine, err = xorm.NewEngine(driver, username+":@/"+database+"?charset="+charset)
	/*
		engine.ShowSQL = showSql     //则会在控制台打印出生成的SQL语句；
		engine.ShowDebug = showDebug //则会在控制台打印调试信息；
		engine.ShowErr = showErr     //则会在控制台打印错误信息；
		engine.ShowWarn = showWarn   //则会在控制台打印警告信息；
		engine.ShowInfo = showInfo
	*/

	if err != nil {
	}
}

func Select() bool {
	return true
}

/*
func GetUrlOne() orm.Params {
	params := []interface{}{1, 2, 3, "http://jream.lu"}
	sql := "SELECT * FROM redirect WHERE redirect_id IN (?, ?, ?) AND long_url = ? "
	a, _ := mysql.Select(params, sql)
	return a[0]
}
*/
