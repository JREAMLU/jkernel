package mredis

import (
	"fmt"

	"github.com/JREAMLU/core/db/redigos"
	"github.com/JREAMLU/jkernel/base/services/atom"
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

//TODO 单独的功能 数据放到service层
func ShortenHMGet(origin []interface{}) (shortens map[string]interface{}, emptys []interface{}, err error) {
	params := append([]interface{}{shortenKey}, origin...)
	conn := redigos.GetRedisClient(REDISSERVER_BASE, true)
	reply, err := redis.Strings(conn.Do("HMGET", params...))
	conn.Close()
	if err != nil {
		return nil, nil, err
	}
	shortens = make(map[string]interface{})
	for k, v := range origin {
		atom.Mu.Lock()
		if reply[k] != "" {
			shortens[v.(string)] = reply[k]
		} else {
			emptys = append(emptys, v.(string))
		}
		atom.Mu.Unlock()
	}
	return shortens, emptys, err
}

func ShortenHMSet(url []interface{}) (reply string, err error) {
	params := append([]interface{}{shortenKey}, url...)
	conn := redigos.GetRedisClient(REDISSERVER_BASE, true)
	reply, err = redis.String(conn.Do("HMSET", params...))
	conn.Close()
	fmt.Println("<<<<<<<<<<", reply, err)
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
