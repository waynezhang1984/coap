package controllers

import (
	"github.com/astaxie/beego"
)

type DebugController struct {
	beego.Controller
}

func (this *DebugController) Get() {
	beego.Info("ok")
	this.Ctx.Output.Body([]byte{'o', 'k'})
}
