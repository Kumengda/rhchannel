package platform

//
//type Platform struct {
//	Host            string
//	latestTimeStamp int64
//	deviceID        string
//	tarantola.BaseCrawler
//}
//
//func NewPlatform(token string, host string, deviceID string) *Platform {
//	authHeader := make(map[string]interface{})
//	authHeader["X-Access-Token"] = token
//	return &Platform{
//		BaseCrawler: tarantola.BaseCrawler{
//			BaseOptions: tarantola.BaseOptions{
//				Headers: authHeader,
//			},
//		},
//		Host:            host,
//		deviceID:        deviceID,
//		latestTimeStamp: time.Now().UnixMilli(),
//	}
//}
//
//func (p *Platform) apiUrl(apiPath string) string {
//	return strings.TrimRight(p.Host, "/") + apiPath
//}
//
//func (p *Platform) Crawl() error {
//	for {
//		p.getMessage()
//		time.Sleep(1 * time.Second)
//	}
//	return nil
//}
//func (p *Platform) getMessage() {
//	resp, err := p.HttpRequest.Get(p.apiUrl(MessageApi) + "?pageNo=1&pageSize=100&deviceCode=" + p.deviceID)
//	if err != nil {
//		fmt.Print(err)
//		return
//	}
//	records := gjson.ParseBytes(resp).Get("data").Get("records").Array()
//	for i := len(records) - 1; i >= 0; i-- {
//		b := records[i]
//		createTime := b.Get("createTime").Int()
//		if createTime > p.latestTimeStamp {
//			devices := gjson.Parse(b.Get("content").String()).Get("devices").Array()[0]
//			if devices.Get("deviceId").String() == p.deviceID {
//				p.PushResult(devices.Get("services").Array()[0])
//				p.latestTimeStamp = createTime
//			}
//		}
//
//	}
//
//}
