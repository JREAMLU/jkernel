// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"base/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	/**
	 * https://base.jream.lu/v1/url/goshorten.json
	 * https://base.jream.lu/v1/url/goexpand.json
	 */
	base := beego.NewNamespace("/v1",
		beego.NSCond(func(ctx *context.Context) bool {
			if ctx.Input.Domain() == beego.AppConfig.String("baseDomain") {
				return true
			}
			return false
		}),
		beego.NSNamespace("url",
			beego.NSRouter("/goshorten.json", &controllers.UrlController{}, "post:GoShorten"),
			beego.NSRouter("/goshorten", &controllers.UrlController{}, "post:GoShorten"),
			beego.NSRouter("/goexpand.json", &controllers.UrlController{}, "post:GoExpand"),
			beego.NSRouter("/goexpand", &controllers.UrlController{}, "post:GoExpand"),
		),
		beego.NSRouter("/object", &controllers.ObjectController{}, "post:Post"),
		beego.NSRouter("/test", &controllers.TesController{}, "post:Post"),
	)

	beego.AddNamespace(base)
}
