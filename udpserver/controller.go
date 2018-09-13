package udpserver

import (
	"strconv"
	"TechPlat/datapipe/udpserver/outputadapter"
	"TechPlat/datapipe/udpserver/protocol"
	"TechPlat/datapipe/util/log"
	"TechPlat/datapipe/config"
)

const (
	timeLayout = "2006-01-02 15:04:05.999999"
)

const (
	PROTOCOL_TXT  = "TXT"
	PROTOCOL_JSON = "JSON"
)

const (
	ADAPTER_FILE  = "FILE"
	ADAPTER_REDIS = "REDIS"
)

var (
	innerLogger *logger.InnerLogger
	protocolHandlerMap map[string]protocol.ProtocolHandler
	outputAdapterMap map[string]outputadapter.OutputAdapter
)

func init() {
	innerLogger = logger.GetInnerLogger()

	protocolHandlerMap = map[string]protocol.ProtocolHandler {
		PROTOCOL_TXT  : protocol.OnHandleTextProtocol,
		PROTOCOL_JSON : protocol.OnHandleJsonProtocol,
	}

	outputAdapterMap = map[string]outputadapter.OutputAdapter {
		ADAPTER_FILE  : outputadapter.OutputFileAdapter,
		ADAPTER_REDIS : outputadapter.OutputRedisAdapter,
	}
}

func StartServer() {
	if !config.CurrentConfig.UdpServer.Enable {
		innerLogger.Debug("未配置UDP日志服务")
		return
	}

	if len(config.CurrentConfig.UdpServer.UDPPorts) <= 0{
		innerLogger.Debug("未配置UDP端口")
		return
	}

	outputAdapters := initOutputAdapter()
	if len(outputAdapters) <= 0 {
		innerLogger.Debug("未配置输出适配器")
		return
	}

	innerLogger.Debug("开始启动UDP日志服务")
	initUdpServer(outputAdapters)
	innerLogger.Debug("完成启动UDP日志服务")
}

func initOutputAdapter() []Adapter {
	var adapters []Adapter
	for _, adpConf := range config.CurrentConfig.OutputAdapters {
		if !adpConf.Enable {
			continue
		}

		if _,ok:=outputAdapterMap[adpConf.Type]; !ok {
			continue
		}

		adapter := Adapter{outputAdapterMap[adpConf.Type], adpConf}
		adapters = append(adapters, adapter)
	}
	return adapters
}

func initUdpServer(outputAdapters []Adapter) {
	for _, portConf := range config.CurrentConfig.UdpServer.UDPPorts {
		if !portConf.Enable {
			innerLogger.Debug(portConf.Name + ", UDP端口未开启")
			continue
		}

		if portConf.Port <= 0 {
			innerLogger.Debug(portConf.Name + ", UDP端口必须大于0")
			continue
		}

		if _,ok:=protocolHandlerMap[portConf.Protocol]; !ok {
			innerLogger.Debug(portConf.Name + ", 未知的处理协议")
			continue
		}

		protocolHandler := protocolHandlerMap[portConf.Protocol]
		server, err := GetNewServer(portConf.Port, protocolHandler, outputAdapters)
		if err != nil {
			innerLogger.Debug(portConf.Name + ", 创建UDP日志服务实例失败, " + err.Error())
			continue
		}

		server.Start()
		innerLogger.Debug("启动UDP日志服务:" + portConf.Name + ", port=" + strconv.Itoa(portConf.Port) + ", protocol=" + portConf.Protocol)
	}
}