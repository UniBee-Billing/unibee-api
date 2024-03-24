// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Refund is the golang structure for table refund.
type Refund struct {
	Id                   int64       `json:"id"                   description:"id"`                                                 // id
	CompanyId            int64       `json:"companyId"            description:"company id"`                                         // company id
	MerchantId           uint64      `json:"merchantId"           description:"merchant id"`                                        // merchant id
	UserId               uint64      `json:"userId"               description:"user_id"`                                            // user_id
	OpenApiId            int64       `json:"openApiId"            description:"open api id"`                                        // open api id
	GatewayId            uint64      `json:"gatewayId"            description:"gateway_id"`                                         // gateway_id
	BizType              int         `json:"bizType"              description:"biz type, copy from payment.biz_type"`               // biz type, copy from payment.biz_type
	ExternalRefundId     string      `json:"externalRefundId"     description:"external_refund_id"`                                 // external_refund_id
	CountryCode          string      `json:"countryCode"          description:"country code"`                                       // country code
	Currency             string      `json:"currency"             description:"currency"`                                           // currency
	PaymentId            string      `json:"paymentId"            description:"relative payment id"`                                // relative payment id
	RefundId             string      `json:"refundId"             description:"refund id (system generate)"`                        // refund id (system generate)
	RefundAmount         int64       `json:"refundAmount"         description:"refund amount, cent"`                                // refund amount, cent
	RefundComment        string      `json:"refundComment"        description:"refund comment"`                                     // refund comment
	Status               int         `json:"status"               description:"status。10-pending，20-success，30-failure, 40-cancel"` // status。10-pending，20-success，30-failure, 40-cancel
	RefundTime           int64       `json:"refundTime"           description:"refund success time"`                                // refund success time
	GmtCreate            *gtime.Time `json:"gmtCreate"            description:"create time"`                                        // create time
	GmtModify            *gtime.Time `json:"gmtModify"            description:"update time"`                                        // update time
	GatewayRefundId      string      `json:"gatewayRefundId"      description:"gateway refund id"`                                  // gateway refund id
	AppId                string      `json:"appId"                description:"app id"`                                             // app id
	RefundCommentExplain string      `json:"refundCommentExplain" description:"refund comment"`                                     // refund comment
	ReturnUrl            string      `json:"returnUrl"            description:"return url after refund success"`                    // return url after refund success
	MetaData             string      `json:"metaData"             description:"meta_data(json)"`                                    // meta_data(json)
	UniqueId             string      `json:"uniqueId"             description:"unique id"`                                          // unique id
	SubscriptionId       string      `json:"subscriptionId"       description:"subscription id"`                                    // subscription id
	CreateTime           int64       `json:"createTime"           description:"create utc time"`                                    // create utc time
	Type                 int         `json:"type"                 description:"1-gateway refund,2-mark refund"`                     // 1-gateway refund,2-mark refund
}
