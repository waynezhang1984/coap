package controllers

import (
	"github.com/astaxie/beego"
)

type LogsetController struct {
	beego.Controller
}

func (this *LogsetController) Get() {
	beego.Info("ok")
	this.Ctx.Output.Body([]byte{'o', 'k'})
}
