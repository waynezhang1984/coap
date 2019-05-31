package controllers

import (
	"encoding/json"
	"os"
	"runtime"

	"github.com/astaxie/beego"
	"github.com/coapprocess/types"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
}

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Post() {

	devid := this.Input().Get("devid")
	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "login.go",
		"line":     line,
		"func":     f.Name(),
		"devid":    devid,
	}).Info("LoginController Post")

	ob := &types.RestLogin{}
	beego.Info(string(this.Ctx.Input.RequestBody))
	err := json.Unmarshal(this.Ctx.Input.RequestBody, ob)
	if err != nil {
		//this.Ctx
		this.Ctx.WriteString(err.Error())
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "login.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("LoginController Post Unmarshal error: ", err)

		rsp := types.RestRspNoData{
			Code: 100001,
		}
		this.Data["json"] = rsp
		return
	}

	pc, _, line, _ = runtime.Caller(0)
	f = runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "login.go",
		"line":     line,
		"func":     f.Name(),
		"ip":       ob.Data.Ip,
		"port":     ob.Data.Port,
		"user":     ob.Data.Username,
		"passwd":   ob.Data.Password,
	}).Info("LoginController Post")
	rsp := types.RestRspNoData{
		Code: 0,
	}
	this.Data["json"] = rsp
	this.Ctx.Output.Body([]byte{'o', 'k'})
}
