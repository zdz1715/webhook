# webhook
通过一个配置文件启动一个webhook转发程序

## 配置参考
[config.yaml.example](config.yaml.example)

## 使用方式
> 没有鉴权方式，推荐内网使用
### 二进制
1. 从[release](https://github.com/zdz1715/webhook/releases)下载二进制包
2. 解压运行
```shell
tar -xzvf webhook-linux-amd64-v0.1.0.tar.gz 
cd bin
webhook --config config.yaml
```

### docker-compose
```yaml
version: "3.1"
services:
  webhook:
    image: zdzserver/webhook:v0.1.0
    container_name: webhook
    ports:
      - "8000:8000"
    environment:
      TZ: "Asia/Shanghai" 
    command:
      - -c
      - /user/local/webhook/config/config.yaml
    volumes:
      - ./config/config.yaml:/user/local/webhook/config/config.yaml
      - ./log:/user/local/webhook/log
    networks:
      webhook:
networks:
  webhook:
```

### kubernetes

```yaml

```

## 接口列表
### /webhook/:uuid
- HTTP请求方式： `GET POST `
#### 请求参数（query）

| 参数     | 必填  | 类型          | 说明            |
|--------|-----|-------------|---------------|
| url    | 否   | string      | 覆盖配置文件的`url`  |


## 自定义消息格式
模版参考：https://pkg.go.dev/text/template

支持函数：https://masterminds.github.io/sprig/

### 额外函数
- `urlEncode`: encode url
- `urlDecode`: decode url





