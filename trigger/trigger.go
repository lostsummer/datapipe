package trigger

import (
	"TechPlat/datapipe/util/redis"
	"github.com/pkg/errors"
)

// SendSignal send trigger signal
func SendSignal(server, queue string, signal string) error{
	redisClient := redisutil.GetRedisClient(server)
	if redisClient == nil{
		return errors.New("no exists trigger server "+ server)
	}
	_, err := redisClient.LPush(queue, signal)
	return err
}
