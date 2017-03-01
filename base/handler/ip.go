package handler

import (
	jcontext "context"
	"fmt"
	"strings"

	"github.com/JREAMLU/core/com"
	"github.com/JREAMLU/jkernel/base/atom"
)

// IPInfo ip info struct
type IPInfo struct {
	IPs []string `json:"ips" valid:"Required"`
}

// IPsInfo get ips info
func IPsInfo(jctx jcontext.Context, ipInfo *IPInfo) (map[string]interface{}, error) {
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
