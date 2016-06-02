package controllers

import (
	"base/services"
	"core/global"

	"github.com/astaxie/beego"

	//"encoding/json"
)

type UrlController struct {
	global.BaseController
}

/**
HEADER:
	Accept-Language: zh-CN
	source: Advanced Rest Client
	version: 1.0
	Secret-Key: ABDEFGHIJKLMNOPQRSTUVWXYZ
	Request-Id: AAAAAAAAAAAAAAAAAAAAAAAAA
	token: ONE-PIECE
	ip: 192.168.1.1

DATA:
{
    "data": {
        "urls": [
            {
                "long_url": "http://o9d.cn",
                "IP": "127.0.0.1"
            },
            {
                "long_url": "http://huiyimei.com",
                "IP": "192.168.1.1"
            }
        ]
    }
}
*/
/**
 *	@auther			jream.lu
 *	@Request		post
 *	@url			https://base.jream.lu/v1/url/goshorten
 *	@Description 	入参rawMetaHeader, rawDataBody raw形式  meta以header信息传递 data以raw json形式传递
 *	@todo 			参数验证, 封装返回
 */
func (r *UrlController) GoShorten() {
	//入参 meta data
	rawMetaHeader := r.Ctx.Input.Context.Request.Header
	rawDataBody := r.Ctx.Input.RequestBody

	//记录参数日志
	beego.Trace("入参body:" + string(rawDataBody))

	//调用servcie方法, 将参数传递过去
	var service services.Url
	httpStatus, shorten := service.GoShorten(rawMetaHeader, rawDataBody)

	r.Ctx.Output.SetStatus(httpStatus)
	r.Data["json"] = shorten
	r.ServeJSON()
}

func (r *UrlController) GoExpand() {
}
