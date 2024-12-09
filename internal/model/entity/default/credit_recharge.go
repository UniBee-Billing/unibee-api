// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CreditRecharge is the golang structure for table credit_recharge.
type CreditRecharge struct {
	Id                int64       `json:"id"                description:"Id"`                                                                      // Id
	UserId            uint64      `json:"userId"            description:"user_id"`                                                                 // user_id
	CreditId          uint64      `json:"creditId"          description:"id of credit account"`                                                    // id of credit account
	RechargeId        string      `json:"rechargeId"        description:"unique recharge id for credit account"`                                   // unique recharge id for credit account
	RechargeStatus    int         `json:"rechargeStatus"    description:"recharge status, 10-in charging，20-recharge success，30-recharge failed"`  // recharge status, 10-in charging，20-recharge success，30-recharge failed
	Currency          string      `json:"currency"          description:"currency"`                                                                // currency
	TotalAmount       int64       `json:"totalAmount"       description:"recharge total amount, cent"`                                             // recharge total amount, cent
	PaymentAmount     int64       `json:"paymentAmount"     description:"the payment amount for recharge"`                                         // the payment amount for recharge
	Name              string      `json:"name"              description:"recharge title"`                                                          // recharge title
	Description       string      `json:"description"       description:"recharge description"`                                                    // recharge description
	PaidTime          int64       `json:"paidTime"          description:"paid time"`                                                               // paid time
	GatewayId         uint64      `json:"gatewayId"         description:"payment gateway id"`                                                      // payment gateway id
	InvoiceId         string      `json:"invoiceId"         description:"invoice_id"`                                                              // invoice_id
	PaymentId         string      `json:"paymentId"         description:"paymentId"`                                                               // paymentId
	TotalRefundAmount int64       `json:"totalRefundAmount" description:"total refund amount,cent"`                                                // total refund amount,cent
	GmtCreate         *gtime.Time `json:"gmtCreate"         description:"create time"`                                                             // create time
	GmtModify         *gtime.Time `json:"gmtModify"         description:"update time"`                                                             // update time
	CreateTime        int64       `json:"createTime"        description:"create utc time"`                                                         // create utc time
	MerchantId        uint64      `json:"merchantId"        description:"merchant id"`                                                             // merchant id
	AccountType       int         `json:"accountType"       description:"type of credit account, 1-main recharge account, 2-promo credit account"` // type of credit account, 1-main recharge account, 2-promo credit account
}
