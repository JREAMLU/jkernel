package controllers

import (
	"github.com/JREAMLU/core/global"
	io "github.com/JREAMLU/core/inout"
	"github.com/JREAMLU/jkernel/base/services"
)

type IPController struct {
	global.BaseController
}

func (r *IPController) Info() {
	data, jctx := io.InputParams(r.Ctx)

	var service services.IP
	httpStatus, shorten := service.IPsInfo(jctx, data)

	r.Ctx.Output.SetStatus(httpStatus)
	r.Data["json"] = shorten
	r.ServeJSON()
}