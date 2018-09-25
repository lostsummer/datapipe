package tasks

import (
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/const/log"
	"TechPlat/datapipe/counter"
	"TechPlat/datapipe/repository/impl"
	"TechPlat/datapipe/trigger"
	"TechPlat/datapipe/util/log"

	"github.com/devfeel/dottask"
)

//将queue获取的数据存入mongodb
func MongoDBHandler(ctx *task.TaskContext) error {
	title := "Tasks:MongoDBHandler"
	taskConf := ctx.TaskData.(*config.TaskInfo)

	if taskConf.Enable == false {
		return nil
	}

	//get queue data
	val, err := getQueueData(taskConf)
	if err != nil {
		logger.Log(title+":getRedisData error -> "+err.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
		return err
	}
	logger.Log(title+":getRedisData -> "+val, taskConf.TaskID, logdefine.LogLevel_Debug)

	mongoHandler := new(impl.MongoHandler)
	mongoName := taskConf.TargetName
	var mongoConf *config.MongoDB
	for _, m := range config.CurrentConfig.MongoDBs.MongoDBList {
		if m.Name == mongoName {
			mongoConf = &m
		}
	}
	if mongoConf != nil {
		mongoHandler.SetConn(mongoConf.ServerUrl, mongoConf.DBName)
	} else {
		logger.Log(title+":GetConfig no "+mongoName+" define", taskConf.TaskID, logdefine.LogLevel_Error)
		return nil
	}
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
		} else {
			logger.Log(title+":SendTriggerSignal success -> ["+val+"]", taskConf.TaskID, logdefine.LogLevel_Debug)
		}
	}

	//deal counter
	if err == nil && taskConf.HasCounter() {
		err = counter.Count(taskConf.CounterServer, taskConf.CounterKey)
		if err != nil {
			logger.Log(title+":Counter error -> "+err.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
		} else {
			logger.Log(title+":Counter success", taskConf.TaskID, logdefine.LogLevel_Debug)

		}
	}
	return nil
}
