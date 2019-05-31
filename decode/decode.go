package msgjson

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"runtime"

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

type DeviceRspMessage struct {
	Data   string
	Method string `json:"method"`
	MsgId  string `json:"msgid"`
}

func KafkaMsgDecode(buf []byte) ([]byte, string, string, error) {

	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "decode.go",
		"line":     line,
		"func":     f.Name(),
	}).Info("KafkaMsgDecode: ", string(buf))

	info := DeviceRspMessage{}
	err := json.Unmarshal(buf, &info)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "decode.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("kafkaMsgDecode error: ", err)
		return []byte{}, "", "", err
	}

	bbuf, err := base64.StdEncoding.DecodeString(info.Data)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "decode.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("base64 decode failure, error=[%v]\n", err)
		return []byte{}, "", "", err
	}
	pc, _, line, _ = runtime.Caller(0)
	f = runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "decode.go",
		"line":     line,
		"func":     f.Name(),
		"method":   info.Method,
		"msgid":    info.MsgId,
		"Data":     string(bbuf),
	}).Info("KafkaMsgDecode")

	return bbuf, info.Method, info.MsgId, nil
}

func PlatReqDecode(buf []byte) (string, string, error) {
	info := types.CommonReq{}
	err := json.Unmarshal(buf, &info)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "decode.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("msgDecode error: ", err)
		return "", "", err
	}

	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "decode.go",
		"line":     line,
		"func":     f.Name(),
		"method":   info.Method,
		"devid":    info.DevId,
		"taskid":   info.Taskid,
	}).Info("PlatReqDecode")
	return info.DevId, info.Taskid, nil
}

func PlatRsp(buf []byte) (string, string, error) {
	rsp := types.CommonRsp{}

	err := json.Unmarshal(buf, &rsp)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "decode.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("msgDecode error: ", err)
		return "", "", err
	}

	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "decode.go",
		"line":     line,
		"func":     f.Name(),
		"method":   rsp.Method,
		"devid":    rsp.DevId,
		"taskid":   rsp.Taskid,
	}).Info("PlatRsp")
	return rsp.DevId, rsp.Taskid, nil
}

func PlatMsgDecode(buf []byte) (string, string, error) {

	devid, taskid, err := PlatRsp(buf)
	if err == nil {
		return devid, taskid, err
	}

	devid, taskid, err = PlatReqDecode(buf)
	if err == nil {
		return devid, taskid, err
	}

	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "decode.go",
		"line":     line,
		"func":     f.Name(),
	}).Info("PlatMsgDecode error: ", err)
	return devid, taskid, err
}

func DeviceReqDecode(buf []byte) (string, string, string, error) {
	info := types.CommonReq{}
	err := json.Unmarshal(buf, &info)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "decode.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("msgDecode error: ", err)
		return "", "", "", err
	}

	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "decode.go",
		"line":     line,
		"func":     f.Name(),
		"method":   info.Method,
		"devid":    info.DevId,
		"taskid":   info.Taskid,
	}).Info("DeviceReqDecode")
	return info.Method, info.DevId, info.Taskid, nil
}

func DeviceRsp(buf []byte) (string, string, string, error) {
	rsp := types.CommonRsp{}

	err := json.Unmarshal(buf, &rsp)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "decode.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("msgDecode error: ", err)
		return "", "", "", err
	}

	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "decode.go",
		"line":     line,
		"func":     f.Name(),
		"method":   rsp.Method,
		"devid":    rsp.DevId,
		"error":    rsp.Code,
		"taskid":   rsp.Taskid,
	}).Info("DeviceRsp")
	return rsp.Method, rsp.DevId, rsp.Taskid, nil
}

func RestDeviceRsp(buf []byte) (string, string, int, error) {
	rsp := types.CommonRsp{}

	err := json.Unmarshal(buf, &rsp)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "decode.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("msgDecode error: ", err)
		return "", "", 1, err
	}

	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "decode.go",
		"line":     line,
		"func":     f.Name(),
		"method":   rsp.Method,
		"devid":    rsp.DevId,
		"error":    rsp.Code,
	}).Info("DeviceRsp")
	return rsp.Method, rsp.DevId, rsp.Code, nil
}

func CoapMsgDecode(buf []byte, isreq bool) (string, string, string, error) {

	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "decode.go",
		"line":     line,
		"func":     f.Name(),
	}).Info("CoapMsgDecode: ", string(buf))

	if isreq {
		method, id, tid, err := DeviceReqDecode(buf)
		return method, id, tid, err
	} else {
		method, id, tid, err := DeviceRsp(buf)
		return method, id, tid, err
	}
}

func CommonReqDecode(buf []byte) (string, string, error) {
	info := types.CommonReq{}
	err := json.Unmarshal(buf, &info)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "decode.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("msgDecode error: ", err)
		return "", "", err
	}

	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "decode.go",
		"line":     line,
		"func":     f.Name(),
		"method":   info.Method,
		"devid":    info.DevId,
	}).Info("CommonReqDecode:")
	return info.Method, info.DevId, nil
}

func CommonRsp(buf []byte) (string, error) {
	rsp := types.CommonRsp{}

	err := json.Unmarshal(buf, &rsp)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "decode.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("msgDecode error: ", err)
		return "", err
	}

	switch rsp.Method {
	case "Login_Rsp":
	case "GetData_Rsp":
	case "GetAlarm_Rsp":
	case "GetConfig_Rsp":
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "decode.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("msgDecode error: ", errors.New("method mismatch!"))
		return "", errors.New("method mismatch!")
	default:
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "decode.go",
			"line":     line,
			"func":     f.Name(),
		}).Info("method: ", rsp.Method)
		return rsp.DevId, nil
	}

	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "decode.go",
		"line":     line,
		"func":     f.Name(),
	}).Info("method: ", rsp.Method)
	return rsp.DevId, nil
}
