// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantAccountOrder is the golang structure of table merchant_account_order for DAO operations like Where/Data.
type MerchantAccountOrder struct {
	g.Meta                        `orm:"table:merchant_account_order, do:true"`
	Id                            interface{} //
	CompanyId                     interface{} //
	MerchantId                    interface{} // 商户ID
	MainId                        interface{} // 结算单id，表merchant_account_main的id
	BizId                         interface{} // 业务ID。可能是payId、refundId
	BizType                       interface{} // 业务ID类型，1-pay,2-refund
	OrderType                     interface{} // 账单类型；1支付；2 退款
	OrigCurrency                  interface{} // 原始货币类型
	OrigTradeFee                  interface{} // 原始货币交易金额。单位：分
	CurrencyRate                  interface{} // 汇率，万分位
	CurrencyRateDataJson          interface{} // 汇率数据JSON结构
	Currency                      interface{} 类型
	TradeFee                      interface{} // 交易金额（正值，退款代表退还金额）。单位：分
	DeductPoint                   interface{} // 服务费扣点，万分位
	DeductFee                     interface{} // 扣点金额（正值，退款代表退服务费）。单位：分
	BillFee                       interface{} // 结算金额（正值，退款代表退还金额）。单位：分
	GmtCreate                     *gtime.Time //
	GmtModify                     *gtime.Time //
	MerchantReference             interface{} // 客户订单号
	ChannelOrderNo                interface{} //
	MerchantOrderNo               interface{} //
	Channel                       interface{} //
	PaymentCurrency               interface{} //
	Authorised                    interface{} //
	Captured                      interface{} //
	CurrencyRateMarkup            interface{} //
	PaymentMethodVariant          interface{} //
	ModificationMerchantReference interface{} //
	MerchantOrderReference        interface{} //
	DelayCaptureTime              interface{} //
	NotificationUrl               interface{} //
	OpenAppId                     interface{} //
}
