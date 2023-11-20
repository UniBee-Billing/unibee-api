// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantAccountOrder is the golang structure for table merchant_account_order.
type MerchantAccountOrder struct {
	Id                            uint64      `json:"id"                            ` //
	CompanyId                     int64       `json:"companyId"                     ` //
	MerchantId                    int64       `json:"merchantId"                    ` // 商户ID
	MainId                        int64       `json:"mainId"                        ` // 结算单id，表merchant_account_main的id
	BizId                         int64       `json:"bizId"                         ` // 业务ID。可能是payId、refundId
	BizType                       int         `json:"bizType"                       ` // 业务ID类型，1-pay,2-refund
	OrderType                     int         `json:"orderType"                     ` // 账单类型；1支付；2 退款
	OrigCurrency                  string      `json:"origCurrency"                  ` // 原始货币类型
	OrigTradeFee                  int64       `json:"origTradeFee"                  ` // 原始货币交易金额。单位：分
	CurrencyRate                  int64       `json:"currencyRate"                  ` // 汇率，万分位
	CurrencyRateDataJson          string      `json:"currencyRateDataJson"          ` // 汇率数据JSON结构
	Currency                      string      `json:"currency"                      ` // 货币类型
	TradeFee                      int64       `json:"tradeFee"                      ` // 交易金额（正值，退款代表退还金额）。单位：分
	DeductPoint                   string      `json:"deductPoint"                   ` // 服务费扣点，万分位
	DeductFee                     int64       `json:"deductFee"                     ` // 扣点金额（正值，退款代表退服务费）。单位：分
	BillFee                       int64       `json:"billFee"                       ` // 结算金额（正值，退款代表退还金额）。单位：分
	GmtCreate                     *gtime.Time `json:"gmtCreate"                     ` //
	GmtModify                     *gtime.Time `json:"gmtModify"                     ` //
	MerchantReference             string      `json:"merchantReference"             ` // 客户订单号
	ChannelOrderNo                string      `json:"channelOrderNo"                ` //
	MerchantOrderNo               string      `json:"merchantOrderNo"               ` //
	Channel                       string      `json:"channel"                       ` //
	PaymentCurrency               string      `json:"paymentCurrency"               ` //
	Authorised                    int64       `json:"authorised"                    ` //
	Captured                      int64       `json:"captured"                      ` //
	CurrencyRateMarkup            string      `json:"currencyRateMarkup"            ` //
	PaymentMethodVariant          string      `json:"paymentMethodVariant"          ` //
	ModificationMerchantReference string      `json:"modificationMerchantReference" ` //
	MerchantOrderReference        string      `json:"merchantOrderReference"        ` //
	DelayCaptureTime              string      `json:"delayCaptureTime"              ` //
	NotificationUrl               string      `json:"notificationUrl"               ` //
	OpenAppId                     string      `json:"openAppId"                     ` //
}
