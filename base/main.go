package main

import (
	"time"

	"github.com/JREAMLU/core/logs"
	_ "github.com/JREAMLU/jkernel/base/routers"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

func init() {
	//timezone set
	time.LoadLocation(beego.AppConfig.String("Timezone"))

	//beego log
	jlogs.InitLogs()
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
	}

	beego.AddFuncMap("i18n", i18n.Tr)

	beego.Run()
}
