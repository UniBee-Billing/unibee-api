# Blockonomics 后端实现需求

## 概述
本文档列出了在 Go 后端中实现 Blockonomics 支付网关所需的功能和数据结构。

## Action 数据结构需求

在 `internal/logic/gateway/api/blockonomics.go` 的 `GatewayNewPayment` 函数中，需要在 `action` 字段中设置以下数据：

```go
action := gjson.New("")
_ = action.Set("blockonomicsAddress", paymentAddress)      // 支付地址
_ = action.Set("blockonomicsApiKey", apiKey)              // API 密钥
_ = action.Set("blockonomicsOrderId", orderId)            // 订单 ID
_ = action.Set("blockonomicsAmount", amount)              // 支付金额
_ = action.Set("blockonomicsCurrency", currency)          // 货币类型 (BTC/USDT)
_ = action.Set("blockonomicsTimeout", timeout)            // 支付超时时间（秒）
```

## 必需的后端功能

### 1. 支付地址生成
- **功能**: 为每个支付生成唯一的加密货币地址
- **API**: 调用 Blockonomics API 创建支付地址
- **实现位置**: `internal/logic/gateway/api/blockonomics.go`
- **参考**: [Blockonomics API - Create Payment Address](https://developers.blockonomics.co/reference/create-or-get-payment-addresspost)

### 2. 订单管理
- **功能**: 创建和管理 Blockonomics 订单
- **API**: 调用 Blockonomics API 创建订单
- **实现位置**: `internal/logic/gateway/api/blockonomics.go`
- **参考**: [Blockonomics API - Order Management](https://developers.blockonomics.co/reference/get-all-ordersget)

### 3. 支付监控
- **功能**: 监控支付状态变化
- **实现方式**: 
  - WebSocket 连接监控
  - 回调通知处理
  - 定期轮询检查
- **实现位置**: `internal/logic/gateway/webhook/blockonomics.go`

### 4. 价格转换
- **功能**: 将法币金额转换为加密货币金额
- **API**: 调用 Blockonomics 价格 API
- **实现位置**: `internal/logic/gateway/api/blockonomics.go`
- **参考**: [Blockonomics API - Fetch Cryptocurrency Price](https://developers.blockonomics.co/reference/fetch-cryptocurrency-priceget)

## 配置需求

### 1. 环境变量
```bash
BLOCKONOMICS_API_KEY=your_api_key_here
BLOCKONOMICS_API_URL=https://www.blockonomics.co/api
BLOCKONOMICS_WEBSOCKET_URL=wss://www.blockonomics.co/api/ws
```

### 2. 网关配置
在 `api/bean/merchant_gateway.go` 中需要支持：
- `GatewayKey`: Blockonomics API Key
- `GatewaySecret`: 可选，用于签名验证
- `Host`: API 基础 URL
- `SubGateway`: 支持的加密货币类型 (BTC, USDT)

## 数据库需求

### 1. 支付记录扩展
在 `api/bean/payment.go` 中可能需要添加：
- `CryptoAddress`: 加密货币地址
- `CryptoAmount`: 加密货币金额
- `CryptoCurrency`: 加密货币类型
- `PaymentTimeout`: 支付超时时间

### 2. 订单状态跟踪
- `Pending`: 等待支付
- `Paid`: 已支付
- `Expired`: 已过期
- `Failed`: 支付失败

## Webhook 处理

### 1. 支付通知
- **端点**: `/payment/gateway_webhook_entry/blockonomics/notifications`
- **功能**: 处理 Blockonomics 的支付通知
- **实现位置**: `internal/logic/gateway/webhook/blockonomics.go`

### 2. 回调验证
- 验证 Webhook 签名的真实性
- 处理重复通知
- 更新支付状态

## 错误处理

### 1. 常见错误场景
- API 调用失败
- 网络连接问题
- 支付超时
- 金额不匹配
- 地址生成失败

### 2. 错误码定义
```go
const (
    BlockonomicsErrorAPIKeyInvalid = "BLOCKONOMICS_API_KEY_INVALID"
    BlockonomicsErrorAddressGeneration = "BLOCKONOMICS_ADDRESS_GENERATION_FAILED"
    BlockonomicsErrorPaymentTimeout = "BLOCKONOMICS_PAYMENT_TIMEOUT"
    BlockonomicsErrorAmountMismatch = "BLOCKONOMICS_AMOUNT_MISMATCH"
)
```

## 测试需求

### 1. 测试环境
- 使用 Blockonomics 测试网络
- 测试 API 密钥
- 模拟支付场景

### 2. 测试用例
- 支付地址生成
- 支付状态监控
- 支付完成处理
- 支付超时处理
- 错误场景处理

## 安全考虑

### 1. API 密钥管理
- 安全存储 API 密钥
- 密钥轮换机制
- 访问日志记录

### 2. 支付验证
- 验证支付金额
- 验证支付地址
- 防止重复支付

## 性能优化

### 1. 缓存策略
- 价格信息缓存
- 支付状态缓存
- API 响应缓存

### 2. 并发处理
- WebSocket 连接池
- 异步支付处理
- 批量状态更新

## 监控和日志

### 1. 关键指标
- 支付成功率
- 平均支付时间
- API 调用延迟
- 错误率统计

### 2. 日志记录
- API 调用日志
- 支付状态变化日志
- 错误日志
- 性能日志

## 部署注意事项

### 1. 环境配置
- 生产环境 API 密钥
- WebSocket 连接配置
- 超时设置

### 2. 依赖服务
- Blockonomics API 可用性
- 网络连接稳定性
- 数据库性能

---

**注意**: 以上需求基于 Blockonomics API 文档和前端实现需求。实际实现时请参考最新的 Blockonomics API 文档进行相应调整。
