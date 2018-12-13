package handlers

import (
	"TechPlat/datapipe/global"
	"encoding/json"
	"fmt"

	"github.com/devfeel/dotweb"
)

type PageRecord struct {
	PageID    string `json:"pageid"`
	PageData  string `json:"pagedata"`
	EnterTime int64  `json:"entertime"`
	QuitTime  int64  `json:"quittime"`
}

type PageRecords struct {
	AppID   string       `json:"appid"`
	GUID    string       `json:"guid"`
	UserID  int64        `json:"userid"`
	Records []PageRecord `json:"records"`
}

type PageRecToQ struct {
	AppID     string `json:"appid"`
	GUID      string `json:"guid"`
	UserID    int64  `json:"userid"`
	PageID    string `json:"pageid"`
	PageData  string `json:"pagedata"`
	EnterTime int64  `json:"entertime"`
	QuitTime  int64  `json:"quittime"`
}

type EventRecord struct {
	PageID     string `json:"pageid"`
	PageData   string `json:"pagedata"`
	EventID    string `json:"eventid"`
	EventData  string `json:"eventdata"`
	ActionTime int64  `json:"actiontime"`
}

type EventRecords struct {
	AppID   string        `json:"appid"`
	GUID    string        `json:"guid"`
	UserID  int64         `json:"userid"`
	Records []EventRecord `json:"records"`
}

type EventRecToQ struct {
	AppID      string `json:"appid"`
	GUID       string `json:"guid"`
	UserID     int64  `json:"userid"`
	PageID     string `json:"pageid"`
	PageData   string `json:"pagedata"`
	EventID    string `json:"eventid"`
	EventData  string `json:"eventdata"`
	ActionTime int64  `json:"actiontime"`
}

type LoginInfo struct {
	AppID      string `json:"appid"`
	GUID       string `json:"guid"`
	UserID     int64  `json:"userid"`
	Event      string `json:"event"`
	LoginTime  int64  `json:"logintime"`
	LogoutTime int64  `json:"logouttime"`
}

const (
	postRecsKey      = "records"
	postLoginInfoKey = "logininfo"
)

func PageRecordsHandle(ctx dotweb.Context) error {
	respstr := respFailed
	defer func() {
		innerLogger.Info("HttpServer::PageRecords")
		ctx.WriteString(respstr)
	}()

	inputJson := ctx.PostFormValue(postRecsKey)
	if inputJson == "" {
		innerLogger.Error("HttpServer::PageRecords " + global.LessParamError.Error())
		return nil
	}

	var pageRecords PageRecords
	err := json.Unmarshal([]byte(inputJson), &pageRecords)
	if err != nil {
		innerLogger.Error("HttpServer::PageRecords fail to parse post json actionData: " +
			err.Error() + "\r\n" + inputJson + "\r\n")
		return nil
	}

	appid := pageRecords.AppID
	guid := pageRecords.GUID
	userid := pageRecords.UserID

	for _, rec := range pageRecords.Records {
		var dataToQ PageRecToQ
		dataToQ.AppID = appid
		dataToQ.GUID = guid
		dataToQ.UserID = userid
		dataToQ.PageID = rec.PageID
		dataToQ.PageData = rec.PageData
		dataToQ.EnterTime = rec.EnterTime
		dataToQ.QuitTime = rec.QuitTime
		if outputJson, err := json.Marshal(dataToQ); err != nil {
			innerLogger.Error("HttpServer::PageRecords " + err.Error())
			return nil
		} else {
			target, err := getImporterTarget("PageRecords")
			if err != nil {
				panic(err)
				return nil
			}
			_, err = target.Push(string(outputJson))
			if err != nil {
				innerLogger.Error("HttpServer::PageRecords push queue data failed!")
				innerLogger.Error(err.Error())
				return nil
			}
		}
	}
	respstr = fmt.Sprintf("%d", len(pageRecords.Records))
	return nil
}

func EventRecordsHandle(ctx dotweb.Context) error {
	respstr := respFailed
	defer func() {
		innerLogger.Info("HttpServer::EventRecords")
		ctx.WriteString(respstr)
	}()

	inputJson := ctx.PostFormValue(postRecsKey)
	if inputJson == "" {
		innerLogger.Error("HttpServer::EventRecords " + global.LessParamError.Error())
		return nil
	}

	var eventRecords EventRecords
	err := json.Unmarshal([]byte(inputJson), &eventRecords)
	if err != nil {
		innerLogger.Error("HttpServer::EventRecords fail to parse post json actionData: " +
			err.Error() + "\r\n" + inputJson + "\r\n")
		return nil
	}

	appid := eventRecords.AppID
	guid := eventRecords.GUID
	userid := eventRecords.UserID

	for _, rec := range eventRecords.Records {
		var dataToQ EventRecToQ
		dataToQ.AppID = appid
		dataToQ.GUID = guid
		dataToQ.UserID = userid
		dataToQ.PageID = rec.PageID
		dataToQ.PageData = rec.PageData
		dataToQ.EventID = rec.EventID
		dataToQ.EventData = rec.EventData
		dataToQ.ActionTime = rec.ActionTime
		if outputJson, err := json.Marshal(dataToQ); err != nil {
			innerLogger.Error("HttpServer::EventRecords " + err.Error())
			return nil
		} else {
			target, err := getImporterTarget("EventRecords")
			if err != nil {
				panic(err)
				return nil
			}
			_, err = target.Push(string(outputJson))
			if err != nil {
				innerLogger.Error("HttpServer::EventRecord push queue data failed!")
				innerLogger.Error(err.Error())
				return nil
			}
		}
	}
	respstr = fmt.Sprintf("%d", len(eventRecords.Records))
	return nil
}

func LoginInfoHandle(ctx dotweb.Context) error {
	respstr := respFailed
	defer func() {
		innerLogger.Info("HttpServer::LoginInfoHandle")
		ctx.WriteString(respstr)
	}()

	inputJson := ctx.PostFormValue(postLoginInfoKey)
	if inputJson == "" {
		innerLogger.Error("HttpServer::LoginInfo " + global.LessParamError.Error())
		return nil
	}

	var loginInfo LoginInfo
	err := json.Unmarshal([]byte(inputJson), &loginInfo)
	if err != nil {
		innerLogger.Error("HttpServer::LoginInfo fail to parse post json actionData: " +
			err.Error() + "\r\n" + inputJson + "\r\n")
		return nil
	}

	appid := loginInfo.AppID
	guid := loginInfo.GUID
	userid := loginInfo.UserID
	event := loginInfo.Event
	logintime := loginInfo.LoginTime
	logouttime := loginInfo.LogoutTime

	dataToQ := make(map[string]interface{})
	dataToQ["appid"] = appid
	dataToQ["guid"] = guid
	dataToQ["userid"] = userid
	dataToQ["event"] = event
	dataToQ["logintime"] = logintime
	if event == "logout" && logouttime > 0 {
		dataToQ["logouttime"] = logouttime
	}

	if outputJson, err := json.Marshal(dataToQ); err != nil {
		innerLogger.Error("HttpServer::LoginInfo " + err.Error())
		return nil
	} else {
		target, err := getImporterTarget("LoginInfo")
		if err != nil {
			panic(err)
			return nil
		}
		_, err = target.Push(string(outputJson))
		if err != nil {
			innerLogger.Error("HttpServer::LoginInfo push queue data failed!")
			innerLogger.Error(err.Error())
			return nil
		}
	}

	respstr = fmt.Sprintf("%d", 0)
	return nil
}
