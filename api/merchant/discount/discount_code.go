package discount

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type ListReq struct {
	g.Meta          `path:"/list" tags:"Discount" method:"get" summary:"Get Discount Code List" dc:"Get discountCode list"`
	DiscountType    []int  `json:"discountType"  dc:"discount_type, 1-percentage, 2-fixed_amount" `
	BillingType     []int  `json:"billingType"  dc:"billing_type, 1-one-time, 2-recurring" `
	Status          []int  `json:"status" dc:"status, 1-editable, 2-active, 3-deactive, 4-expire, 10-archive" `
	Code            string `json:"code" dc:"Filter Code"  `
	SearchKey       string `json:"searchKey" dc:"Search Key, code or name"  `
	Currency        string `json:"currency" dc:"Filter Currency"  `
	SortField       string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType        string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int    `json:"page"  dc:"Page, Start 0" `
	Count           int    `json:"count"  dc:"Count Of Per Page" `
	CreateTimeStart int64  `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd   int64  `json:"createTimeEnd" dc:"CreateTimeEnd" `
}

type ListRes struct {
	Discounts []*detail.MerchantDiscountCodeDetail `json:"discounts" dc:"Discount Object List"`
	Total     int                                  `json:"total" dc:"Total"`
}

type DetailReq struct {
	g.Meta `path:"/detail" tags:"Discount" method:"get,post" summary:"Get Merchant Discount Detail"`
	Id     uint64 `json:"id"                 dc:"The discount's Id" v:"required"`
}

type DetailRes struct {
	Discount *detail.MerchantDiscountCodeDetail `json:"discount" dc:"Discount Object"`
}

type NewReq struct {
	g.Meta             `path:"/new" tags:"Discount" method:"post" summary:"New Discount Code" dc:"Create a new discount code, code can used in onetime or subscription purchase to make discount"`
	Code               string                 `json:"code" dc:"The discount's unique code, customize by merchant" v:"required"`
	Name               *string                `json:"name"              dc:"The discount's name"`                                                                                                                                                                                                                                                                    // name
	BillingType        int                    `json:"billingType"       dc:"The billing type of the discount code, 1-one-time, 2-recurring, define the situation the code can be used, the code of one-time billing_type can used for all situation that effect only once, the code of recurring billing_tye can only used for subscription purchase"  v:"required"` // billing_type, 1-one-time, 2-recurring
	DiscountType       int                    `json:"discountType"      dc:"The discount type of the discount code, 1-percentage, 2-fixed_amount, the discountType of code, the discountPercentage will be effect when discountType is percentage, the discountAmount and currency will be effect when discountTYpe is fixed_amount"  v:"required"`                  // discount_type, 1-percentage, 2-fixed_amount
	DiscountAmount     int64                  `json:"discountAmount"    dc:"The discount amount of the discount code, available when discount_type is fixed_amount"`                                                                                                                                                                                                 // amount of discount, available when discount_type is fixed_amount
	DiscountPercentage int64                  `json:"discountPercentage" dc:"The discount percentage of discount code, 100=1%, available when discount_type is percentage"`                                                                                                                                                                                          // percentage of discount, 100=1%, available when discount_type is percentage
	Currency           string                 `json:"currency"          dc:"The discount currency of discount code, available when discount_type is fixed_amount"`                                                                                                                                                                                                   // currency of discount, available when discount_type is fixed_amount
	CycleLimit         int                    `json:"cycleLimit"         dc:"The count limitation of subscription cycle, each subscription is valid separately , 0-no limit"`                                                                                                                                                                                        // the count limitation of subscription cycle , 0-no limit
	StartTime          *int64                 `json:"startTime"         dc:"The start time of discount code can effect, utc time"  v:"required"`                                                                                                                                                                                                                     // start of discount available utc time
	EndTime            *int64                 `json:"endTime"           dc:"The end time of discount code can effect, utc time"  v:"required"`
	PlanApplyType      *int                   `json:"planApplyType"      description:"plan apply type, 0-apply for all, 1-apply for plans specified, 2-exclude for plans specified"`
	PlanIds            []int64                `json:"planIds"  dc:"Ids of plan which discount code can effect, default effect all plans if not set" `
	Quantity           *uint64                `json:"quantity"           description:"Quantity of code, default 0, set 0 to disable quantity management"`
	Advance            *bool                  `json:"advance"            description:"AdvanceConfig, 0-false,1-true, will enable all advance config if set true"` // AdvanceConfig,  0-false,1-true, will enable all advance config if set 1
	UserScope          *int                   `json:"userScope"  dc:"AdvanceConfig, Apply user scope,0-for all, 1-for only new user, 2-for only renewals, renewals is upgrade&downgrade&renew"`
	UpgradeOnly        *bool                  `json:"upgradeOnly"  dc:"AdvanceConfig, true or false, will forbid for all except same interval upgrade action if set true" `
	UpgradeLongerOnly  *bool                  `json:"upgradeLongPlanOnly"  dc:"AdvanceConfig, true or false, will forbid for all except upgrade to longer plan if set true" `
	UserLimit          *int                   `json:"userLimit"         dc:"AdvanceConfig, The limit of every customer can apply, the recurring apply not involved, 0-unlimited"`
	Metadata           map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}

type NewRes struct {
	Discount *bean.MerchantDiscountCode `json:"discount" dc:"Discount Object"`
}

type EditReq struct {
	g.Meta             `path:"/edit" tags:"Discount" method:"post" summary:"Edit Discount Code" dc:"Edit the discount code before activate"`
	Id                 uint64                 `json:"id"                 dc:"The discount's Id" v:"required"`
	Name               *string                `json:"name"              dc:"The discount's name"`                                                                                                                                                                                                                                                      // name
	BillingType        int                    `json:"billingType"       dc:"The billing type of the discount code, 1-one-time, 2-recurring, define the situation the code can be used, the code of one-time billing_type can used for all situation that effect only once, the code of recurring billing_tye can only used for subscription purchase"` // billing_type, 1-one-time, 2-recurring
	DiscountType       int                    `json:"discountType"      dc:"The discount type of the discount code, 1-percentage, 2-fixed_amount, the discountType of code, the discountPercentage will be effect when discountType is percentage, the discountAmount and currency will be effect when discountTYpe is fixed_amount"`                  // discount_type, 1-percentage, 2-fixed_amount
	DiscountAmount     int64                  `json:"discountAmount"    dc:"The discount amount of the discount code, available when discount_type is fixed_amount"`                                                                                                                                                                                   // amount of discount, available when discount_type is fixed_amount
	DiscountPercentage int64                  `json:"discountPercentage" dc:"The discount percentage of discount code, 100=1%, available when discount_type is percentage"`                                                                                                                                                                            // percentage of discount, 100=1%, available when discount_type is percentage
	Currency           string                 `json:"currency"          dc:"The discount currency of discount code, available when discount_type is fixed_amount"`                                                                                                                                                                                     // currency of discount, available when discount_type is fixed_amount
	CycleLimit         int                    `json:"cycleLimit"         dc:"The count limitation of subscription cycle，each subscription is valid separately, 0-no limit"`
	StartTime          *int64                 `json:"startTime"         dc:"The start time of discount code can effect, editable after activate, utc time"`
	EndTime            *int64                 `json:"endTime"           dc:"The end time of discount code can effect, editable after activate, utc time"`
	PlanApplyType      *int                   `json:"planApplyType"      description:"plan apply type, 0-apply for all, 1-apply for plans specified, 2-exclude for plans specified"`
	PlanIds            []int64                `json:"planIds"  dc:"Ids of plan which discount code can effect, default effect all plans if not set" `
	Quantity           *uint64                `json:"quantity"           description:"Quantity of code, default 0, set 0 to disable quantity management"`
	Advance            *bool                  `json:"advance"            description:"AdvanceConfig, 0-false,1-true, will enable all advance config if set true"` // AdvanceConfig,  0-false,1-true, will enable all advance config if set 1
	UserScope          *int                   `json:"userScope"  dc:"AdvanceConfig, Apply user scope,0-for all, 1-for only new user, 2-for only renewals, renewals is upgrade&downgrade&renew"`
	UpgradeOnly        *bool                  `json:"upgradeOnly"  dc:"AdvanceConfig, true or false, will forbid for all except same interval upgrade action if set true" `
	UpgradeLongerOnly  *bool                  `json:"upgradeLongPlanOnly"  dc:"AdvanceConfig, true or false, will forbid for all except upgrade to longer plan if set true" `
	UserLimit          *int                   `json:"userLimit"         dc:"AdvanceConfig, The limit of every customer can apply, the recurring apply not involved, 0-unlimited"`
	Metadata           map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}

type EditRes struct {
	Discount *bean.MerchantDiscountCode `json:"discount" dc:"Discount Object"`
}

type DeleteReq struct {
	g.Meta `path:"/delete" tags:"Discount" method:"post" summary:"Delete Discount Code" dc:"Delete discount code before activate"`
	Id     uint64 `json:"id"                 description:"The discount's Id" v:"required"`
}

type DeleteRes struct {
}

type ActivateReq struct {
	g.Meta `path:"/activate" tags:"Discount" method:"post" summary:"Activate Discount Code" dc:"Activate discount code, the discount code can only effect to payment or subscription after activated"`
	Id     uint64 `json:"id"                 description:"The discount's Id" v:"required"`
}

type ActivateRes struct {
}

type DeactivateReq struct {
	g.Meta `path:"/deactivate" tags:"Discount" method:"post" summary:"Deactivate Discount Code" dc:"Deactivate discount code"`
	Id     uint64 `json:"id"                 description:"The discount's Id" v:"required"`
}

type DeactivateRes struct {
}

type UserDiscountListReq struct {
	g.Meta          `path:"/user_discount_list" tags:"User Discount" method:"get" summary:"Get User Discount Code List" dc:"Get user discountCode list"`
	Id              uint64   `json:"id"                 description:"The discount's Id" v:"required"`
	UserIds         []uint64 `json:"userIds" dc:"Filter UserIds Default All" `
	Email           string   `json:"email" dc:"Filter Email Default All" `
	PlanIds         []uint64 `json:"planIds" dc:"Filter PlanIds Default All" `
	SortField       string   `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType        string   `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int      `json:"page"  dc:"Page, Start 0" `
	Count           int      `json:"count"  dc:"Count Of Per Page" `
	CreateTimeStart int64    `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd   int64    `json:"createTimeEnd" dc:"CreateTimeEnd" `
}

type UserDiscountListRes struct {
	UserDiscounts []*detail.MerchantUserDiscountCodeDetail `json:"userDiscounts" dc:"User Discount Object List"`
	Total         int                                      `json:"total" dc:"Total"`
}

type PlanApplyPreviewReq struct {
	g.Meta         `path:"/plan_apply_preview" tags:"User Discount" method:"post" summary:"Plan Apply Preview" dc:"Check discount can apply to plan, Only check rules about plan，the actual usage is subject to the subscription interface"`
	Code           string `json:"code" dc:"The discount's unique code, customize by merchant" v:"required"`
	PlanId         int64  `json:"planId" dc:"The id of plan which code to apply, either planId or externalPlanId is needed"`
	ExternalPlanId string `json:"externalPlanId" dc:"The externalId of plan which code to apply, either planId or externalPlanId is needed"`
	//SubscriptionId     string `json:"subscriptionId"            description:"SubscriptionId"`
	IsUpgrade                  *bool  `json:"isUpgrade"            description:"IsUpgrade"`
	IsChangeToSameIntervalPlan *bool  `json:"isChangeToSameIntervalPlan"  description:"IsChangeToSameIntervalPlan"`
	IsChangeToLongPlan         *bool  `json:"isChangeToLongPlan"  description:"IsChangeToLongPlan"`
	Email                      string `json:"email"  description:"Email"`
}

type PlanApplyPreviewRes struct {
	Valid          bool                       `json:"valid" dc:"The apply preview result, true or false" `
	DiscountAmount int64                      `json:"discountAmount" dc:"The discount amount can apply to plan" `
	DiscountCode   *bean.MerchantDiscountCode `json:"discountCode" dc:"The discount code object" `
	FailureReason  string                     `json:"failureReason" dc:"The apply preview failure reason" `
}

type QuantityIncrementReq struct {
	g.Meta `path:"/quantity_increment" tags:"Discount" method:"post" summary:"Quantity Increment" dc:"Increase discount code quantity, if original quantity is 0, increase greater than 0 will enable quantity control"`
	Id     uint64 `json:"id"                 description:"The discount's Id" v:"required"`
	Amount uint64 `json:"amount" dc:"The discount quantity amount to increase, should greater than 0" `
}

type QuantityIncrementRes struct {
	DiscountCode *bean.MerchantDiscountCode `json:"discountCode" dc:"The discount code object" `
}

type QuantityDecrementReq struct {
	g.Meta `path:"/decrease_quantity" tags:"Discount" method:"post" summary:"Quantity Decrement" dc:"Decrease discount code quantity, the quantity after decreased should greater than 0, the action may disable quantity control if quantity decrease to 0 or lower than quantityUsed after decreased"`
	Id     uint64 `json:"id"                 description:"The discount's Id" v:"required"`
	Amount uint64 `json:"amount" dc:"The discount quantity Amount to decrease, should greater than 0" `
}

type QuantityDecrementRes struct {
	DiscountCode *bean.MerchantDiscountCode `json:"discountCode" dc:"The discount code object" `
}
