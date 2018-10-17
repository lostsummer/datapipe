package handlers

import (
	"TechPlat/datapipe/global"
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/util/redis"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/devfeel/dotweb"
)

type IncrData struct {
	Category string `json:"category"`
	AppID    string `json:"appid"`
	Key      string `json:"key"`
	GlobalID string `json:"globalid"`
	//Increase int    `json:"increase"` //将来去除，为上报端版本兼容暂时保留
	Uid  string `json:"uid"`
	Time int64  `json:"time"`
}

type EnterData struct {
	Category string `json:"category"`
	AppID    string `json:"appid"`
	Key      string `json:"key"`
	GlobalID string `json:"globalid"`
	Uid      string `json:"uid"`
	Time     int64  `json:"time"`
}

type QueueElem EnterData

func dateStr() string {
	return time.Now().Format("20060102")
}

func counterIncrBy(accConf *config.Accumulator, incr *IncrData) (int64, error) {
	serverUrl := accConf.ServerUrl
	category, appid, key := incr.Category, incr.AppID, incr.Key
	incrVal := 1
	redisField := fmt.Sprintf("%s:%s", appid, key)
	redisKey := fmt.Sprintf("%s:%s:%s", accConf.ToCounter, category, dateStr())
	redisClient := redisutil.GetRedisClient(serverUrl)
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
	return redisClient.HIncrBy(redisKey, redisField, incrVal)
}

func counterSet(accConf *config.Accumulator, enter *EnterData, val int64) error {
	serverUrl := accConf.ServerUrl
	category, appid, key := enter.Category, enter.AppID, enter.Key
	redisKey := fmt.Sprintf("%s:%s:%s", accConf.ToCounter, category, dateStr())
	redisField := fmt.Sprintf("%s:%s", appid, key)
	redisClient := redisutil.GetRedisClient(serverUrl)
	if redisClient == nil {
		return global.GetRedisError
	}
	defer func() {
		if p := recover(); p != nil {
			if s, ok := p.(string); ok {
				innerLogger.Error(s)
			}
		}
	}()
	return redisClient.HSet(redisKey, redisField, fmt.Sprintf("%d", val))
}

func addToSet(accConf *config.Accumulator, enter *EnterData) (int64, error) {
	serverUrl := accConf.ServerUrl
	category, appid, key, globalid := enter.Category, enter.AppID, enter.Key, enter.GlobalID
	redisKey := fmt.Sprintf("%s:%s:%s:%s:%s", accConf.ToSet, category, dateStr(), appid, key)
	redisClient := redisutil.GetRedisClient(serverUrl)
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

	err := redisClient.SAdd(redisKey, globalid)
	if err != nil {
		return -1, err
	}
	return redisClient.SCard(redisKey)
}
func PVCounter(ctx dotweb.Context) error {
	respstr := respFailed
	defer func() {
		innerLogger.Info("HttpServer::PVCounter")
		ctx.WriteString(respstr)
	}()

	accConf, err := getAccumulatorConf("PVCounter")
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::PVCounter " + err.Error())
		return nil
	}

	datajson := ctx.PostFormValue(postActionDataKey)
	if datajson == "" {
		innerLogger.Error("HttpServer::PVCounter " + global.LessParamError.Error())
		respstr = respFailed
		return nil
	}

	var incr IncrData
	err = json.Unmarshal([]byte(datajson), &incr)
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::PVCounter fail to parse post json incr data: " +
			err.Error() + "\r\n" + datajson + "\r\n")
		return nil
	} else if incr.Category == "" || incr.AppID == "" || incr.Key == "" {
		respstr = respFailed
		innerLogger.Error("HttpServer::PVCounter json data less fields:" +
			"\r\n" + datajson + "\r\n")
		return nil
	}

	ret, err := counterIncrBy(accConf, &incr)

	if ret > 0 && err == nil {
		respstr = strconv.FormatInt(ret, 10)
	} else {
		innerLogger.Error("HttpServer::PVCounter incr failed!")
		if err != nil {
			innerLogger.Error(err.Error())
		}
		respstr = respFailed
	}

	importerConf, err := getImporterConf("Counter")
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::PVCounter " + err.Error())
		return nil
	}

	qelem := getCounterQElem(&incr)
	if qelem == nil {
		innerLogger.Error("HttpServer::PVCounter incr -> qelem failed")
		return nil
	}
	if qelem.GlobalID == "" {
		return nil //前一版接口无此字段，兼容接收但不推队列
	}
	qelem.Time = getNowUnixSec() // 以服务器时间为准
	if data, err := json.Marshal(qelem); err != nil {
		innerLogger.Error(err.Error())
		respstr = respFailed
	} else {
		_, err = pushQueueData(importerConf, string(data))
		if err != nil {
			innerLogger.Error(err.Error())
		}
	}
	return nil
}

func UVCounter(ctx dotweb.Context) error {
	respstr := respFailed
	defer func() {
		innerLogger.Info("HttpServer::UVCounter")
		ctx.WriteString(respstr)
	}()

	accConf, err := getAccumulatorConf("UVCounter")
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::UVCounter " + err.Error())
		return nil
	}

	datajson := ctx.PostFormValue(postActionDataKey)
	if datajson == "" {
		innerLogger.Error("HttpServer::UVCounter " + global.LessParamError.Error())
		respstr = respFailed
		return nil
	}

	var enter EnterData
	err = json.Unmarshal([]byte(datajson), &enter)
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::UVCounter fail to parse post json enter data: " +
			err.Error() + "\r\n" + datajson + "\r\n")
		return nil
	} else if enter.Category == "" || enter.AppID == "" || enter.Key == "" || enter.GlobalID == "" {
		respstr = respFailed
		innerLogger.Error("HttpServer::UVCounter json data less fields:" +
			"\r\n" + datajson + "\r\n")
		return nil
	}

	scard, err := addToSet(accConf, &enter)

	if err != nil || scard < 0 {
		innerLogger.Error("HttpServer::UVCounter saddToSet failed!")
		if err != nil {
			innerLogger.Error(err.Error())
		}
		respstr = respFailed
		return nil
	}

	err = counterSet(accConf, &enter, scard)
	if err != nil {
		innerLogger.Error("HttpServer::UVCounter counterSet failed!")
		if err != nil {
			innerLogger.Error(err.Error())
		}
		respstr = respFailed
		return nil
	}

	//respstr = strconv.FormatInt(scard, 10)
	respstr = fmt.Sprintf("%d", scard)

	return nil
}

func getCounterQElem(data interface{}) *QueueElem {
	var e QueueElem
	if incr, ok := data.(*IncrData); ok {
		e = QueueElem{
			incr.Category,
			incr.AppID,
			incr.Key,
			incr.GlobalID,
			incr.Uid,
			incr.Time,
		}
		return &e
	} else if enter, ok := data.(*EnterData); ok {
		e = QueueElem{
			enter.Category,
			enter.AppID,
			enter.Key,
			enter.GlobalID,
			enter.Uid,
			enter.Time,
		}
		return &e
	} else {
		return nil
	}
}
