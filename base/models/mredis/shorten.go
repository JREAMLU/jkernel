package mredis

import (
	"github.com/JREAMLU/core/db/redigos"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
)

const (
	REDISSERVER_BASE = "base"
)

var (
	shortenKey = beego.AppConfig.String("redis.key.shorten")
	expandKey  = beego.AppConfig.String("redis.key.expand")
)

func ShortenHSet(origin string, short string) (reply int64, err error) {
	conn := redigos.GetRedisClient(REDISSERVER_BASE, true)
	reply, err = redis.Int64(conn.Do("HSET", shortenKey, origin, short))
	conn.Close()
	return reply, err
}

func ShortenHGet(origin string) (reply string, err error) {
	conn := redigos.GetRedisClient(REDISSERVER_BASE, true)
	reply, err = redis.String(conn.Do("HGET", shortenKey, origin))
	conn.Close()
	return reply, err
}

func ExpandHSet(short string, origin string) (reply int64, err error) {
	conn := redigos.GetRedisClient(REDISSERVER_BASE, true)
	reply, err = redis.Int64(conn.Do("HSET", expandKey, short, origin))
	conn.Close()
	return reply, err
}

func ExpandHGet(short string) (reply string, err error) {
	conn := redigos.GetRedisClient(REDISSERVER_BASE, true)
	reply, err = redis.String(conn.Do("HGET", expandKey, short))
	conn.Close()
	return reply, err
}
