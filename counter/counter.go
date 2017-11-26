package counter

import (
	"errors"
	"TechPlat/tongjiservice/util/redis"
	"time"
	"sync"
)

const(
	timeFormatDaily = "20060102"
)

var(
	messageChan chan *countInfo
	countInfoPool *sync.Pool
)

type(
	countInfo struct{
		Server string
		Key string
	}
)

func init(){
	messageChan = make(chan *countInfo, 10000)
	countInfoPool = &sync.Pool{
		New : func() interface{}{
			return new(countInfo)
		},
	}

	//start count deal
	go doCount()
}

// Count counter with server and key
func Count(server, key string) error{
	info:=countInfoPool.Get().(*countInfo)
	info.reset(server,key)
	messageChan <- info
	return nil
}


func (info *countInfo) reset(server, key string){
	info.Server = server
	info.Key = key
}

func doCount() error{
	for{
		info :=<- messageChan
		redisClient := redisutil.GetRedisClient(info.Server)
		if redisClient == nil{
			return errors.New("no exists Counter server "+ info.Server)
		}
		_, err:= redisClient.INCR(createDailyCounterKey(info.Key))
		countInfoPool.Put(info)
		return err
	}

}

func createDailyCounterKey(rawKey string) string{
	return rawKey + "_" + time.Now().Format(timeFormatDaily)
}