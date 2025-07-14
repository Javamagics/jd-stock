package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/zhuweitung/jd-stock/message"
	"github.com/zhuweitung/jd-stock/models"
	"github.com/zhuweitung/jd-stock/utils"

	"github.com/go-co-op/gocron"
)

// 定时任务函数
func task() {
	config := utils.GetConfig()
	utils.QueryStock(config.SkuInfos)
}

// 获取配置
func getConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	config := utils.GetConfig()

	// 将配置转换为JSON
	jsonData, err := json.Marshal(config)
	if err != nil {
		http.Error(w, "JSON转换失败", http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// 保存配置
func saveConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 读取请求体
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "读取请求失败", http.StatusBadRequest)
		return
	}

	// 解析JSON到Config结构体
	var config models.Config
	if err := json.Unmarshal(body, &config); err != nil {
		http.Error(w, "JSON解析失败", http.StatusBadRequest)
		return
	}

	// 保存到文件
	err = utils.SaveConfig(&config)
	if err != nil {
		http.Error(w, "配置写入文件失败", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("配置保存成功"))
}

func main() {

	// 加载配置
	cfg, err := utils.LoadConfig("config/config.yaml")
	if err != nil {
		log.Printf("加载配置文件失败: %v", err)
		return
	}

	// 打印解析后的配置
	log.Println("=============配置信息=============")
	log.Printf("|| 间隔执行：%d分钟\n", utils.GetEveryMinutes())
	log.Printf("|| 库存省份：%v\n", cfg.Provinces)
	log.Printf("|| 库存地区编码：%v\n", cfg.AddressCodes)
	log.Printf("|| 监控商品：%v\n", cfg.SkuInfos)
	log.Printf("|| 查询延迟：%d毫秒\n", utils.GetDelay())
	log.Printf("|| 启用通知：%v\n", cfg.EnableNotify)
	if cfg.EnableNotify {
		log.Printf("|| 通知间隔：%d分钟\n", cfg.NotifyInterval)
		log.Printf("|| 通知方式：%s\n", cfg.NotifyType)
		var sender message.Sender
		sender, err = utils.GetSender()
		if err != nil {
			log.Printf("%v", err)
			return
		}
		log.Printf("|| %s配置：%v\n", sender.GetName(), sender)
	}
	log.Printf("|| 当前版本：v1.0.5\n")
	log.Println("================================")

	// 启动时发送测试通知，验证消息渠道是否正常
	utils.SendMessage("【测试】京东库存监控服务启动成功")

	// 加载地区编码
	err = utils.LoadAreaCodes()
	if err != nil {
		log.Printf("加载地区编码失败: %v", err)
		return
	}

	// 设置静态文件服务
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	// 注册配置保存接口
	http.HandleFunc("/api/get-config", getConfig)
	http.HandleFunc("/api/save-config", saveConfig)

	// 启动HTTP服务器
	go func() {
		log.Printf("HTTP服务器启动在 :9090 端口")
		if err := http.ListenAndServe(":9090", nil); err != nil {
			log.Fatal("HTTP服务器启动失败:", err)
		}
	}()

	// 初始化 gocron 调度器
	scheduler := gocron.NewScheduler(time.Local)

	// 每隔 5 分钟执行一次任务
	scheduler.Every(utils.GetEveryMinutes()).Minutes().Do(task)

	// 启动调度器（异步运行）
	scheduler.StartAsync()

	// 阻止主协程退出
	select {}
}
