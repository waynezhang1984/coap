package msginfo

//	"github.com/mqttprocess/types"

type MsgData struct {
	PubType string
	Data    []byte
	IsReq   bool
	Url     string
}

type RestReqMsg struct {
	Method string
	Devid  string
}

type RestRspMsg struct {
	Method string
	Devid  string
	Error  int
}

var (
	DataChan         = make(chan MsgData, 10)
	RestGetDataChan  = make(chan []byte)
	RestGetOtherChan = make(chan []byte)
)

var RestDeviceBuf []string
var RestReqBuf []RestReqMsg

func GetDataChan() chan MsgData {
	return DataChan
}

func GetRestDataChan() chan []byte {
	return RestGetDataChan
}

func GetRestOtherChan() chan []byte {
	return RestGetOtherChan
}

func GetRestDeviceBuf() []string {
	return RestDeviceBuf
}

func SetRestDeviceBuf(id string) {
	RestDeviceBuf = append(RestDeviceBuf, id)
}

func DelRestDeviceBuf(id string) {
	for k, v := range RestDeviceBuf {
		if v == id {
			RestDeviceBuf = append(RestDeviceBuf[:k], RestDeviceBuf[k+1:]...)
			break
		}
	}
}

func GetRestReqBuf() []RestReqMsg {
	return RestReqBuf
}

func SetRestReqBuf(id RestReqMsg) {
	RestReqBuf = append(RestReqBuf, id)
}
