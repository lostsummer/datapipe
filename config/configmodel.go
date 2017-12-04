package config

import (
	"encoding/xml"
)

const (
	Target_MongoDB = "mongodb"
	Target_Http    = "http"
	Target_Kafka    = "kafka"
)

//代理配置信息
type AppConfig struct {
	XMLName xml.Name   `xml:"config"`
	Log     Log        `xml:"log"`
	Redis   Redis      `xml:"redis"`
	MongoDB MongoDB    `xml:"mongodb"`
	Kafka 	*Kafka    `xml:"kafka"`
	Tasks   []TaskInfo `xml:"tasks>task"`
	TaskMap map[string]*TaskInfo
}

//Redis配置
type Redis struct {
	KeyCommonPre string `xml:"keycommonpre,attr"`
}

//MongodDB配置
type MongoDB struct {
	ServerUrl string `xml:"serverurl,attr"`
	DBName   string `xml:"dbname,attr"`
}

//kafka配置
type Kafka struct {
	ServerUrl string `xml:"serverurl,attr"`
}

//log配置
type Log struct {
	FilePath string `xml:"filepath,attr"`
}

//数据模板
type TaskInfo struct {
	TaskID          string `xml:"taskid,attr"`
	FromServer    string `xml:"fromserver,attr"`
	FromQueue      string `xml:"fromqueue,attr"`
	TargetValue     string `xml:"targetvalue,attr"`
	TargetType      string `xml:"targettype,attr"`
	TriggerServer string `xml:"triggerserver,attr"`
	TriggerQueue   string `xml:"triggerqueue,attr"`
	CounterServer   string `xml:"counterserver,attr"` //计数器server
	CounterKey   string `xml:"counterkey,attr"` //计数器key
}

// HasTrigger 检查是否存在触发器配置
func (t *TaskInfo) HasTrigger() bool{
	return t.TriggerServer != "" && t.TriggerQueue != ""
}

// HasCounter 检查是否存在计数器配置
func (t *TaskInfo) HasCounter() bool{
	return t.CounterServer != "" && t.CounterKey != ""
}