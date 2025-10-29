P# Pay Client SDK

支付客户端SDK，提供与支付服务端的交互能力，支持gRPC和HTTP两种通信协议。

## 目录
- [功能特性](#功能特性)
- [安装](#安装)
- [快速开始](#快速开始)
  - [HTTP客户端](#http客户端)
  - [gRPC客户端](#grpc客户端)
- [使用说明](#使用说明)
  - [配置参数](#配置参数)
  - [创建收款订单](#创建收款订单)
  - [创建付款订单](#创建付款订单)
  - [查询订单](#查询订单)
  - [其他功能](#其他功能)
- [API参考](#api参考)
- [协议](#协议)

## 功能特性

- 支持gRPC和HTTP两种通信协议
- 提供统一的接口访问支付服务
- 自动处理数据加密解密
- 完整的日志记录功能
- 支持以下业务功能：
  - 虚拟账户创建
  - 收款订单管理
  - 付款订单管理
  - 订单查询
  - 支付通道查询
  - 商户余额查询

## 安装

```bash
go get github.com/XingMenTech/XMPAY-SDK-GO
```

## 快速开始

### HTTP客户端

```go
import "github.com/XingMenTech/XMPAY-SDK-GO"

config := &client.Config{
    ApiUrl:       "http://localhost:8080",
    AccessId:     "your_access_id",
    AccessKey:    "your_access_key",
    InId:         "receive_channel_id",
    OutId:        "payment_channel_id",
    InNotifyUrl:  "http://yourdomain.com/notify/receive",
    OutNotifyUrl: "http://yourdomain.com/notify/payment",
}

httpClient := client.NewHttpClient(config, nil)
```

### gRPC客户端

```go
config := &client.Config{
    ApiUrl:       "localhost:9090",
    AccessId:     "your_access_id",
    AccessKey:    "your_access_key",
    InId:         "receive_channel_id",
    OutId:        "payment_channel_id",
    InNotifyUrl:  "http://yourdomain.com/notify/receive",
    OutNotifyUrl: "http://yourdomain.com/notify/payment",
}

grpcClient, err := client.NewGrpcClient(config, nil)
if err != nil {
    // 处理错误
}
defer grpcClient.Close()
```

## 使用说明

### 配置参数

| 参数 | 类型 | 描述 |
|------|------|------|
| ApiUrl | string | API地址 (HTTP客户端为完整基础URL，gRPC客户端为主机名和端口) |
| AccessId | string | 访问ID，用于身份验证 |
| AccessKey | string | 访问密钥，用于数据加密 |
| InId | string | 收款通道ID |
| OutId | string | 代付通道ID |
| InNotifyUrl | string | 收款回调地址 |
| OutNotifyUrl | string | 代付回调地址 |

### 创建收款订单

```go
param := &client.ReceiveParam{
    OrderParam: client.OrderParam{
        OrderNo:   "ORDER123456",
        Amount:    10000, // 单位：分
        Uid:       "user123",
        Ip:        "192.168.1.1",
        Name:      "张三",
        Phone:     "13800138000",
        Email:     "zhangsan@example.com",
        IdNum:     "110101199003076598",
        Pid:       1,
        NotifyUrl: "http://yourdomain.com/notify/custom",
        Subject:   "商品标题",
        Body:      "商品描述",
    },
    ReturnUrl: "http://yourdomain.com/return",
}

resp, err := client.CreateReceive(param)
```

### 创建付款订单

```go
param := &client.OutParam{
    OrderParam: client.OrderParam{
        OrderNo:   "OUT123456",
        Amount:    10000, // 单位：分
        Uid:       "user123",
        Ip:        "192.168.1.1",
        Name:      "李四",
        Phone:     "13800138001",
        Email:     "lisi@example.com",
        IdNum:     "110101199003079856",
        Pid:       2,
        NotifyUrl: "http://yourdomain.com/notify/out",
        Subject:   "提现",
        Body:      "用户提现",
    },
    BankNo:   "6222001234567890123",
    BankCode: "ICBC",
    BankName: "工商银行",
    Mode:     "3", // 银行卡
}

resp, err := client.CreateOut(param)
```

### 查询订单

```go
// 查询收款订单
resp, err := client.QueryReceive("merchant_order_no", "platform_order_no")

// 查询付款订单
resp, err := client.QueryOut("merchant_order_no", "platform_order_no")
```

### 其他功能

```go
// 查询可用支付通道
channels, err := client.Channel(pb.ORDER_TYPE_RECEIVE)

// 查询商户余额
balance, err := client.Balance()

// 创建虚拟账户
virtualResp, err := client.CreateVirtual(virtualParam)
```

## API参考

该SDK提供了以下主要接口：

- `CreateVirtual`: 创建虚拟账户
- `CreateReceive`: 创建收款订单
- `QueryReceive`: 查询收款订单
- `CreateOut`: 创建付款订单
- `QueryOut`: 查询付款订单
- `Channel`: 查询支付通道
- `Balance`: 查询商户余额

## 协议

本项目采用MIT协议。