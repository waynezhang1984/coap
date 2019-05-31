package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"strings"

	"encoding/json"
	"io/ioutil"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/astaxie/beego"
	"github.com/coapprocess/coap"
	"github.com/coapprocess/configure"
	"github.com/coapprocess/controllers"
	"github.com/coapprocess/decode"
	"github.com/coapprocess/eureka"
	"github.com/coapprocess/kafka"

	"github.com/coapprocess/msginfo"
	"github.com/coapprocess/redis"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

var (
	brokerList    = flag.String("kafka", os.Getenv("KAFKA_PEERS"), "The comma separated list of brokers in the Kafka cluster")
	ctopic        = flag.String("ctopic", "", "REQUIRED: the topic to consume")
	cpartitions   = flag.String("partitions", "all", "The partitions to consume, can be 'all' or comma-separated numbers")
	coffset       = flag.String("offset", "newest", "The offset to start with. Can be `oldest`, `newest`")
	verbose       = flag.Bool("verbose", false, "Whether to turn on sarama logging")
	tlsEnabled    = flag.Bool("tls-enabled", false, "Whether to enable TLS")
	tlsSkipVerify = flag.Bool("tls-skip-verify", false, "Whether skip TLS server cert verification")
	tlsClientCert = flag.String("tls-client-cert", "", "Client cert for client authentication (use with -tls-enabled and -tls-client-key)")
	tlsClientKey  = flag.String("tls-client-key", "", "Client key for client authentication (use with tls-enabled and -tls-client-cert)")

	cbufferSize = flag.Int("buffer-size", 256, "The buffer size of the message channel.")

	ptopic       = flag.String("ptopic", "", "REQUIRED: the topic to produce to")
	pkey         = flag.String("key", "", "The key of the message to produce. Can be empty.")
	pvalue       = flag.String("value", "", "REQUIRED: the value of the message to produce. You can also provide the value on stdin.")
	ppartitioner = flag.String("partitioner", "", "The partitioning scheme to use. Can be `hash`, `manual`, or `random`")
	ppartition   = flag.Int("partition", -1, "The partition to produce to.")
	pshowMetrics = flag.Bool("metrics", false, "Output metrics on successful publish to stderr")
	psilent      = flag.Bool("silent", false, "Turn off printing the message's topic, partition, and offset to stdout")

	eurekaurls  = flag.String("eureka", "192.168.101.189:9000", "eureka server ip")
	redisurls   = flag.String("redis", "192.168.101.189:6379", "eureka server ip")
	MSG         = "coap_msgid"
	URL         = "coap_url"
	hostname, _ = os.Hostname()
)

type CfgPara struct {
	Kafka    string `json:"kafka"`
	Reqtopic string `json:"reqtopic"`
	Rsptopic string `json:"rsptopic"`
	Eureka   string `json:"eureka"`
	Redis    string `json:"redis"`
}

func main() {
	fmt.Println("welcome coapprocess")
	flag.Parse()

	initConfig()
	eurekaStart()

	go func() {
		coapx.StartServer()
	}()

	go func() {
		for {
			kafka.ConsumeMsg(configure.GetConsumerCfg())
		}
	}()

	go func() {
		for {
			select {
			case v := <-msginfo.GetDataChan():
				if v.PubType == "coap" {
					coapMsgHandle(v)

				} else if v.PubType == "kafka" {

					kafkaMsgHandle(v)

				} else if v.PubType == "restapi" {
					restapiHandle(v)
				}
			}
		}

	}()

	beegoRun()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
	<-done
}

func coapMsgHandle(v msginfo.MsgData) {
	method, id, taskid, err := msgjson.CoapMsgDecode(v.Data, v.IsReq)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "device.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("CoapMsgDecode Error: ", err)
		return
	}

	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "device.go",
		"line":     line,
		"func":     f.Name(),
		"method":   method,
		"devid":    id,
		"taskid":   taskid,
		"msgdata":  v.Data,
	}).Info("coapMsgHandle")

	if v.IsReq {
		setUrl(id, v.Url)
	}

	msgid := getMsgid(method)
	msgbyte := msgjson.CoapMsgEncode(v.Data, method, msgid)
	kafka.ProduceMsg(configure.GetProducerCfg(), string(msgbyte))

}

func kafkaMsgHandle(v msginfo.MsgData) {
	msgdata, method, msgid, err := msgjson.KafkaMsgDecode(v.Data)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "device.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("KafkaMsgDecode Error: ", err)
		return
	}

	devid, taskid, err := msgjson.PlatMsgDecode(msgdata)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "device.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("PlatMsgDecode Error: ", err)
		return
	}

	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "device.go",
		"line":     line,
		"func":     f.Name(),
		"msgid":    msgid,
		"method":   method,
		"devid":    devid,
		"taskid":   taskid,
		"msgdata":  msgdata,
	}).Info("kafkaMsgHandle")

	setMsgid(method, msgid)

	if strings.Contains(method, "_Rsp") {
		key := "coap_rsp" + "_" + hostname + "_" + method + "_" + devid + "_" + taskid
		fmt.Println("key: ", key)
		redisclient.SaveRedisByte(key, msgdata)
	} else {
		pubUrl := getUrl(devid)
		buf, err := coapx.PubCoapMsg(msgdata, method, pubUrl)
		if err == nil {
			if len(buf) == 0 {
				pc, _, line, _ := runtime.Caller(0)
				f := runtime.FuncForPC(pc)
				log.WithFields(log.Fields{
					"filename": "device.go",
					"line":     line,
					"func":     f.Name(),
				}).Error("buf is nil: ", buf)
			}
			msg := msginfo.MsgData{
				PubType: "coap",
				Data:    buf,
				IsReq:   false,
				Url:     "",
			}
			coapMsgHandle(msg)
		}
	}

}

func restapiHandle(v msginfo.MsgData) {
	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "device.go",
		"line":     line,
		"func":     f.Name(),
	}).Info("restapiHandle: ", v)

	pubv := strings.Split(v.Url, ":")
	buf, err := coapx.PubCoapMsg(v.Data, pubv[2], pubv[0])
	if err == nil {
		method, id, _, err := msgjson.CoapMsgDecode(buf, false)
		if err != nil {
			pc, _, line, _ := runtime.Caller(0)
			f := runtime.FuncForPC(pc)
			log.WithFields(log.Fields{
				"filename": "device.go",
				"line":     line,
				"func":     f.Name(),
			}).Error("CoapMsgDecode Error: ", err)
		} else {
			key := "coaprestrsp_" + method + "_" + id
			if method == "GetData_Rsp" || method == "CheckTime_Rsp" || method == "Rest_Rsp" || method == "Ctrl_Rsp" {
				redisclient.SaveRedisByte(key, buf)
			}
		}
	}
}

func getUrl(devid string) string {
	key := URL + "_" + hostname + "_" + devid

	log.Info("devurlkey: ", key)
	url, err := redisclient.ReadRedisString(key)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "device.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("ReadRedisString Err: ", err)
		return ""
	}

	log.Info("devurl: ", url)
	return url
}

func setUrl(devid string, puburl string) {

	// key := URL + "_" + hostname + "_" + devid
	// log.Info("devurlkey: ", key)
	// _, err := redisclient.ReadRedisString(key)
	// if err != nil {
	// 	fmt.Println(puburl)
	// 	redisclient.SaveRedisString(key, puburl)
	// 	if err != nil {
	// 		pc, _, line, _ := runtime.Caller(0)
	// 		f := runtime.FuncForPC(pc)
	// 		log.WithFields(log.Fields{
	// 			"filename": "device.go",
	// 			"line":     line,
	// 			"func":     f.Name(),
	// 		}).Error("SaveRedisString Err: ", err)
	// 	}
	// }

	key := URL + "_" + hostname + "_" + devid
	log.Info("devurlkey: ", key)
	fmt.Println(puburl)
	err := redisclient.SaveRedisString(key, puburl)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "device.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("SaveRedisString Err: ", err)
	}

}

func getMsgid(method string) string {
	msgid := "1"
	err := errors.New("error")
	if method == "GetAlarm_Rsp" || method == "GetConfig_Rsp" || method == "SetConfig_Rsp" || method == "GetSysInfo_Rsp" {
		methodtmp := strings.TrimRight(method, "_Rsp")

		key := MSG + "_" + hostname + "_" + methodtmp
		log.Info("msgidkey: ", key)
		msgid, err = redisclient.ReadRedisString(key)
		if err != nil {
			pc, _, line, _ := runtime.Caller(0)
			f := runtime.FuncForPC(pc)
			log.WithFields(log.Fields{
				"filename": "device.go",
				"line":     line,
				"func":     f.Name(),
			}).Error("redis get failed: ", err)
		}
		redisclient.DelRedis(key)
	}
	log.Info("msgid: ", msgid)
	return msgid
}

func setMsgid(method, msgid string) {

	if strings.Contains(method, "_Rsp") {
		return
	}

	key := MSG + "_" + hostname + "_" + method
	log.Info("msgidkey: ", key)
	redisclient.SaveRedisString(key, msgid)
}

func beegoRun() {
	beego.Router("/healthcheck", &controllers.HealthController{})
	beego.Router("/debug", &controllers.DebugController{})
	beego.Router("/loglevelset", &controllers.LogsetController{})

	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/heartbeat", &controllers.HbController{})
	beego.Router("/checktime", &controllers.CkController{})
	beego.Router("/getValue", &controllers.GetvalueController{})
	beego.Router("/setValue", &controllers.SetvalueController{})
	beego.Router("/reset", &controllers.ResetController{})

	beego.Router("/getCfg", &controllers.GetcfgController{})
	beego.Router("/setCfg", &controllers.SetcfgController{})
	beego.Router("/getAlarm", &controllers.GetalarmController{})

	beego.Run()
}

func initConfig() {
	cfg := &CfgPara{}
	data, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.WithFields(log.Fields{
			"filename": "device.go",
		}).Error("ReadFile Error: ", err)
		return
	}

	err = json.Unmarshal(data, cfg)
	if err != nil {
		log.WithFields(log.Fields{
			"filename": "device.go",
		}).Error("Unmarshal Error: ", err)
		return
	}

	configure.GetLocalIp()

	cp := configure.ProducerCfg{
		Topic:         cfg.Reqtopic,
		BrokerList:    cfg.Kafka,
		Key:           *pkey,
		Value:         "test",
		Partitioner:   *ppartitioner,
		Partition:     *ppartition,
		Verbose:       *verbose,
		ShowMetrics:   *pshowMetrics,
		Silent:        *psilent,
		TlsEnabled:    *tlsEnabled,
		TlsSkipVerify: *tlsSkipVerify,
		TlsClientCert: *tlsClientCert,
		TlsClientKey:  *tlsClientKey,
	}
	configure.SetProducerCfg(cp)

	cc := configure.ConsumerCfg{
		Topic:         cfg.Rsptopic,
		BrokerList:    cfg.Kafka,
		Partitions:    *cpartitions,
		Offset:        *coffset,
		Verbose:       *verbose,
		TlsEnabled:    *tlsEnabled,
		TlsSkipVerify: *tlsSkipVerify,
		TlsClientCert: *tlsClientCert,
		TlsClientKey:  *tlsClientKey,
		BufferSize:    *cbufferSize,
	}
	configure.SetConsumerCfg(cc)

	eu := configure.EurekaServer{
		Url: cfg.Eureka,
	}

	ru := configure.RedisServer{
		Url: cfg.Redis,
	}

	configure.SetEurekaIp(eu)
	configure.SetRedisIp(ru)
	configure.PrintCfg()
}

func eurekaStart() {
	eureka.InitEurekaClient()
	eureka.RegEureka()
	eureka.SyncEureka()
	eureka.HeartBeatEureka()
}
