package handler

import (
	jcontext "context"
	"strings"

	"github.com/JREAMLU/jkernel/base/atom"
	"github.com/JREAMLU/jkernel/base/models/mredis"
)

// URLExpand shorten to expand struct
type URLExpand struct {
	Shorten []string `json:"shorten" valid:"Required"`
}

// NewURLExpand return *URLExpand
func NewURLExpand() *URLExpand {
	return &URLExpand{}
}

// Expand expand handler
func (ue *URLExpand) Expand(jctx jcontext.Context) (list map[string]interface{}, err error) {
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
