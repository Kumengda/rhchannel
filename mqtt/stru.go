package mqtt

type MsgType string

const (
	TestType MsgType = "TEST"
)

type Tdata struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`

	TimeStamp int64 `json:"timestamp"`
}

type Paramm struct {
	DeviceId string `json:"deviceId "`
	Event    string `json:"event"`
	Data     Tdata  `json:"data"`
}

type ReportMessage struct {
	Mid       int64   `json:"mid"`
	Type      MsgType `json:"type"`
	Timestamp int64   `json:"timestamp"`
	DeviceId  string  `json:"deviceId"`
	Expire    int     `json:"expire"`
	Param     Paramm  `json:"param"`
}
