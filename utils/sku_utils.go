package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/zhuweitung/jd-stock/models"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var (
	// 正则表达式，用于提取 JSON 数据部分
	skuJsonPattern = regexp.MustCompile(`\((\{.*})\)`)
)

// QueryStock 查询库存
func QueryStock(customSkuInfos []models.CustomSkuInfo) {
	config := GetConfig()

	provinceNames := config.Provinces
	addressCodes := config.AddressCodes

	var areaCodeCombinations []string
	if addressCodes != nil && len(addressCodes) > 0 {
		// 只查询用户在配置中显式指定的地区编码
		areaCodeCombinations = addressCodes
	} else {
		// 如果未配置具体地区编码，则退化为随机抽取各省一条组合编码
		areaCodeCombinations = GetRandomCodeCombination(provinceNames)
	}
	stockAreaNames := make(map[string][]string)

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", "https://api.m.jd.com/stocks", nil)
	if err != nil {
		log.Printf("请求创建失败: %v", err)
		return
	}

	for index, areaCodeCombination := range areaCodeCombinations {
		q := req.URL.Query()
		q.Add("type", "getstocks")
		q.Add("skuIds", getSkuIds(customSkuInfos))
		q.Add("appid", "item-v3")
		q.Add("functionId", "pc_stocks")
		q.Add("callback", "jQuery111107584463972365898_1729065548044")
		q.Add("area", areaCodeCombination)
		q.Add("_", fmt.Sprint(time.Now().UnixMilli()))
		req.URL.RawQuery = q.Encode()
		req.Header.Set("User-Agent", GetConfig().Ua)

		// 发送请求
		resp, _ := client.Do(req)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		contentType := resp.Header.Get("Content-Type")
		var response string
		if strings.Contains(strings.ToLower(contentType), "gbk") {
			response, _ = convertGBKToUTF8(body)
		} else {
			response = string(body)
		}
		groups := skuJsonPattern.FindStringSubmatch(response)
		areaName := ""
		if contains(areaCodeCombination, addressCodes) {
			// 自定义地区获取完整地址
			areaName = GetAreaNameByCodeCombination(areaCodeCombination)
		} else {
			// 获取省份
			area, err := GetAreaByID(strings.Split(areaCodeCombination, "_")[0])
			if err != nil {
				log.Printf("%v", err)
				return
			}
			areaName = area.Name
		}

		if len(groups) > 1 {
			var skuInfoMap map[string]models.SkuInfo
			if err := json.Unmarshal([]byte(groups[1]), &skuInfoMap); err != nil {
				log.Printf("%s：查询异常，response=%s", areaName, response)
				continue
			}

			for _, customSkuInfo := range customSkuInfos {
				skuId := customSkuInfo.Id
				skuInfo, ok := skuInfoMap[skuId]
				if !ok {
					continue
				}
				stockStateName := skuInfo.StockStateName
				isPurchase := skuInfo.IsPurchase
				skuState := skuInfo.SkuState

				purchaseStr := "可购买"
				if !isPurchase {
					purchaseStr = "不可购买"
				}

				// 记录更详细的调试信息，包含 skuState
				log.Printf("[%s] %s %s：%s %s skuState=%d", skuId, customSkuInfo.Name, areaName, stockStateName, purchaseStr, skuState)

				if skuState == 1 {
					stockAreaNames[skuId] = append(stockAreaNames[skuId], areaName)
				}
			}
		} else {
			log.Printf("%s：查询异常，response=%s", areaName, response)
		}

		if index != len(areaCodeCombinations)-1 {
			time.Sleep(time.Duration(GetDelay()) * time.Millisecond)
		}
	}

	var messages []string
	for _, customSkuInfo := range customSkuInfos {
		skuId := customSkuInfo.Id
		areaNames := stockAreaNames[skuId]
		intersection := getIntersection(provinceNames, areaNames)
		var line string
		if len(areaNames) > 0 {
			showNames := areaNames
			if len(intersection) > 0 {
				// 如果在重点省份中，就只显示重点省份
				showNames = intersection
			}
			line = fmt.Sprintf("商品 [%s] %s 在 %s 地区有现货！", skuId, customSkuInfo.Name, strings.Join(showNames, "、"))
		}
		if line != "" {
			if customSkuInfo.Link != "" {
				// Server酱支持 Markdown，追加购买链接
				line = fmt.Sprintf("%s [立即购买](%s)", line, customSkuInfo.Link)
			}
			messages = append(messages, line+"\n")
		}
	}
	if len(messages) > 0 {
		message := strings.Join(messages, "\n")
		log.Printf("%s", message)
		SendMessage(message)
	} else {
		log.Printf("商品 %v 无货或已下架...", customSkuInfos)
	}
}

// 获取商品ids
func getSkuIds(skuInfos []models.CustomSkuInfo) string {
	if len(skuInfos) == 0 {
		return ""
	}
	var skuIds []string
	for _, skuInfo := range skuInfos {
		skuIds = append(skuIds, skuInfo.Id)
	}
	return strings.Join(skuIds, ",")
}

// 获取两个字符串数组的交集
func getIntersection(arr1, arr2 []string) []string {
	// 创建一个映射用于存储数组元素
	set := make(map[string]struct{})
	for _, item := range arr1 {
		set[item] = struct{}{} // 使用空结构体占位，节省内存
	}
	var intersection []string
	for _, item := range arr2 {
		if _, exists := set[item]; exists {
			intersection = append(intersection, item) // 如果存在于第一个数组中，则加入交集
		}
	}
	return intersection
}

// 将 GBK 编码的字节数组转换为 UTF-8 字符串
func convertGBKToUTF8(gbkData []byte) (string, error) {
	reader := transform.NewReader(strings.NewReader(string(gbkData)), simplifiedchinese.GBK.NewDecoder())
	utf8Data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(utf8Data), nil
}
