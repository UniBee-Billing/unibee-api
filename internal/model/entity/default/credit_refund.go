// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CreditRefund is the golang structure for table credit_refund.
type CreditRefund struct {
	Id                     int64       `json:"id"                     description:"Id"`                                                                      // Id
	UserId                 uint64      `json:"userId"                 description:"user_id"`                                                                 // user_id
	CreditId               uint64      `json:"creditId"               description:"id of credit account"`                                                    // id of credit account
	Currency               string      `json:"currency"               description:"currency"`                                                                // currency
	InvoiceId              string      `json:"invoiceId"              description:"invoice_id"`                                                              // invoice_id
	CreditPaymentId        string      `json:"creditPaymentId"        description:"credit refund id"`                                                        // credit refund id
	CreditRefundId         string      `json:"creditRefundId"         description:"credit refund id"`                                                        // credit refund id
	ExternalCreditRefundId string      `json:"externalCreditRefundId" description:"external credit refund id"`                                               // external credit refund id
	RefundAmount           int64       `json:"refundAmount"           description:"total refund amount,cent"`                                                // total refund amount,cent
	RefundTime             int64       `json:"refundTime"             description:"refund time"`                                                             // refund time
	Name                   string      `json:"name"                   description:"recharge transaction title"`                                              // recharge transaction title
	Description            string      `json:"description"            description:"recharge transaction description"`                                        // recharge transaction description
	GmtCreate              *gtime.Time `json:"gmtCreate"              description:"create time"`                                                             // create time
	GmtModify              *gtime.Time `json:"gmtModify"              description:"update time"`                                                             // update time
	CreateTime             int64       `json:"createTime"             description:"create utc time"`                                                         // create utc time
	MerchantId             uint64      `json:"merchantId"             description:"merchant id"`                                                             // merchant id
	AccountType            int         `json:"accountType"            description:"type of credit account, 1-main recharge account, 2-promo credit account"` // type of credit account, 1-main recharge account, 2-promo credit account
}
