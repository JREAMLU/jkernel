package services

import (
	jcontext "context"
	"hash/crc32"
	"net/http"
	"strings"
	"time"

	"github.com/JREAMLU/core/com"
	io "github.com/JREAMLU/core/inout"
	"github.com/JREAMLU/jkernel/base/models/mentity"
	"github.com/JREAMLU/jkernel/base/models/mmysql"
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

var ip string

func GetParams(url Url) Url {
	return url
}

func (r *Url) Valid(v *validation.Validation) {}

func (r *Url) GoShorten(jctx jcontext.Context, data map[string]interface{}) (httpStatus int, output io.Output) {
	ffjson.Unmarshal(data["body"].([]byte), r)
	ip = data["headermap"].(http.Header)["X-Forwarded-For"][0]
	ch, err := io.InputParamsCheck(jctx, data, &r.Data)
	if err != nil {
		return http.StatusExpectationFailed, io.Fail(
			ch.Message,
			"DATAPARAMSILLEGAL",
			ch.RequestID,
		)
	}

	list := shorten(jctx, r)

	var datalist entity.DataList
	datalist.List = list
	datalist.Total = len(list)

	return http.StatusCreated, io.Suc(
		datalist,
		ch.RequestID,
	)
}

func shorten(jctx jcontext.Context, r *Url) map[string]interface{} {
	list := make(map[string]interface{})

	for _, val := range r.Data.Urls {
		shortUrl := atom.GetShortenUrl(val.LongURL)

		short, err := setDB(jctx, val.LongURL, shortUrl)
		if err != nil {
			beego.Info(jctx.Value("requestID").(string), ":", "setDB error: ", err)
		}

		list[val.LongURL] = beego.AppConfig.String("ShortenDomain") + short
	}

	return list
}

func setDB(jctx jcontext.Context, origin string, short string) (string, error) {
	reply, err := mredis.ShortenHGet(origin)
	if err != nil && err.Error() != "redigo: nil returned" {
		return "", err
	}
	if reply == "" {
		var redirect mentity.Redirect
		redirect.LongUrl = origin
		redirect.ShortUrl = short
		redirect.LongCrc = uint64(crc32.ChecksumIEEE([]byte(origin)))
		redirect.ShortCrc = uint64(crc32.ChecksumIEEE([]byte(short)))
		redirect.Status = 1
		redirect.CreatedByIP = uint64(com.Ip2Int(ip))
		redirect.UpdateByIP = uint64(com.Ip2Int(ip))
		redirect.CreateAT = uint64(time.Now().Unix())
		redirect.UpdateAT = uint64(time.Now().Unix())

		_, err := mmysql.ShortenIn(redirect)
		if err != nil {
			beego.Error(jctx.Value("requestID").(string), ":", "setDB error: ", err)
			return "", err
		}

		_, err = mredis.ShortenHSet(origin, short)
		if err != nil {
			return "", err
		}
		mredis.ExpandHSet(short, origin)
		return short, nil
	}
	return reply, nil
}

func (r *Url) GoExpand(jctx jcontext.Context, data map[string]interface{}) (httpStatus int, output io.Output) {
	var ue UrlExpand
	ffjson.Unmarshal([]byte(data["querystrjson"].(string)), &ue)

	ch, err := io.InputParamsCheck(jctx, data, ue)
	if err != nil {
		return http.StatusExpectationFailed, io.Fail(
			ch.Message,
			"DATAPARAMSILLEGAL",
			ch.RequestID,
		)
	}

	list := expand(jctx, &ue)

	var datalist entity.DataList
	datalist.List = list
	datalist.Total = len(list)

	return http.StatusCreated, io.Suc(
		datalist,
		ch.RequestID,
	)
}

func expand(jctx jcontext.Context, ue *UrlExpand) map[string]interface{} {
	list := make(map[string]interface{})
	shortens := ue.Shorten[0]
	for _, shorten := range strings.Split(shortens, ",") {
		reply, err := mredis.ExpandHGet(shorten)
		if err != nil {
			beego.Info(jctx.Value("requestID").(string), ":", "expand error: ", err)
		}
		atom.Mu.Lock()
		list[shorten] = reply
		atom.Mu.Unlock()
	}
	return list
}
