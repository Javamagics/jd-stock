package models

import (
	"encoding/json"

	"github.com/zhuweitung/jd-stock/message"
)

// AreaInfo 地区信息
type AreaInfo struct {
	ID    json.Number `json:"id"`    // 地区ID
	PID   json.Number `json:"pid"`   // 父级ID
	Level int         `json:"level"` // 层级
	Name  string      `json:"name"`  // 地区名称
}

// SkuInfo 商品信息
type SkuInfo struct {
	FreshEdi       *string `json:"freshEdi"`
	SidDely        string  `json:"sidDely"`
	Channel        int     `json:"channel"`
	Rid            *string `json:"rid"`
	Sid            string  `json:"sid"`
	DcId           string  `json:"dcId"`
	IsPurchase     bool    `json:"IsPurchase"`
	Eb             string  `json:"eb"`
	Ec             string  `json:"ec"`
	StockState     int     `json:"StockState"`
	Ab             string  `json:"ab"`
	CanAddCart     string  `json:"canAddCart"`
	Ac             string  `json:"ac"`
	Ad             string  `json:"ad"`
	Ae             string  `json:"ae"`
	SkuState       int     `json:"skuState"`
	PopType        int     `json:"PopType"`
	Af             string  `json:"af"`
	Ag             string  `json:"ag"`
	StockStateName string  `json:"StockStateName"`
	M              string  `json:"m"`
	Rfg            int     `json:"rfg"`
	ArrivalDate    string  `json:"ArrivalDate"`
	V              string  `json:"v"`
	Rn             int     `json:"rn"`
	Dc             string  `json:"dc"`
}

// Config 定义模型结构体，映射 YAML 配置文件
type Config struct {
	EveryMinutes      int                       `yaml:"everyMinutes" json:"everyMinutes"`     // 每隔N分钟执行
	Provinces         []string                  `yaml:"provinces" json:"provinces"`           // 库存省份
	AddressCodes      []string                  `yaml:"addressCodes" json:"addressCodes"`     // 库存地区编码
	SkuInfos          []CustomSkuInfo           `yaml:"skuInfos" json:"skuInfos"`             // 商品信息列表
	Delay             int                       `yaml:"delay" json:"delay"`                   // 查询延迟（毫秒）
	Ua                string                    `yaml:"ua" json:"ua"`                         // 用户代理字符串
	EnableNotify      bool                      `yaml:"enableNotify" json:"enableNotify"`     // 是否启用通知
	NotifyInterval    int                       `yaml:"notifyInterval" json:"notifyInterval"` // 通知间隔（分钟），0表示允许重复提醒
	NotifyType        string                    `yaml:"notifyType" json:"notifyType"`         // 通知方式
	DingtalkBotSender message.DingtalkBotSender `yaml:"dingtalkBot" json:"dingtalkBotSender"` // 钉钉机器人配置
	QyWechatBotSender message.QyWechatBotSender `yaml:"qyWechatBot" json:"qyWechatBotSender"` // 企业微信机器人配置
	PushPlusSender    message.PushPlusSender    `yaml:"pushPlus" json:"pushPlusSender"`       // PushPlus配置
	ServerChanSender  message.ServerChanSender  `yaml:"serverChan" json:"serverChanSender"`   // Server酱配置
	BarkSender        message.BarkSender        `yaml:"bark" json:"barkSender"`               // bark配置
}

// CustomSkuInfo 自定义商品信息
type CustomSkuInfo struct {
	Id   string `yaml:"id" json:"id"`     // 商品id（用于查询库存）
	Name string `yaml:"name" json:"name"` // 商品名称（用于提示）
	Link string `yaml:"link" json:"link"` // 商品跳转链接（用于通知）
}

// MessageCache 通知缓存
type MessageCache struct {
	Content string `json:"Content"` // 通知内容（md5）
	Time    string `json:"time"`    // 通知时间
}
