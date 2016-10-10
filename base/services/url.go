package services

import (
	io "github.com/JREAMLU/core/inout"
	"github.com/JREAMLU/jkernel/base/models/mredis"
	"github.com/JREAMLU/jkernel/base/services/atom"
	"github.com/JREAMLU/jkernel/base/services/entity"
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

func (r *Url) Valid(v *validation.Validation) {}

func (r *Url) GoShorten(data map[string]interface{}) (httpStatus int, output io.Output) {
	ffjson.Unmarshal(data["body"].([]byte), r)

	ch, err := io.InputParamsCheck(data, &r.Data)
	if err != nil {
		return io.EXPECTATION_FAILED, io.Fail(
			ch.Message,
			"DATAPARAMSILLEGAL",
			ch.RequestID,
		)
	}

	list := shorten(r)

	var datalist entity.DataList
	datalist.List = list
	datalist.Total = len(list)

	return io.OK, io.Suc(
		datalist,
		ch.RequestID,
	)
}

func shorten(r *Url) map[string]interface{} {
	list := make(map[string]interface{})
	params_map := make(map[string]interface{})
	params := []map[string]interface{}{}

	for _, val := range r.Data.Urls {
		shortUrl := atom.GetShortenUrl(val.LongURL)

		short, err := setCache(val.LongURL, shortUrl)
		if err != nil {
			beego.Trace("setCache error: ", err)
		}

		list[val.LongURL] = beego.AppConfig.String("ShortenDomain") + short

		params_map["long_url"] = val.LongURL
		params_map["short_url"] = short
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
	beego.Trace(io.Rid + ":" + "持久化")

	return list
}

func setCache(origin string, short string) (string, error) {
	reply, err := mredis.ShortenHGet(origin)
	if err != nil && err.Error() != "redigo: nil returned" {
		return "", err
	}

	if reply == "" {
		_, err = mredis.ShortenHSet(origin, short)
		if err != nil {
			return "", err
		}
		mredis.ExpandHSet(short, origin)
		return short, nil
	}

	return reply, nil
}

func GoExpand() {

}
