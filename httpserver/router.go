package httpserver

import (
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/httpserver/handlers"

	"strings"

	"github.com/devfeel/dotweb"
)

//记录路由handler关联
type routeInfo struct {
	HttpMethod  func(string, dotweb.HttpHandle)
	Route       string
	HandlerFunc dotweb.HttpHandle
}

//根据routeInfo结构信息绑定路由handler
func (this *routeInfo) bound() {
	this.HttpMethod(this.Route, this.HandlerFunc)
	this.HttpMethod(strings.ToLower(this.Route), this.HandlerFunc)
	return
}

//根据importer配置路由初化
func InitRoute(server *dotweb.DotWeb) {
	//默认主页
	server.HttpServer.GET("/", handlers.Index)
	server.HttpServer.GET("/test", handlers.Test)

	//根据importer, accumulator开关绑定route
	imptMap := config.CurrentConfig.ImporterMap
	accMap := config.CurrentConfig.AccumulatorMap

	//importer_name : routeInfo
	//原php程序全方法支持, 另外发现web客户端发数据存在方法乱用问题
	//(例如对 /soft 用 POST) 所以对路由全方法支持
	var imptIDRouteMap = map[string]routeInfo{
		"PageClick":     {server.HttpServer.Any, "/Page/PageClick", handlers.PageClick},
		"PageView":      {server.HttpServer.Any, "/Page/PageView", handlers.PageView},
		"ADView":        {server.HttpServer.Any, "/Page/AdView", handlers.ADView},
		"ADClick":       {server.HttpServer.Any, "/Page/ADClick", handlers.ADClick},
		"WebData":       {server.HttpServer.Any, "/Data/WebData", handlers.WebData},
		"AppData":       {server.HttpServer.Any, "/Data/AppData", handlers.AppData},
		"PayLog":        {server.HttpServer.Any, "/PayLog/PayLog", handlers.PayLog},
		"UserLog":       {server.HttpServer.Any, "/UserLog/UserLog", handlers.UserLog},
		"Soft":          {server.HttpServer.Any, "/Soft", handlers.Soft},
		"SoftEncrypt":   {server.HttpServer.Any, "/Soft/EncryptInfo", handlers.SoftEncrypt},
		"SoftActionLog": {server.HttpServer.Any, "/Soft/ActionLog", handlers.SoftActionLog},
		"ActLog":        {server.HttpServer.Any, "/ActLog", handlers.SoftActionLog},
		"FrontEndLog":   {server.HttpServer.Any, "/FrontEnd/Log", handlers.FrontEndLog},
		"LiveDuration":  {server.HttpServer.Any, "/LiveDuration/Data", handlers.LiveDuration},
		"PageRecords":   {server.HttpServer.Any, "/Mobile/Page", handlers.PageRecordsHandle},
		"EventRecords":  {server.HttpServer.Any, "/Mobile/Event", handlers.EventRecordsHandle},
		"LoginInfo":     {server.HttpServer.Any, "/Mobile/LoginInfo", handlers.LoginInfoHandle},
		"AsoRank":       {server.HttpServer.Any, "/Mobile/AsoRank", handlers.AsoRankHandle},
		"HPSoftLog":     {server.HttpServer.Any, "/Gaojia/SoftLog", handlers.HPSoftLogHandle},
	}
	var accIDRounteMap = map[string]routeInfo{
		"PVCounter": {server.HttpServer.Any, "/Counter/PV", handlers.PVCounter},
		"UVCounter": {server.HttpServer.Any, "/Counter/UV", handlers.UVCounter},
	}
	for id := range imptMap {
		if route, exist := imptIDRouteMap[id]; exist {
			route.bound()
		}
	}
	for id := range accMap {
		if route, exist := accIDRounteMap[id]; exist {
			route.bound()
		}
	}
}
