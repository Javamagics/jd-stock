package message

import (
	"fmt"
)

// ServerChanSender Server酱配置
type ServerChanSender struct {
	// SendKeys 支持配置一个或多个 sendKey
	// YAML 中既可以写成 sendKey: SCTxxx，也可以写成 sendKey:
	//   - SCTxxx
	//   - SCTyyy
	SendKeys []string `yaml:"sendKey" json:"sendKey"`
}

func (_ ServerChanSender) GetName() string {
	return "Server酱"
}

// Send 发送通知
func (sender ServerChanSender) Send(msg string) error {
	if len(sender.SendKeys) == 0 {
		return fmt.Errorf("server酱配置缺失")
	}
	body := map[string]interface{}{
		"title": "京东库存监控",
		"desp":  msg,
	}
	var lastErr error
	for _, key := range sender.SendKeys {
		if key == "" {
			continue
		}
		url := fmt.Sprintf("https://sctapi.ftqq.com/%s.send", key)
		if err := SendPost(sender.GetName(), url, body); err != nil {
			lastErr = err
		}
	}
	return lastErr
}
