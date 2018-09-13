package outputadapter

import (
	"errors"
	"TechPlat/datapipe/util/redis"
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/util/log"
)

func OutputRedisAdapter(conf config.OutputAdapter, appid string, logstr string)  {
	server, queue := conf.Url, conf.ToQueue
	pushQueueData(server, queue, logstr)
}

func pushQueueData(server, queue, val string) (int64, error) {
	redisClient := redisutil.GetRedisClient(server)
	if redisClient == nil {
		return -1, errors.New("get rediscli failed")
	}

	defer func() {
		if p := recover(); p != nil {
			if s, ok := p.(string); ok {
				logger.GetInnerLogger().Error(s)
			}
		}
	}()
	return redisClient.LPush(queue, val)
}