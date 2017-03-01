package controllers

import (
	"github.com/JREAMLU/core/global"
	io "github.com/JREAMLU/core/inout"
	"github.com/JREAMLU/jkernel/base/services"
)

// IPController ip struct
type IPController struct {
	global.BaseController
}

// Info ip info controller
func (r *IPController) Info() {
	data, jctx := io.InputParams(r.Ctx)

	service := services.NewIP()
	httpStatus, shorten := service.IPsInfo(jctx, data)

	r.Ctx.Output.SetStatus(httpStatus)
	r.Data["json"] = shorten
	r.ServeJSON()
}
