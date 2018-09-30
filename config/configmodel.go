package config

import (
	"encoding/xml"
)

const (
	Target_MongoDB = "mongodb"
	Target_Http    = "http"
	Target_Kafka   = "kafka"
)

//代理配置信息
type AppConfig struct {
	XMLName        xml.Name        `xml:"config"`
	Log            Log             `xml:"log"`
	Redis          Redis           `xml:"redis"`
	MongoDBs       MongoDBs        `xml:"mongodbs"`
	Kafka          Kafka           `xml:"kafka"`
	HttpServer     HttpServer      `xml:"httpserver"`
	Tasks          []TaskInfo      `xml:"tasks>task"`
	UdpServer      UDPServer       `xml:"udpserver"`
	OutputAdapters []OutputAdapter `xml:"outputadapter>adapter"`
	ImporterMap    map[string]*Importer
	AccumulatorMap map[string]*Accumulator
	TaskMap        map[string]*TaskInfo
}

//Redis配置
type Redis struct {
	KeyCommonPre string `xml:"keycommonpre,attr"`
}

//MongodDB配置
type MongoDB struct {
	Name      string `xml:"name,attr"`
	ServerUrl string `xml:"serverurl,attr"`
	DBName    string `xml:"dbname,attr"`
}
type MongoDBs struct {
	MongoDBList []MongoDB `xml:"mongodb"`
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
	TaskID        string `xml:"taskid,attr"`
	Enable        bool   `xml:"enable,attr"`
	FromServer    string `xml:"fromserver,attr"`
	FromQueue     string `xml:"fromqueue,attr"`
	TargetValue   string `xml:"targetvalue,attr"`
	TargetType    string `xml:"targettype,attr"`
	TargetName    string `xml:"targetname,attr"`
	TargetFilter  string `xml:"targetfilter,attr"`
	TriggerServer string `xml:"triggerserver,attr"`
	TriggerQueue  string `xml:"triggerqueue,attr"`
	TriggerFilter string `xml:"triggerfilter,attr"`
	CounterServer string `xml:"counterserver,attr"` //计数器server
	CounterKey    string `xml:"counterkey,attr"`    //计数器key
}

type HttpServer struct {
	Enable       bool          `xml:"enable,attr"`
	Importers    []Importer    `xml:"importer"`
	Accumulators []Accumulator `xml:"accumulator"`
}

type Importer struct {
	Name       string `xml:"name,attr"`
	Enable     bool   `xml:"enable,attr"`
	ServerType string `xml:"servertype,attr"`
	ServerUrl  string `xml:"serverurl,attr"`
	ToQueue    string `xml:"toqueue,attr"`
}

type Accumulator struct {
	Name       string `xml:"name,attr"`
	Enable     bool   `xml:"enable,attr"`
	ServerType string `xml:"servertype,attr"`
	ServerUrl  string `xml:"serverurl,attr"`
	ToCounter  string `xml:"tocounter,attr"`
	ToSet      string `xml:"toset,attr"`
}

type UDPServer struct {
	Enable   bool          `xml:"enable,attr"`
	UDPPorts []UDPPortInfo `xml:"server"`
}

type UDPPortInfo struct {
	Enable         bool   `xml:"enable,attr"`
	Name           string `xml:"name,attr"`
	Port           int    `xml:"port,attr"`
	Protocol       string `xml:"protocol,attr"`
	Outputadapters string `xml:"outputadapters,attr"`
}

type OutputAdapter struct {
	Name    string `xml:"name,attr"`
	Type    string `xml:"type,attr"`
	Url     string `xml:"url,attr"`
	ToQueue string `xml:"toqueue,attr"`
}

// 检查是否存在触发器配置
func (t *TaskInfo) HasTrigger() bool {
	return t.TriggerServer != "" && t.TriggerQueue != ""
}

// 检查是否存在触发器过滤字段配置
func (t *TaskInfo) HasTriggerFilter() bool {
	return t.TriggerFilter != ""
}

// 检查是否存在写入目标过滤字段配置
func (t *TaskInfo) HasTargetFilter() bool {
	return t.TargetFilter != ""
}

// 检查是否存在计数器配置
func (t *TaskInfo) HasCounter() bool {
	return t.CounterServer != "" && t.CounterKey != ""
}
