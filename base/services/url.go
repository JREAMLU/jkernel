package services

import (
	"fmt"
	"strings"

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

type UrlExpand struct {
	Shorten []string `json:"shorten" valid:"Required"`
}

func GetParams(url Url) Url {
	return url
}

func (r *Url) Valid(v *validation.Validation) {}

func (r *Url) GoShorten(data map[string]interface{}) (httpStatus int, output io.Output) {
	ffjson.Unmarshal(data["body"].([]byte), r)

	ch, err := io.InputParamsCheck(data, &r.Data)
	if err != nil {
		return io.BAD_REQUEST, io.Fail(
			ch.Message,
			"DATAPARAMSILLEGAL",
			ch.RequestID,
		)
	}

	list := shorten(r)

	var datalist entity.DataList
	datalist.List = list
	datalist.Total = len(list)

	return io.CREATED, io.Suc(
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

		short, err := setDB(val.LongURL, shortUrl)
		if err != nil {
			beego.Trace("setDB error: ", err)
		}

		list[val.LongURL] = beego.AppConfig.String("ShortenDomain") + short

		atom.Mu.Lock()
		params_map["long_url"] = val.LongURL
		params_map["short_url"] = short
		params_map["long_crc"] = 1
		params_map["short_crc"] = 1
		params_map["status"] = 1
		params_map["created_by_ip"] = 123
		params_map["updated_by_ip"] = 123
		params_map["created_at"] = 456
		params_map["updated_at"] = 456
		atom.Mu.Unlock()
		params = append(params, params_map)
	}

	//持久化到mysql
	beego.Trace(io.Rid + ":" + "持久化")

	return list
}

func setDB(origin string, short string) (string, error) {
	reply, err := mredis.ShortenHGet(origin)
	if err != nil && err.Error() != "redigo: nil returned" {
		return "", err
	}
	if reply == "" {
		//@todo mysql
		_, err = mredis.ShortenHSet(origin, short)
		if err != nil {
			return "", err
		}
		mredis.ExpandHSet(short, origin)
		return short, nil
	}
	return reply, nil
}

func (r *Url) GoExpand(data map[string]interface{}) (httpStatus int, output io.Output) {
	var ue UrlExpand
	ffjson.Unmarshal([]byte(data["querystrjson"].(string)), &ue)

	ch, err := io.InputParamsCheck(data, ue)
	if err != nil {
		return io.BAD_REQUEST, io.Fail(
			ch.Message,
			"DATAPARAMSILLEGAL",
			ch.RequestID,
		)
	}

	list := expand(&ue)

	var datalist entity.DataList
	datalist.List = list
	datalist.Total = len(list)

	return io.OK, io.Suc(
		datalist,
		ch.RequestID,
	)
}

func expand(ue *UrlExpand) map[string]interface{} {
	list := make(map[string]interface{})
	shortens := ue.Shorten[0]
	for _, shorten := range strings.Split(shortens, ",") {
		reply, err := mredis.ExpandHGet(shorten)
		if err != nil {
			beego.Trace("expand error: ", err)
		}
		fmt.Println("<<<<<", reply)
		atom.Mu.Lock()
		list[shorten] = reply
		atom.Mu.Unlock()
	}
	return list
}
