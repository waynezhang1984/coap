package configure

import (
	"net"
	"os"
	"runtime"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
}

type ProducerCfg struct {
	Topic         string
	BrokerList    string
	Key           string
	Value         string
	Partitioner   string
	Partition     int
	Verbose       bool
	ShowMetrics   bool
	Silent        bool
	TlsEnabled    bool
	TlsSkipVerify bool
	TlsClientCert string
	TlsClientKey  string
}

type ConsumerCfg struct {
	Topic         string
	BrokerList    string
	Partitions    string
	Offset        string
	Verbose       bool
	TlsEnabled    bool
	TlsSkipVerify bool
	TlsClientCert string
	TlsClientKey  string
	BufferSize    int
}

type EurekaServer struct {
	Url string
}

type RedisServer struct {
	Url string
}

var (
	pCfg      ProducerCfg
	CCfg      ConsumerCfg
	HostIp    string
	EurekaUrl EurekaServer
	RedisUrl  RedisServer
)

func SetProducerCfg(p ProducerCfg) {
	pCfg = p
}

func SetConsumerCfg(c ConsumerCfg) {
	CCfg = c
}

func GetProducerCfg() ProducerCfg {
	return pCfg
}

func GetConsumerCfg() ConsumerCfg {
	return CCfg
}

func GetIpPool() (map[string]string, error) {
	ips := make(map[string]string)

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range interfaces {
		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			return nil, err
		}

		addresses, err := byName.Addrs()
		for _, v := range addresses {
			ips[byName.Name] = v.String()
			pc, _, line, _ := runtime.Caller(0)
			f := runtime.FuncForPC(pc)
			log.WithFields(log.Fields{
				"filename": "configure.go",
				"line":     line,
				"func":     f.Name(),
				"ip":       v.String(),
			}).Info("PrintCfg")
		}
	}

	return ips, nil

}

func GetLocalIp() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Exit(1)
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				pc, _, line, _ := runtime.Caller(0)
				f := runtime.FuncForPC(pc)
				log.WithFields(log.Fields{
					"filename": "configure.go",
					"line":     line,
					"func":     f.Name(),
					"ip":       ipnet.IP.String(),
				}).Info("PrintCfg")
				HostIp = ipnet.IP.String()
			}
		}
	}
}

func GetHostIp() string {
	return HostIp
}

func SetEurekaIp(url EurekaServer) {
	EurekaUrl = url
}

func GetEurekaUrl() string {

	return EurekaUrl.Url
}

func SetRedisIp(url RedisServer) {
	RedisUrl = url
}

func GetRedisUrl() string {

	return RedisUrl.Url
}

func PrintCfg() {
	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename":       "configure.go",
		"line":           line,
		"func":           f.Name(),
		"brokers":        pCfg.BrokerList,
		"consumer topic": CCfg.Topic,
		"producer topic": pCfg.Topic,
		"eureka server":  EurekaUrl.Url,
		"redis server":   RedisUrl.Url,
	}).Info("PrintCfg")
}
