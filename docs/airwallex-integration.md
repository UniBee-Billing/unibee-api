# Airwallex 网关集成

## 概述

Airwallex 是一个全球支付网关，支持多种支付方式和货币。本文档描述了如何在 Unibee API 中集成 Airwallex 网关。

## 功能特性

- 支持信用卡支付
- 全球货币支持
- 实时汇率转换
- 安全的支付处理
- 支持退款和取消操作

## 配置参数

### 必需参数

- `gateway_key`: Airwallex API Key
- `gateway_secret`: Airwallex Client Secret

### 可选参数

- `currency`: 默认货币（默认为 USD）
- `logo`: 网关 Logo URL
- `display_name`: 显示名称

## API 端点

- 生产环境: `https://api.airwallex.com`
- 测试环境: `https://api.airwallex.com`

## 支付流程

1. **创建支付意图**: 调用 `GatewayNewPayment` 方法
2. **处理支付**: 用户完成支付
3. **确认支付**: 调用 `GatewayCapture` 方法
4. **退款处理**: 如需要，调用 `GatewayRefund` 方法

## 支持的操作

- ✅ 创建用户账户
- ✅ 创建支付
- ✅ 确认支付
- ✅ 取消支付
- ✅ 退款处理
- ✅ 查询支付状态
- ✅ 查询商户余额
- ✅ 支付方式管理

## 不支持的操作

- ❌ 加密货币交易
- ❌ 银行转账

## 测试

运行测试以确保集成正常工作：

```bash
go test ./test/airwallex_test.go
```

## 注意事项

1. 确保在生产环境中使用有效的 API 密钥
2. 测试环境中的 API 调用不会产生实际费用
3. 所有金额都以分为单位
4. 支持多种货币，但建议使用 USD 作为基础货币

## 错误处理

常见的错误包括：
- API 密钥无效
- 网络连接问题
- 支付金额超出限制
- 货币不支持

## 更多信息

- [Airwallex 官方文档](https://www.airwallex.com/docs)
- [API 参考](https://www.airwallex.com/docs/api)
- [支持联系方式](https://www.airwallex.com/contact)
