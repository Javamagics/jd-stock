### 说明

实现思路来自于[京东有货地区|库存查询](https://cll.name/jd/stockjd.htm)

通过指定的地区编码和商品id定时查询库存

### 特性
+ 支持多商品
+ 支持多地区
+ 支持钉钉群聊机器人通知
+ 支持企业微信机器人
+ 支持PushPlus
+ 支持Server酱
+ 支持Bark
+ 支持网页编辑配置

### 待实现

- [x] java镜像太大，使用go重构
- [x] 实现更多通知方式

### 使用方式

#### 下载配置文件

```bash
mkdir config
wget -O config/config.yaml https://github.com/zhuweitung/jd-stock/raw/refs/heads/go/config/config.yaml.example

# 国内
wget -O config/config.yaml https://fastly.jsdelivr.net/gh/zhuweitung/jd-stock@go/config/config.yaml.example
```

#### 修改配置文件

```yml
everyMinutes: 5 # 每隔N分钟执行
provinces: # 库存省份，省份有货后通知
  - 江苏
  - 浙江
  - 上海
addressCodes: # 库存地区编码，在 https://cll.name/jd/jdloc.html 中查询完整的地区编码
  - 2_2830_0_0
  - 15_1213_3411_0
skuInfos: # 监控商品信息
  - id: 100014150579
    name: 蓝漂XPLUS会员联名款 抽纸4层100抽*20包
delay: 5500 # 每次查询延迟（毫秒），建议设置大些，防止触发风控
ua: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36 Edg/129.0.0.0
enableNotify: false # 启用通知
notifyInterval: 720 # 通知间隔（分钟），防止频繁发送相同通知，0表示允许重复提醒
notifyType: "dingtalk_bot" # 通知方式（dingtalk_bot、qy_wechat_bot、push_plus、server_chan、bark）
dingtalkBot: # 钉钉机器人
  token:
  secret:
qyWechatBot: # 企业微信机器人
  key:
pushPlus: # PushPlus
  token:
serverChan: # Server酱
  sendKey:
bark: # bark
  apiUrl:
  deviceKey:
```

库存地区编码`addressCodes`参数从[京东地区编码](https://cll.name/jd/jdloc.html)查询获得

库存省份`provinces`与库存地区编码`addressCodes`参数两者选其一



#### docker（二选一）

```bash
docker run -d --name jd-stock -p 8080:8080 -v ./config:/app/config zhuweitung/jd-stock:latest
```

#### docker-compose（二选一）

下载`docker-compose.yml`

```bash
wget -O docker-compose.yml https://github.com/zhuweitung/jd-stock/raw/refs/heads/go/docker-compose.yml

# 国内
wget -O docker-compose.yml https://fastly.jsdelivr.net/gh/zhuweitung/jd-stock@go/docker-compose.yml
```

启动

```bash
docker-compose up -d
```

#### 网页编辑

访问`ip:8080`端口

### 侵删

