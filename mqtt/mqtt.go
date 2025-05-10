package mqtt

import (
	"database/sql"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql" // 导入 MySQL 驱动
	"github.com/tidwall/gjson"
	"net"
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
			fmt.Println(string(jsonBytes))
			res := gjson.ParseBytes(jsonBytes)
			fmt.Println(res.Get("time").String())
			fmt.Println(res.Get("type").String())
			fmt.Println(res.Get("sip").String())
			fmt.Println(res.Get("tip").String())
			m.insert(res.Get("time").String(), res.Get("type").String(), res.Get("sip").String(), res.Get("tip").String(), msg.DeviceId)
			//jsonBytes, _ := json.Marshal(msg)
			//jsonBytes = append(jsonBytes, '\n')
			//m.conn.Write(jsonBytes)
			//token := m.client.Publish(MessageTopic, 0, false, jsonBytes)
			//token.Wait()
			time.Sleep(1 * time.Second)
		}
	}()
}
func (m *MyMqttServer) insert(time string, attackType string, sip string, tip string, agentid string) {
	insertSQL := `
		INSERT INTO flow_warning (time, type,sip, tip,agentid)
		VALUES (?, ?, ?, ?,?)
	`
	_, err := m.db.Exec(insertSQL, time, attackType, sip, tip, agentid)
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
