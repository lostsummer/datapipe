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
//180615 update: 原php程序全方法支持, 另外发现web客户端发数据存在方法乱用问题
//(例如对 /soft 用 POST) 所以无奈对路由全方法支持
func (this *routeInfo) bound(server *dotweb.DotWeb) {
	route := this.Route
	handlerFunc := this.HandlerFunc
	switch this.HttpMethod {
	case "GET":
		//server.HttpServer.GET(route, handlerFunc)
		fallthrough
	case "POST":
		//server.HttpServer.POST(route, handlerFunc)
		fallthrough
	default:
		server.HttpServer.Any(route, handlerFunc)
	}
	return
}

//importer_name : routeInfo
var routeMap map[string]routeInfo = map[string]routeInfo{
	"PageClick":     {"GET", "/page/pageclick", handlers.PageClick},
	"PageView":      {"GET", "/page/pageview", handlers.PageView},
	"ADView":        {"GET", "/page/adview", handlers.ADView},
	"ADClick":       {"GET", "/page/adclick", handlers.ADClick},
	"WebData":       {"GET", "/data/webdata", handlers.WebData},
	"AppData":       {"POST", "/data/appdata", handlers.AppData},
	"PayLog":        {"GET", "/paylog/paylog", handlers.PayLog},
	"UserLog":       {"POST", "/userlog/userlog", handlers.UserLog},
	"Soft":          {"GET", "/soft", handlers.Soft},
	"SoftActionLog": {"POST", "/soft/actionlog", handlers.SoftActionLog},
	"ActLog":        {"POST", "/actlog", handlers.SoftActionLog}, // same as SoftActionLog, i don't know why
	"FrontEndLog":   {"POST", "/frontend/log", handlers.FrontEndLog},
	"LiveDuration":  {"GET", "/liveduration/data", handlers.LiveDuration},
}

//没解决URL大小写敏感的问题, 客户端脚本中都是大小写错落
var routeMap2 map[string]routeInfo = map[string]routeInfo{
	"PageClick":     {"GET", "/Page/PageClick", handlers.PageClick},
	"PageView":      {"GET", "/Page/PageView", handlers.PageView},
	"ADView":        {"GET", "/Page/AdView", handlers.ADView},
	"ADClick":       {"GET", "/Page/ADClick", handlers.ADClick},
	"WebData":       {"GET", "/Data/WebData", handlers.WebData},
	"AppData":       {"POST", "/Data/AppData", handlers.AppData},
	"PayLog":        {"GET", "/PayLog/PayLog", handlers.PayLog},
	"UserLog":       {"POST", "/UserLog/UserLog", handlers.UserLog},
	"Soft":          {"GET", "/Soft", handlers.Soft},
	"SoftActionLog": {"POST", "/Soft/ActionLog", handlers.SoftActionLog},
	"ActLog":        {"POST", "/ActLog", handlers.SoftActionLog},
	"FrontEndLog":   {"POST", "/FrontEnd/Log", handlers.FrontEndLog},
	"LiveDuration":  {"GET", "/LiveDuration/Data", handlers.LiveDuration},
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
			if route, exist := routeMap2[importerInfo.Name]; exist {
				route.bound(server)
			}
		}
	}
}
