package httpserver

import (
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/httpserver/handlers"

	"github.com/devfeel/dotweb"
)

//记录路由handler关联
type routeInfo struct {
	HttpMethod  string
	Route       string
	HandlerFunc func(dotweb.Context) error
}

//根据routeInfo结构信息绑定路由handler
func (this *routeInfo) bound(server *dotweb.DotWeb) {
	route := this.Route
	handlerFunc := this.HandlerFunc
	switch this.HttpMethod {
	case "GET":
		server.HttpServer.GET(route, handlerFunc)
	case "POST":
		server.HttpServer.POST(route, handlerFunc)
	}
	return
}

//importer_name : routeInfo
var routeMap map[string]routeInfo = map[string]routeInfo{
	"PageClick":     routeInfo{"GET", "/page/pageclick", handlers.PageClick},
	"PageView":      routeInfo{"GET", "/page/pageview", handlers.PageView},
	"WebData":       routeInfo{"GET", "/data/webdata", handlers.WebData},
	"AppData":       routeInfo{"GET", "/data/appdata", handlers.AppData},
	"PayLog":        routeInfo{"GET", "/paylog/paylog", handlers.PayLog},
	"UserLog":       routeInfo{"GET", "/userlog/userlog", handlers.UserLog},
	"Soft":          routeInfo{"GET", "/soft", handlers.Soft},
	"SoftActionLog": routeInfo{"POST", "/soft/actionlog", handlers.SoftActionLog},
	"FrontEndLog":   routeInfo{"POST", "/frontend/log", handlers.FrontEndLog},
	"LiveDuration":  routeInfo{"GET", "/liveduration/data", handlers.LiveDuration},
}

//根据importer配置路由初化
func InitRoute(server *dotweb.DotWeb) {
	//默认主页
	server.HttpServer.GET("/", handlers.Index)

	//根据importer开关绑定route
	importers := config.CurrentConfig.HttpServer.Importers
	for _, importerInfo := range importers {
		if importerInfo.Enable {
			if route, exist := routeMap[importerInfo.Name]; exist {
				route.bound(server)
			}
		}
	}
}
