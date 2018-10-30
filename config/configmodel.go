package config

import (
	"TechPlat/datapipe/queue"
	"encoding/xml"
)

const (
	Target_MongoDB = "mongodb"
	Target_Http    = "http"
	Target_Kafka   = "kafka"
	Target_Redis   = "redis"
)

//代理配置信息
type AppConfig struct {
	XMLName           xml.Name        `xml:"config"`
	Log               Log             `xml:"log"`
	MongoDBs          []MongoDB       `xml:"mongodbs>mongodb"`
	Redises           []Redis         `xml:"redises>redis"`
	Kafkas            []Kafka         `xml:"kafkas>kafka"`
	HTTPs             []HTTP          `xml:"https>http"`
	HttpServer        HttpServer      `xml:"httpserver"`
	Tasks             []TaskInfo      `xml:"tasks>task"`
	UdpServer         UDPServer       `xml:"udpserver"`
	OutputAdapters    []OutputAdapter `xml:"outputadapter>adapter"`
	ImporterMap       map[string]*Importer
	AccumulatorMap    map[string]*Accumulator
	TaskMap           map[string]*TaskInfo
	RedisMap          map[string]*Redis
	MongoDBMap        map[string]*MongoDB
	KafkaMap          map[string]*Kafka
	HTTPMap           map[string]*HTTP
	ImporterTargetMap map[string]queue.Target
	TaskSourceMap     map[string]queue.Source
	TaskTargetMap     map[string]queue.Target
	TaskTriggerMap    map[string]queue.Target
}

//Redis配置
type Redis struct {
	ID  string `xml:"id,attr"`
	URL string `xml:"url,attr"`
	DB  string `xml:"db,attr"`
}

//MongodDB配置
type MongoDB Redis

//kafka配置
type Kafka struct {
	ID  string `xml:"id,attr"`
	URL string `xml:"url,attr"`
}

type HTTP Kafka

//log配置
type Log struct {
	FilePath string `xml:"filepath,attr"`
}

//数据模板
type TaskInfo struct {
	ID      string       `xml:"id,attr"`
	Enable  bool         `xml:"enable,attr"`
	Source  Queue        `xml:"source"`
	Target  Queue        `xml:"target"`
	Trigger Queue        `xml:"trigger"`
	Counter Task_Counter `xml:"counter"`
}

type Queue struct {
	Type   string `xml:"type,attr"`
	ID     string `xml:"id,attr"`
	Queue  string `xml:"queue,attr"`
	Filter string `xml:"filter,attr"`
}

type Task_Counter struct {
	Type string `xml:"type,attr"`
	ID   string `xml:"id,attr"`
	Key  string `xml:"key,attr"`
}

type HttpServer struct {
	Enable       bool          `xml:"enable,attr"`
	Importers    []Importer    `xml:"importer"`
	Accumulators []Accumulator `xml:"accumulator"`
}

type Importer struct {
	ID     string `xml:"id,attr"`
	Enable bool   `xml:"enable,attr"`
	Target Queue  `xml:"target"`
}

type Accumulator struct {
	ID     string             `xml:"id,attr"`
	Enable bool               `xml:"enable,attr"`
	Target Accumulator_Target `xml:"target"`
}

type Accumulator_Target struct {
	Type    string `xml:"type,attr"`
	ID      string `xml:"id,attr"`
	Counter string `xml:"counter,attr"`
	Set     string `xml:"set,attr"`
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
	return t.Trigger.Queue != ""
}

// 检查是否存在触发器过滤字段配置
func (t *TaskInfo) HasTriggerFilter() bool {
	return t.Trigger.Filter != ""
}

// 检查是否存在写入目标过滤字段配置
func (t *TaskInfo) HasTargetFilter() bool {
	return t.Target.Filter != ""
}

// 检查是否存在计数器配置
func (t *TaskInfo) HasCounter() bool {
	return t.Counter.Key != ""
}

func (tc *Task_Counter) GetServer() string {
	r, ok := CurrentConfig.RedisMap[tc.ID]
	if !ok {
		return ""
	}
	return r.URL
}
