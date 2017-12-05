package tasks

import (
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/util/log"
	"TechPlat/datapipe/util/redis"
	"github.com/pkg/errors"
)

var (
	innerLogger *logger.InnerLogger
)

var(
	NotConfigError = errors.New("no exists such config info")
)

func init() {
	innerLogger = logger.GetInnerLogger()
}





// getQueueData get queue data from redis
func getQueueData(taskConf *config.TaskInfo) (string, error){
	redisClient := redisutil.GetRedisClient(taskConf.FromServer)
	if redisClient == nil{
		return "", NotConfigError
	}
	val, err := redisClient.BRPop(taskConf.FromQueue)
	return val, err
}