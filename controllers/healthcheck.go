package controllers

import (
	//"fmt"

	"github.com/astaxie/beego"
)

type HealthController struct {
	beego.Controller
}

func (this *HealthController) Get() {
	beego.Info("ok")
	this.Ctx.Output.Body([]byte{'o', 'k'})
}
