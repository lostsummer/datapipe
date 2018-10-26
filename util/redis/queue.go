package redisutil

import (
	"TechPlat/datapipe/global"
)

type Queue struct {
	Server string
	DB     int
	Key    string
}

// get queue data from redis
func (q *Queue) Pop() (string, error) {
	redisClient := GetRedisClient(q.Server)
	if redisClient == nil {
		return "", global.GetRedisError
	}
	if q.DB > 0 {
		if err := redisClient.Select(q.DB); err != nil {
			return "", err
		}
	}
	val, err := redisClient.BRPop(q.Key)
	return val, err
}

// add data to redis queue
func (q *Queue) Push(val string) (int64, error) {
	redisClient := GetRedisClient(q.Server)
	if redisClient == nil {
		return -1, global.GetRedisError
	}
	if q.DB > 0 {
		if err := redisClient.Select(q.DB); err != nil {
			return -1, err
		}
	}
	n, err := redisClient.LPush(q.Key, val)
	return n, err
}
