package impl

import (
	"emoney/tongjiservice/const/log"
	"emoney/tongjiservice/util/log"
	"encoding/json"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoHandler struct {
}

var (
	monogDBConn string
	defaultDB   string
)

func (da *MongoHandler) SetConn(conn, dbName string) {
	monogDBConn = conn
	defaultDB = dbName
}

/*新增Json数据插入指定的Collection
* Author: Panxinming
* LastUpdateTime: 2016-10-17 10:00
* 如果失败，则返回具体的error，成功则返回nil
 */
func (da *MongoHandler) InsertJsonData(collectionName, jsonData string) error {
	session, err := mgo.Dial(monogDBConn)
	if err != nil {
		logger.Log("repository.impl.MongoHandler::Dial["+monogDBConn+"]失败 - "+err.Error(), logdefine.LogTarget_MongoDB, logdefine.LogLevel_Error)
		return err
		//panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(defaultDB).C(collectionName)
	var f interface{}
	data := []byte(jsonData)
	//特殊处理，需要先json序列化，成为标准json，再进一步转为bosn，才能成功
	err_jsonunmar := json.Unmarshal(data, &f)
	if err_jsonunmar != nil {
		logger.Log("repository.impl.MongoHandler::InsertJsonData::json.Unmarshal["+collectionName+"]["+jsonData+"]失败 - "+err_jsonunmar.Error(), logdefine.LogTarget_MongoDB, logdefine.LogLevel_Error)
		return err_jsonunmar
	}
	bjsondata, err_jsonmar := json.Marshal(f)
	if err_jsonmar != nil {
		logger.Log("repository.impl.MongoHandler::InsertJsonData::json.Marshal["+collectionName+"]["+jsonData+"]失败 - "+err_jsonmar.Error(), logdefine.LogTarget_MongoDB, logdefine.LogLevel_Error)
		return err_jsonmar
	}

	var bf interface{}
	bsonerr := bson.UnmarshalJSON(bjsondata, &bf)
	if bsonerr != nil {
		logger.Log("repository.impl.MongoHandler::InsertJsonData::bson.UnmarshalJSON["+collectionName+"]["+jsonData+"]失败 - "+bsonerr.Error(), logdefine.LogTarget_MongoDB, logdefine.LogLevel_Error)
		return bsonerr
	}
	err = c.Insert(&bf)
	if err != nil {
		logger.Log("repository.impl.MongoHandler::InsertJsonData["+collectionName+"]["+jsonData+"]失败 - "+err.Error(), logdefine.LogTarget_MongoDB, logdefine.LogLevel_Error)
		return err
	}
	return nil
}
