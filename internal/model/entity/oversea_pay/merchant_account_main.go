// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantAccountMain is the golang structure for table merchant_account_main.
type MerchantAccountMain struct {
	Id          int64       `json:"id"          ` // 主键ID
	CompanyId   int64       `json:"companyId"   ` //
	MerchantId  int64       `json:"merchantId"  ` // 商户ID
	Currency    string      `json:"currency"    ` // 货币，“SGD” “MYR” “PHP” “IDR” “THB” 与付款金额关联的货币。指定三个字母的ISO 4217货币代码
	TotalTrade  int64       `json:"totalTrade"  ` // 交易金额总计
	TotalRefund int64       `json:"totalRefund" ` // 退款金额总计
	TotalCut    int64       `json:"totalCut"    ` // 服务扣点金额总计
	TotalSend   int64       `json:"totalSend"   ` // 结算金额总计
	Year        int         `json:"year"        ` // 结算key-年
	Month       int         `json:"month"       ` // 结算key-月
	Day         int         `json:"day"         ` // 结算key-日
	GmtCreate   *gtime.Time `json:"gmtCreate"   ` //
	GmtModify   *gtime.Time `json:"gmtModify"   ` //
	Statistic   int         `json:"statistic"   ` // 统计使用
	IsDeleted   int64       `json:"isDeleted"   ` // 是否删除，0-未删除，1-已删除
}
