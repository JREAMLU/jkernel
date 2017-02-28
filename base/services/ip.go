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

// IP ip service
type IP struct {
}

// Valid valid struct
func (r *IP) Valid(v *validation.Validation) {}

// IPsInfo is list info
func (r *IP) IPsInfo(jctx jcontext.Context, data map[string]interface{}) (httpStatus int, output io.Output) {
	var ipInfo entity.IPInfo
	ffjson.Unmarshal([]byte(data["querystrjson"].(string)), &ipInfo)

	ch, err := io.InputParamsCheck(jctx, data, ipInfo)
	if err != nil {
		return http.StatusExpectationFailed, io.Fail(ch.Message, "DATAPARAMSILLEGAL", jctx.Value("requestID").(string))
	}

	list, err := handler.IPsInfo(jctx, &ipInfo)
	if err != nil {
		beego.Info(jctx.Value("requestID").(string), ":", "getIPsInfo error: ", err)
		return http.StatusExpectationFailed, io.Fail(i18n.Tr(global.Lang, "ip.IPSINFOILLEGAL"), "LOGICILLEGAL", jctx.Value("requestID").(string))
	}

	var datalist entity.DataList
	datalist.List = list
	datalist.Total = len(list)

	return http.StatusCreated, io.Suc(datalist, ch.RequestID)
}
