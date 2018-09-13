package udpserver

import (
	"net"
	"strconv"
	"TechPlat/datapipe/udpserver/protocol"
	"TechPlat/datapipe/udpserver/outputadapter"
	"TechPlat/datapipe/config"
)

type Adapter struct {
	Adapter outputadapter.OutputAdapter
	Conf config.OutputAdapter
}

type Server struct {
	conn *net.UDPConn 						  //UDP连接
	logs chan protocol.LogInfo 				  //日志消息
	protocolHandler protocol.ProtocolHandler  //协议处理函数
	outputAdapters []Adapter				  //输出适配器
}

func GetNewServer(udpPort int, protocolHandler protocol.ProtocolHandler, outputAdapters []Adapter) (Server, error) {
	var currentUdpPort string
	if udpPort > 0 {
		currentUdpPort = ":" + strconv.Itoa(udpPort)
	}

	var s Server
	udpAddr, err := net.ResolveUDPAddr("udp4", currentUdpPort)
	if err == nil {
		s.protocolHandler = protocolHandler
		s.outputAdapters = outputAdapters
		s.logs = make(chan protocol.LogInfo, 100000)
		s.conn, err = net.ListenUDP("udp", udpAddr)
	}
	return s, err
}

func (s *Server) Start() {
	go func(){
		for {
			s.handleLog()
		}
	}()

	go func(){
		for {
			s.readLog()
		}
	}()
}

func (s *Server) readLog() {
	var buf [4096]byte
	n, addr, err := s.conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}

	if (n>0 && len(buf)>0){
		log := s.protocolHandler(buf, n)
		log.Clientip = addr.IP.String()

		s.logs <- log
	}
}

func (s *Server) handleLog() {
	log := <-s.logs
	logStr := "[" + log.Clientip + "] [" + log.Appid + "] [" + log.Level + "] [" + log.Message + "]"

	for _, adp := range s.outputAdapters {
		adp.Adapter(adp.Conf, log.Appid, logStr)
	}
}