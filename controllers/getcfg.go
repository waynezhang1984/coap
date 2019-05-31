package controllers

import (
	"encoding/json"
	"time"

	"os"
	"runtime"

	"github.com/astaxie/beego"
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

type GetcfgController struct {
	beego.Controller
}

func (this *GetcfgController) Post() {
	devid := this.Input().Get("devid")
	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "getcfg.go",
		"line":     line,
		"func":     f.Name(),
		"devid":    devid,
	}).Info("GetcfgController Post")

	reqmsg := msginfo.RestReqMsg{
		Method: "GetConfig",
		Devid:  devid,
	}
	msginfo.SetRestDeviceBuf(devid)
	msginfo.SetRestReqBuf(reqmsg)

	ob := &types.RestGetcfg{}

	err := json.Unmarshal(this.Ctx.Input.RequestBody, ob)
	if err != nil {

		this.Ctx.WriteString(err.Error())
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "getcfg.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("GetcfgController Get Unmarshal error: ", err)
		rsp := types.RestRspNoData{
			Code: 100001,
		}
		this.Data["json"] = rsp
		return
	}

	pc, _, line, _ = runtime.Caller(0)
	f = runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "getcfg.go",
		"line":     line,
		"func":     f.Name(),
		"ip":       ob.Data.Ip,
		"port":     ob.Data.Port,
		"user":     ob.Data.Username,
		"passwd":   ob.Data.Password,
		"id":       ob.Data.Id,
		"attrid":   ob.Data.AttrId,
	}).Info("GetcfgController Post")

	dd := types.DeviceData{
		Id:     ob.Data.Id,
		AttrId: ob.Data.AttrId,
	}
	method := "GetConfig"

	ddr := types.CfgQueryReq{
		Method: "GetConfig",
		Taskid: "1",
		DevId:  devid,
		Data:   dd,
	}

	v, err := json.Marshal(&ddr)
	if err != nil {
		this.Ctx.WriteString(err.Error())
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "getcfg.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("GetcfgController Post Marshal error: ", err)
	}

	key := URL + "_" + hostname + "_" + devid
	url, err := redisclient.ReadRedisString(key)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "getcfg.go",
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

		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "getcfg.go",
			"line":     line,
			"func":     f.Name(),
		}).Info("getcfg rsp: ", string(buf))

		drsp := types.CfgQueryRsp{}
		err = json.Unmarshal(buf, &drsp)
		if err != nil {
			pc, _, line, _ := runtime.Caller(0)
			f := runtime.FuncForPC(pc)
			log.WithFields(log.Fields{
				"filename": "getvalue.go",
				"line":     line,
				"func":     f.Name(),
			}).Error("GetcfgController Post Marshal error: ", err)

			rsp := types.RestRspNoData{
				Code: 1,
			}
			this.Data["json"] = &rsp
			this.ServeJSON()
			return
		}

		rsp := types.RestRspGetCfg{
			Code: 0,
			Data: drsp.Data,
		}
		this.Data["json"] = &rsp
		this.ServeJSON()
		return
	} else {

		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "getvalue.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("getcfg Timeout!")
		rsp := types.RestRspNoData{
			Code: 100001,
		}
		this.Data["json"] = &rsp
		this.ServeJSON()
		return
	}

}
