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

// NewIP return *ip
func NewIP() *IP {
	return &IP{}
}

// Valid valid struct
func (r *IP) Valid(v *validation.Validation) {}

// IPsInfo is list info
func (r *IP) IPsInfo(jctx jcontext.Context, data map[string]interface{}) (httpStatus int, output io.Output) {
	return ipInfo(jctx, data)
}

func ipInfo(jctx jcontext.Context, data map[string]interface{}) (httpStatus int, output io.Output) {
	ipHandler := handler.NewIPInfo()
	ffjson.Unmarshal([]byte(data["querystrjson"].(string)), &ipHandler)

	ch, err := io.InputParamsCheck(jctx, data, ipHandler)
	if err != nil {
		return http.StatusExpectationFailed, io.Fail(ch.Message, "DATAPARAMSILLEGAL", jctx.Value("requestID").(string))
	}

	list, err := ipHandler.IPsInfo(jctx)
	if err != nil {
		beego.Info(jctx.Value("requestID").(string), ":", "getIPsInfo error: ", err)
		return http.StatusExpectationFailed, io.Fail(i18n.Tr(global.Lang, "ip.IPSINFOILLEGAL"), "LOGICILLEGAL", jctx.Value("requestID").(string))
	}

	datalist := entity.NewDataList()
	datalist.List = list
	datalist.Total = len(list)

	return http.StatusOK, io.Suc(datalist, ch.RequestID)

}
