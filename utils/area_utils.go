package utils

import (
	"encoding/json"
	"fmt"
	"github.com/zhuweitung/jd-stock/models"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	provinces = []models.AreaInfo{
		{ID: json.Number("2"), PID: json.Number("0"), Level: 0, Name: "上海"},
		{ID: json.Number("12"), PID: json.Number("0"), Level: 0, Name: "江苏"},
		{ID: json.Number("15"), PID: json.Number("0"), Level: 0, Name: "浙江"},
		{ID: json.Number("14"), PID: json.Number("0"), Level: 0, Name: "安徽"},
		{ID: json.Number("19"), PID: json.Number("0"), Level: 0, Name: "广东"},
		{ID: json.Number("20"), PID: json.Number("0"), Level: 0, Name: "广西"},
		{ID: json.Number("16"), PID: json.Number("0"), Level: 0, Name: "福建"},
		{ID: json.Number("23"), PID: json.Number("0"), Level: 0, Name: "海南"},
		{ID: json.Number("1"), PID: json.Number("0"), Level: 0, Name: "北京"},
		{ID: json.Number("5"), PID: json.Number("0"), Level: 0, Name: "河北"},
		{ID: json.Number("11"), PID: json.Number("0"), Level: 0, Name: "内蒙古"},
		{ID: json.Number("13"), PID: json.Number("0"), Level: 0, Name: "山东"},
		{ID: json.Number("6"), PID: json.Number("0"), Level: 0, Name: "山西"},
		{ID: json.Number("3"), PID: json.Number("0"), Level: 0, Name: "天津"},
		{ID: json.Number("17"), PID: json.Number("0"), Level: 0, Name: "湖北"},
		{ID: json.Number("18"), PID: json.Number("0"), Level: 0, Name: "湖南"},
		{ID: json.Number("7"), PID: json.Number("0"), Level: 0, Name: "河南"},
		{ID: json.Number("21"), PID: json.Number("0"), Level: 0, Name: "江西"},
		{ID: json.Number("8"), PID: json.Number("0"), Level: 0, Name: "辽宁"},
		{ID: json.Number("10"), PID: json.Number("0"), Level: 0, Name: "黑龙江"},
		{ID: json.Number("9"), PID: json.Number("0"), Level: 0, Name: "吉林"},
		{ID: json.Number("22"), PID: json.Number("0"), Level: 0, Name: "四川"},
		{ID: json.Number("4"), PID: json.Number("0"), Level: 0, Name: "重庆"},
		{ID: json.Number("25"), PID: json.Number("0"), Level: 0, Name: "云南"},
		{ID: json.Number("24"), PID: json.Number("0"), Level: 0, Name: "贵州"},
		{ID: json.Number("26"), PID: json.Number("0"), Level: 0, Name: "西藏"},
		{ID: json.Number("27"), PID: json.Number("0"), Level: 0, Name: "陕西"},
		{ID: json.Number("30"), PID: json.Number("0"), Level: 0, Name: "宁夏"},
		{ID: json.Number("28"), PID: json.Number("0"), Level: 0, Name: "甘肃"},
		{ID: json.Number("29"), PID: json.Number("0"), Level: 0, Name: "青海"},
		{ID: json.Number("31"), PID: json.Number("0"), Level: 0, Name: "新疆"},
	}
	areas    []models.AreaInfo
	areaOnce sync.Once
	// 正则表达式，用于提取 JSON 数据部分
	areaJsonPattern = regexp.MustCompile(`\((\[.*])\)`)
)

// LoadAreaCodes 加载地区编码
func LoadAreaCodes() error {
	var err error
	areaOnce.Do(func() {
		_areas, _err := fetch()
		if _err != nil {
			err = _err
		}
		areas = _areas
	})
	return err
}

// 从文件或接口获取地区数据，并发处理查询。
func fetch() ([]models.AreaInfo, error) {
	filepath := "data/area_code.json"
	if _, err := os.Stat(filepath); err == nil {
		log.Println("读取现有的地区编码数据...")
		data, err := os.ReadFile(filepath)
		if err != nil {
			return nil, err
		}
		var _areas []models.AreaInfo
		if err := json.Unmarshal(data, &_areas); err != nil {
			return nil, err
		}
		return _areas, nil
	}

	log.Println("未找到地区编码数据，开始从接口获取...")
	var (
		mu      sync.Mutex
		_areas  []models.AreaInfo
		wg      sync.WaitGroup
		errChan = make(chan error, len(provinces))
	)

	wg.Add(len(provinces))
	for _, province := range provinces {
		go func(province models.AreaInfo) {
			defer wg.Done()
			_areas = append(_areas, province)
			children, err := fetchChildren(province, 0)
			if err != nil {
				errChan <- err
				return
			}
			mu.Lock()
			_areas = append(_areas, children...)
			mu.Unlock()
		}(province)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan // 返回第一个错误
	}

	data, _ := json.Marshal(_areas)
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return nil, err
	}
	return _areas, nil
}

// 递归获取下级地区编码
func fetchChildren(currentArea models.AreaInfo, level int) ([]models.AreaInfo, error) {
	if level >= 3 { // 限制递归深度
		return nil, nil
	}

	req, err := http.NewRequest("GET", "https://fts.jd.com/area/get", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("fid", currentArea.ID.String())
	q.Add("callback", "jQuery1111047169012038874314_1729066415663")
	q.Add("_", fmt.Sprint(time.Now().UnixMilli()))
	req.URL.RawQuery = q.Encode()
	req.Header.Set("User-Agent", GetConfig().Ua)

	client := &http.Client{Timeout: 10 * time.Second}
	log.Printf("查询地区编码：%s %s\n", currentArea.ID.String(), currentArea.Name)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	matches := areaJsonPattern.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		return nil, fmt.Errorf("无法解析响应中的JSON数据")
	}
	jsonData := matches[1]

	var _areas []models.AreaInfo
	if err := json.Unmarshal([]byte(jsonData), &_areas); err != nil {
		return nil, err
	}

	var children []models.AreaInfo
	for _, area := range _areas {
		area.PID = currentArea.ID
		area.Level = level + 1
		children = append(children, area)
		_children, err := fetchChildren(area, level+1)
		if err != nil {
			log.Println(err)
			continue
		}
		if _children != nil {
			children = append(children, _children...)
		}
	}

	return children, nil
}

// GetRandomCodeCombination 获取随机地区编码组合（各省）
func GetRandomCodeCombination(provinceNames []string) []string {
	var codes []string
	for _, province := range provinces {
		if provinceNames != nil && !contains(province.Name, provinceNames) {
			continue
		}
		_codes := getRandomChildrenCode(province)
		_codesLength := len(_codes)
		if _codesLength < 4 {
			for i := 0; i < 4-_codesLength; i++ {
				_codes = append(_codes, "0")
			}
		}
		codes = append(codes, strings.Join(_codes, "_"))
	}
	return codes
}

// 获取随机下级地区编码
func getRandomChildrenCode(area models.AreaInfo) []string {
	var codes []string
	codes = append(codes, area.ID.String())
	children := GetChildAreasByID(area.ID.String())
	if children == nil {
		return codes
	}
	// 使用当前时间戳作为随机数种子
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(children)) // 获取随机索引
	_codes := getRandomChildrenCode(children[randomIndex])
	codes = append(codes, _codes...)
	return codes
}

// 判断集合是否包含字符串
func contains(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// GetAreaByName 根据地区名称查找地区信息
func GetAreaByName(name string) (*models.AreaInfo, error) {
	for _, area := range areas {
		if area.Name == name {
			return &area, nil // 返回找到的地区信息
		}
	}
	return nil, fmt.Errorf("未找到名称为 %s 的地区", name)
}

// GetAreaByID 根据地区id查找地区信息
func GetAreaByID(id string) (*models.AreaInfo, error) {
	for _, area := range areas {
		if area.ID.String() == id {
			return &area, nil // 返回找到的地区信息
		}
	}
	return nil, fmt.Errorf("未找到id为 %s 的地区", id)
}

// GetChildAreasByID 根据地区ID获取下级地区
func GetChildAreasByID(id string) []models.AreaInfo {
	var children []models.AreaInfo
	for _, area := range areas {
		if area.PID.String() == id {
			children = append(children, area)
		}
	}
	if len(children) == 0 {
		return nil
	}
	return children
}

// GetAreaNameByCodeCombination 根据组合地区编码获取地区名称
func GetAreaNameByCodeCombination(areaCodeCombination string) string {
	codes := strings.Split(areaCodeCombination, "_")
	var names []string
	for _, code := range codes {
		if code != "0" {
			area, err := GetAreaByID(code)
			if err != nil {
				log.Printf("%v", err)
			} else {
				areaName := area.Name
				names = append(names, areaName)
			}
		}
	}
	return strings.Join(names, "")
}
