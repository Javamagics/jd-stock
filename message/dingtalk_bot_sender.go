package message

import (
	"fmt"
	"github.com/wanghuiyt/ding"
	"log"
)

// DingtalkBotSender 钉钉机器人配置
type DingtalkBotSender struct {
	Token  string `yaml:"token" json:"token"`
	Secret string `yaml:"secret" json:"secret"`
}

func (_ DingtalkBotSender) GetName() string {
	return "钉钉机器人"
}

// Send 发送通知
func (sender DingtalkBotSender) Send(msg string) error {
	if sender.Token == "" || sender.Secret == "" {
		return fmt.Errorf("钉钉机器人配置缺失")
	}
	d := ding.Webhook{
		AccessToken: sender.Token,
		Secret:      sender.Secret,
	}
	err := d.SendMessageText(msg)
	if err != nil {
		log.Printf("钉钉机器人通知发送异常：%v", err)
		return err
	}
	return nil
}
