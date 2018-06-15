package impl

import (
	"TechPlat/datapipe/const/log"
	"TechPlat/datapipe/util/log"
	"encoding/json"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sync"
)

var mgoSessionPool *sync.Map
const MaxSessionPoolLimit = 50

func init(){
	mgoSessionPool = new(sync.Map)
}

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
	session, err := getSessionCopy(monogDBConn)
	if err != nil {
		logger.Log("repository.impl.MongoHandler::Dial["+monogDBConn+"]失败 - "+err.Error(), logdefine.LogTarget_MongoDB, logdefine.LogLevel_Error)
		return err
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

// getSessionCopy get seesion copy with conn from pool
func getSessionCopy(conn string) (*mgo.Session, error){
	data, isOk:=mgoSessionPool.Load(conn)
	if isOk{
		session, isSuccess := data.(mgo.Session)
		if isSuccess{
			return session.Clone(), nil
		}
	}

	session, err := mgo.Dial(conn)
	if err!= nil{
		logger.Log("repository.impl.getSessionCopy::Dial["+conn+"]失败 - "+err.Error(), logdefine.LogTarget_MongoDB, logdefine.LogLevel_Error)
		return nil, err
	}else{
		session.SetPoolLimit(MaxSessionPoolLimit)
		mgoSessionPool.Store(conn, session)
		return session.Clone(), nil
	}
}