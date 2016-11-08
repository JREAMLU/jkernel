package services

import (
	jcontext "context"
	"fmt"
	"net/http"

	io "github.com/JREAMLU/core/inout"
	"github.com/JREAMLU/jkernel/base/services/entity"
	"github.com/astaxie/beego/validation"
	"github.com/pquerna/ffjson/ffjson"
)

type IP struct {
}

type IPInfo struct {
	IPs []string `json:"ips" valid:"Required"`
}

func (r *IP) Valid(v *validation.Validation) {}

func (r *IP) IPsInfo(jctx jcontext.Context, data map[string]interface{}) (httpStatus int, output io.Output) {
	ffjson.Unmarshal(data["body"].([]byte), r)
	var ipInfo IPInfo
	ffjson.Unmarshal([]byte(data["querystrjson"].(string)), &ipInfo)

	ch, err := io.InputParamsCheck(jctx, data, ipInfo)
	if err != nil {
		return http.StatusExpectationFailed, io.Fail(ch.Message, "DATAPARAMSILLEGAL", ch.RequestID)
	}

	fmt.Println("<<<<<<<<<<<", ipInfo.IPs)
	list := getIPsInfo(jctx, &ipInfo)

	var datalist entity.DataList
	datalist.List = list
	datalist.Total = len(list)

	return http.StatusCreated, io.Suc(
		datalist,
		ch.RequestID,
	)
}

func getIPsInfo(jctx jcontext.Context, ipInfo *IPInfo) map[string]interface{} {
	// fmt.Println("<<<<<<<<", ipInfo.IPs[0], []string{"127.0.0.1", "119.75.218.70"})
	// ip, err := com.Query(ipInfo.IPs, "memory")
	// fmt.Println(ip, err)
	return nil
}
