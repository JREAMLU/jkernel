package services

import (
	"encoding/json"
	"fmt"
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
	Data DataParams `valid:"Required"`
}

type DataParams struct {
	Urls      []UrlsParams `json:"urls" valid:"Required"`
	Timestamp int64        `json:"timestamp" valid:"Required"`
	Sign      string       `json:"sign" valid:"Required"`
}

type UrlsParams struct {
	LongUrl string `json:"long_url" valid:"Required"`
	IP      string `json:"ip" valid:"IP"`
}

type dataList struct {
	Total int                    `json:"total"`
	List  map[string]interface{} `json:"list"`
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
func (r *Url) GoShorten(rawMetaHeader map[string][]string, rawDataBody []byte) (httpStatus int, output interface{}) {
	//将传递过来多json raw解析到struct
	ffjson.Unmarshal(rawDataBody, r)

	//日志
	beego.Trace("Url json解析:", r)

	//参数验证
	checked, err := inout.InputParamsCheck(rawMetaHeader, rawDataBody, &r.Data)
	if err != nil {
		return inout.EXPECTATION_FAILED, inout.OutputFail(checked.Message, "DATAPARAMSILLEGAL", checked.RequestID)
	}

	//进行shorten
	var list = make(map[string]interface{})
	var params_map = make(map[string]interface{})
	params := []map[string]interface{}{}
	for _, val := range r.Data.Urls {
		shortUrl := atom.GetShortenUrl(val.LongUrl, beego.AppConfig.String("ShortenDomain"))
		list[val.LongUrl] = shortUrl

		params_map["long_url"] = val.LongUrl
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

	var data dataList
	data.List = list
	data.Total = len(list)

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
				checked.MetaCheckResult["Request-Id"],
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
				checked.MetaCheckResult["Request-Id"],
			},
		},
	)

	res, err := async.GoAsyncRequest(addFunc, 2)
	fmt.Println(res, err)
	fmt.Println("=================================", "res: ", res["a"][0], "err: ", res["a"][1])
	fmt.Println("=================================", "res: ", res["b"][0], "err: ", res["a"][1])

	//持久化到mysql

	// a := models.GetUrlOne()
	// fmt.Println("AAAAAAAAAA", a["long_url"])

	return inout.OK, inout.OutputSuccess(data, checked.MetaCheckResult["Request-Id"])
}
