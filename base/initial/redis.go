package initial

import (
	"github.com/JREAMLU/core/db/redigos"
	"github.com/astaxie/beego"
)

func InitRedis() {
	err := redigos.LoadRedisConfig(beego.AppConfig.String("redis.file"))
	if err != nil {
		beego.Error("init redis error: ", err)
	}
}
