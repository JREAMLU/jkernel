package services

import (
	jcontext "context"
	"hash/crc32"
	"net/http"
	"strings"
	"time"

	"github.com/JREAMLU/core/com"
	"github.com/JREAMLU/core/db/mysql"
	"github.com/JREAMLU/core/global"
	io "github.com/JREAMLU/core/inout"
	"github.com/JREAMLU/jkernel/base/models/mentity"
	"github.com/JREAMLU/jkernel/base/models/mmysql"
	"github.com/JREAMLU/jkernel/base/models/mredis"
	"github.com/JREAMLU/jkernel/base/services/atom"
	"github.com/JREAMLU/jkernel/base/services/entity"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/beego/i18n"
	"github.com/pquerna/ffjson/ffjson"
)

type Url struct {
	Meta struct {
		Auth string
	} `json:"meta" valid:"Required"`
	Data struct {
		Urls []interface{} `json:"urls" valid:"Required"`
		IP   string        `json:"ip" valid:"IP"`
	} `json:"data" valid:"Required"`
}

type UrlExpand struct {
	Shorten []string `json:"shorten" valid:"Required"`
}

var ip string

const (
	DELETE = 0
	NORMAL = 1
)

func GetParams(url Url) Url {
	return url
}

func (r *Url) Valid(v *validation.Validation) {}

func (r *Url) GoShorten(jctx jcontext.Context, data map[string]interface{}) (httpStatus int, output io.Output) {
	ffjson.Unmarshal(data["body"].([]byte), r)
	ip = data["headermap"].(http.Header)["X-Forwarded-For"][0]
	ch, err := io.InputParamsCheck(jctx, data, &r.Data)
	if err != nil {
		beego.Info(jctx.Value("requestID").(string), ":", "goShorten error: ", err)
		return http.StatusExpectationFailed, io.Fail(ch.Message, "DATAPARAMSILLEGAL", jctx.Value("requestID").(string))
	}

	if len(r.Data.Urls) > 10 {
		beego.Info(jctx.Value("requestID").(string), ":", "goShorten error: ", err)
		return http.StatusExpectationFailed, io.Fail(i18n.Tr(global.Lang, "url.NUMBERLIMIT"), "DATAPARAMSILLEGAL", jctx.Value("requestID").(string))
	}

	list, err := shorten(r)
	if err != nil {
		beego.Info(jctx.Value("requestID").(string), ":", "goShorten error: ", err)
		return http.StatusExpectationFailed, io.Fail(i18n.Tr(global.Lang, "url.SHORTENILLEGAL"), "LOGICILLEGAL", jctx.Value("requestID").(string))
	}

	var datalist entity.DataList
	datalist.List = list
	datalist.Total = len(list)

	return http.StatusCreated, io.Suc(datalist, ch.RequestID)
}

func shorten(r *Url) (map[string]interface{}, error) {
	list, err := setDB(r)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func setDB(r *Url) (map[string]interface{}, error) {
	var shortenMap = make(map[string]interface{})
	reply, err := mredis.ShortenHMGet(r.Data.Urls)
	if err != nil {
		return nil, err
	}
	exist, notExistLongCRCList, notExistMapList := splitExistOrNot(r, reply)
	if len(notExistLongCRCList) == 0 && len(notExistMapList) == 0 {
		return exist, nil
	}

	existShortListInDB, err := mmysql.GetShortens(notExistLongCRCList)
	if err != nil {
		return nil, err
	}
	existQueue, existQueueShortenList, existQueueExpandList, notExistMapList := getShortenData(existShortListInDB, notExistMapList)
	if len(existQueue) == 0 {
		return exist, nil
	}

	x, err := mysql.GetXS(mmysql.BASE)
	if err != nil {
		return nil, err
	}
	tx := x.Begin()
	if len(notExistMapList) > 0 {
		err = mmysql.ShortenInBatch(notExistMapList, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if len(existQueueShortenList) > 0 {
		_, err = mredis.ShortenHMSet(existQueueShortenList)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if len(existQueueExpandList) > 0 {
		_, err = mredis.ExpandHMSet(existQueueExpandList)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()

	shortenMap = com.MapMerge(exist, existQueue)

	return shortenMap, nil
}

func splitExistOrNot(r *Url, reply []string) (exist map[string]interface{}, notExistLongCRCList []uint64, notExistMapList []mentity.Redirect) {
	exist = make(map[string]interface{})
	var redirect mentity.Redirect
	for key, url := range r.Data.Urls {
		if reply[key] != "" {
			atom.Mu.Lock()
			exist[url.(string)] = reply[key]
			atom.Mu.Unlock()
		} else {
			longCrc := uint64(crc32.ChecksumIEEE([]byte(url.(string))))
			shortUrl := atom.GetShortenUrl(url.(string))
			shortCrc := uint64(crc32.ChecksumIEEE([]byte(shortUrl)))
			notExistLongCRCList = append(notExistLongCRCList, longCrc)
			redirect.LongUrl = url.(string)
			redirect.ShortUrl = shortUrl
			redirect.LongCrc = longCrc
			redirect.ShortCrc = shortCrc
			redirect.Status = NORMAL
			redirect.CreatedByIP = uint64(com.Ip2Int(ip))
			redirect.UpdateByIP = uint64(com.Ip2Int(ip))
			redirect.CreateAT = uint64(time.Now().Unix())
			redirect.UpdateAT = uint64(time.Now().Unix())
			notExistMapList = append(notExistMapList, redirect)
		}
	}
	return exist, notExistLongCRCList, notExistMapList
}

func getShortenData(existShortListInDB []mentity.Redirect, notExistMapList []mentity.Redirect) (existQueue map[string]interface{}, existQueueShortenList []interface{}, existQueueExpandList []interface{}, notExistMapLists []mentity.Redirect) {
	existQueue = make(map[string]interface{})
	for _, existShortListInDBVal := range existShortListInDB {
		atom.Mu.Lock()
		existQueue[existShortListInDBVal.LongUrl] = existShortListInDBVal.ShortUrl
		atom.Mu.Unlock()
		existQueueShortenList = append(existQueueShortenList, existShortListInDBVal.LongUrl)
		existQueueShortenList = append(existQueueShortenList, existShortListInDBVal.ShortUrl)
		existQueueExpandList = append(existQueueExpandList, existShortListInDBVal.ShortUrl)
		existQueueExpandList = append(existQueueExpandList, existShortListInDBVal.LongUrl)
		for k, notExistMapListVal := range notExistMapList {
			if existShortListInDBVal.LongUrl == notExistMapListVal.LongUrl {
				notExistMapList = append(notExistMapList[:k], notExistMapList[k+1:]...)
			}
		}
	}
	for _, notExistMapListVal := range notExistMapList {
		atom.Mu.Lock()
		existQueue[notExistMapListVal.LongUrl] = notExistMapListVal.ShortUrl
		atom.Mu.Unlock()
		existQueueShortenList = append(existQueueShortenList, notExistMapListVal.LongUrl)
		existQueueShortenList = append(existQueueShortenList, notExistMapListVal.ShortUrl)
		existQueueExpandList = append(existQueueExpandList, notExistMapListVal.ShortUrl)
		existQueueExpandList = append(existQueueExpandList, notExistMapListVal.LongUrl)
	}
	return existQueue, existQueueShortenList, existQueueExpandList, notExistMapList
}

func (r *Url) GoExpand(jctx jcontext.Context, data map[string]interface{}) (httpStatus int, output io.Output) {
	var ue UrlExpand
	ffjson.Unmarshal([]byte(data["querystrjson"].(string)), &ue)

	ch, err := io.InputParamsCheck(jctx, data, ue)
	if err != nil {
		beego.Info(jctx.Value("requestID").(string), ":", "goExpand error: ", err)
		return http.StatusExpectationFailed, io.Fail(ch.Message, "DATAPARAMSILLEGAL", jctx.Value("requestID").(string))
	}

	list, err := expand(jctx, &ue)
	if err != nil {
		beego.Info(jctx.Value("requestID").(string), ":", "goExpand error: ", err)
		return http.StatusExpectationFailed, io.Fail(i18n.Tr(global.Lang, "url.EXPANDILLEGAL"), "LOGICILLEGAL", jctx.Value("requestID").(string))
	}

	var datalist entity.DataList
	datalist.List = list
	datalist.Total = len(list)

	return http.StatusCreated, io.Suc(datalist, ch.RequestID)
}

func expand(jctx jcontext.Context, ue *UrlExpand) (list map[string]interface{}, err error) {
	shortenList := shortenList(ue.Shorten[0])
	expandList, err := mredis.ExpandHMGet(shortenList)
	if err != nil {
		return nil, err
	}
	list = getExpandData(shortenList, expandList)
	return list, nil
}

func shortenList(shortens string) []interface{} {
	shortensListStr := strings.Split(shortens, ",")
	var shortenList = make([]interface{}, len(shortensListStr))
	for k, shroten := range shortensListStr {
		shortenList[k] = shroten
	}
	return shortenList
}

func getExpandData(shortenList []interface{}, expandList []string) map[string]interface{} {
	list := make(map[string]interface{})
	for key, url := range expandList {
		atom.Mu.Lock()
		list[shortenList[key].(string)] = url
		atom.Mu.Unlock()
	}
	return list
}
