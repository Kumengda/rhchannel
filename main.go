package main

import (
	"bufio"
	"fmt"
	"github.com/Kumengda/Tarantola/request"
	"github.com/Kumengda/Tarantola/tarantola"
	"github.com/Kumengda/rhchannel/mqtt"
	"github.com/Kumengda/rhchannel/platform"
	"github.com/tidwall/gjson"
	"os"
	"time"
)

func main() {
	go func() {
		plantformClient := platform.NewPlatform("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjE3NDAxNTY5MzAsInVzZXJuYW1lIjoiYWRtaW4ifQ.6TmhQwJfkA7nuUdOQuaEjhFbOSGVKfumqVaE-yztVnU", "http://192.168.10.20", "1001204907901807")
		plantformClient.SetDataProcessFunc(func(crawlRes interface{}, httpRequest *request.HttpRequest) error {
			fmt.Print(crawlRes.(gjson.Result).Get("data").Get("msg"))
			return nil
		})

		t1 := tarantola.NewTarantola()
		t1.AddCrawler(plantformClient)
		t1.MonoCrawl(1)
	}()

	server, err := mqtt.NewMyMqttServer("192.168.10.101", 1883, 10*time.Second, "1001204907901807")
	if err != nil {
		panic(err)
		return
	}

	server.Start()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		msg := scanner.Text()
		server.Push(mqtt.TestType, mqtt.Tdata{
			Type:      "testData",
			Msg:       msg, // 发送用户输入的消息
			TimeStamp: time.Now().UnixMilli(),
		})

	}

}
