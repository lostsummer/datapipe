package handlers

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
	"time"

	"TechPlat/datapipe/config"
	"TechPlat/datapipe/util/log"
	"TechPlat/datapipe/util/redis"

	"github.com/devfeel/dotweb"
	"github.com/pkg/errors"
)

const (
	respSucceed        = "1"
	respFailed         = "-1"
	timeLayout         = "2006-01-02 03:04:05"
	tjgidCookieName    = "em_tongji_cookie_globalid"
	tjfvtCookieName    = "em_tongji_cookie_firstvisittime"
	comgidCookieName   = "tongji_globalid"
	tjCookieDomain     = "tongji.emoney.cn"
	comCookieDomain    = "emoney.cn"
	cookieValidSeconds = 311040000
	postActionDataKey  = "ActionData"
)

var (
	NotConfigError = errors.New("not exists such config info")
	LessParamError = errors.New("less param")
	GetRedisError  = errors.New("get rediscli failed")
	innerLogger    *logger.InnerLogger
	agentOS        map[string]string = map[string]string{
		"Windows CE":     "Windows CE",
		"iPhone":         "iPhone",
		"Android":        "Android",
		"Windows NT 10":  "Windows 10",
		"Windows NT 6.1": "Windows 7",
		"Windows NT 6.0": "Windows Vista",
		"Windows NT 5.2": "Windows 2003",
		"Windows NT 5.1": "Windows XP",
		"Windows NT 5.0": "Windows 2000",
	}
	agentBrowser map[string]string = map[string]string{
		"GreenBrowser":    "GreenBrowser",
		"NetCaptor":       "NetCaptor",
		"TencentTraveler": "TencentTraveler",
		"TheWorld":        "TheWorld",
		"MAXTHON":         "MAXTHON",
		"MyIE":            "MyIE",
		"MSIE 10":         "IE 10",
		"MSIE 9":          "IE 9",
		"MSIE 8":          "IE 8",
		"MSIE 7":          "IE 7",
		"MSIE 6":          "IE 6",
		"MSIE 5.5":        "IE 5.5",
		"Netscape":        "Netscape",
		"Chrome":          "Chrome",
		"Firefox":         "Firefox",
		//"Safari":         "Safari",       因为现在Chrome的UA中包含"Safari",所以对于Safari要特殊判断
		"Opera": "Opera",
		"R4EA":  "R4EA",
		"UP":    "UP",
		"UCWEB": "UCWEB",
	}
)

func init() {
	innerLogger = logger.GetInnerLogger()
}

func getImporterConf(name string) (*config.Importer, error) {
	impMap := config.CurrentConfig.ImporterMap
	importerInfo, exist := impMap[name]
	if exist {
		return importerInfo, nil
	} else {
		return nil, NotConfigError
	}
}

func getMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func createGuid() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	md5str := getMd5String(base64.URLEncoding.EncodeToString(b))
	return strings.Join([]string{
		strings.ToUpper(md5str[0:8]),
		strings.ToUpper(md5str[8:12]),
		strings.ToUpper(md5str[12:16]),
		strings.ToUpper(md5str[16:20]),
		strings.ToUpper(md5str[20:32])}, "-")

}

func getGlobalID(ctx dotweb.Context) string {
	var gid string
	if cookie, err := ctx.ReadCookie(tjgidCookieName); err == nil {
		gid = cookie.Value
	} else if cookie, err := ctx.ReadCookie(comgidCookieName); err == nil {
		gid = cookie.Value
	} else {
		gid = createGuid()
		SetCookieNameValueDomain(ctx, tjgidCookieName, gid, tjCookieDomain)
		SetCookieNameValueDomain(ctx, comgidCookieName, gid, comCookieDomain)
	}
	return gid
}

func SetCookieNameValueDomain(ctx dotweb.Context, name string, value string, domain string) {
	cookie := &http.Cookie{
		Name:    name,
		Value:   value,
		Expires: time.Now().Add(cookieValidSeconds * time.Second),
		Domain:  domain,
	}
	ctx.SetCookie(cookie)
}

func getFirstVistTime(ctx dotweb.Context) string {
	var fvt string
	if cookie, err := ctx.ReadCookie(tjfvtCookieName); err == nil {
		fvt = cookie.Value
	} else {
		fvt = time.Now().Format(timeLayout)
		cookie = &http.Cookie{
			Name:    tjfvtCookieName,
			Value:   fvt,
			Expires: time.Now().Add(cookieValidSeconds * time.Second),
			Domain:  tjCookieDomain,
		}
		ctx.SetCookie(cookie)
	}
	return fvt
}

func getNowFormatTime() string {
	return time.Now().Format(timeLayout)
}

func getClientIP(ctx dotweb.Context) string {
	if ip := ctx.Request().Header.Get("Remote_addr"); ip == "" {
		return ctx.Request().RemoteAddr
	} else {
		return ip
	}
}

func getFullUrl(ctx dotweb.Context) string {
	return "http://" + ctx.Request().Host + ctx.Request().Url()
}

func getUserAgent(ctx dotweb.Context) string {
	return ctx.Request().UserAgent()
}

func getAgentOS(ua string) string {
	for k, v := range agentOS {
		if strings.Contains(ua, k) {
			return v
		}
	}
	return "other"
}

func getAgentBrowser(ua string) string {
	if strings.Contains(ua, "Safari") && !strings.Contains(ua, "Chrome") {
		return "Safari"
	}
	for k, v := range agentBrowser {
		if strings.Contains(ua, k) {
			return v
		}
	}
	return "other"
}

// pushQueueData push data to redis queue
func pushQueueData(importerConf *config.Importer, val string) (int64, error) {
	server, queue := importerConf.ServerUrl, importerConf.ToQueue
	return pushQueueDataToSQ(server, queue, val)
}

// args expose server and queue
func pushQueueDataToSQ(server, queue, val string) (int64, error) {
	redisClient := redisutil.GetRedisClient(server)
	if redisClient == nil {
		return -1, GetRedisError
	}
	defer func() {
		if p := recover(); p != nil {
			if s, ok := p.(string); ok {
				innerLogger.Error(s)
			}
		}
	}()
	return redisClient.LPush(queue, val)
}
