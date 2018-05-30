// redisclient
package redisutil

import (
	"sync"

	"github.com/garyburd/redigo/redis"
)

type RedisClient struct {
	pool    *redis.Pool
	Address string
}

var (
	redisMap map[string]*RedisClient
	mapMutex *sync.RWMutex
)

const (
	defaultTimeout = 60 * 10 //默认10分钟
)

func init() {
	redisMap = make(map[string]*RedisClient)
	mapMutex = new(sync.RWMutex)
}

// 重写生成连接池方法
func newPool(redisIP string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:   5,
		MaxActive: 20, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisIP)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

//获取指定Address的RedisClient
func GetRedisClient(address string) *RedisClient {
	mapMutex.RLock()
	redis, mok := redisMap[address]
	mapMutex.RUnlock()
	if !mok {
		mapMutex.Lock()
		if r, ok := redisMap[address]; ok {
			redis = r
		} else {
			redis = &RedisClient{Address: address, pool: newPool(address)}
			redisMap[address] = redis
		}
		mapMutex.Unlock()
	}
	return redis
}

//获取指定key的内容
func (rc *RedisClient) Get(key string) (string, error) {
	// 从连接池里面获得一个连接
	conn := rc.pool.Get()
	// 连接完关闭，其实没有关闭，是放回池里，也就是队列里面，等待下一个重用
	defer conn.Close()
	reply, errDo := conn.Do("GET", key)
	if errDo == nil && reply == nil {
		return "", nil
	}
	val, err := redis.String(reply, errDo)
	return val, err
}

//获取指定hashset的内容
func (rc *RedisClient) HGet(hashID string, field string) (string, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	reply, errDo := conn.Do("HGET", hashID, field)
	if errDo == nil && reply == nil {
		return "", nil
	}
	val, err := redis.String(reply, errDo)
	return val, err
}

//对存储在指定key的数值执行原子的加1操作
func (rc *RedisClient) INCR(key string) (int, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	reply, errDo := conn.Do("INCR", key)
	if errDo == nil && reply == nil {
		return 0, nil
	}
	val, err := redis.Int(reply, errDo)
	return val, err
}

//获取指定hashset的所有内容
func (rc *RedisClient) HGetAll(hashID string) (map[string]string, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	reply, err := redis.StringMap(conn.Do("HGetAll", hashID))
	return reply, err
}

//设置指定hashset的内容
func (rc *RedisClient) HSet(hashID string, field string, val string) error {
	conn := rc.pool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", hashID, field, val)
	return err
}

//删除，并获得该列表中的最后一个元素，或阻塞，直到有一个可用
func (rc *RedisClient) BRPop(key string) (string, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.StringMap(conn.Do("BRPOP", key, defaultTimeout))
	if err != nil {
		return "", err
	} else {
		return val[key], nil
	}
}

//将所有指定的值插入到存于 key 的列表的头部
func (rc *RedisClient) LPush(key string, val string) (int64, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	ret, err := redis.Int64(conn.Do("LPUSH", key, val))
	if err != nil {
		return -1, err
	} else {
		return ret, nil
	}
}

//设置指定key的内容
func (rc *RedisClient) Set(key string, val string) (string, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.String(conn.Do("SET", key, val))
	return val, err
}
