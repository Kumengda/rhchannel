package mqtt

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
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
}

func NewMyMqttServer(host string, port int, timeout time.Duration, deviceID string) (*MyMqttServer, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
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
	}, nil
}

func (m *MyMqttServer) Start() {
	go func() {
		for {
			fmt.Println("[+]start receive message")
			msg := m.messageChan.receive()
			fmt.Println("[+]receive message")
			fmt.Println(msg)
			jsonBytes, _ := json.Marshal(msg)
			m.conn.Write(jsonBytes)
			//token := m.client.Publish(MessageTopic, 0, false, jsonBytes)
			//token.Wait()
			time.Sleep(1 * time.Second)
		}
	}()
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
