package services

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"

	"github.com/JREAMLU/core/inout"
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

	GenerateSign(rawDataBody, 1466490032, "kkkkk")

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

	// var batch int64
	// batch = time.Now().Unix()
	// fmt.Println("=======================", batch)

	//持久化到mysql

	// a := models.GetUrlOne()
	// fmt.Println("AAAAAAAAAA", a["long_url"])

	return inout.OK, inout.OutputSuccess(data, checked.MetaCheckResult["Request-Id"])
}

func GenerateSign(requestData []byte, requestTime int64, secretKey string) string {
	//ksort
	var rdata map[string]interface{}
	json.Unmarshal([]byte(requestData), &rdata)

	sorted_keys := make([]string, 0)
	for k, _ := range rdata {
		sorted_keys = append(sorted_keys, k)
	}

	sort.Strings(sorted_keys)

	// var str string
	//
	// for _, k := range sorted_keys {
	// 	fmt.Printf("k=%v, v=%v\n", k, rdata[k])
	// 	str = str + k + rdata[k]
	// }

	// fmt.Println("==================", str)

	str := Serialize(rdata)
	fmt.Println("rrrrrr", str)

	return ""
}

func Serialize(data interface{}) interface{} {
	fmt.Println("----", data, "---", reflect.TypeOf(data).Kind())
	var str string
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(data)
		for i := 0; i < s.Len(); i++ {
			serial := Serialize(s.Index(i).Interface())
			if reflect.TypeOf(serial).Kind() == reflect.Float64 {
				serial = strconv.Itoa(int(serial.(float64)))
			}
			str = str + strconv.Itoa(i) + serial.(string)
		}
		return str
	case reflect.Map:
		s := reflect.ValueOf(data)

		//ksort
		// sorted_keys := make([]string, 0)
		// for k, _ := range rdata {
		// 	sorted_keys = append(sorted_keys, k)
		// }
		// sort.Strings(sorted_keys)
		// for _, k := range sorted_keys {
		// 	fmt.Printf("k=%v, v=%v\n", k, rdata[k])
		// 	str = str + k + rdata[k]
		// }

		keys := s.MapKeys()
		//ksort
		sorted_keys := make([]string, 0)
		for _, key := range keys {
			sorted_keys = append(sorted_keys, key.Interface().(string))
		}
		sort.Strings(sorted_keys)
		for _, key := range sorted_keys {
			serial := Serialize(s.MapIndex(reflect.ValueOf(key)).Interface())
			if reflect.TypeOf(serial).Kind() == reflect.Float64 {
				serial = strconv.Itoa(int(serial.(float64)))
			}
			str = str + key + serial.(string)
		}
		//     for _, key := range keys {
		//         serial := Serialize(s.MapIndex(reflect.ValueOf(key.String())).Interface(), true)
		//         if reflect.TypeOf(serial).Kind() == reflect.Float64 {
		//             serial = strconv.Itoa(int(serial.(float64)))
		//         }
		//         str = str + key.String() + serial.(string)
		//     }
		// }
		return str
	}

	fmt.Println("pppppp")
	// fmt.Println(data)
	// switch data.(type) {
	// case []interface{}:
	// case map[string]interface{}:
	// 	for k, v := range data {
	//
	// 	}
	// }
	return data
}
