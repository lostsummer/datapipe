package httpserver

import (
	"TechPlat/datapipe/config"

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
	return err
}

func RestartServer() error {
	if srv != nil {
		srv.Close()
	}
	return StartServer()
}
