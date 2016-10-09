package initial

import (
	"github.com/JREAMLU/core/db/redigos"
	"github.com/astaxie/beego"
)

func InitRedis() {
	redigos.LoadRedisConfig(beego.AppConfig.String("redis.file"))
}
