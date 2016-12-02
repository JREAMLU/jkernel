package initial

import (
	"time"

	"github.com/astaxie/beego"
)

// InitTimezone init timezone
func InitTimezone() {
	time.LoadLocation(beego.AppConfig.String("Timezone"))
}
