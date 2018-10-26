package queue

import (
	"TechPlat/datapipe/util/redis"
)

type RedisTarget struct {
	Server string
	DB     int
	Key    string
}

type RedisSource RedisTarget

func (r *RedisTarget) Push(val string) (int64, error) {
	q := redisutil.Queue{
		r.Server,
		r.DB,
		r.Key,
	}
	return q.Push(val)
}

func (r *RedisSource) Pop() (string, error) {
	q := redisutil.Queue{
		r.Server,
		r.DB,
		r.Key,
	}
	return q.Pop()
}
