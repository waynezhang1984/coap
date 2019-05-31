package types

const (
	DsTopic     string = "BW_MQTT_DEVICE_V1"
	Device2Plat string = "mqtt_topic_access_process"
	Plat2Device string = "mqtt_topic_process_access"
)

type DeviceInfo struct {
	DeviceType  string `json:"type"`
	DeviceId    int    `json:"id"`
	DeviceName  string `json:"name"`
	DeviceDesc  string `json:"desc"`
	DeviceChild []DeviceInfo
}

type LoginReq struct {
	Method string       `json:"method"`
	Taskid string       `json:"taskId"`
	DevId  string       `json:"devId"`
	Data   []DeviceInfo `json:"data"`
}

type Period struct {
	HeartBeat  int `json:"heartbeatPeriod"`
	DataReport int `json:"dataReportPeriod"`
}

type LoginRsp struct {
	Method string `json:"method"`
	Error  int    `json:"code"`
	Taskid string `json:"taskId"`
	DevId  string `json:"devId"`
	Data   Period `json:"data"`
}

type HeartBeatReq struct {
	Method string `json:"method"`
	Taskid string `json:"taskId"`
	DevId  string `json:"devId"`
}

type HeartBeatRsp struct {
	Method string `json:"method"`
	Error  int    `json:"code"`
	Taskid string `json:"taskId"`
	DevId  string `json:"devId"`
}

type DeviceValue struct {
	AttrId string `json:"attrId"`
	Value  string `json:"value"`
	Time   string `json:"time"`
}

type DevicePara struct {
	Id     string        `json:"id"`
	Values []DeviceValue `json:"values"`
}

type DataReport struct {
	Method string       `json:"method"`
	Taskid string       `json:"taskId"`
	DevId  string       `json:"devId"`
	Data   []DevicePara `json:"data"`
}

type DataReportRsp struct {
	Method string `json:"method"`
	Error  int    `json:"code"`
	Taskid string `json:"taskId"`
	DevId  string `json:"devId"`
}

type AlertData struct {
	Id        string `json:"id"`
	AttrId    string `json:"attrId"`
	Timestamp string `json:"time"`
	Flag      string `json:"flag"`
	Value     string `json:"value"`
	Desc      string `json:"desc"`
}

type AlertReport struct {
	Method string    `json:"method"`
	Taskid string    `json:"taskId"`
	DevId  string    `json:"devId"`
	Data   AlertData `json:"data"`
}

type AlertReportRsp struct {
	Method string `json:"method"`
	Error  int    `json:"code"`
	Taskid string `json:"taskId"`
	DevId  string `json:"devId"`
}

type CfgInfo struct {
	AttrId      string  `json:"attrId"`
	Min         string  `json:"min"`
	Max         string  `json:"max"`
	NormalValue string  `json:"normalValue"`
	Kv          float64 `json:"k"`
	Bv          float64 `json:"b"`
}

type CfgInfos struct {
	Id   string    `json:"id"`
	Cfgs []CfgInfo `json:"cfgs"`
}

type CfgReport struct {
	Method string     `json:"method"`
	Taskid string     `json:"taskId"`
	DevId  string     `json:"devId"`
	Data   []CfgInfos `json:"data"`
}

type CfgReportRsp struct {
	Method string `json:"method"`
	Error  int    `json:"code"`
	Taskid string `json:"taskId"`
	DevId  string `json:"devId"`
}

type ControlData struct {
	Id     string `json:"id"`
	AttrId string `json:"attrId"`
	Value  string `json:"value"`
}

type ControlInfoReq struct {
	Method string      `json:"method"`
	Taskid string      `json:"taskId"`
	DevId  string      `json:"devId"`
	Data   ControlData `json:"data"`
}

type ControlInfoRsp struct {
	Method string `json:"method"`
	Error  int    `json:"code"`
	Taskid string `json:"taskId"`
	DevId  string `json:"devId"`
}

type DeviceData struct {
	Id     string `json:"id"`
	AttrId string `json:"attrId"`
}

type DeviceDatasReq struct {
	Method string       `json:"method"`
	Taskid string       `json:"taskId"`
	DevId  string       `json:"devId"`
	Data   []DeviceData `json:"data"`
}

type DeviceDatasRsp struct {
	Method string       `json:"method"`
	Error  int          `json:"code"`
	Taskid string       `json:"taskId"`
	DevId  string       `json:"devId"`
	Data   []DevicePara `json:"data"`
}

type AlertMeta struct {
	Type      string `json:"type"`
	Id        string `json:"id"`
	StartTime string `json:"sTime"`
	EndTime   string `json:"eTime"`
}

type AlertQueryReq struct {
	Method string    `json:"method"`
	Taskid string    `json:"taskId"`
	DevId  string    `json:"devId"`
	Data   AlertMeta `json:"data"`
}

type AlarmInfo struct {
	AttrId    string `json:"attrId"`
	AlarmTime string `json:"alarmTime"`
	ClearTime string `json:"clearTime"`
	Flag      string `json:"flag"`
	Value     string `json:"value"`
	Desc      string `json:"desc"`
}

type AlertRspMeta struct {
	Id     string      `json:"id"`
	Alarms []AlarmInfo `json:"alarms"`
}

type AlertQueryRsp struct {
	Method string         `json:"method"`
	Error  int            `json:"code"`
	Taskid string         `json:"taskId"`
	DevId  string         `json:"devId"`
	Data   []AlertRspMeta `json:"data"`
}

type CfgQueryReq struct {
	Method string     `json:"method"`
	Taskid string     `json:"taskId"`
	DevId  string     `json:"devId"`
	Data   DeviceData `json:"data"`
}

type CfgQueryRsp struct {
	Method string     `json:"method"`
	Error  int        `json:"code"`
	Taskid string     `json:"taskId"`
	DevId  string     `json:"devId"`
	Data   []CfgInfos `json:"data"`
}

type CfgSetInfo struct {
	Id          string  `json:"id"`
	AttrId      string  `json:"attrId"`
	Min         float64 `json:"min"`
	Max         float64 `json:"max"`
	NormalValue float64 `json:"normalValue"`
	Kv          float64 `json:"k"`
	Bv          float64 `json:"b"`
}

type CfgSetReq struct {
	Method string     `json:"method"`
	Taskid string     `json:"taskId"`
	DevId  string     `json:"devId"`
	Data   CfgSetInfo `json:"data"`
}

type CfgSetRsp struct {
	Method string `json:"method"`
	Error  int    `json:"code"`
	Taskid string `json:"taskId"`
	DevId  string `json:"devId"`
}

type UpdateInfo struct {
	Type       string `json:"type"`
	ServerIp   string `json:"serverIp"`
	ServerPort string `json:"serverPort"`
	FileName   string `json:"fileName"`
	FilePath   string `json:"filePath"`
	FileSize   string `json:"fileSize"`
	Version    string `json:"version"`
}

type DeviceUpdateReq struct {
	Method string     `json:"method"`
	Taskid string     `json:"taskId"`
	DevId  string     `json:"devId"`
	Data   UpdateInfo `json:"data"`
}

type DeviceUpdateRsp struct {
	Method string `json:"method"`
	Error  int    `json:"code"`
	Taskid string `json:"taskId"`
	DevId  string `json:"devId"`
}

type DeviceUpdNotifyReq struct {
	Method string        `json:"method"`
	Taskid string        `json:"taskId"`
	DevId  string        `json:"devId"`
	Error  int           `json:"code"`
	Data   UpdateVersion `json:"data"`
}

type UpdateVersion struct {
	Version string `json:"version"`
}

type DeviceUpdaNotifyRsp struct {
	Method string `json:"method"`
	Error  int    `json:"code"`
	Taskid string `json:"taskId"`
	DevId  string `json:"devId"`
}

type CheckTimeReq struct {
	Method string `json:"method"`
	Taskid string `json:"taskId"`
	DevId  string `json:"devId"`
	Data   string `json:"data"`
}

type CheckTimeRsp struct {
	Method string `json:"method"`
	Error  int    `json:"code"`
	Taskid string `json:"taskId"`
	DevId  string `json:"devId"`
}

type ResetInfo struct {
	Type  string `json:"type"`
	Delay int    `json:"delay"`
}

type ResetReq struct {
	Method string    `json:"method"`
	Taskid string    `json:"taskId"`
	DevId  string    `json:"devId"`
	Data   ResetInfo `json:"data"`
}

type ResetRsp struct {
	Method string `json:"method"`
	Error  int    `json:"code"`
	Taskid string `json:"taskId"`
	DevId  string `json:"devId"`
}

type CommonRsp struct {
	Method string `json:"method"`
	Code   int    `json:"code"`
	Taskid string `json:"taskId"`
	DevId  string `json:"devId"`
}

type DeviceMessage struct {
	Method  string `json:"method"`
	MsgId   string `json:"msgid"`
	Data    []byte
	Msgfrom string `json:"msgfrom"`
}

type CommonReq struct {
	Method string `json:"method"`
	Taskid string `json:"taskId"`
	DevId  string `json:"devId"`
}

type RestData struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RestLogin struct {
	Data RestData
}

type RestHeatbeat struct {
	Data RestData
}

type RestChecktimeData struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Time     string `json:"time"`
}

type RestChecktime struct {
	Data RestChecktimeData
}

type RestGetvalueData struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Id       string `json:"id"`
	AttrId   string `json:"attrid"`
}

type RestGetvalue struct {
	Data RestGetvalueData
}

type RestGetcfgData struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Id       string `json:"id"`
	AttrId   string `json:"attrid"`
}

type RestGetcfg struct {
	Data RestGetcfgData
}

type RestGetalarmData struct {
	Ip        string `json:"ip"`
	Port      string `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Type      string `json:"type"`
	Id        string `json:"id"`
	StartTime string `json:"sTime"`
	EndTime   string `json:"eTime"`
}

type RestGetalarm struct {
	Data RestGetalarmData
}

type RestSetvalueData struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Id       string `json:"id"`
	AttrId   string `json:"attrid"`
	Value    string `json:"value"`
}

type RestSetvalue struct {
	Data RestSetvalueData
}

type RestSetcfgData struct {
	Ip          string  `json:"ip"`
	Port        string  `json:"port"`
	Username    string  `json:"username"`
	Password    string  `json:"password"`
	Id          string  `json:"id"`
	AttrId      string  `json:"attrId"`
	Min         float64 `json:"min"`
	Max         float64 `json:"max"`
	NormalValue float64 `json:"normalValue"`
	Kv          float64 `json:"k"`
	Bv          float64 `json:"b"`
}

type RestSetcfg struct {
	Data RestSetcfgData
}

type RestResetData struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"`
	Delay    int    `json:"delay"`
}

type RestReset struct {
	Data RestResetData
}

type RestRspNoData struct {
	Code int `json:"error"`
}

type RestRspGetValue struct {
	Code int          `json:"code"`
	Data []DevicePara `json:"data"`
}

type RestRspGetAlarm struct {
	Code int            `json:"code"`
	Data []AlertRspMeta `json:"data"`
}

type RestRspGetCfg struct {
	Code int        `json:"code"`
	Data []CfgInfos `json:"data"`
}

type RestDevice struct {
	DevId  string
	Method string
	IsRsp  bool
}
