package atom

import (
	"github.com/JREAMLU/core/com"
	"github.com/astaxie/beego"
)

var IP string

func init() {
	IP, _ = com.ExternalIP()
}

func GetHeader(requestID string) map[string]string {
	var header = map[string]string{
		"Content-Type":    beego.AppConfig.String("Content-Type"),
		"Accept":          beego.AppConfig.String("Accept"),
		"Accept-Language": beego.AppConfig.String("lang.default"),
		"source":          beego.AppConfig.String("appname"),
		"Request-ID":      requestID,
		"ip":              IP,
	}

	return header
}
