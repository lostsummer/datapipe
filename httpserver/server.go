package httpserver

import (
	"TechPlat/datapipe/config"
	"net/http"

	"github.com/devfeel/dotweb"
)

var srv *(dotweb.DotWeb) = nil

func StartServer() error {
	//初始化DotServer
	srv = dotweb.New()

	//设置日志目录

	srv.SetLogPath(config.CurrentConfig.Log.FilePath)
	srv.SetEnabledLog(true)
	srv.HttpServer.SetEnabledGzip(true)

	//设置路由
	InitRoute(srv)
	port := config.CurrentConfig.HttpServer.HttpPort
	err := srv.StartServer(port)
	if err == http.ErrServerClosed {
		return nil
	} else {
		return err
	}
}

func RestartServer() {
	if srv != nil {
		srv.Close()
	}
}
