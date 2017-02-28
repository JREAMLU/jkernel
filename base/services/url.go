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
	entity.URLShorten
}

// GetParams get params
func GetParams(url URL) URL {
	return url
}

// Valid valid struct
func (r *URL) Valid(v *validation.Validation) {}

// GoShorten shorten url
func (r *URL) GoShorten(jctx jcontext.Context, data map[string]interface{}) (httpStatus int, output io.Output) {
	ffjson.Unmarshal(data["body"].([]byte), r)
	r.FromIP = data["headermap"].(http.Header)["X-Forwarded-For"][0]
	ch, err := io.InputParamsCheck(jctx, data, &r.Data)
	if err != nil {
		beego.Info(jctx.Value("requestID").(string), ":", "goShorten error: ", err)
		return http.StatusExpectationFailed, io.Fail(ch.Message, "DATAPARAMSILLEGAL", jctx.Value("requestID").(string))
	}

	if len(r.Data.URLs) > 10 {
		beego.Info(jctx.Value("requestID").(string), ":", "goShorten error: ", err)
		return http.StatusExpectationFailed, io.Fail(i18n.Tr(global.Lang, "url.NUMBERLIMIT"), "DATAPARAMSILLEGAL", jctx.Value("requestID").(string))
	}

	var us entity.URLShorten
	us.Data = r.Data
	us.Meta = r.Meta
	list, err := handler.Shorten(&us)
	if err != nil {
		beego.Info(jctx.Value("requestID").(string), ":", "goShorten error: ", err)
		return http.StatusExpectationFailed, io.Fail(i18n.Tr(global.Lang, "url.SHORTENILLEGAL"), "LOGICILLEGAL", jctx.Value("requestID").(string))
	}

	var datalist entity.DataList
	datalist.List = list
	datalist.Total = len(list)

	return http.StatusCreated, io.Suc(datalist, ch.RequestID)
}

// GoExpand expand shorten url
func (r *URL) GoExpand(jctx jcontext.Context, data map[string]interface{}) (httpStatus int, output io.Output) {
	var ue entity.URLExpand
	ffjson.Unmarshal([]byte(data["querystrjson"].(string)), &ue)

	ch, err := io.InputParamsCheck(jctx, data, ue)
	if err != nil {
		beego.Info(jctx.Value("requestID").(string), ":", "goExpand error: ", err)
		return http.StatusExpectationFailed, io.Fail(ch.Message, "DATAPARAMSILLEGAL", jctx.Value("requestID").(string))
	}

	list, err := handler.Expand(jctx, &ue)
	if err != nil {
		beego.Info(jctx.Value("requestID").(string), ":", "goExpand error: ", err)
		return http.StatusExpectationFailed, io.Fail(i18n.Tr(global.Lang, "url.EXPANDILLEGAL"), "LOGICILLEGAL", jctx.Value("requestID").(string))
	}

	var datalist entity.DataList
	datalist.List = list
	datalist.Total = len(list)

	return http.StatusCreated, io.Suc(datalist, ch.RequestID)
}
