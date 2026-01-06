# SberPay 支付网关集成文档

## 概述

SberPay是俄罗斯Sberbank的移动支付解决方案，现已集成到UniBee支付系统中。本文档详细说明了SberPay的集成方式、配置要求和API使用。

## 功能特性

- ✅ 支付创建和查询
- ✅ 支付状态同步
- ✅ 退款处理（全额和部分退款）
- ✅ Webhook通知处理
- ✅ 签名验证（HMAC-SHA256）
- ❌ 支付方式管理（SberPay不支持）
- ❌ 余额查询（SberPay不支持）

## 配置要求

### 1. 网关配置

在商户后台配置SberPay网关时，需要提供以下信息：

- **网关名称**: `sberpay`
- **API密钥**: SberPay项目名称
- **API密钥**: SberPay API密钥
- **Webhook密钥**: 用于验证webhook签名的密钥

### 2. 环境变量

确保以下环境变量已正确配置：

```bash
# SberPay API配置
SBERPAY_API_KEY=your_api_key
SBERPAY_API_SECRET=your_api_secret
SBERPAY_WEBHOOK_SECRET=your_webhook_secret
```

## API集成

### 支付创建

```go
// 创建SberPay支付
paymentData := map[string]interface{}{
    "amount":      1000, // 金额（分）
    "currency":    "RUB", // 货币
    "order_id":    "order_123", // 订单ID
    "description": "Payment for order order_123", // 描述
    "return_url":  "https://your-site.com/success", // 成功回调URL
    "fail_url":    "https://your-site.com/fail", // 失败回调URL
}
```

### 支付查询

```go
// 查询支付状态
response, err := sberPayAPI.GetPaymentDetail(paymentId)
```

### 退款处理

```go
// 创建退款
refundData := map[string]interface{}{
    "payment_id": "payment_123", // 原支付ID
    "amount":     500, // 退款金额（分）
    "currency":   "RUB", // 货币
    "reason":     "Customer request", // 退款原因
}
```

## Webhook处理

### Webhook事件类型

SberPay支持以下webhook事件：

1. **payment.finished** - 支付完成
2. **refund.finished** - 退款完成

### 签名验证

所有webhook请求都使用HMAC-SHA256进行签名验证：

```go
func verifyWebhookSignature(r *http.Request, gateway *entity.MerchantGateway, body []byte) bool {
    signature := r.Header.Get("X-SberPay-Signature")
    h := hmac.New(sha256.New, []byte(gateway.WebhookSecret))
    h.Write(body)
    expectedSignature := hex.EncodeToString(h.Sum(nil))
    return signature == expectedSignature
}
```

### Webhook处理流程

1. 接收webhook请求
2. 验证签名
3. 解析事件类型
4. 更新支付/退款状态
5. 发送状态变更消息到RedisMQ

## 状态映射

### 支付状态映射

| SberPay状态 | UniBee状态 | 描述 |
|------------|-----------|------|
| pending | PaymentCreated (10) | 待处理 |
| succeeded | PaymentSuccess (20) | 成功 |
| failed | PaymentFailed (30) | 失败 |
| cancelled | PaymentCancelled (40) | 已取消 |

### 退款状态映射

| SberPay状态 | UniBee状态 | 描述 |
|------------|-----------|------|
| pending | RefundCreated (10) | 待处理 |
| succeeded | RefundSuccess (20) | 成功 |
| failed | RefundFailed (30) | 失败 |
| cancelled | RefundCancelled (40) | 已取消 |

## 错误处理

### 常见错误

1. **API认证失败**: 检查API密钥和密钥是否正确
2. **签名验证失败**: 检查webhook密钥是否正确
3. **支付不存在**: 检查支付ID是否正确
4. **金额错误**: 确保金额格式正确（以分为单位）

### 错误响应格式

```json
{
    "success": false,
    "error": "错误描述"
}
```

## 测试

### 测试环境

SberPay提供测试环境用于开发和测试：

- **测试API地址**: `https://test.proxy.bank131.ru/api/v1`
- **测试账户**: 联系SberPay获取测试账户

### 测试用例

1. **支付创建测试**
2. **支付状态查询测试**
3. **退款创建测试**
4. **Webhook接收测试**

## 部署注意事项

1. **HTTPS要求**: 生产环境必须使用HTTPS
2. **Webhook URL**: 确保webhook URL可公开访问
3. **超时设置**: API调用超时设置为30秒
4. **日志记录**: 所有API调用都会记录到日志中

## 监控和日志

### 日志记录

所有SberPay相关的API调用都会记录到以下位置：

- **API调用日志**: `internal/logic/gateway/api/log/`
- **Webhook日志**: 系统日志

### 监控指标

建议监控以下指标：

- API调用成功率
- Webhook接收成功率
- 支付成功率
- 退款成功率

## 支持

如有问题，请联系：

- **技术文档**: https://developer.131.ru/en/payments/payment-sber-pay/
- **API支持**: 联系SberPay技术支持
- **UniBee支持**: 联系UniBee技术支持

## 更新日志

### v1.0.0 (2024-01-XX)
- 初始SberPay集成
- 支持支付和退款功能
- 支持webhook处理
- 支持签名验证 