package atom

import (
	"encoding/json"

	"github.com/JREAMLU/core/curl"
	"github.com/JREAMLU/core/sign"
	"github.com/astaxie/beego"
)

func RequestGetAyiName(requestParams map[string]interface{}, timestamp int64) (string, error) {
	raw, _ := json.Marshal(requestParams)
	sign := sign.GenerateSign(raw, timestamp, beego.AppConfig.String("sign.secretKey"))
	requestParams["sign"] = sign
	rawSign, _ := json.Marshal(requestParams)

	res, err := curl.RollingCurl(
		curl.Requests{
			Method: "POST",
			UrlStr: "http://localhost/study/rest/put.php",
			Header: map[string]string{
				"Content-Type": "application/json;charset=utf-8;",
				"Accept":       "application/json",
			},
			Raw: string(rawSign),
		},
	)
	if err != nil {
		return "", err
	}

	return res, nil
}
