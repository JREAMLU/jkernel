package services

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/JREAMLU/core/async"
	io "github.com/JREAMLU/core/inout"
	"github.com/JREAMLU/core/sign"
	"github.com/JREAMLU/jkernel/example/services/atom"
	"github.com/JREAMLU/jkernel/example/services/entity"
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
func (r *Url) GoShorten(data map[string]interface{}) (httpStatus int, output io.Output) {
	//将传递过来多json raw解析到struct
	ffjson.Unmarshal(data["body"].([]byte), r)

	//参数验证
	ch, err := io.InputParamsCheck(data, &r.Data)
	if err != nil {
		return http.StatusExpectationFailed, io.Fail(
			ch.Message,
			"DATAPARAMSILLEGAL",
			ch.RequestID,
		)
	}

	//shorten && indb
	list := shorten(r)

	var datalist entity.DataList
	datalist.List = list
	datalist.Total = len(list)

	//请求其他接口
	request(data)

	return http.StatusAccepted, io.Suc(
		datalist,
		ch.RequestID,
	)
}

func shorten(r *Url) map[string]interface{} {
	list := make(map[string]interface{})
	params_map := make(map[string]interface{})
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

	//持久化到mysql
	beego.Info(io.Jctx.Value("requestID").(string) + ":" + "持久化")

	return list
}

func request(data map[string]interface{}) {
	token := data["headermap"].(http.Header)["Token"][0]
	ts, _ := strconv.ParseInt(data["headermap"].(http.Header)["Timestamp"][0], 10, 64)
	beego.Info(token, ts)
	requestParams := make(map[string]interface{})
	rdata := make(map[string]interface{})
	urls := []map[string]string{}
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
			Logo:    "a",
			Handler: atom.RequestGetAyiName,
			Params: []interface{}{
				requestParams,
				timestamp,
				io.Jctx.Value("requestID").(string),
			},
		},
	)
	addFunc = append(
		addFunc,
		async.AddFunc{
			Logo:    "b",
			Handler: atom.RequestGetAyiName,
			Params: []interface{}{
				requestParams,
				timestamp,
				io.Jctx.Value("requestID").(string),
			},
		},
	)

	async.GoAsyncRequest(addFunc, 2)

}
