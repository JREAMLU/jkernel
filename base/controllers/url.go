package controllers

import (
	"github.com/JREAMLU/core/global"
	io "github.com/JREAMLU/core/inout"
	"github.com/JREAMLU/jkernel/base/services"
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
	data, jctx := io.InputParams(r.Ctx)

	var service services.Url
	httpStatus, shorten := service.GoShorten(jctx, data)

	r.Ctx.Output.SetStatus(httpStatus)
	r.Data["json"] = shorten
	r.ServeJSON()
}

func (r *UrlController) GoExpand() {
	data, jctx := io.InputParams(r.Ctx)

	var service services.Url
	httpStatus, expand := service.GoExpand(jctx, data)

	r.Ctx.Output.SetStatus(httpStatus)
	r.Data["json"] = expand
	r.ServeJSON()
}
