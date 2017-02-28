package handler

import (
	jcontext "context"
	"fmt"
	"strings"

	"github.com/JREAMLU/core/com"
	"github.com/JREAMLU/jkernel/base/atom"
	"github.com/JREAMLU/jkernel/base/entity"
)

// IPsInfo get ips info
func IPsInfo(jctx jcontext.Context, ipInfo *entity.IPInfo) (map[string]interface{}, error) {
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
