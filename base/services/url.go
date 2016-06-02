package services

import (

	//"encoding/json"

	"base/services/atom"
	"fmt"

	"core/inout"

	"github.com/pquerna/ffjson/ffjson"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

type Url struct {
	Data DataParams `valid:"Required"`
}

type DataParams struct {
	Urls []UrlsParams `json:"urls" valid:"Required"`
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
	var u Url
	ffjson.Unmarshal(rawDataBody, &u)
	// ffjson.Unmarshal(rawDataBody, &u.Data.Urls)

	//日志
	fmt.Println("Url json解析:", u)

	//参数验证
	checked, err := inout.InputParamsCheck(rawMetaHeader, &u.Data)
	if err != nil {
		return inout.EXPECTATION_FAILED, inout.OutputFail(checked.Message, "DATAPARAMSILLEGAL", checked.RequestID)
	}

	//进行shorten
	var list = make(map[string]interface{})
	var params_map = make(map[string]interface{})
	params := []map[string]interface{}{}
	for _, val := range u.Data.Urls {
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

	//持久化到mysql

	// a := models.GetUrlOne()
	// fmt.Println("AAAAAAAAAA", a["long_url"])

	return inout.OK, inout.OutputSuccess(data, checked.MetaCheckResult["Request-Id"])
}
