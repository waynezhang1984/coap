package msgjson

import (
	"encoding/json"

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

func CoapMsgEncode(buf []byte, md string, id string) []byte {
	info := types.DeviceMessage{
		Method:  md,
		MsgId:   id,
		Data:    buf,
		Msgfrom: "coap_v100",
	}

	msgbyte, err := json.Marshal(info)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "encode.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("CoapMsgEncode error: ", err)
		return []byte{}
	}
	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "encode.go",
		"line":     line,
		"func":     f.Name(),
	}).Info("CoapMsgEncode: ", string(msgbyte))
	return msgbyte
}
