package services

import (
	"encoding/json"
	"time"

	"github.com/JREAMLU/core/async"
	"github.com/JREAMLU/core/inout"
	"github.com/JREAMLU/core/sign"
	"github.com/JREAMLU/jkernel/base/services/atom"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/pquerna/ffjson/ffjson"
)

type Url struct {
	Data struct {
		Urls []struct {
			LongURL string `json:"long_url" valid:"Required"`
			IP      string `json:"ip" valid:"IP"`
		} `json:"urls" valid:"Required"`
		Timestamp int    `json:"timestamp" valid:"Required"`
		Sign      string `json:"sign" valid:"Required"`
	} `json:"data" valid:"Required"`
}

func GetParams(url Url) Url {
	return url
}

/**
 *	自定义多valid
 */
func (r *Url) Valid(v *validation.Validation) {}

/**
 *	@auther		jream.lu
 *	@intro		原始链接=>短链接
 *	@logic
 *	@todo		参数验证抽出去
 *	@params		params []byte	参数
 *	@return 	httpStatus lice
 */
func (r *Url) GoShorten(data map[string]interface{}) (httpStatus int, output inout.Output) {
	//将传递过来多json raw解析到struct
	ffjson.Unmarshal(data["body"].([]byte), r)

	//参数验证
	checked, err := inout.InputParamsCheck(data, &r.Data)
	if err != nil {
		return inout.EXPECTATION_FAILED, inout.OutputFail(
			checked.Message,
			"DATAPARAMSILLEGAL",
			checked.RequestID,
		)
	}

	//进行shorten
	var list = make(map[string]interface{})
	var params_map = make(map[string]interface{})
	params := []map[string]interface{}{}
	for _, val := range r.Data.Urls {
		shortUrl := atom.GetShortenUrl(val.LongURL, beego.AppConfig.String("ShortenDomain"))
		list[val.LongURL] = shortUrl

		params_map["long_url"] = val.LongURL
		params_map["short_url"] = shortUrl
		params_map["long_crc"] = 1
		params_map["short_crc"] = 1
		params_map["status"] = 1
		params_map["created_by_ip"] = 123
		params_map["updated_by_ip"] = 123
		params_map["created_at"] = 456
		params_map["updated_at"] = 456
		params = append(params, params_map)
	}

	var datalist atom.DataList
	datalist.List = list
	datalist.Total = len(list)

	//请求其他接口
	var requestParams = make(map[string]interface{})
	var rdata = make(map[string]interface{})
	var urls []map[string]string
	timestamp := time.Now().Unix()
	urls = append(urls, map[string]string{"long_url": "http://o9d.cn", "IP": "127.0.0.1"})
	urls = append(urls, map[string]string{"long_url": "http://huiyimei.net", "IP": "192.168.1.1"})
	rdata["urls"] = urls
	rdata["timestamp"] = timestamp
	requestParams["data"] = rdata
	raw, _ := json.Marshal(requestParams)
	sign := sign.GenerateSign(raw, timestamp, beego.AppConfig.String("sign.secretKey"))
	requestParams["sign"] = sign
	// atom.RequestGetAyiName(requestParams, timestamp)

	var addFunc async.MultiAddFunc
	addFunc = append(
		addFunc,
		async.AddFunc{
			Name:    "a",
			Handler: atom.RequestGetAyiName,
			Params: []interface{}{
				requestParams,
				timestamp,
				checked.CheckRes["Request-Id"],
			},
		},
	)
	addFunc = append(
		addFunc,
		async.AddFunc{
			Name:    "b",
			Handler: atom.RequestGetAyiName,
			Params: []interface{}{
				requestParams,
				timestamp,
				checked.CheckRes["Request-Id"],
			},
		},
	)

	async.GoAsyncRequest(addFunc, 2)

	//持久化到mysql
	beego.Trace(inout.Rid + ":" + "持久化")

	return inout.OK, inout.OutputSuccess(
		datalist,
		checked.RequestID,
	)
}
