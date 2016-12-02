package initial

import (
	"github.com/JREAMLU/core/db/redigos"
	"github.com/astaxie/beego"
)

// InitRedis init redis
func InitRedis() {
	err := redigos.LoadRedisConfig(beego.AppConfig.String("redis.file"))
	if err != nil {
		beego.Error("init redis error: ", err)
		panic("init redis error")
	}
}
