package queue

import (
	"TechPlat/datapipe/repository/impl"
)

type MongoDBTarget struct {
	URL        string
	DB         string
	Collection string
}

func (m *MongoDBTarget) Push(val string) (int64, error) {
	mongoHandler := new(impl.MongoHandler)
	mongoHandler.SetConn(m.URL, m.DB)
	//insert data to mongo
	err := mongoHandler.InsertJsonData(m.Collection, val)
	if err != nil {
		return -1, err
	}
	return 1, nil
}
