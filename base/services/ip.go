package services

import (
	jcontext "context"
	"fmt"
	"net/http"
	"strings"

	"github.com/JREAMLU/core/com"
	io "github.com/JREAMLU/core/inout"
	"github.com/JREAMLU/jkernel/base/services/atom"
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
	var ipInfo IPInfo
	ffjson.Unmarshal([]byte(data["querystrjson"].(string)), &ipInfo)

	ch, err := io.InputParamsCheck(jctx, data, ipInfo)
	if err != nil {
		return http.StatusExpectationFailed, io.Fail(ch.Message, "DATAPARAMSILLEGAL", jctx.Value("requestID").(string))
	}

	list, err := getIPsInfo(jctx, &ipInfo)
	if err != nil {
		return http.StatusExpectationFailed, io.Fail(err.Error(), "LOGICILLEGAL", jctx.Value("requestID").(string))
	}

	var datalist entity.DataList
	datalist.List = list
	datalist.Total = len(list)

	return http.StatusCreated, io.Suc(datalist, ch.RequestID)
}

func getIPsInfo(jctx jcontext.Context, ipInfo *IPInfo) (map[string]interface{}, error) {
	var ipList []string
	ips := strings.Split(ipInfo.IPs[0], ",")
	for _, ip := range ips {
		fmt.Println(ip)
		ipList = append(ipList, ip)
	}
	ip, err := com.Query(ipList, "memory")
	if err != nil {
		return nil, err
	}
	var ipResult = make(map[string]interface{})
	for k, v := range ip {
		atom.Mu.Lock()
		ipResult[k] = map[string]interface{}{
			"cityID":   v.CityId,
			"country":  v.Country,
			"region":   v.Region,
			"province": v.Province,
			"city":     v.City,
			"isp":      v.ISP,
		}
		atom.Mu.Unlock()
	}
	return ipResult, nil
}
