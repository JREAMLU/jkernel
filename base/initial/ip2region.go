package initial

import (
	"github.com/JREAMLU/core/com"
	"github.com/astaxie/beego"
)

// InitIP2Region init ip to region project
func InitIP2Region() {
	err := com.InitIP2Region(beego.AppConfig.String("ip.path"))
	if err != nil {
		beego.Error("init ip2region error: ", err)
		panic("init ip2region error")
	}
}
