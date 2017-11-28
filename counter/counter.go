package counter

import (
	"errors"
	"TechPlat/datapipe/util/redis"
	"TechPlat/datapipe/global"
	"time"
	"sync"
	"github.com/devfeel/dottask"
)

const(
	timeFormatDaily = "20060102"
)
const(
	QueueTaskName = "CounterQueueTask"
	QueueSize = 10000
)

var(
	countInfoPool *sync.Pool
	queueTask *task.QueueTask
)

type(
	countInfo struct{
		Server string
		Key string
	}
)

func (info *countInfo) reset(server, key string){
	info.Server = server
	info.Key = key
}

// InitCounter init counter basic object
func StartCounter() error{
	//init pool
	countInfoPool = &sync.Pool{
		New : func() interface{}{
			return new(countInfo)
		},
	}

	//load queue task
	t, exists := global.TaskService.GetTask(QueueTaskName)
	if !exists{
		return errors.New("not exists queue task")
	}else{
		queueTask = t.(*task.QueueTask)
	}
	return nil
}

// Count do count with server and key
func Count(server, key string) error{
	info:=countInfoPool.Get().(*countInfo)
	info.reset(server,key)
	queueTask.EnQueue(info)
	return nil
}

// DealMessage deal count message
func DealMessage(ctx *task.TaskContext) error{
	info, ok := ctx.Message.(*countInfo)
	if !ok{
		return errors.New("message is not legal type")
	}

	redisClient := redisutil.GetRedisClient(info.Server)
	if redisClient == nil{
		return errors.New("no exists Counter server "+ info.Server)
	}
	_, err:= redisClient.INCR(createDailyCounterKey(info.Key))
	countInfoPool.Put(info)
	return err
}

func createDailyCounterKey(rawKey string) string{
	return rawKey + "_" + time.Now().Format(timeFormatDaily)
}