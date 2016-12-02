package initial

import (
	"github.com/JREAMLU/core/db/mysql"
	"github.com/astaxie/beego"
)

// InitMysql init mysql
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

	var gconfs mysql.GormConfs
	err = gconfs.InitGorms(beego.AppConfig.String("mysql.file"))
	if err != nil {
		beego.Error("init mysqls error: ", err)
		panic("init mysqls error")
	}
}
