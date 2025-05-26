package mqtt

import (
	"database/sql"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql" // 导入 MySQL 驱动
	"github.com/tidwall/gjson"
	"net"
	"strconv"
	"strings"
	"time"
)

type MyMqttServer struct {
	host          string
	port          int
	timeout       time.Duration
	client        mqtt.Client
	clientOptions *mqtt.ClientOptions
	topic         string
	messageChan   *messageQueue
	deviceID      string
	conn          net.Conn
	db            *sql.DB
}

func NewMyMqttServer(host string, port int, dbname, dbuser, dbpass string, timeout time.Duration, deviceID string, isSql bool) (*MyMqttServer, error) {
	var conn net.Conn
	var err error
	var db *sql.DB
	if isSql {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbuser, dbpass, host, port, dbname)
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			return nil, err
		}
	} else {
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			return nil, err
		}
	}

	//broker := fmt.Sprintf("tcp://%s:%d", host, port)
	//opts := mqtt.NewClientOptions()
	//opts.SetAutoReconnect(true)
	//opts.AddBroker(broker)
	//opts.SetConnectTimeout(timeout)
	//opts.SetClientID(randomClientID())
	//client := mqtt.NewClient(opts)
	//if token := client.Connect(); token.Wait() && token.Error() != nil {
	//	return nil, token.Error()
	//}
	return &MyMqttServer{
		host:    host,
		port:    port,
		timeout: timeout,
		//clientOptions: opts,
		messageChan: newMessageQueue(1),
		deviceID:    deviceID,
		topic:       MessageTopic,
		//client:        client,
		conn: conn,
		db:   db,
	}, nil
}

func (m *MyMqttServer) Start() {
	go func() {
		for {
			fmt.Println("[+]start receive message")
			msg := m.messageChan.receive().(ReportMessage)
			fmt.Println("[+]receive message")
			jsonBytes, _ := json.Marshal(msg.Param.Data.Msg)
			unquoted, _ := strconv.Unquote(string(jsonBytes))
			if strings.Contains(unquoted, "网络流量异常告警") {
				res := gjson.Parse(unquoted)
				m.insertFlowWarning(res.Get("time").String(),
					res.Get("type").String(),
					res.Get("attacktype").String(),
					res.Get("sip").String(),
					res.Get("tip").String(),
					res.Get("tport").String(),
					res.Get("sport").String(),
					res.Get("schme").String(),
					res.Get("source").String(),
					res.Get("link_type").String(),
					msg.DeviceId)
			} else {
				res := gjson.Parse(unquoted)
				m.insertProcessWarning(res.Get("time").String(),
					res.Get("type").String(),
					res.Get("detail").String(),
					res.Get("path").String(),
					res.Get("hash").String(),
					msg.DeviceId,
				)
			}

			//jsonBytes, _ := json.Marshal(msg)
			//jsonBytes = append(jsonBytes, '\n')
			//m.conn.Write(jsonBytes)
			//token := m.client.Publish(MessageTopic, 0, false, jsonBytes)
			//token.Wait()
			time.Sleep(1 * time.Second)
		}
	}()
}
func (m *MyMqttServer) insertFlowWarning(time string, ttype, attackType, sip, tip, tport, sport, schem, source, linktype, agetid string) {
	insertSQL := `
		INSERT INTO flow_warning (time, type,attackType,sip, tip,tport,sport,agetid,schem,source,linktype)
		VALUES (?, ?, ?, ?,?,?,?,?,?,?,?)
	`
	_, err := m.db.Exec(insertSQL, time, attackType, sip, tip, agetid)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (m *MyMqttServer) insertProcessWarning(time string, ttype, detail, path, hash, agetid string) {
	insertSQL := `
		INSERT INTO terminal_warning (time, type,detail,path, hash,agentid)
		VALUES (?,?,?,?,?,?)
	`
	_, err := m.db.Exec(insertSQL, time, ttype, detail, path, hash, agetid)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (m *MyMqttServer) Push(mtype MsgType, mdata Tdata) {
	m.messageChan.send(ReportMessage{
		Mid:       generateRandom10DigitNumber(),
		Type:      "AGENT_" + mtype,
		Timestamp: mdata.TimeStamp,
		DeviceId:  m.deviceID,
		Expire:    -1,
		Param: Paramm{
			DeviceId: m.deviceID,
			Event:    "resource_alarm",
			Data:     mdata,
		},
	})
}
