package redis

import (
	"go-websocket/configs"
	"strconv"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

var redisPool *redis.Pool
var mu sync.Mutex

func GetInstance() *redis.Pool {
	if redisPool == nil {
		//创建时才加锁
		mu.Lock()
		defer mu.Unlock()
		redisPool = newPool()
	}
	return redisPool
}

func newPool() *redis.Pool {
	host := configs.RedisHost
	port := configs.RedisPort
	db, _ := strconv.Atoi(configs.RedisChannel)
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host+":"+port, redis.DialDatabase(db))
			if nil != err {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:         100,                //最大空闲连接数
		MaxActive:       4000,              //最大连接数
		IdleTimeout:     180 * time.Second, //空闲连接超时时间
		Wait:            true,              //过最大连接，是报错，还是等待
		MaxConnLifetime: 0,
	}

	return pool
}

//获取key值
func Get(key string) (string, error) {
	rs := GetInstance().Get()
	defer rs.Close()
	reply, err := rs.Do("GET", key)
	if reply == nil {
		return "", nil
	} else {
		return redis.String(reply, err)
	}
}

//设置key值
func Set(key string, value string) (interface{}, error) {
	rs := GetInstance().Get()
	defer rs.Close()
	return rs.Do("SET", key, value)
}

//设置key值加过期时间
func SetWithSurvivalTime(key string, value string, survivalTime int) (interface{}, error) {
	if survivalTime == 0 {
		return Set(key, value)
	}

	rs := GetInstance().Get()
	defer rs.Close()
	return rs.Do("SET", key, value, "EX", survivalTime)
}

//设置key的过期时间
func SetSurvivalTime(key string, survivalTime int) (interface{}, error) {
	if survivalTime < 0 {
		return nil, nil
	}
	rs := GetInstance().Get()
	defer rs.Close()
	return rs.Do("expire", key, survivalTime)
}

//删除key值
func Del(key string) (interface{}, error) {
	rs := GetInstance().Get()
	defer rs.Close()
	return rs.Do("DEL", key)
}

//向集合里添加元素
func SetAdd(key, value string) (interface{}, error) {
	rs := GetInstance().Get()
	defer rs.Close()
	return rs.Do("SADD", key, value)
}

//删除集合里的某个元素
func DelSetKey(key, member string) (interface{}, error) {
	rs := GetInstance().Get()
	defer rs.Close()
	return rs.Do("SREM", key, member)
}

//返回集合里的元素列表
func SMEMBERS(key string) ([]string, error) {
	rs := GetInstance().Get()
	defer rs.Close()
	reply, err := rs.Do("SMEMBERS", key)
	if reply == nil {
		return []string{}, nil
	} else {
		return redis.Strings(reply, err)
	}
}

//集合里的元素个数
func SCARD(key string) (int64, error) {
	rs := GetInstance().Get()
	defer rs.Close()
	numString, err := rs.Do("SCARD", key)
	return numString.(int64), err
}

//判断成员元素是否是集合的成员
func SISMEMBER(key string, value string) (bool, error) {
	rs := GetInstance().Get()
	defer rs.Close()
	numInt, err := rs.Do("SISMEMBER", key, value)
	if err != nil {
		return false, err
	}
	return numInt.(int64) > 0, err
}
