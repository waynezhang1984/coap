package coapx

import (
	"encoding/json"
	"fmt"

	// "net"
	"os"
	"runtime"
	"time"

	"github.com/coapprocess/msginfo"
	"github.com/coapprocess/redis"
	"github.com/coapprocess/types"

	// "github.com/dustin/go-coap"
	cp "github.com/go-ocf/go-coap"
	log "github.com/sirupsen/logrus"
)

func init() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

var (
	hostname, _ = os.Hostname()
	coaprsp     = "coap_rsp"
)

func handleLogin(w cp.ResponseWriter, req *cp.Request) {

	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "coaps.go",
		"line":     line,
		"func":     f.Name(),
		"code":     req.Msg.Code(),
		"type":     req.Msg.Type(),
		"payload":  string(req.Msg.Payload()),
		"path":     req.Msg.Path(),
		"UDPAddr":  req.Client.RemoteAddr(),
	}).Info("Got message in handleLogin")

	msg := msginfo.MsgData{
		PubType: "coap",
		Data:    req.Msg.Payload(),
		IsReq:   true,
		Url:     req.Client.RemoteAddr().String(),
	}

	msginfo.GetDataChan() <- msg

	creq := types.CommonReq{}
	err := json.Unmarshal(req.Msg.Payload(), &creq)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "coaps.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("handleLogin Unmarshal Error  : ", err)
		return
	}
	devid := creq.DevId
	taskid := creq.Taskid
	repcnt := 0
	for {
		select {
		case <-time.After(200 * time.Millisecond):

			rediskey := coaprsp + "_" + hostname + "_" + "Login_Rsp" + "_" + devid + "_" + taskid
			pc, _, line, _ := runtime.Caller(0)
			f := runtime.FuncForPC(pc)
			log.WithFields(log.Fields{
				"filename": "coaps.go",
				"line":     line,
				"func":     f.Name(),
				"rediskey": rediskey,
			}).Info("ReadRedisByte rediskey")

			redisv, err := redisclient.ReadRedisByte(rediskey)
			if err == nil {
				pc, _, line, _ := runtime.Caller(0)
				f := runtime.FuncForPC(pc)
				log.WithFields(log.Fields{
					"filename": "coaps.go",
					"line":     line,
					"func":     f.Name(),
					"redisv":   redisv,
				}).Info("ReadRedisByte redisv")
				redisclient.DelRedis(rediskey)

				if req.Msg.IsConfirmable() {

					w.SetContentFormat(cp.AppJSON)

					if _, err := w.Write(redisv); err != nil {
						pc, _, line, _ := runtime.Caller(0)
						f := runtime.FuncForPC(pc)
						log.WithFields(log.Fields{
							"filename": "coaps.go",
							"line":     line,
							"func":     f.Name(),
							"rediskey": rediskey,
						}).Error("send response error!")
					}
					pc, _, line, _ := runtime.Caller(0)
					f := runtime.FuncForPC(pc)
					log.WithFields(log.Fields{
						"filename": "coaps.go",
						"line":     line,
						"func":     f.Name(),
						"rediskey": rediskey,
					}).Info("Transmitting from handleLogin")
					w.NewResponse(cp.Content)
				}
				return
			}

			if repcnt >= 5 {
				if req.Msg.IsConfirmable() {

					w.SetContentFormat(cp.AppJSON)
					w.SetCode(cp.InternalServerError)
					if _, err := w.Write([]byte("error!")); err != nil {
						pc, _, line, _ := runtime.Caller(0)
						f := runtime.FuncForPC(pc)
						log.WithFields(log.Fields{
							"filename": "coaps.go",
							"line":     line,
							"func":     f.Name(),
						}).Error("Cannot send response  : ", err)

					}
					//w.NewResponse(cp.InternalServerError)
				}
				pc, _, line, _ := runtime.Caller(0)
				f := runtime.FuncForPC(pc)
				log.WithFields(log.Fields{
					"filename": "coaps.go",
					"line":     line,
					"func":     f.Name(),
				}).Error("send response timeout!")
				return
			}
			repcnt++
		}
	}

}

func handleHeartBeat(w cp.ResponseWriter, req *cp.Request) {

	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "coaps.go",
		"line":     line,
		"func":     f.Name(),
		"code":     req.Msg.Code(),
		"type":     req.Msg.Type(),
		"payload":  string(req.Msg.Payload()),
		"path":     req.Msg.Path(),
		"UDPAddr":  req.Client.RemoteAddr(),
	}).Info("Got message in handleHeartBeat")

	msg := msginfo.MsgData{
		PubType: "coap",
		Data:    req.Msg.Payload(),
		IsReq:   true,
		Url:     req.Client.RemoteAddr().String(),
	}

	msginfo.GetDataChan() <- msg

	creq := types.CommonReq{}
	err := json.Unmarshal(req.Msg.Payload(), &creq)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "coaps.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("handleHeartBeat Unmarshal Error  : ", err)
	}
	devid := creq.DevId
	taskid := creq.Taskid
	repcnt := 0
	for {
		select {
		case <-time.After(200 * time.Millisecond):

			rediskey := coaprsp + "_" + hostname + "_" + "Heartbeat_Rsp" + "_" + devid + "_" + taskid
			pc, _, line, _ := runtime.Caller(0)
			f := runtime.FuncForPC(pc)
			log.WithFields(log.Fields{
				"filename": "coaps.go",
				"line":     line,
				"func":     f.Name(),
				"rediskey": rediskey,
			}).Info("ReadRedisByte rediskey")

			redisv, err := redisclient.ReadRedisByte(rediskey)
			if err == nil {
				pc, _, line, _ := runtime.Caller(0)
				f := runtime.FuncForPC(pc)
				log.WithFields(log.Fields{
					"filename": "coaps.go",
					"line":     line,
					"func":     f.Name(),
					"redisv":   redisv,
				}).Info("ReadRedisByte redisv")
				redisclient.DelRedis(rediskey)

				if req.Msg.IsConfirmable() {

					w.SetContentFormat(cp.AppJSON)

					if _, err := w.Write(redisv); err != nil {
						pc, _, line, _ := runtime.Caller(0)
						f := runtime.FuncForPC(pc)
						log.WithFields(log.Fields{
							"filename": "coaps.go",
							"line":     line,
							"func":     f.Name(),
							"rediskey": rediskey,
						}).Error("send response error!")
					}
					pc, _, line, _ := runtime.Caller(0)
					f := runtime.FuncForPC(pc)
					log.WithFields(log.Fields{
						"filename": "coaps.go",
						"line":     line,
						"func":     f.Name(),
						"rediskey": rediskey,
					}).Info("Transmitting from handleHeartBeat")

					w.NewResponse(cp.Content)
				}
				return
			}

			if repcnt >= 5 {
				if req.Msg.IsConfirmable() {

					w.SetContentFormat(cp.AppJSON)
					w.SetCode(cp.InternalServerError)
					log.Printf("Transmitting from A")
					if _, err := w.Write([]byte("error!")); err != nil {
						pc, _, line, _ := runtime.Caller(0)
						f := runtime.FuncForPC(pc)
						log.WithFields(log.Fields{
							"filename": "coaps.go",
							"line":     line,
							"func":     f.Name(),
						}).Error("Cannot send response  : ", err)
					}
					//w.NewResponse(cp.InternalServerError)
				}
				pc, _, line, _ := runtime.Caller(0)
				f := runtime.FuncForPC(pc)
				log.WithFields(log.Fields{
					"filename": "coaps.go",
					"line":     line,
					"func":     f.Name(),
				}).Error("send response timeout!")
				return
			}
			repcnt++
		}
	}

}

func handleDataReport(w cp.ResponseWriter, req *cp.Request) {

	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "coaps.go",
		"line":     line,
		"func":     f.Name(),
		"code":     req.Msg.Code(),
		"type":     req.Msg.Type(),
		"payload":  string(req.Msg.Payload()),
		"path":     req.Msg.Path(),
		"UDPAddr":  req.Client.RemoteAddr(),
	}).Info("Got message in handleDataReport")

	msg := msginfo.MsgData{
		PubType: "coap",
		Data:    req.Msg.Payload(),
		IsReq:   true,
		Url:     req.Client.RemoteAddr().String(),
	}

	msginfo.GetDataChan() <- msg

	creq := types.CommonReq{}
	err := json.Unmarshal(req.Msg.Payload(), &creq)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "coaps.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("handleDataReport Unmarshal Error  : ", err)
	}
	devid := creq.DevId
	taskid := creq.Taskid
	repcnt := 0
	for {
		select {
		case <-time.After(200 * time.Millisecond):

			rediskey := coaprsp + "_" + hostname + "_" + "DataReport_Rsp" + "_" + devid + "_" + taskid
			pc, _, line, _ := runtime.Caller(0)
			f := runtime.FuncForPC(pc)
			log.WithFields(log.Fields{
				"filename": "coaps.go",
				"line":     line,
				"func":     f.Name(),
				"rediskey": rediskey,
			}).Info("ReadRedisByte rediskey")

			redisv, err := redisclient.ReadRedisByte(rediskey)
			if err == nil {
				pc, _, line, _ := runtime.Caller(0)
				f := runtime.FuncForPC(pc)
				log.WithFields(log.Fields{
					"filename": "coaps.go",
					"line":     line,
					"func":     f.Name(),
					"redisv":   redisv,
				}).Info("ReadRedisByte redisv")
				redisclient.DelRedis(rediskey)

				if req.Msg.IsConfirmable() {

					w.SetContentFormat(cp.AppJSON)
					//log.Printf("Transmitting from A")
					if _, err := w.Write(redisv); err != nil {
						pc, _, line, _ := runtime.Caller(0)
						f := runtime.FuncForPC(pc)
						log.WithFields(log.Fields{
							"filename": "coaps.go",
							"line":     line,
							"func":     f.Name(),
						}).Error("send response error  : ", err)
					}
					pc, _, line, _ := runtime.Caller(0)
					f := runtime.FuncForPC(pc)
					log.WithFields(log.Fields{
						"filename": "coaps.go",
						"line":     line,
						"func":     f.Name(),
						"rediskey": rediskey,
					}).Info("Transmitting from handleDataReport")
					w.NewResponse(cp.Content)
				}
				return
			}

			if repcnt >= 5 {
				if req.Msg.IsConfirmable() {

					w.SetContentFormat(cp.AppJSON)
					w.SetCode(cp.InternalServerError)
					if _, err := w.Write([]byte("error!")); err != nil {
						pc, _, line, _ := runtime.Caller(0)
						f := runtime.FuncForPC(pc)
						log.WithFields(log.Fields{
							"filename": "coaps.go",
							"line":     line,
							"func":     f.Name(),
						}).Error("Cannot send response  : ", err)
					}
					//w.NewResponse(cp.InternalServerError)
				}
				pc, _, line, _ := runtime.Caller(0)
				f := runtime.FuncForPC(pc)
				log.WithFields(log.Fields{
					"filename": "coaps.go",
					"line":     line,
					"func":     f.Name(),
				}).Error("send response timeout!")
				return
			}
			repcnt++
		}
	}

}

func handleAlarmReport(w cp.ResponseWriter, req *cp.Request) {

	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "coaps.go",
		"line":     line,
		"func":     f.Name(),
		"code":     req.Msg.Code(),
		"type":     req.Msg.Type(),
		"payload":  string(req.Msg.Payload()),
		"path":     req.Msg.Path(),
		"UDPAddr":  req.Client.RemoteAddr(),
	}).Info("Got message in handleAlarmReport")

	msg := msginfo.MsgData{
		PubType: "coap",
		Data:    req.Msg.Payload(),
		IsReq:   true,
		Url:     req.Client.RemoteAddr().String(),
	}

	msginfo.GetDataChan() <- msg

	creq := types.CommonReq{}
	err := json.Unmarshal(req.Msg.Payload(), &creq)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "coaps.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("handleAlarmReport Unmarshal Error  : ", err)
	}
	devid := creq.DevId
	taskid := creq.Taskid
	repcnt := 0
	for {
		select {
		case <-time.After(200 * time.Millisecond):

			rediskey := coaprsp + "_" + hostname + "_" + "AlarmReport_Rsp" + "_" + devid + "_" + taskid
			pc, _, line, _ := runtime.Caller(0)
			f := runtime.FuncForPC(pc)
			log.WithFields(log.Fields{
				"filename": "coaps.go",
				"line":     line,
				"func":     f.Name(),
				"rediskey": rediskey,
			}).Info("ReadRedisByte rediskey")

			redisv, err := redisclient.ReadRedisByte(rediskey)
			if err == nil {
				pc, _, line, _ := runtime.Caller(0)
				f := runtime.FuncForPC(pc)
				log.WithFields(log.Fields{
					"filename": "coaps.go",
					"line":     line,
					"func":     f.Name(),
					"redisv":   redisv,
				}).Info("ReadRedisByte redisv")
				redisclient.DelRedis(rediskey)

				if req.Msg.IsConfirmable() {

					w.SetContentFormat(cp.AppJSON)
					log.Printf("Transmitting from A")
					if _, err := w.Write(redisv); err != nil {
						pc, _, line, _ := runtime.Caller(0)
						f := runtime.FuncForPC(pc)
						log.WithFields(log.Fields{
							"filename": "coaps.go",
							"line":     line,
							"func":     f.Name(),
						}).Error("send response error  : ", err)
					}
					pc, _, line, _ := runtime.Caller(0)
					f := runtime.FuncForPC(pc)
					log.WithFields(log.Fields{
						"filename": "coaps.go",
						"line":     line,
						"func":     f.Name(),
						"rediskey": rediskey,
					}).Info("Transmitting from handleAlarmReport")
					w.NewResponse(cp.Content)
				}
				return
			}

			if repcnt >= 5 {
				if req.Msg.IsConfirmable() {

					w.SetContentFormat(cp.AppJSON)
					w.SetCode(cp.InternalServerError)
					log.Printf("Transmitting from A")
					if _, err := w.Write([]byte("error!")); err != nil {
						log.Printf("Cannot send response: %v", err)
					}
					//w.NewResponse(cp.InternalServerError)
				}
				pc, _, line, _ := runtime.Caller(0)
				f := runtime.FuncForPC(pc)
				log.WithFields(log.Fields{
					"filename": "coaps.go",
					"line":     line,
					"func":     f.Name(),
				}).Error("send response timeout!")
				return
			}
			repcnt++
		}
	}

}

func handleConfigReport(w cp.ResponseWriter, req *cp.Request) {

	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "coaps.go",
		"line":     line,
		"func":     f.Name(),
		"code":     req.Msg.Code(),
		"type":     req.Msg.Type(),
		"payload":  string(req.Msg.Payload()),
		"path":     req.Msg.Path(),
		"UDPAddr":  req.Client.RemoteAddr(),
	}).Info("Got message in handleConfigReport")

	msg := msginfo.MsgData{
		PubType: "coap",
		Data:    req.Msg.Payload(),
		IsReq:   true,
		Url:     req.Client.RemoteAddr().String(),
	}

	msginfo.GetDataChan() <- msg

	creq := types.CommonReq{}
	err := json.Unmarshal(req.Msg.Payload(), &creq)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "coaps.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("handleConfigReport Unmarshal Error  : ", err)
	}
	devid := creq.DevId
	taskid := creq.Taskid
	repcnt := 0
	for {
		select {
		case <-time.After(200 * time.Millisecond):

			rediskey := coaprsp + "_" + hostname + "_" + "ConfigReport_Rsp" + "_" + devid + "_" + taskid
			pc, _, line, _ := runtime.Caller(0)
			f := runtime.FuncForPC(pc)
			log.WithFields(log.Fields{
				"filename": "coaps.go",
				"line":     line,
				"func":     f.Name(),
				"rediskey": rediskey,
			}).Info("ReadRedisByte rediskey")

			redisv, err := redisclient.ReadRedisByte(rediskey)
			if err == nil {
				pc, _, line, _ := runtime.Caller(0)
				f := runtime.FuncForPC(pc)
				log.WithFields(log.Fields{
					"filename": "coaps.go",
					"line":     line,
					"func":     f.Name(),
					"redisv":   redisv,
				}).Info("ReadRedisByte redisv")
				fmt.Println(string(redisv))
				redisclient.DelRedis(rediskey)

				if req.Msg.IsConfirmable() {

					w.SetContentFormat(cp.AppJSON)
					log.Printf("Transmitting from A")
					if _, err := w.Write(redisv); err != nil {
						pc, _, line, _ := runtime.Caller(0)
						f := runtime.FuncForPC(pc)
						log.WithFields(log.Fields{
							"filename": "coaps.go",
							"line":     line,
							"func":     f.Name(),
						}).Error("send response error  : ", err)
					}
					pc, _, line, _ := runtime.Caller(0)
					f := runtime.FuncForPC(pc)
					log.WithFields(log.Fields{
						"filename": "coaps.go",
						"line":     line,
						"func":     f.Name(),
						"rediskey": rediskey,
					}).Info("Transmitting from handleConfigReport")
					w.NewResponse(cp.Content)
				}
				return
			}

			if repcnt >= 5 {
				if req.Msg.IsConfirmable() {

					w.SetContentFormat(cp.AppJSON)
					w.SetCode(cp.InternalServerError)
					//w.NewResponse(cp.InternalServerError)
				}
				pc, _, line, _ := runtime.Caller(0)
				f := runtime.FuncForPC(pc)
				log.WithFields(log.Fields{
					"filename": "coaps.go",
					"line":     line,
					"func":     f.Name(),
				}).Error("send response timeout!")
				return
			}
			repcnt++
		}
	}

}

func handleUpgradeNotify(w cp.ResponseWriter, req *cp.Request) {

	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "coaps.go",
		"line":     line,
		"func":     f.Name(),
		"code":     req.Msg.Code(),
		"type":     req.Msg.Type(),
		"payload":  string(req.Msg.Payload()),
		"path":     req.Msg.Path(),
		"UDPAddr":  req.Client.RemoteAddr(),
	}).Info("Got message in handleUpgradeNotify")

	msg := msginfo.MsgData{
		PubType: "coap",
		Data:    req.Msg.Payload(),
		IsReq:   true,
		Url:     req.Client.RemoteAddr().String(),
	}

	msginfo.GetDataChan() <- msg

	creq := types.CommonReq{}
	err := json.Unmarshal(req.Msg.Payload(), &creq)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "coaps.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("handleUpgradeNotify Unmarshal Error  : ", err)
	}
	devid := creq.DevId
	taskid := creq.Taskid
	repcnt := 0
	for {
		select {
		case <-time.After(200 * time.Millisecond):

			rediskey := coaprsp + "_" + hostname + "_" + "UpgradeNotify_Rsp" + "_" + devid + "_" + taskid
			pc, _, line, _ := runtime.Caller(0)
			f := runtime.FuncForPC(pc)
			log.WithFields(log.Fields{
				"filename": "coaps.go",
				"line":     line,
				"func":     f.Name(),
				"rediskey": rediskey,
			}).Info("ReadRedisByte rediskey")

			redisv, err := redisclient.ReadRedisByte(rediskey)
			if err == nil {
				pc, _, line, _ := runtime.Caller(0)
				f := runtime.FuncForPC(pc)
				log.WithFields(log.Fields{
					"filename": "coaps.go",
					"line":     line,
					"func":     f.Name(),
					"redisv":   redisv,
				}).Info("ReadRedisByte redisv")
				fmt.Println(string(redisv))
				redisclient.DelRedis(rediskey)

				if req.Msg.IsConfirmable() {

					w.SetContentFormat(cp.AppJSON)
					//log.Printf("Transmitting from A")
					if _, err := w.Write(redisv); err != nil {
						pc, _, line, _ := runtime.Caller(0)
						f := runtime.FuncForPC(pc)
						log.WithFields(log.Fields{
							"filename": "coaps.go",
							"line":     line,
							"func":     f.Name(),
						}).Error("send response error  : ", err)
					}
					pc, _, line, _ := runtime.Caller(0)
					f := runtime.FuncForPC(pc)
					log.WithFields(log.Fields{
						"filename": "coaps.go",
						"line":     line,
						"func":     f.Name(),
						"rediskey": rediskey,
					}).Info("Transmitting from handleUpgradeNotify")
					w.NewResponse(cp.Content)
				}
				return
			}

			if repcnt >= 5 {
				if req.Msg.IsConfirmable() {

					w.SetContentFormat(cp.AppJSON)
					w.SetCode(cp.InternalServerError)
					//w.NewResponse(cp.InternalServerError)
				}
				pc, _, line, _ := runtime.Caller(0)
				f := runtime.FuncForPC(pc)
				log.WithFields(log.Fields{
					"filename": "coaps.go",
					"line":     line,
					"func":     f.Name(),
				}).Error("send response timeout!")
				return
			}
			repcnt++
		}
	}

}

func StartServer() {
	mux := cp.NewServeMux()
	mux.Handle("/Login", cp.HandlerFunc(handleLogin))
	mux.Handle("/Heartbeat", cp.HandlerFunc(handleHeartBeat))
	mux.Handle("/DataReport", cp.HandlerFunc(handleDataReport))
	mux.Handle("/AlarmReport", cp.HandlerFunc(handleAlarmReport))
	mux.Handle("/ConfigReport", cp.HandlerFunc(handleConfigReport))
	mux.Handle("/UpgradeNotify", cp.HandlerFunc(handleUpgradeNotify))
	//mux.Handle("/Login", coap.FuncHandler(handleLogin))
	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "coap.go",
		"line":     line,
		"func":     f.Name(),
	}).Error(cp.ListenAndServe(":5683", "udp", mux))
	//log.Fatal(coap.ListenAndServe(":5688", "udp", mux))
}
