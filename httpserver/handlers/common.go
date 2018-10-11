package handlers

import (
	"TechPlat/datapipe/global"
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
)

const (
	respSucceed        = "1"
	respFailed         = "-1"
	timeLayout         = "2006-01-02 15:04:05"
	tjgidCookieName    = "em_tongji_cookie_globalid"
	tjfvtCookieName    = "em_tongji_cookie_firstvisittime"
	comgidCookieName   = "tongji_globalid"
	tjCookieDomain     = "tongji.emoney.cn"
	comCookieDomain    = "emoney.cn"
	cookieValidSeconds = 311040000
	postActionDataKey  = "ActionData"
)

var (
	innerLogger *logger.InnerLogger

	agentOS = map[string]string{
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

	agentBrowser = map[string]string{
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

	/* 可获取客户端IP地址的http头
	** 后三个可能在nginx给php-fpm中反向代理中有添加
	** 但在nginx做负载均衡时并没有
	** 实测nginx做负载均衡时主要靠前两个
	 */
	clientIPHeader = []string{
		"X-Forwarded-For",
		"X-Real-IP",
		"Proxy-Client-IP",
		"WL-Proxy-Client-IP",
		"HTTP_X_FORWARDED_FOR",
		"HTTP_CLIENT_IP",
		"REMOTE_ADDR",
	}
)

func init() {
	innerLogger = logger.GetInnerLogger()
}

func getImporterConf(name string) (*config.Importer, error) {
	impMap := config.CurrentConfig.ImporterMap
	impConf, exist := impMap[name]
	if exist {
		return impConf, nil
	} else {
		return nil, global.NotConfigError
	}
}

func getAccumulatorConf(name string) (*config.Accumulator, error) {
	accMap := config.CurrentConfig.AccumulatorMap
	accConf, exist := accMap[name]
	if exist {
		return accConf, nil
	} else {
		return nil, global.NotConfigError
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
	md5str := strings.ToUpper(
		getMd5String(base64.URLEncoding.EncodeToString(b)))

	return strings.Join(
		[]string{
			md5str[0:8],
			md5str[8:12],
			md5str[12:16],
			md5str[16:20],
			md5str[20:32],
		}, "-")

}

func getGlobalID(ctx dotweb.Context) string {
	var gid string
	if cookie, err := ctx.ReadCookie(tjgidCookieName); err == nil {
		gid = cookie.Value
	} else if cookie, err := ctx.ReadCookie(comgidCookieName); err == nil {
		gid = cookie.Value
	} else {
		gid = createGuid()
		SetCookieNameValueDomainPath(ctx, tjgidCookieName, gid, tjCookieDomain, "/")
		SetCookieNameValueDomainPath(ctx, comgidCookieName, gid, comCookieDomain, "/")
	}
	return gid
}

func SetCookieNameValueDomainPath(ctx dotweb.Context, name string, value string, domain string, path string) {
	cookie := &http.Cookie{
		Name:    name,
		Value:   value,
		Expires: time.Now().Add(cookieValidSeconds * time.Second),
		Domain:  domain,
		Path:    path,
	}
	ctx.SetCookie(cookie)
}

func getFirstVistTime(ctx dotweb.Context) string {
	var fvt string
	if cookie, err := ctx.ReadCookie(tjfvtCookieName); err == nil {
		fvt, _ = URLUnescape(cookie.Value)
	} else {
		fvt = time.Now().Format(timeLayout)
		SetCookieNameValueDomainPath(ctx, tjfvtCookieName, URLEscape(fvt), tjCookieDomain, "/")
	}
	return fvt
}

func getNowFormatTime() string {
	return time.Now().Format(timeLayout)
}

func getNowUnixSec() int64 {
	return time.Now().Unix()
}

func getClientIP(ctx dotweb.Context) string {
	for _, name := range clientIPHeader {
		ip := ctx.Request().Header.Get(name)
		if ip != "" {
			// 经过多级反向代理 X-Forwarded-For 的值可能是多个逗号分隔的IP地址
			// 取第一个
			if strings.Index(ip, ",") > -1 {
				ip = strings.Split(ip, ",")[0]
			}
			return ip
		}
	}
	return ctx.Request().RemoteIP()
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
		return -1, global.GetRedisError
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

// 因为php框架对写入cookie的字串Escape处理(主要原因是cookie无法直接存中文字符,
// 其实对于我们写入的时间格式完全不必要),
// 为了兼容线上以及写入的客户端cookie，也需要对first_visit_time进行escape,
// 以下部分代码扒自标准库net/url并做了精简
func URLEscape(s string) string {
	spaceCount, hexCount := 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c) {
			if c == ' ' {
				spaceCount++
			} else {
				hexCount++
			}
		}
	}

	if spaceCount == 0 && hexCount == 0 {
		return s
	}

	t := make([]byte, len(s)+2*hexCount)
	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == ' ':
			t[j] = '+'
			j++
		case shouldEscape(c):
			t[j] = '%'
			t[j+1] = "0123456789ABCDEF"[c>>4]
			t[j+2] = "0123456789ABCDEF"[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}

func shouldEscape(c byte) bool {
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return false
	}
	switch c {
	case '-', '_', '.', '~': // §2.3 Unreserved characters (mark)
		return false

	case '$', '&', '+', ',', '/', ':', ';', '=', '?', '@': // §2.2 Reserved characters (reserved)
		return true

	}
	return true
}

func URLUnescape(s string) (string, error) {
	// Count %, check that they're well-formed.
	n := 0
	hasPlus := false
	for i := 0; i < len(s); {
		switch s[i] {
		case '%':
			n++
			if i+2 >= len(s) || !ishex(s[i+1]) || !ishex(s[i+2]) {
				s = s[i:]
				if len(s) > 3 {
					s = s[:3]
				}
				return "", global.EscapeError
			}
			i += 3
		case '+':
			hasPlus = true
			i++
		default:
			i++
		}
	}

	if n == 0 && !hasPlus {
		return s, nil
	}

	t := make([]byte, len(s)-2*n)
	j := 0
	for i := 0; i < len(s); {
		switch s[i] {
		case '%':
			t[j] = unhex(s[i+1])<<4 | unhex(s[i+2])
			j++
			i += 3
		case '+':
			t[j] = ' '
			j++
			i++
		default:
			t[j] = s[i]
			j++
			i++
		}
	}
	return string(t), nil
}

func unhex(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}

func ishex(c byte) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	}
	return false
}
