package handlers

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"net/http"
	"time"

	"TechPlat/datapipe/config"
	"TechPlat/datapipe/util/log"
	"TechPlat/datapipe/util/redis"

	"github.com/devfeel/dotweb"
	"github.com/pkg/errors"
)

const (
	respSucceed  = "1"
	respFailed   = "-1"
	timeLayout   = "2006-01-02 03:04:05"
	gidCookieKey = "em_tongji_cookie_globalid"
	fvtCookieKey = "em_tongji_cookie_firstvisittime"
	postDataKey  = "ActionData"
)

var (
	innerLogger *logger.InnerLogger
)

var (
	NotConfigError = errors.New("not exists such config info")
	LessParamError = errors.New("less param")
	GetRedisError  = errors.New("get rediscli failed")
)

func init() {
	innerLogger = logger.GetInnerLogger()
}

func getImporterInfo(id string) (*config.ImporterInfo, error) {
	impMap := config.CurrentConfig.ImporterMap
	importerInfo, exist := impMap[id]
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
	return getMd5String(base64.URLEncoding.EncodeToString(b))
}

func getGlobalID(ctx dotweb.Context) string {
	var gid string
	if cookie, err := ctx.ReadCookie(gidCookieKey); err == nil {
		gid = cookie.Value
	} else {
		gid = createGuid()
		cookie = &http.Cookie{
			Name:    gidCookieKey,
			Value:   gid,
			Expires: time.Now().Add(311040000),
			Domain:  "tongji.emoney.cn",
		}
		ctx.SetCookie(cookie)
	}
	return gid
}

func getFirstVistTime(ctx dotweb.Context) string {
	var fvt string
	if cookie, err := ctx.ReadCookie(fvtCookieKey); err == nil {
		fvt = cookie.Value
	} else {
		fvt = time.Now().Format(timeLayout)
		cookie = &http.Cookie{
			Name:    fvtCookieKey,
			Value:   fvt,
			Expires: time.Now().Add(311040000),
			Domain:  "tongji.emoney.cn",
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

// pushQueueData push data to redis queue
func pushQueueData(importerConf *config.ImporterInfo, val string) (int64, error) {
	server := importerConf.ToServer
	queue := importerConf.ToQueue
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
