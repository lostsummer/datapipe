package endpoint

import (
	"testing"
)

var postStr = `{"App":"10168","ClickData":"6","ClickKey":"month","ClickRemark":"null","ClientIP":"192.168.6.95","FirstVisitTime":"2018-10-30 11:29:03","GlobalID":"6C46EEBB-6223-CA39-C115-C49E8A10002D","HtmlType":"null","Module":"summaryreport","PageUrl":"http://192.168.8.215/account-analysis-ui/m/summaryreport","Remark":"","UserAgent":"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36","WriteTime":"2018-10-30 11:29:03"}`

func Test_Push(t *testing.T) {
	h := HttpTarget{"http://192.168.8.215/EMoney.JG.OpenApi/api/Tongji/PageClick"}
	_, err := h.Push(postStr)
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log("sucess to post data")
	}
}
