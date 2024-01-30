// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantOperationLog is the golang structure for table merchant_operation_log.
type MerchantOperationLog struct {
	Id                 uint64      `json:"id"                 description:"主键id"`                                                              // 主键id
	CompanyId          int64       `json:"companyId"          description:"公司ID"`                                                              // 公司ID
	MerchantId         int64       `json:"merchantId"         description:"merchantId"`                                                        // merchantId
	UserId             int64       `json:"userId"             description:"操作userId，系统自动操作可能没有"`                                               // 操作userId，系统自动操作可能没有
	OptAccount         string      `json:"optAccount"         description:"操作账号"`                                                              // 操作账号
	ClientType         int         `json:"clientType"         description:"操作渠道 0:云店后台 1:云管家app 2:Java服务 3:小程序"`                               // 操作渠道 0:云店后台 1:云管家app 2:Java服务 3:小程序
	BizType            int         `json:"bizType"            description:"操作业务 0:菜单 1:商品 2:门店 3:订单 4:账号|会员 5:优惠券转赠中 6:优惠券转赠领取成功 7:优惠券转赠自动取消"` // 操作业务 0:菜单 1:商品 2:门店 3:订单 4:账号|会员 5:优惠券转赠中 6:优惠券转赠领取成功 7:优惠券转赠自动取消
	OptTarget          string      `json:"optTarget"          description:"操作对象"`                                                              // 操作对象
	OptContent         string      `json:"optContent"         description:"操作内容"`                                                              // 操作内容
	OptCreate          *gtime.Time `json:"optCreate"          description:"操作发生时间"`                                                            // 操作发生时间
	IsDelete           int         `json:"isDelete"           description:"0-UnDeleted，1-Deleted"`                                             // 0-UnDeleted，1-Deleted
	GmtCreate          *gtime.Time `json:"gmtCreate"          description:"创建时间"`                                                              // 创建时间
	GmtModify          *gtime.Time `json:"gmtModify"          description:"修改时间"`                                                              // 修改时间
	QueryportRequestId string      `json:"queryportRequestId" description:"queryport请求Id，可在request_security_log查询请求信息"`                        // queryport请求Id，可在request_security_log查询请求信息
	ServerType         int         `json:"serverType"         description:"操作终端，参看 message-api包 OperationLogServerTypeEnum的code"`              // 操作终端，参看 message-api包 OperationLogServerTypeEnum的code
	ServerTypeDesc     string      `json:"serverTypeDesc"     description:"操作终端描述，参看 message-api包 OperationLogServerTypeEnum的desc"`            // 操作终端描述，参看 message-api包 OperationLogServerTypeEnum的desc
}
