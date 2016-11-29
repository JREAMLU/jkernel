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

func ShortenHMGet(origin []interface{}) (reply []string, err error) {
	params := append([]interface{}{shortenKey}, origin...)
	conn := redigos.GetRedisClient(REDISSERVER_BASE, true)
	reply, err = redis.Strings(conn.Do("HMGET", params...))
	conn.Close()
	return reply, err
}

func ShortenHMSet(originShort []interface{}) (reply string, err error) {
	params := append([]interface{}{shortenKey}, originShort...)
	conn := redigos.GetRedisClient(REDISSERVER_BASE, true)
	reply, err = redis.String(conn.Do("HMSET", params...))
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

func ExpandHMSet(shortOrigin []interface{}) (reply string, err error) {
	params := append([]interface{}{expandKey}, shortOrigin...)
	conn := redigos.GetRedisClient(REDISSERVER_BASE, true)
	reply, err = redis.String(conn.Do("HMSET", params...))
	conn.Close()
	return reply, err
}

func ExpandHMGet(short []interface{}) (reply []string, err error) {
	params := append([]interface{}{expandKey}, short...)
	conn := redigos.GetRedisClient(REDISSERVER_BASE, true)
	reply, err = redis.Strings(conn.Do("HMGET", params...))
	conn.Close()
	return reply, err
}
