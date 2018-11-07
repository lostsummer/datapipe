package endpoint

import (
	"TechPlat/datapipe/util/http"
)

type Http struct {
	URL string
}

func (h *Http) Push(val string) (int64, error) {
	//retBody, _, _, httpErr := httputil.HttpPost(h.URL, val, "")
	_, _, _, httpErr := httputil.HttpPost(h.URL, val, "")
	if httpErr != nil {
		return -1, httpErr
	}
	return 1, nil
}
