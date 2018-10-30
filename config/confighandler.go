package config

import (
	"TechPlat/datapipe/queue"
	"TechPlat/datapipe/util/log"
	"encoding/xml"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
)

const ()

var (
	CurrentConfig  *AppConfig
	CurrentBaseDir string
	innerLogger    *logger.InnerLogger

	mutex *sync.RWMutex
)

func init() {
	//初始化读写锁
	mutex = new(sync.RWMutex)
	innerLogger = logger.GetInnerLogger()
}

func SetBaseDir(baseDir string) {
	CurrentBaseDir = baseDir
}

//初始化配置文件
func InitConfig(configFile string) *AppConfig {
	innerLogger.Info("AppConfig::InitConfig 配置文件[" + configFile + "]开始...")
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		innerLogger.Warn("AppConfig::InitConfig 配置文件[" + configFile + "]无法解析 - " + err.Error())
		os.Exit(1)
	}

	var result AppConfig
	CurrentConfig = &result
	err = xml.Unmarshal(content, &result)
	if err != nil {
		innerLogger.Warn("AppConfig::InitConfig 配置文件[" + configFile + "]解析失败 - " + err.Error())
		os.Exit(1)
	}

	result.RedisMap = make(map[string]*Redis)
	for k, v := range result.Redises {
		result.RedisMap[v.ID] = &result.Redises[k]
		innerLogger.Info("AppConfig::InitConfig Load Redises => " + v.ID + "," + v.URL + "," + v.DB)
	}

	result.MongoDBMap = make(map[string]*MongoDB)
	for k, v := range result.MongoDBs {
		result.MongoDBMap[v.ID] = &result.MongoDBs[k]
		innerLogger.Info("AppConfig::InitConfig Load MongoDBs => " + v.ID + "," + v.URL + "," + v.DB)
	}

	result.KafkaMap = make(map[string]*Kafka)
	for k, v := range result.Kafkas {
		result.KafkaMap[v.ID] = &result.Kafkas[k]
		innerLogger.Info("AppConfig::InitConfig Load Kafkas => " + v.ID + v.URL)
	}

	result.HTTPMap = make(map[string]*HTTP)
	for k, v := range result.HTTPs {
		result.HTTPMap[v.ID] = &result.HTTPs[k]
		innerLogger.Info("AppConfig::InitConfig Load HTTPs => " + v.ID + v.URL)
	}

	result.TaskMap = make(map[string]*TaskInfo)
	for k, v := range result.Tasks {
		if v.Enable {
			result.TaskMap[v.ID] = &result.Tasks[k]
			innerLogger.Info("AppConfig::InitConfig Load Task => " + v.ID + "," + v.Target.Type + "," + v.Target.ID + "," + v.Trigger.ID + "," + v.Trigger.Queue)
		}
	}

	result.ImporterMap = make(map[string]*Importer)
	if result.HttpServer.Enable {
		for k, v := range result.HttpServer.Importers {
			if v.Enable {
				result.ImporterMap[v.ID] = &result.HttpServer.Importers[k]
				innerLogger.Info("AppConfig::InitConfig Load Importer => " + v.ID + "," + v.Target.Type + "," + v.Target.ID + "," + v.Target.Queue)
			}
		}
	}

	result.AccumulatorMap = make(map[string]*Accumulator)
	if result.HttpServer.Enable {
		for k, v := range result.HttpServer.Accumulators {
			if v.Enable {
				result.AccumulatorMap[v.ID] = &result.HttpServer.Accumulators[k]
				innerLogger.Info("AppConfig::InitConfig Load Accumulator => " + v.ID + "," + v.Target.Type + "," + v.Target.ID + "," + v.Target.Counter + v.Target.Set)
			}
		}
	}

	result.ImporterTargetMap = make(map[string]queue.Target)
	for k, v := range result.ImporterMap {
		t := getImptTarget(v)
		if t != nil {
			result.ImporterTargetMap[k] = t
			innerLogger.Info("AppConfig::InitConfig Load ImporterTargetMap => " + k)
		}
	}

	result.TaskSourceMap = make(map[string]queue.Source)
	for k, v := range result.TaskMap {
		t := getTaskSource(v)
		if t != nil {
			result.TaskSourceMap[k] = t
			innerLogger.Info("AppConfig::InitConfig Load TaskSourceMap => " + k)
		}
	}

	result.TaskTargetMap = make(map[string]queue.Target)
	for k, v := range result.TaskMap {
		t := getTaskTarget(v)
		if t != nil {
			result.TaskTargetMap[k] = t
			innerLogger.Info("AppConfig::InitConfig Load TaskTargetMap => " + k)
		}
	}

	result.TaskTriggerMap = make(map[string]queue.Target)
	for k, v := range result.TaskMap {
		t := getTaskTrigger(v)
		if t != nil {
			result.TaskTriggerMap[k] = t
			innerLogger.Info("AppConfig::InitConfig Load TaskTriggerMap => " + k)
		}
	}

	innerLogger.Info("AppConfig::InitConfig 配置文件[" + configFile + "]完成")

	return CurrentConfig
}

func getImptTarget(imptConf *Importer) queue.Target {
	switch imptConf.Target.Type {
	case Target_Redis:
		return getRedisTarget(&imptConf.Target)
	default:
		return nil // http importer 暂无连接redis之外类型target的必要
	}

}

func getTaskTarget(taskConf *TaskInfo) queue.Target {
	switch taskConf.Target.Type {
	case Target_Redis:
		return getRedisTarget(&taskConf.Target)
	case Target_MongoDB:
		return getMongoDBTarget(&taskConf.Target)
	case Target_Kafka:
		return getKafkaTarget(&taskConf.Target)
	case Target_Http:
		return getHTTPTarget(&taskConf.Target)
	default:
		return nil
	}

}

func getTaskSource(taskConf *TaskInfo) queue.Source {
	switch taskConf.Source.Type {
	case Target_Redis:
		return getRedisSource(&taskConf.Source)
	default:
		return nil // http importer 暂无连接redis之外类型target的必要
	}

}

func getTaskTrigger(taskConf *TaskInfo) queue.Target {
	getRedisTrigger := getRedisTarget
	switch taskConf.Trigger.Type {
	case Target_Redis:
		return getRedisTrigger(&taskConf.Trigger)
	default:
		return nil // trigger 目前设计仅有redis类型
	}
}

func getRedisSource(qConf *Queue) queue.Source {
	if qConf == nil {
		return nil
	}
	if r, exist := CurrentConfig.RedisMap[qConf.ID]; exist {
		db, err := strconv.Atoi(r.DB)
		if err != nil {
			db = 0 // 目前线上只使用db0, 所以对于配置错误暂时做db=0处理
		}
		url := r.URL
		q := qConf.Queue
		return &queue.RedisSource{
			url,
			db,
			q,
		}
	}
	return nil
}

func getRedisTarget(qConf *Queue) queue.Target {
	if qConf == nil {
		return nil
	}
	if r, exist := CurrentConfig.RedisMap[qConf.ID]; exist {
		db, err := strconv.Atoi(r.DB)
		if err != nil {
			db = 0 // 目前线上只使用db0, 所以对于配置错误暂时做db=0处理
		}
		url := r.URL
		q := qConf.Queue
		return &queue.RedisTarget{
			url,
			db,
			q,
		}
	}
	return nil
}

func getMongoDBTarget(qConf *Queue) queue.Target {
	if qConf == nil {
		return nil
	}
	if m, exist := CurrentConfig.MongoDBMap[qConf.ID]; exist {
		url := m.URL
		db := m.DB
		q := qConf.Queue
		return &queue.MongoDBTarget{
			url,
			db,
			q,
		}
	}
	return nil
}

func getKafkaTarget(qConf *Queue) queue.Target {
	if qConf == nil {
		return nil
	}
	if k, exist := CurrentConfig.KafkaMap[qConf.ID]; exist {
		url := k.URL
		q := qConf.Queue
		return &queue.KafkaTarget{
			url,
			q, //kafka topic
		}
	}
	return nil
}

func getHTTPTarget(qConf *Queue) queue.Target {
	if qConf == nil {
		return nil
	}
	if h, exist := CurrentConfig.HTTPMap[qConf.ID]; exist {
		url := h.URL
		return &queue.HttpTarget{
			url,
		}
	}
	return nil
}
