package controllers

import (
	"encoding/json"
	"time"

	// "fmt"
	"os"
	"runtime"

	// "time"

	"github.com/astaxie/beego"
	"github.com/coapprocess/coap"

	//	"github.com/coapprocess/mqttclient"
	"github.com/coapprocess/msginfo"
	"github.com/coapprocess/redis"
	"github.com/coapprocess/types"
	log "github.com/sirupsen/logrus"
)

func init() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

type GetvalueController struct {
	beego.Controller
}

func (this *GetvalueController) Post() {
	devid := this.Input().Get("devid")
	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "getvalue.go",
		"line":     line,
		"func":     f.Name(),
		"devid":    devid,
	}).Info("GetvalueController Post")

	reqmsg := msginfo.RestReqMsg{
		Method: "GetData",
		Devid:  devid,
	}
	msginfo.SetRestDeviceBuf(devid)
	msginfo.SetRestReqBuf(reqmsg)

	ob := &types.RestGetvalue{}

	err := json.Unmarshal(this.Ctx.Input.RequestBody, ob)
	if err != nil {

		this.Ctx.WriteString(err.Error())
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "getvalue.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("GetvalueController Get Unmarshal error: ", err)
		rsp := types.RestRspNoData{
			Code: 100001,
		}
		this.Data["json"] = rsp
		return
	}

	pc, _, line, _ = runtime.Caller(0)
	f = runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "getvalue.go",
		"line":     line,
		"func":     f.Name(),
		"ip":       ob.Data.Ip,
		"port":     ob.Data.Port,
		"user":     ob.Data.Username,
		"passwd":   ob.Data.Password,
		"id":       ob.Data.Id,
		"attrid":   ob.Data.AttrId,
	}).Info("GetvalueController Post")

	// dd := types.DeviceData{
	// 	Id:     ob.Data.Id,
	// 	AttrId: ob.Data.AttrId,
	// }
	method := "GetData"
	var dda []types.DeviceData
	if ob.Data.Id != "" || ob.Data.AttrId != "" {
		dd := types.DeviceData{
			Id:     ob.Data.Id,
			AttrId: ob.Data.AttrId,
		}
		dda = append(dda, dd)
	}
	//dda = append(dda, dd)
	ddr := types.DeviceDatasReq{
		Method: "GetData",
		Taskid: "1",
		DevId:  devid,
		Data:   dda,
	}

	v, err := json.Marshal(&ddr)
	if err != nil {
		this.Ctx.WriteString(err.Error())
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "getvalue.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("GetvalueController Post Marshal error: ", err)
	}

	key := URL + "_" + hostname + "_" + devid
	url, err := redisclient.ReadRedisString(key)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "getvalue.go",
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
			"filename": "getvalue.go",
			"line":     line,
			"func":     f.Name(),
		}).Info("getvalue rsp: ", string(buf))

		drsp := types.DeviceDatasRsp{}
		err = json.Unmarshal(buf, &drsp)
		if err != nil {
			pc, _, line, _ := runtime.Caller(0)
			f := runtime.FuncForPC(pc)
			log.WithFields(log.Fields{
				"filename": "getvalue.go",
				"line":     line,
				"func":     f.Name(),
			}).Error("GetvalueController Post Marshal error: ", err)

			rsp := types.RestRspNoData{
				Code: 1,
			}
			this.Data["json"] = &rsp
			this.ServeJSON()
			return
		}

		rsp := types.RestRspGetValue{
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
		}).Error("Getvalue Timeout!")
		rsp := types.RestRspNoData{
			Code: 100001,
		}
		this.Data["json"] = &rsp
		this.ServeJSON()
		return
	}

	// msg := msginfo.MsgData{
	// 	PubType: "restapi",
	// 	Data:    v,
	// 	IsReq:   false,
	// 	Url:     url + ":" + method,
	// }
	// msginfo.GetDataChan() <- msg

	// rediskey := "coaprestrsp_" + "GetData_Rsp" + "_" + devid
	// redisv, err := redisclient.ReadRedisByte(rediskey)
	// if err == nil {
	// 	pc, _, line, _ := runtime.Caller(0)
	// 	f := runtime.FuncForPC(pc)
	// 	log.WithFields(log.Fields{
	// 		"filename": "getvalue.go",
	// 		"line":     line,
	// 		"func":     f.Name(),
	// 	}).Info("getvalue rsp: ", string(redisv))

	// 	drsp := types.DeviceDatasRsp{}
	// 	err = json.Unmarshal(redisv, &drsp)
	// 	if err != nil {
	// 		pc, _, line, _ := runtime.Caller(0)
	// 		f := runtime.FuncForPC(pc)
	// 		log.WithFields(log.Fields{
	// 			"filename": "getvalue.go",
	// 			"line":     line,
	// 			"func":     f.Name(),
	// 		}).Error("GetvalueController Post Marshal error: ", err)

	// 		rsp := types.RestRspNoData{
	// 			Code: 1,
	// 		}
	// 		this.Data["json"] = rsp
	// 		this.ServeJSON()
	// 		return
	// 	}
	// 	redisclient.DelRedis(rediskey)
	// 	rsp := types.RestRspGetValue{
	// 		Code: 0,
	// 		Data: drsp.Data[0],
	// 	}
	// 	this.Data["json"] = rsp
	// 	this.ServeJSON()
	// 	return

	// }

	// repcnt := 0
	// for {
	// 	select {
	// 	case <-time.After(200 * time.Millisecond):

	// 		redisv, err := redisclient.ReadRedisByte(rediskey)
	// 		if err == nil {
	// 			pc, _, line, _ := runtime.Caller(0)
	// 			f := runtime.FuncForPC(pc)
	// 			log.WithFields(log.Fields{
	// 				"filename": "getvalue.go",
	// 				"line":     line,
	// 				"func":     f.Name(),
	// 			}).Info("getvalue rsp: ", string(redisv))

	// 			drsp := types.DeviceDatasRsp{}
	// 			err = json.Unmarshal(redisv, &drsp)
	// 			if err != nil {
	// 				pc, _, line, _ := runtime.Caller(0)
	// 				f := runtime.FuncForPC(pc)
	// 				log.WithFields(log.Fields{
	// 					"filename": "getvalue.go",
	// 					"line":     line,
	// 					"func":     f.Name(),
	// 				}).Error("GetvalueController Post Marshal error: ", err)

	// 				rsp := types.RestRspNoData{
	// 					Code: 1,
	// 				}
	// 				this.Data["json"] = &rsp
	// 				this.ServeJSON()
	// 				return
	// 			}
	// 			redisclient.DelRedis(rediskey)
	// 			rsp := types.RestRspGetValue{
	// 				Code: 0,
	// 				Data: drsp.Data[0],
	// 			}
	// 			this.Data["json"] = &rsp
	// 			this.ServeJSON()
	// 			return
	// 		}

	// 		if repcnt >= 5 {
	// 			redisclient.DelRedis(rediskey)
	// 			rsp := types.RestRspNoData{
	// 				Code: 100001,
	// 			}
	// 			this.Data["json"] = &rsp
	// 			this.ServeJSON()
	// 			return
	// 		}
	// 		repcnt++
	// 	}
	// }

}
