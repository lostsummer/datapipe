package config

import (
	"TechPlat/datapipe/util/log"
	"encoding/xml"
	"io/ioutil"
	"os"
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
	err = xml.Unmarshal(content, &result)
	if err != nil {
		innerLogger.Warn("AppConfig::InitConfig 配置文件[" + configFile + "]解析失败 - " + err.Error())
		os.Exit(1)
	}
	result.TaskMap = make(map[string]*TaskInfo)
	for k, v := range result.Tasks {
		result.TaskMap[v.TaskID] = &result.Tasks[k]
		innerLogger.Info("AppConfig::InitConfig Load Task => " + v.TaskID + "," + v.TargetType + "," + v.TargetValue + "," + v.TriggerServer + "," + v.TriggerQueue)
	}

	result.ImporterMap = make(map[string]*Importer)
	if result.HttpServer.Enable {
		for k, v := range result.HttpServer.Importers {
			if v.Enable {
				result.ImporterMap[v.Name] = &result.HttpServer.Importers[k]
				innerLogger.Info("AppConfig::InitConfig Load Importer => " + v.Name + "," + v.ServerType + "," + v.ServerUrl + "," + v.ToQueue)
			}
		}
	}

	CurrentConfig = &result

	innerLogger.Info("AppConfig::InitConfig 配置文件[" + configFile + "]完成")

	return CurrentConfig
}

// GetKafkaServerUrl return kafka server info
func GetKafkaServerUrl() string {
	if CurrentConfig == nil {
		return ""
	}
	return CurrentConfig.Kafka.ServerUrl
}
