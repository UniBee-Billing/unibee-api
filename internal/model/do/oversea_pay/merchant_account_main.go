// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantAccountMain is the golang structure of table merchant_account_main for DAO operations like Where/Data.
type MerchantAccountMain struct {
	g.Meta      `orm:"table:merchant_account_main, do:true"`
	Id          interface{} // 主键ID
	CompanyId   interface{} //
	MerchantId  interface{} // 商户ID
	Currency    interface{} //，“SGD” “MYR” “PHP” “IDR” “THB” 与付款金额关联的货币。指定三个字母的ISO 4217货币代码
	TotalTrade  interface{} // 交易金额总计
	TotalRefund interface{} // 退款金额总计
	TotalCut    interface{} // 服务扣点金额总计
	TotalSend   interface{} // 结算金额总计
	Year        interface{} // 结算key-年
	Month       interface{} // 结算key-月
	Day         interface{} // 结算key-日
	GmtCreate   *gtime.Time //
	GmtModify   *gtime.Time //
	Statistic   interface{} // 统计使用
	IsDeleted   interface{} // 是否删除，0-未删除，1-已删除
}
