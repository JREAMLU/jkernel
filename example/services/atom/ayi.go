package atom

import (
	"encoding/json"

	"github.com/JREAMLU/core/curl"
)

func RequestGetAyiName(requestParams map[string]interface{}, timestamp int64, requestID string) (string, error) {
	rawSign, _ := json.Marshal(requestParams)

	res, err := curl.RollingCurl(
		curl.Requests{
			Method: "POST",
			UrlStr: "http://localhost/study/rest/put.php",
			Header: GetHeader(requestID),
			Raw:    string(rawSign),
		},
	)
	if err != nil {
		return "", err
	}

	return res, nil
}
