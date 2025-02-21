package message

import (
	"fmt"
)

// PushPlusSender PushPlus配置
type PushPlusSender struct {
	Token string `yaml:"token" json:"token"`
}

func (_ PushPlusSender) GetName() string {
	return "PushPlus"
}

// Send 发送通知
func (sender PushPlusSender) Send(msg string) error {
	if sender.Token == "" {
		return fmt.Errorf("PushPlus配置缺失")
	}
	body := map[string]interface{}{
		"token":    sender.Token,
		"content":  msg,
		"template": "txt",
	}
	return SendPost(sender.GetName(), "https://www.pushplus.plus/send", body)
}
