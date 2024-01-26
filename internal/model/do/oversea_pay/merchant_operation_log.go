// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantOperationLog is the golang structure of table merchant_operation_log for DAO operations like Where/Data.
type MerchantOperationLog struct {
	g.Meta             `orm:"table:merchant_operation_log, do:true"`
	Id                 interface{} // 主键id
	CompanyId          interface{} // 公司ID
	MerchantId         interface{} // merchantId
	UserId             interface{} // 操作userId，系统自动操作可能没有
	OptAccount         interface{} // 操作账号
	ClientType         interface{} // 操作渠道 0:云店后台 1:云管家app 2:Java服务 3:小程序
	BizType            interface{} // 操作业务 0:菜单 1:商品 2:门店 3:订单 4:账号|会员 5:优惠券转赠中 6:优惠券转赠领取成功 7:优惠券转赠自动取消
	OptTarget          interface{} // 操作对象
	OptContent         interface{} // 操作内容
	OptCreate          *gtime.Time // 操作发生时间
	IsDelete           interface{} // 是否删除 1-删除，0-未删除
	GmtCreate          *gtime.Time // 创建时间
	GmtModify          *gtime.Time // 修改时间
	QueryportRequestId interface{} // queryport请求Id，可在request_security_log查询请求信息
	ServerType         interface{} // 操作终端，参看 message-api包 OperationLogServerTypeEnum的code
	ServerTypeDesc     interface{} // 操作终端描述，参看 message-api包 OperationLogServerTypeEnum的desc
}
