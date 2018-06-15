package impl

import (
	"testing"
	"strconv"
)

func TestMongoHandler_InsertJsonData(t *testing.T) {
	for i:=0;i<100;i++ {
		mongoHandler := new(MongoHandler)
		mongoHandler.SetConn("mongodb://192.168.8.178:27017/emoney_tongji", "emoney_tongji")
		//insert data to mongo
		err := mongoHandler.InsertJsonData("capped_tongji_userlog", "{\"UserName\":\"111\"}")
		if err != nil {
			t.Error("InsertJsonData " + strconv.Itoa(i) + " error -> " + err.Error())
		} else {
			t.Log("InsertJsonData " + strconv.Itoa(i) + " success ->")
		}
	}
}
