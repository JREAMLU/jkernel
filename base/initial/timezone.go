package initial

import (
	"time"

	"github.com/astaxie/beego"
)

func InitTimezone() {
	time.LoadLocation(beego.AppConfig.String("Timezone"))
}
