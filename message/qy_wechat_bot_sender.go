package message

import (
	"fmt"
)

var (
	qyWechatBotWebhookUrl = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s"
)

// QyWechatBotSender 企业微信机器人配置
type QyWechatBotSender struct {
	Key string `yaml:"key" json:"key"`
}

func (_ QyWechatBotSender) GetName() string {
	return "企业微信机器人"
}

// Send 发送通知
func (sender QyWechatBotSender) Send(msg string) error {
	if sender.Key == "" {
		return fmt.Errorf("企业微信机器人配置缺失")
	}
	body := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": msg,
		},
	}
	return SendPost(sender.GetName(), fmt.Sprintf(qyWechatBotWebhookUrl, sender.Key), body)
}
