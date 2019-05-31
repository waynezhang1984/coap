package controllers

import (
	"encoding/json"
	"os"
	"runtime"

	"time"

	"github.com/astaxie/beego"
	// "github.com/coapprocess/decode"

	"github.com/coapprocess/coap"
	"github.com/coapprocess/msginfo"
	"github.com/coapprocess/redis"
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

type SetvalueController struct {
	beego.Controller
}

func (this *SetvalueController) Post() {
	devid := this.Input().Get("devid")
	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "setvalue.go",
		"line":     line,
		"func":     f.Name(),
		"devid":    devid,
	}).Info("SetvalueController Post")

	reqmsg := msginfo.RestReqMsg{
		Method: "Ctrl",
		Devid:  devid,
	}
	msginfo.SetRestDeviceBuf(devid)
	msginfo.SetRestReqBuf(reqmsg)

	ob := &types.RestSetvalue{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, ob)
	if err != nil {
		//this.Ctx
		this.Ctx.WriteString(err.Error())
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "setvalue.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("SetvalueController Post Unmarshal error: ", err)
		rsp := types.RestRspNoData{
			Code: 100001,
		}
		this.Data["json"] = rsp
		return
	}

	pc, _, line, _ = runtime.Caller(0)
	f = runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "setvalue.go",
		"line":     line,
		"func":     f.Name(),
		"ip":       ob.Data.Ip,
		"port":     ob.Data.Port,
		"user":     ob.Data.Username,
		"passwd":   ob.Data.Password,
		"id":       ob.Data.Id,
		"attrid":   ob.Data.AttrId,
		"value":    ob.Data.Value,
	}).Info("SetvalueController Post")

	cd := types.ControlData{
		Id:     ob.Data.Id,
		AttrId: ob.Data.AttrId,
		Value:  ob.Data.Value,
	}
	method := "Ctrl"
	cir := types.ControlInfoReq{
		Method: "Ctrl",
		Taskid: "1",
		DevId:  devid,
		Data:   cd,
	}

	v, err := json.Marshal(&cir)
	if err != nil {
		this.Ctx.WriteString(err.Error())
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "setvalue.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("SetvalueController Post Marshal error: ", err)
	}

	key := URL + "_" + hostname + "_" + devid
	url, err := redisclient.ReadRedisString(key)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)

		log.WithFields(log.Fields{
			"filename": "setvalue.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("ReadRedisString Err: ", err)
	}
	log.Info("devurl: ", key)
	repcnt := 0
	buf, err := coapx.PubCoapMsg(v, method, url)
	for err != nil {

		select {
		case <-time.After(200 * time.Second):

			buf, err = coapx.PubCoapMsg(v, method, url)
			if err == nil {
				repcnt = 0
				return
			}

			if repcnt >= 5 {
				rsp := types.RestRspNoData{
					Code: 100001,
				}
				this.Data["json"] = &rsp
				this.ServeJSON()
				return
			}
			repcnt++
		}

	}

	if repcnt == 0 {

		crsp := types.CommonRsp{}
		err := json.Unmarshal(buf, &crsp)
		if err != nil {
			rsp := types.RestRspNoData{
				Code: 100001,
			}
			this.Data["json"] = &rsp
			this.ServeJSON()
			return
		}

		rsp := types.RestRspNoData{
			Code: crsp.Code,
		}
		this.Data["json"] = &rsp
		this.ServeJSON()
		return
	} else {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "setvalue.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("Setvalue Timeout!")
		rsp := types.RestRspNoData{
			Code: 100001,
		}
		this.Data["json"] = &rsp
		this.ServeJSON()
		return
	}

}
