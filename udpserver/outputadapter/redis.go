package outputadapter

import (
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/queue"
)

func OutputRedisAdapter(conf config.OutputAdapter, appid string, logstr string) {
	q := &queue.Queue{
		conf.Url,
		0,
		conf.ToQueue,
	}
	q.Push(logstr)
}
