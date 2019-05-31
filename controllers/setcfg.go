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
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

type SetcfgController struct {
	beego.Controller
}

func (this *SetcfgController) Post() {
	devid := this.Input().Get("devid")
	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "setcfg.go",
		"line":     line,
		"func":     f.Name(),
		"devid":    devid,
	}).Info("SetcfgController Post")

	reqmsg := msginfo.RestReqMsg{
		Method: "SetConfig",
		Devid:  devid,
	}
	msginfo.SetRestDeviceBuf(devid)
	msginfo.SetRestReqBuf(reqmsg)

	ob := &types.RestSetcfg{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, ob)
	if err != nil {
		//this.Ctx
		this.Ctx.WriteString(err.Error())
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "setcfg.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("SetcfgController Post Unmarshal error: ", err)
		rsp := types.RestRspNoData{
			Code: 100001,
		}
		this.Data["json"] = rsp
		return
	}

	pc, _, line, _ = runtime.Caller(0)
	f = runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename":    "setcfg.go",
		"line":        line,
		"func":        f.Name(),
		"ip":          ob.Data.Ip,
		"port":        ob.Data.Port,
		"user":        ob.Data.Username,
		"passwd":      ob.Data.Password,
		"id":          ob.Data.Id,
		"attrid":      ob.Data.AttrId,
		"min":         ob.Data.Min,
		"max":         ob.Data.Max,
		"normalvalue": ob.Data.NormalValue,
		"kv":          ob.Data.Kv,
		"bv":          ob.Data.Bv,
	}).Info("SetcfgController Post")

	cd := types.CfgSetInfo{
		Id:          ob.Data.Id,
		AttrId:      ob.Data.AttrId,
		Min:         ob.Data.Min,
		Max:         ob.Data.Max,
		NormalValue: ob.Data.NormalValue,
		Kv:          ob.Data.Kv,
		Bv:          ob.Data.Bv,
	}
	method := "SetConfig"
	cir := types.CfgSetReq{
		Method: "SetConfig",
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
			"filename": "setcfg.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("SetcfgController Post Marshal error: ", err)
	}

	key := URL + "_" + hostname + "_" + devid
	url, err := redisclient.ReadRedisString(key)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)

		log.WithFields(log.Fields{
			"filename": "setcfg.go",
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
			"filename": "setcfg.go",
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
