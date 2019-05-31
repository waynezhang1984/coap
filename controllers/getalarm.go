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

type GetalarmController struct {
	beego.Controller
}

func (this *GetalarmController) Post() {
	devid := this.Input().Get("devid")
	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "getalarm.go",
		"line":     line,
		"func":     f.Name(),
		"devid":    devid,
	}).Info("GetalarmController Post")

	reqmsg := msginfo.RestReqMsg{
		Method: "GetAlarm",
		Devid:  devid,
	}
	msginfo.SetRestDeviceBuf(devid)
	msginfo.SetRestReqBuf(reqmsg)

	ob := &types.RestGetalarm{}

	err := json.Unmarshal(this.Ctx.Input.RequestBody, ob)
	if err != nil {

		this.Ctx.WriteString(err.Error())
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "getalarm.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("GetalarmController Get Unmarshal error: ", err)
		rsp := types.RestRspNoData{
			Code: 100001,
		}
		this.Data["json"] = rsp
		return
	}

	pc, _, line, _ = runtime.Caller(0)
	f = runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "getalarm.go",
		"line":     line,
		"func":     f.Name(),
		"ip":       ob.Data.Ip,
		"port":     ob.Data.Port,
		"user":     ob.Data.Username,
		"passwd":   ob.Data.Password,
		"id":       ob.Data.Id,
		"type":     ob.Data.Type,
		"sTime":    ob.Data.StartTime,
		"eTime":    ob.Data.EndTime,
	}).Info("GetalarmController Post")

	am := types.AlertMeta{
		Type:      ob.Data.Type,
		Id:        ob.Data.Id,
		StartTime: ob.Data.StartTime,
		EndTime:   ob.Data.EndTime,
	}
	method := "GetAlarm"
	aqr := types.AlertQueryReq{
		Method: "GetAlarm",
		Taskid: "1",
		DevId:  devid,
		Data:   am,
	}

	v, err := json.Marshal(&aqr)
	if err != nil {
		this.Ctx.WriteString(err.Error())
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "getalarm.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("GetalarmController Post Marshal error: ", err)
	}

	key := URL + "_" + hostname + "_" + devid
	url, err := redisclient.ReadRedisString(key)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "getalarm.go",
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
			"filename": "getalarm.go",
			"line":     line,
			"func":     f.Name(),
		}).Info("getalarm rsp: ", string(buf))

		arsp := types.AlertQueryRsp{}
		err = json.Unmarshal(buf, &arsp)
		if err != nil {
			pc, _, line, _ := runtime.Caller(0)
			f := runtime.FuncForPC(pc)
			log.WithFields(log.Fields{
				"filename": "getalarm.go",
				"line":     line,
				"func":     f.Name(),
			}).Error("GetalarmController Post Marshal error: ", err)

			rsp := types.RestRspNoData{
				Code: 1,
			}
			this.Data["json"] = &rsp
			this.ServeJSON()
			return
		}

		rsp := types.RestRspGetAlarm{
			Code: 0,
			Data: arsp.Data,
		}
		this.Data["json"] = &rsp
		this.ServeJSON()
		return
	} else {

		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "getalarm.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("getalarm Timeout!")
		rsp := types.RestRspNoData{
			Code: 100001,
		}
		this.Data["json"] = &rsp
		this.ServeJSON()
		return
	}

}
