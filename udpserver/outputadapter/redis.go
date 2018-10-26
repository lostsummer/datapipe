package outputadapter

import (
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/util/redis"
)

func OutputRedisAdapter(conf config.OutputAdapter, appid string, logstr string) {
	q := &redisutil.Queue{
		conf.Url,
		0,
		conf.ToQueue,
	}
	q.Push(logstr)
}
