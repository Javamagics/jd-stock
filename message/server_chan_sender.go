package message

import (
	"fmt"
)

// ServerChanSender Server酱配置
type ServerChanSender struct {
	SendKey string `yaml:"sendKey" json:"sendKey"`
}

func (_ ServerChanSender) GetName() string {
	return "Server酱"
}

// Send 发送通知
func (sender ServerChanSender) Send(msg string) error {
	if sender.SendKey == "" {
		return fmt.Errorf("server酱配置缺失")
	}
	body := map[string]interface{}{
		"title": "京东库存监控",
		"desp":  msg,
	}
	url := fmt.Sprintf("https://sctapi.ftqq.com/%s.send", sender.SendKey)
	return SendPost(sender.GetName(), url, body)
}
