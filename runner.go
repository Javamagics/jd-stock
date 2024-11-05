package main

import (
	"github.com/zhuweitung/jd-stock/message"
	"github.com/zhuweitung/jd-stock/utils"
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

// 定时任务函数
func task() {
	config := utils.GetConfig()
	utils.QueryStock(config.SkuInfos)
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

	// 加载地区编码
	err = utils.LoadAreaCodes()
	if err != nil {
		log.Printf("加载地区编码失败: %v", err)
		return
	}

	// 初始化 gocron 调度器
	scheduler := gocron.NewScheduler(time.Local)

	// 每隔 5 分钟执行一次任务
	scheduler.Every(utils.GetEveryMinutes()).Minutes().Do(task)

	// 启动调度器（异步运行）
	scheduler.StartAsync()

	// 阻止主协程退出
	select {}
}
