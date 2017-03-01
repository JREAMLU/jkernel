package handler

import (
	jcontext "context"
	"hash/crc32"
	"time"

	"github.com/JREAMLU/core/com"
	"github.com/JREAMLU/core/db/mysql"
	"github.com/JREAMLU/jkernel/base/atom"
	"github.com/JREAMLU/jkernel/base/models/mentity"
	"github.com/JREAMLU/jkernel/base/models/mmysql"
	"github.com/JREAMLU/jkernel/base/models/mredis"
)

const (
	// DELETE delete status
	DELETE = 0
	// NORMAL normal status
	NORMAL = 1
)

// URLShorten url shorten
type URLShorten struct {
	Meta struct {
		Auth string
	} `json:"meta" valid:"Required"`
	Data struct {
		URLs []interface{} `json:"urls" valid:"Required"`
		IP   string        `json:"ip" valid:"IP"`
	} `json:"data" valid:"Required"`
	FromIP string
}

// NewURLShorten return *URLShorten
func NewURLShorten() *URLShorten {
	return &URLShorten{}
}

// Shorten shroten handler
func (r *URLShorten) Shorten(jctx jcontext.Context) (map[string]interface{}, error) {
	list, err := setDB(r)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func setDB(r *URLShorten) (map[string]interface{}, error) {
	var shortenMap = make(map[string]interface{})
	reply, err := mredis.ShortenHMGet(r.Data.URLs)
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

func splitExistOrNot(r *URLShorten, reply []string) (exist map[string]interface{}, notExistLongCRCList []uint64, notExistMapList []mentity.Redirect) {
	exist = make(map[string]interface{})
	var redirect mentity.Redirect
	for key, url := range r.Data.URLs {
		if reply[key] != "" {
			atom.Mu.Lock()
			exist[url.(string)] = reply[key]
			atom.Mu.Unlock()
		} else {
			longCrc := uint64(crc32.ChecksumIEEE([]byte(url.(string))))
			shortURL := atom.GetShortenURL(url.(string))
			shortCrc := uint64(crc32.ChecksumIEEE([]byte(shortURL)))
			notExistLongCRCList = append(notExistLongCRCList, longCrc)
			redirect.LongURL = url.(string)
			redirect.ShortURL = shortURL
			redirect.LongCrc = longCrc
			redirect.ShortCrc = shortCrc
			redirect.Status = NORMAL
			redirect.CreatedByIP = uint64(com.IP2Int(r.FromIP))
			redirect.UpdateByIP = uint64(com.IP2Int(r.FromIP))
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
		existQueue[existShortListInDBVal.LongURL] = existShortListInDBVal.ShortURL
		atom.Mu.Unlock()
		existQueueShortenList = append(existQueueShortenList, existShortListInDBVal.LongURL)
		existQueueShortenList = append(existQueueShortenList, existShortListInDBVal.ShortURL)
		existQueueExpandList = append(existQueueExpandList, existShortListInDBVal.ShortURL)
		existQueueExpandList = append(existQueueExpandList, existShortListInDBVal.LongURL)
		for k, notExistMapListVal := range notExistMapList {
			if existShortListInDBVal.LongURL == notExistMapListVal.LongURL {
				notExistMapList = append(notExistMapList[:k], notExistMapList[k+1:]...)
			}
		}
	}
	for _, notExistMapListVal := range notExistMapList {
		atom.Mu.Lock()
		existQueue[notExistMapListVal.LongURL] = notExistMapListVal.ShortURL
		atom.Mu.Unlock()
		existQueueShortenList = append(existQueueShortenList, notExistMapListVal.LongURL)
		existQueueShortenList = append(existQueueShortenList, notExistMapListVal.ShortURL)
		existQueueExpandList = append(existQueueExpandList, notExistMapListVal.ShortURL)
		existQueueExpandList = append(existQueueExpandList, notExistMapListVal.LongURL)
	}
	return existQueue, existQueueShortenList, existQueueExpandList, notExistMapList
}
