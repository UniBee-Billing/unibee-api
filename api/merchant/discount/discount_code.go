package discount

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type ListReq struct {
	g.Meta `path:"/list" tags:"Role" method:"get" summary:"Get Merchant Role List"`
}

type ListRes struct {
	MerchantDiscountCodes []*bean.MerchantDiscountCodeSimplify `json:"merchantDiscountCodes" dc:"MerchantDiscountCodes"`
}

type NewReq struct {
	g.Meta             `path:"/new" tags:"DiscountCode" method:"post" summary:"New Merchant Discount Code"`
	Code               string `json:"code" dc:"Code" v:"required"`
	Name               string `json:"name"              description:"name"`                                                                        // name
	BillingType        int    `json:"billingType"       description:"billing_type, 1-one-time, 2-recurring"  v:"required"`                         // billing_type, 1-one-time, 2-recurring
	DiscountType       int    `json:"discountType"      description:"discount_type, 1-percentage, 2-fixed_amount"  v:"required"`                   // discount_type, 1-percentage, 2-fixed_amount
	DiscountAmount     int64  `json:"discountAmount"    description:"amount of discount, available when discount_type is fixed_amount"`            // amount of discount, available when discount_type is fixed_amount
	DiscountPercentage int64  `json:"discountPercentage" description:"percentage of discount, 100=1%, available when discount_type is percentage"` // percentage of discount, 100=1%, available when discount_type is percentage
	Currency           string `json:"currency"          description:"currency of discount, available when discount_type is fixed_amount"`          // currency of discount, available when discount_type is fixed_amount
	UserLimit          int    `json:"userLimit"         description:"the limit of every user apply, 0-unlimited"`                                  // the limit of every user apply, 0-unlimited
	SubscriptionLimit  int    `json:"subscriptionLimit" description:"the limit of every subscription apply, 0-unlimited"`                          // the limit of every subscription apply, 0-unlimited
	StartTime          int64  `json:"startTime"         description:"start of discount available utc time"  v:"required"`                          // start of discount available utc time
	EndTime            int64  `json:"endTime"           description:"end of discount available utc time"  v:"required"`
}

type NewRes struct {
}

type EditReq struct {
	g.Meta             `path:"/edit" tags:"DiscountCode" method:"post" summary:"Edit Merchant Discount Code"`
	Code               string `json:"code" dc:"Code" v:"required"`
	Name               string `json:"name"              description:"name"`                                                                        // name
	BillingType        int    `json:"billingType"       description:"billing_type, 1-one-time, 2-recurring"`                                       // billing_type, 1-one-time, 2-recurring
	DiscountType       int    `json:"discountType"      description:"discount_type, 1-percentage, 2-fixed_amount"`                                 // discount_type, 1-percentage, 2-fixed_amount
	DiscountAmount     int64  `json:"discountAmount"    description:"amount of discount, available when discount_type is fixed_amount"`            // amount of discount, available when discount_type is fixed_amount
	DiscountPercentage int64  `json:"discountPercentage" description:"percentage of discount, 100=1%, available when discount_type is percentage"` // percentage of discount, 100=1%, available when discount_type is percentage
	Currency           string `json:"currency"          description:"currency of discount, available when discount_type is fixed_amount"`          // currency of discount, available when discount_type is fixed_amount
	UserLimit          int    `json:"userLimit"         description:"the limit of every user apply, 0-unlimited"`                                  // the limit of every user apply, 0-unlimited
	SubscriptionLimit  int    `json:"subscriptionLimit" description:"the limit of every subscription apply, 0-unlimited"`                          // the limit of every subscription apply, 0-unlimited
	StartTime          int64  `json:"startTime"         description:"start of discount available utc time"`                                        // start of discount available utc time
	EndTime            int64  `json:"endTime"           description:"end of discount available utc time"`
}

type EditRes struct {
}

type DeleteReq struct {
	g.Meta `path:"/delete" tags:"DiscountCode" method:"post" summary:"Delete Merchant Discount Code"`
	Code   string `json:"code" dc:"Code" v:"required"`
}

type DeleteRes struct {
}
