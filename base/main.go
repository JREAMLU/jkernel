package main

import (
	_ "github.com/JREAMLU/jkernel/base/routers"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
	}

	beego.AddFuncMap("i18n", i18n.Tr)

	beego.Run()
}
