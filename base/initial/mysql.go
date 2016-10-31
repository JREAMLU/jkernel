package initial

import (
	"git.corp.plu.cn/phpgo/core/mysql"
	"github.com/astaxie/beego"
)

func InitMysql() {
	var gconf mysql.GormConf
	gconf.Driver = beego.AppConfig.String("mysql.driver")
	gconf.Setting = beego.AppConfig.String("mysql.setting")
	gconf.SingularTable, _ = beego.AppConfig.Bool("mysql.singulartable")
	gconf.LogMode, _ = beego.AppConfig.Bool("mysql.logmode")
	err := gconf.InitGorm()
	if err != nil {
		beego.Error("init mysql error: ", err)
		panic("init mysql error")
	}
}
