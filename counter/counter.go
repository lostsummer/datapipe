package counter

import (
	"errors"
	"emoney/tongjiservice/util/redis"
	"time"
)

const(
	timeFormatDaily = "20060102"
)

// Count counter with server and key
func Count(server, key string) error{
	redisClient := redisutil.GetRedisClient(server)
	if redisClient == nil{
		return errors.New("no exists Counter server "+ server)
	}
	_, err:= redisClient.INCR(createDailyCounterKey(key))
	return err
}

func createDailyCounterKey(rawKey string) string{
	return rawKey + "_" + time.Now().Format(timeFormatDaily)
}