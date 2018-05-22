package httpserver

import (
	"TechPlat/datapipe/httpserver/handlers"

	"github.com/devfeel/dotweb"
)

func InitRoute(server *dotweb.DotWeb) {
	server.HttpServer.GET("/", handlers.Index)
	server.HttpServer.GET("/page/pageclick", handlers.PageClick)
	server.HttpServer.GET("/page/pageview", handlers.PageView)
	server.HttpServer.GET("/data/webdata", handlers.WebData)
	server.HttpServer.GET("/data/appdata", handlers.AppData)
	server.HttpServer.GET("/paylog/paylog", handlers.PayLog)
	server.HttpServer.GET("/userlog/userlog", handlers.UserLog)
	server.HttpServer.GET("/soft", handlers.Soft)
	server.HttpServer.POST("/soft/actionlog", handlers.SoftActionLog)
	server.HttpServer.POST("/frontend/log", handlers.FrontEndLog)
	server.HttpServer.GET("/liveduration/data", handlers.LiveDuration)
}
