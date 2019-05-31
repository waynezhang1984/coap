package eureka

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/ArthurHlt/go-eureka-client/eureka"
	"github.com/coapprocess/configure"
	log "github.com/sirupsen/logrus"
)

var (
	AppList        *eureka.Applications
	EurekaClient   *eureka.Client
	EurekaInstance *eureka.InstanceInfo
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
}
func GetEurekaApps() *eureka.Applications {
	return AppList
}

func InitEurekaClient() {

	eurekaurl := configure.GetHostIp()
	pc, _, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "eureka.go",
		"line":     line,
		"func":     f.Name(),
	}).Info("EurekaClient IP: ", eurekaurl)

	serverip := configure.GetEurekaUrl()
	serverip = "http://" + serverip + "/eureka"
	pc, _, line, _ = runtime.Caller(0)
	f = runtime.FuncForPC(pc)
	log.WithFields(log.Fields{
		"filename": "eureka.go",
		"line":     line,
		"func":     f.Name(),
	}).Info("EurekaServer IP: ", serverip)

	EurekaClient = eureka.NewClient([]string{
		serverip, //From a spring boot based eureka server
		// add others servers here
	})
	//hostname := eurekaurl + ":9528:1.0.0.RELEASE"
	EurekaInstance = eureka.NewInstanceInfo(eurekaurl, "coapprocess-v100", eurekaurl, 9528, 30, false) //Create a new instance to register
	EurekaInstance.Metadata = &eureka.MetaData{
		Map: make(map[string]string),
	}
	EurekaInstance.Metadata.Map["management.port"] = "9528" //add metadata for example

}

func RegEureka() {
	repcnt := 0
	err := EurekaClient.RegisterInstance("coapprocess-v100", EurekaInstance) // Register new instance in your eureka(s)
	for err != nil {
		select {
		case <-time.After(10 * time.Second):

			err = EurekaClient.RegisterInstance("coapprocess-v100", EurekaInstance)
			if err == nil {
				repcnt = 0
				return
			}

			if repcnt >= 5 {
				return
			}
			repcnt++
		}
	}

	if repcnt > 0 {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "eureka.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("register error: ", err)
		return
	}

	applications, err1 := EurekaClient.GetApplications()
	if err1 != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "eureka.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("GetApplications error: ", err1)
	} else {
		AppList = applications
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		for _, app := range applications.Applications {
			//fmt.Println(app.Name)
			log.WithFields(log.Fields{
				"filename": "eureka.go",
				"line":     line,
				"func":     f.Name(),
			}).Info("appname: ", app.Name)
		}
	}

	_, err = EurekaClient.GetApplication(EurekaInstance.App)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "eureka.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("GetApplications error: ", err)
	}

	_, err = EurekaClient.GetInstance(EurekaInstance.App, EurekaInstance.HostName)
	if err != nil {
		pc, _, line, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		log.WithFields(log.Fields{
			"filename": "eureka.go",
			"line":     line,
			"func":     f.Name(),
		}).Error("GetInstance error: ", err)
	}

	EurekaClient.SendHeartbeat(EurekaInstance.App, EurekaInstance.HostName) // say to eureka that your app is alive (here you must send heartbeat before 30 sec)

}

func SyncEureka() {
	ticker := time.NewTicker(30 * time.Second)
	quit := make(chan int)

	go func() {

		for {
			select {
			case <-ticker.C:
				//fmt.Println("ticker .")
				applications, err := EurekaClient.GetApplications() // Retrieves all applications from eureka server(s)
				if err == nil {
					AppList = applications
					pc, _, line, _ := runtime.Caller(0)
					f := runtime.FuncForPC(pc)
					for _, app := range applications.Applications {
						log.WithFields(log.Fields{
							"filename": "eureka.go",
							"line":     line,
							"func":     f.Name(),
						}).Info("appname: ", app.Name)
					}
				}

			case <-quit:
				fmt.Println("work well .")
				ticker.Stop()
				return
			}
		}

	}()
}

func HeartBeatEureka() {
	ticker := time.NewTicker(30 * time.Second)
	quit := make(chan int)

	go func() {

		for {
			select {
			case <-ticker.C:
				//fmt.Println("ticker .")
				err := EurekaClient.SendHeartbeat(EurekaInstance.App, EurekaInstance.HostName)
				if err != nil {
					fmt.Println("SendHeartbeat error!", err)
				}
			case <-quit:
				fmt.Println("work well .")
				ticker.Stop()
				return
			}
		}

	}()
}
