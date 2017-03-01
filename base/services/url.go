package services

import (
	jcontext "context"
	"net/http"

	"github.com/JREAMLU/core/global"
	io "github.com/JREAMLU/core/inout"
	"github.com/JREAMLU/jkernel/base/entity"
	"github.com/JREAMLU/jkernel/base/handler"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/beego/i18n"
	"github.com/pquerna/ffjson/ffjson"
)

// URL url message struct
type URL struct {
	// entity.URLShorten
}

// NewURL return &url
func NewURL() *URL {
	return &URL{}
}

// GetParams get params
func GetParams(url URL) URL {
	return url
}

// Valid valid struct
func (r *URL) Valid(v *validation.Validation) {}

// GoShorten shorten url
func (r *URL) GoShorten(jctx jcontext.Context, data map[string]interface{}) (httpStatus int, output io.Output) {
	return shorten(jctx, data)
}

// GoExpand expand shorten url
func (r *URL) GoExpand(jctx jcontext.Context, data map[string]interface{}) (httpStatus int, output io.Output) {
	return expand(jctx, data)
}

func shorten(jctx jcontext.Context, data map[string]interface{}) (httpStatus int, output io.Output) {
	shortenHandler := handler.NewURLShorten()
	ffjson.Unmarshal(data["body"].([]byte), shortenHandler)
	shortenHandler.FromIP = data["headermap"].(http.Header)["X-Forwarded-For"][0]
	ch, err := io.InputParamsCheck(jctx, data, &shortenHandler.Data)
	if err != nil {
		beego.Info(jctx.Value("requestID").(string), ":", "goShorten error: ", err)
		return http.StatusExpectationFailed, io.Fail(ch.Message, "DATAPARAMSILLEGAL", jctx.Value("requestID").(string))
	}

	if len(shortenHandler.Data.URLs) > 10 {
		beego.Info(jctx.Value("requestID").(string), ":", "goShorten error: ", err)
		return http.StatusExpectationFailed, io.Fail(i18n.Tr(global.Lang, "url.NUMBERLIMIT"), "DATAPARAMSILLEGAL", jctx.Value("requestID").(string))
	}

	list, err := shortenHandler.Shorten(jctx)
	if err != nil {
		beego.Info(jctx.Value("requestID").(string), ":", "goShorten error: ", err)
		return http.StatusExpectationFailed, io.Fail(i18n.Tr(global.Lang, "url.SHORTENILLEGAL"), "LOGICILLEGAL", jctx.Value("requestID").(string))
	}

	datalist := entity.NewDataList()
	datalist.List = list
	datalist.Total = len(list)

	return http.StatusCreated, io.Suc(datalist, ch.RequestID)
}

func expand(jctx jcontext.Context, data map[string]interface{}) (httpStatus int, output io.Output) {
	expandHandler := handler.NewURLExpand()
	ffjson.Unmarshal([]byte(data["querystrjson"].(string)), &expandHandler)

	ch, err := io.InputParamsCheck(jctx, data, expandHandler)
	if err != nil {
		beego.Info(jctx.Value("requestID").(string), ":", "goExpand error: ", err)
		return http.StatusExpectationFailed, io.Fail(ch.Message, "DATAPARAMSILLEGAL", jctx.Value("requestID").(string))
	}

	list, err := expandHandler.Expand(jctx)
	if err != nil {
		beego.Info(jctx.Value("requestID").(string), ":", "goExpand error: ", err)
		return http.StatusExpectationFailed, io.Fail(i18n.Tr(global.Lang, "url.EXPANDILLEGAL"), "LOGICILLEGAL", jctx.Value("requestID").(string))
	}

	datalist := entity.NewDataList()
	datalist.List = list
	datalist.Total = len(list)

	return http.StatusOK, io.Suc(datalist, ch.RequestID)
}
