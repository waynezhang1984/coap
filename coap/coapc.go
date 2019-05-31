package coapx

import (
	"bytes"
	"math/rand"
	"os"
	"runtime"
	"strings"

	"github.com/dustin/go-coap"
	cp "github.com/go-ocf/go-coap"
	log "github.com/sirupsen/logrus"
)

func init() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func RandUInt16(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}

func PubCoapMsg1(value []byte, method string, ip string) ([]byte, error) {
	var md string

	//ipa := strings.Split(ip, ":")
	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "coapc.go",
		"line":     line,
		"func":     f.Name(),
	}).Info("ip: ", ip)

	switch method {
	case "Login_Rsp":
	case "Heartbeat_Rsp":
	case "DataReport_Rsp":
	case "AlarmReport_Rsp":
	case "ConfigReport_Rsp":
		md = strings.TrimRight(method, "_Rsp")
	default:
		md = method
	}

	msgId := uint16(RandUInt16(1, 1000))
	pc, _, line, _ = runtime.Caller(0)
	f = runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "coapc.go",
		"line":     line,
		"func":     f.Name(),
	}).Info("MessageID: ", msgId)

	pc, _, line, _ = runtime.Caller(0)
	f = runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "coapc.go",
		"line":     line,
		"func":     f.Name(),
	}).Info("method: ", md)
	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.POST,
		MessageID: msgId,
		Payload:   value,
	}

	path := "/" + md

	req.SetOption(coap.ETag, "weetag")
	req.SetOption(coap.MaxAge, 3)
	req.SetPathString(path)

	c, err := coap.Dial("udp", ip)
	if err != nil {

		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "coapc.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("Error dialing: %v", err)
		return []byte{}, err
	}

	pc, _, line, _ = runtime.Caller(0)
	f = runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "coapc.go",
		"line":     line,
		"func":     f.Name(),
	}).Info("req: ", req)

	rv, err := c.Send(req)
	if err != nil {

		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "coapc.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("Error sending request: ", err)
		return []byte{}, err
	}

	if rv != nil {

		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename":  "coapc.go",
			"line":      line,
			"func":      f.Name(),
			"type":      rv.Type,
			"code":      rv.Code,
			"messageid": rv.MessageID,
			"payload":   rv.Payload,
		}).Info("Response: ", rv)

	}

	return rv.Payload, err
}

func PubCoapMsg(value []byte, method string, ip string) ([]byte, error) {
	var md string

	//ipa := strings.Split(ip, ":")
	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "coapc.go",
		"line":     line,
		"func":     f.Name(),
	}).Info("ip: ", ip)

	switch method {
	case "Login_Rsp":
	case "Heartbeat_Rsp":
	case "DataReport_Rsp":
	case "AlarmReport_Rsp":
	case "ConfigReport_Rsp":
		md = strings.TrimRight(method, "_Rsp")
	default:
		md = method
	}

	msgId := uint16(RandUInt16(1, 1000))
	pc, _, line, _ = runtime.Caller(0)
	f = runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "coapc.go",
		"line":     line,
		"func":     f.Name(),
	}).Info("MessageID: ", msgId)

	pc, _, line, _ = runtime.Caller(0)
	f = runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "coapc.go",
		"line":     line,
		"func":     f.Name(),
	}).Info("method: ", md)

	// url := ipa[0] + ":5683"
	// pc, _, line, _ = runtime.Caller(0)
	// f = runtime.FuncForPC(pc)
	// log.WithFields(log.Fields{
	// 	"filename": "coapc.go",
	// 	"line":     line,
	// 	"func":     f.Name(),
	// }).Info("url: ", url)

	path := "/" + md

	BlockWiseTransfer := true
	BlockWiseTransferSzx := cp.BlockWiseSzx1024
	c := cp.Client{Net: "udp", BlockWiseTransfer: &BlockWiseTransfer, BlockWiseTransferSzx: &BlockWiseTransferSzx}
	con, err := c.Dial(ip)
	if err != nil {
		//log.Fatalf("Error dialing: %v", err)
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "coapc.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("Error dialing: ", err)
		return []byte{}, err
	}

	body := bytes.NewReader(value)
	resp, err := con.Post(path, cp.TextPlain, body)
	if err != nil {
		//log.Fatalf("Error sending request: %v", err)
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "coapc.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("Error Post request: ", err)
		return []byte{}, err
	}

	if resp != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename":  "coapc.go",
			"line":      line,
			"func":      f.Name(),
			"type":      resp.Type(),
			"code":      resp.Code(),
			"messageid": resp.MessageID(),
			"payload":   string(resp.Payload()),
		}).Info("Response")
		return resp.Payload(), nil
	}

	return []byte{}, err
}
