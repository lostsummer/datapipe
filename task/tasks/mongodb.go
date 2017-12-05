package tasks

import (
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/const/log"
	"TechPlat/datapipe/util/log"
	"TechPlat/datapipe/repository/impl"
	"github.com/devfeel/dottask"
	"TechPlat/datapipe/trigger"
	"TechPlat/datapipe/counter"
)


//将queue获取的数据存入mongodb
func MongoDBHandler(ctx *task.TaskContext) error{
	title:= "Tasks:MongoDBHandler"
	taskConf := ctx.TaskData.(*config.TaskInfo)

	//get queue data
	val, err := getQueueData(taskConf)
	if err != nil{
		logger.Log(title+":getRedisData error -> "+err.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
		return err
	}
	logger.Log(title+":getRedisData -> "+val, taskConf.TaskID, logdefine.LogLevel_Debug)


	mongoHandler := new(impl.MongoHandler)
	mongoHandler.SetConn(config.CurrentConfig.MongoDB.ServerUrl, config.CurrentConfig.MongoDB.DBName)
	//insert data to mongo
	err = mongoHandler.InsertJsonData(taskConf.TargetValue, val)
	if err != nil {
		logger.Log(title+":InsertJsonData ["+val+"] error -> "+err.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
	} else {
		logger.Log(title+":InsertJsonData success ->["+val+"]", taskConf.TaskID, logdefine.LogLevel_Debug)
	}

	//deal trigger
	if err == nil && taskConf.HasTrigger() {
		err = trigger.SendSignal(taskConf.TriggerServer, taskConf.TriggerQueue, val)
		if err != nil {
			logger.Log(title+":SendTriggerSignal error -> ["+val+"] "+err.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
		}else{
			logger.Log(title+":SendTriggerSignal success -> ["+val+"]", taskConf.TaskID, logdefine.LogLevel_Debug)
		}
	}

	//deal counter
	if err == nil && taskConf.HasCounter() {
		err = counter.Count(taskConf.CounterServer, taskConf.CounterKey)
		if err != nil {
			logger.Log(title+":Counter error -> "+err.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
		}else{
			logger.Log(title+":Counter success", taskConf.TaskID, logdefine.LogLevel_Debug)

		}
	}
	return nil
}

