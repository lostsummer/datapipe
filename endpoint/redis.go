package endpoint

import (
	"TechPlat/datapipe/util/redis"
)

type Redis struct {
	Server string
	DB     int
	Key    string
}

//type Redis Redis

func (r *Redis) Push(val string) (int64, error) {
	q := redisutil.Queue{
		r.Server,
		r.DB,
		r.Key,
	}
	return q.Push(val)
}

func (r *Redis) Pop() (string, error) {
	q := redisutil.Queue{
		r.Server,
		r.DB,
		r.Key,
	}
	return q.Pop()
}
