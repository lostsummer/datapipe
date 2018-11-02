package httpserver

import (
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/util/log"
	"net/http"

	"github.com/devfeel/dotweb"
	dotconfig "github.com/devfeel/dotweb/config"
)

const (
	RunEnv_Flag       = "RunEnv"
	RunEnv_Develop    = "develop"
	RunEnv_Test       = "test"
	RunEnv_Production = "production"
)

const (
	msgSrvNotConfig = "notconfig"
	msgSrvClosed    = "closed"
	msgSrvError     = "error"
)

var (
	RunEnv     string
	ConfigPath string
)

var (
	innerLogger *logger.InnerLogger
	srv         *(dotweb.DotWeb)
	wait        chan int
)

func init() {
	innerLogger = logger.GetInnerLogger()
	srv = nil
	wait = make(chan int)
}

func StartServer() {
	srvStatus := make(chan string)
	runSrv := func() {
		if config.CurrentConfig.HttpServer.Enable == false {
			srv = nil
			srvStatus <- msgSrvNotConfig
			return
		}
		srvConfig := dotconfig.MustInitConfig(ConfigPath + "/dotweb.conf")
		srv = dotweb.ClassicWithConf(srvConfig)
		srv.UseRequestLog()
		if RunEnv == RunEnv_Develop {
			srv.SetDevelopmentMode()
		}
		InitRoute(srv)
		innerLogger.Debug("httpserver.StartServer => ")
		if err := srv.Start(); err == http.ErrServerClosed {
			innerLogger.Debug("httpserver.StartServer recieve close signal")
			srvStatus <- msgSrvClosed
		} else {
			innerLogger.Error("httpserver.StartServer " + err.Error())
			srvStatus <- msgSrvError
		}
		return
	}

	go runSrv()

	for s := range srvStatus {
		switch s {
		case msgSrvClosed:
			go runSrv()
		case msgSrvError:
			return
		case msgSrvNotConfig:
			innerLogger.Debug("httpserver.StartServer httpserver.enable=false, but also block")
			<-wait
			go runSrv()
		}
	}
}

func RestartServer() {
	innerLogger.Debug("httpserver.RestartServer")
	if srv == nil { //上一次启动没有配置打开
		wait <- 1
	} else { //已经打开
		srv.Close()
	}
}
