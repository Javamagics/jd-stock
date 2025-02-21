package message

import "fmt"

// BarkSender Bark配置
type BarkSender struct {
	ApiUrl    string `yaml:"apiUrl" json:"apiUrl"`       // API地址
	DeviceKey string `yaml:"deviceKey" json:"deviceKey"` // 设备Key
}

func (_ BarkSender) GetName() string {
	return "Bark"
}

// Send 发送通知
func (sender BarkSender) Send(msg string) error {
	if sender.ApiUrl == "" || sender.DeviceKey == "" {
		return fmt.Errorf("Bark配置缺失")
	}
	body := map[string]interface{}{
		"title": "京东库存监控",
		"body":  msg,
	}
	return SendPost(sender.GetName(), fmt.Sprintf("%s/%s", sender.ApiUrl, sender.DeviceKey), body)
}
