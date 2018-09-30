package tasks

import (
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/const/log"
	"TechPlat/datapipe/repository/impl"
	"TechPlat/datapipe/util/log"

	"github.com/devfeel/dottask"
)

type MongoDBPusher PusherBase

const (
	mongodbLogTitle = "Tasks:MongDBHandler"
)

func (m MongoDBPusher) LogTitle() string {
	return m.Title
}

func (m MongoDBPusher) Push(taskConf *config.TaskInfo, val string) error {
	title := m.Title
	mongoHandler := new(impl.MongoHandler)
	mongoName := taskConf.TargetName
	var mongoConf *config.MongoDB
	for _, m := range config.CurrentConfig.MongoDBs.MongoDBList {
		if m.Name == mongoName {
			mongoConf = &m
		}
	}

	// 今后配置检查提到外部
	if mongoConf != nil {
		mongoHandler.SetConn(mongoConf.ServerUrl, mongoConf.DBName)
	} else {
		logger.Log(title+":GetConfig no "+mongoName+" define", taskConf.TaskID, logdefine.LogLevel_Error)
		return NotConfigError
	}
	//insert data to mongo
	err := mongoHandler.InsertJsonData(taskConf.TargetValue, val)
	if err != nil {
		logger.Log(title+":InsertJsonData ["+val+"] error -> "+err.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
		return err
	}
	logger.Log(title+":InsertJsonData success ->["+val+"]", taskConf.TaskID, logdefine.LogLevel_Debug)
	return nil
}

//将queue获取的数据存入mongodb
func MongoDBHandler(ctx *task.TaskContext) error {
	handler(ctx, MongoDBPusher{mongodbLogTitle})
	return nil
}
