package discount

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type ListReq struct {
	g.Meta          `path:"/list" tags:"Discount" method:"get" summary:"DiscountCodeList" dc:"Get discountCode list"`
	DiscountType    []int  `json:"discountType"  dc:"discount_type, 1-percentage, 2-fixed_amount" `
	BillingType     []int  `json:"billingType"  dc:"billing_type, 1-one-time, 2-recurring" `
	Status          []int  `json:"status" dc:"status, 1-editable, 2-active, 3-deactive, 4-expire" `
	Code            string `json:"code" dc:"Filter Code"  `
	Currency        string `json:"currency" dc:"Filter Currency"  `
	SortField       string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType        string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int    `json:"page"  dc:"Page, Start 0" `
	Count           int    `json:"count"  dc:"Count Of Per Page" `
	CreateTimeStart int64  `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd   int64  `json:"createTimeEnd" dc:"CreateTimeEnd" `
}

type ListRes struct {
	Discounts []*bean.MerchantDiscountCodeSimplify `json:"discounts" dc:"Discount Object List"`
	Total     int                                  `json:"total" dc:"Total"`
}

type DetailReq struct {
	g.Meta `path:"/detail" tags:"Discount" method:"get,post" summary:"Merchant Discount Detail"`
	Id     uint64 `json:"id"                 dc:"The discount's Id" v:"required"`
}

type DetailRes struct {
	Discount *detail.MerchantDiscountCodeDetail `json:"discount" dc:"Discount Object"`
}

type NewReq struct {
	g.Meta             `path:"/new" tags:"Discount" method:"post" summary:"NewDiscountCode" dc:"Create a new discount code, code can used in onetime or subscription purchase to make discount"`
	Code               string `json:"code" dc:"The discount's unique code, customize by merchant" v:"required"`
	Name               string `json:"name"              dc:"The discount's name"`                                                                                                                                                                                                                                                                    // name
	BillingType        int    `json:"billingType"       dc:"The billing type of the discount code, 1-one-time, 2-recurring, define the situation the code can be used, the code of one-time billing_type can used for all situation that effect only once, the code of recurring billing_tye can only used for subscription purchase"  v:"required"` // billing_type, 1-one-time, 2-recurring
	DiscountType       int    `json:"discountType"      dc:"The discount type of the discount code, 1-percentage, 2-fixed_amount, the discountType of code, the discountPercentage will be effect when discountType is percentage, the discountAmount and currency will be effect when discountTYpe is fixed_amount"  v:"required"`                  // discount_type, 1-percentage, 2-fixed_amount
	DiscountAmount     int64  `json:"discountAmount"    dc:"The discount amount of the discount code, available when discount_type is fixed_amount"`                                                                                                                                                                                                 // amount of discount, available when discount_type is fixed_amount
	DiscountPercentage int64  `json:"discountPercentage" dc:"The discount percentage of discount code, 100=1%, available when discount_type is percentage"`                                                                                                                                                                                          // percentage of discount, 100=1%, available when discount_type is percentage
	Currency           string `json:"currency"          dc:"The discount currency of discount code, available when discount_type is fixed_amount"`                                                                                                                                                                                                   // currency of discount, available when discount_type is fixed_amount
	//UserLimit          int    `json:"userLimit"         dc:"The limit of every customer can effect, 0-unlimited"`                                                                                                                                                                                                                                    // the limit of every user apply, 0-unlimited
	CycleLimit int                    `json:"cycleLimit"         dc:"The count limitation of subscription cycle, each subscription is valid separately , 0-no limit"` // the count limitation of subscription cycle , 0-no limit
	StartTime  int64                  `json:"startTime"         dc:"The start time of discount code can effect, utc time"  v:"required"`                              // start of discount available utc time
	EndTime    int64                  `json:"endTime"           dc:"The end time of discount code can effect, utc time"  v:"required"`
	PlanIds    []int64                `json:"planIds"  dc:"Ids of plan which discount code can effect, default effect all plans if not set" `
	Metadata   map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}

type NewRes struct {
	Discount *bean.MerchantDiscountCodeSimplify `json:"discount" dc:"Discount Object"`
}

type EditReq struct {
	g.Meta             `path:"/edit" tags:"Discount" method:"post" summary:"EditDiscountCode" dc:"Edit the discount code before activate"`
	Id                 uint64 `json:"id"                 dc:"The discount's Id" v:"required"`
	Name               string `json:"name"              dc:"The discount's name"`                                                                                                                                                                                                                                                      // name
	BillingType        int    `json:"billingType"       dc:"The billing type of the discount code, 1-one-time, 2-recurring, define the situation the code can be used, the code of one-time billing_type can used for all situation that effect only once, the code of recurring billing_tye can only used for subscription purchase"` // billing_type, 1-one-time, 2-recurring
	DiscountType       int    `json:"discountType"      dc:"The discount type of the discount code, 1-percentage, 2-fixed_amount, the discountType of code, the discountPercentage will be effect when discountType is percentage, the discountAmount and currency will be effect when discountTYpe is fixed_amount"`                  // discount_type, 1-percentage, 2-fixed_amount
	DiscountAmount     int64  `json:"discountAmount"    dc:"The discount amount of the discount code, available when discount_type is fixed_amount"`                                                                                                                                                                                   // amount of discount, available when discount_type is fixed_amount
	DiscountPercentage int64  `json:"discountPercentage" dc:"The discount percentage of discount code, 100=1%, available when discount_type is percentage"`                                                                                                                                                                            // percentage of discount, 100=1%, available when discount_type is percentage
	Currency           string `json:"currency"          dc:"The discount currency of discount code, available when discount_type is fixed_amount"`                                                                                                                                                                                     // currency of discount, available when discount_type is fixed_amount
	//UserLimit          int    `json:"userLimit"         dc:"The limit of every user effect, 0-unlimited"`                                                                                                                                                                                                                              // the limit of every user apply, 0-unlimited
	CycleLimit int                    `json:"cycleLimit"         dc:"The count limitation of subscription cycle，each subscription is valid separately, 0-no limit"` // the count limitation of subscription cycle , 0-no limit
	StartTime  int64                  `json:"startTime"         dc:"The start time of discount code can effect, utc time"`                                          // start of discount available utc time
	EndTime    int64                  `json:"endTime"           dc:"The end time of discount code can effect, utc time"`
	PlanIds    []int64                `json:"planIds"  dc:"Ids of plan which discount code can effect, default effect all plans if not set" `
	Metadata   map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}

type EditRes struct {
	Discount *bean.MerchantDiscountCodeSimplify `json:"discount" dc:"Discount Object"`
}

type DeleteReq struct {
	g.Meta `path:"/delete" tags:"Discount" method:"post" summary:"DeleteDiscountCode" dc:"Delete discount code before activate"`
	Id     uint64 `json:"id"                 description:"The discount's Id" v:"required"`
}

type DeleteRes struct {
}

type ActivateReq struct {
	g.Meta `path:"/activate" tags:"Discount" method:"post" summary:"ActivateDiscountCode" dc:"Activate discount code, the discount code can only effect to payment or subscription after activated"`
	Id     uint64 `json:"id"                 description:"The discount's Id" v:"required"`
}

type ActivateRes struct {
}

type DeactivateReq struct {
	g.Meta `path:"/deactivate" tags:"Discount" method:"post" summary:"DeactivateDiscountCode" dc:"Deactivate discount code"`
	Id     uint64 `json:"id"                 description:"The discount's Id" v:"required"`
}

type DeactivateRes struct {
}

type UserDiscountListReq struct {
	g.Meta          `path:"/user_discount_list" tags:"Discount" method:"get" summary:"UserDiscountCodeList" dc:"Get user discountCode list"`
	Id              uint64 `json:"id"                 description:"The discount's Id" v:"required"`
	SortField       string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType        string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int    `json:"page"  dc:"Page, Start 0" `
	Count           int    `json:"count"  dc:"Count Of Per Page" `
	CreateTimeStart int64  `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd   int64  `json:"createTimeEnd" dc:"CreateTimeEnd" `
}

type UserDiscountListRes struct {
	UserDiscounts []*detail.MerchantUserDiscountCodeDetail `json:"userDiscounts" dc:"User Discount Object List"`
	Total         int                                      `json:"total" dc:"Total"`
}
