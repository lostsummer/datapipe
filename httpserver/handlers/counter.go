package handlers

import (
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
	Increase int    `json:"increase"`
	Time     int64  `json:"time"`
}

type IPData struct {
	Category string `json:"category"`
	AppID    string `json:"appid"`
	Key      string `json:"key"`
	IP       string `json:"ip"`
	Time     int64  `json:"time"`
}

func dateStr() string {
	return time.Now().Format("20060102")
}

func counterIncrBy(accConf *config.Accumulator, incr *IncrData) (int64, error) {
	serverUrl := accConf.ServerUrl
	category, appid, key, val := incr.Category, incr.AppID, incr.Key, incr.Increase
	redisField := fmt.Sprintf("%s:%s", appid, key)
	redisKey := fmt.Sprintf("%s:%s:%s", accConf.ToCounter, category, dateStr())
	redisClient := redisutil.GetRedisClient(serverUrl)
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
	return redisClient.HIncrBy(redisKey, redisField, val)
}

func counterSet(accConf *config.Accumulator, ipData *IPData, val int64) error {
	serverUrl := accConf.ServerUrl
	category, appid, key := ipData.Category, ipData.AppID, ipData.Key
	redisKey := fmt.Sprintf("%s:%s:%s", accConf.ToCounter, category, dateStr())
	redisField := fmt.Sprintf("%s:%s", appid, key)
	redisClient := redisutil.GetRedisClient(serverUrl)
	if redisClient == nil {
		return GetRedisError
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

func addToSet(accConf *config.Accumulator, ipData *IPData) (int64, error) {
	serverUrl := accConf.ServerUrl
	category, appid, key, ip := ipData.Category, ipData.AppID, ipData.Key, ipData.IP
	redisKey := fmt.Sprintf("%s:%s:%s:%s:%s", accConf.ToSet, category, dateStr(), appid, key)
	redisClient := redisutil.GetRedisClient(serverUrl)
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

	err := redisClient.SAdd(redisKey, ip)
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
		innerLogger.Error("HttpServer::PVCounter " + LessParamError.Error())
		respstr = respFailed
		return nil
	}

	var incrData IncrData
	err = json.Unmarshal([]byte(datajson), &incrData)
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::PVCounter fail to parse post json incrData: " +
			err.Error() + "\r\n" + datajson + "\r\n")
		return nil
	}

	ret, err := counterIncrBy(accConf, &incrData)

	if ret > 0 && err == nil {
		respstr = strconv.FormatInt(ret, 10)
	} else {
		innerLogger.Error("HttpServer::PVCounter incr failed!")
		if err != nil {
			innerLogger.Error(err.Error())
		}
		respstr = respFailed
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
		innerLogger.Error("HttpServer::UVCounter " + LessParamError.Error())
		respstr = respFailed
		return nil
	}

	var ipData IPData
	err = json.Unmarshal([]byte(datajson), &ipData)
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::UVCounter fail to parse post json incrData: " +
			err.Error() + "\r\n" + datajson + "\r\n")
		return nil
	}

	scard, err := addToSet(accConf, &ipData)

	if err != nil || scard < 0 {
		innerLogger.Error("HttpServer::UVCounter saddToSet failed!")
		if err != nil {
			innerLogger.Error(err.Error())
		}
		respstr = respFailed
		return nil
	}

	err = counterSet(accConf, &ipData, scard)
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
