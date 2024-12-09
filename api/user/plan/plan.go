package plan

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	merhcnatPlan "unibee/api/bean/detail"
)

type ListReq struct {
	g.Meta     `path:"/list" tags:"User-Plan" method:"get,post" summary:"Plan List"`
	ProductIds []int64 `json:"productIds"  dc:"filter id list of product, will use default product if not specified " `
	Type       []int   `json:"type" dc:"Default All，,1-main plan，2-addon plan" `
	Currency   string  `json:"currency" dc:"Currency"  `
	SearchKey  string  `json:"searchKey" dc:"Search Key, plan name or description"  `
	Page       int     `json:"page"  dc:"Page, Start 0" `
	Count      int     `json:"count"  dc:"Count Of Per Page" `
}
type ListRes struct {
	Plans []*merhcnatPlan.PlanDetail `json:"plans" dc:"Plan Detail"`
	Total int                        `json:"total" dc:"Total"`
}

type CodeApplyPreviewReq struct {
	g.Meta         `path:"/code_apply_preview" tags:"User-Plan" method:"post" summary:"CodeApplyPreview" dc:"Check discount can apply to plan, Only check rules about plan，the actual usage is subject to the subscription interface"`
	Code           string `json:"code" dc:"The discount's unique code, customize by merchant" v:"required"`
	PlanId         int64  `json:"planId" dc:"The id of plan which code to apply, either planId or externalPlanId is needed"`
	ExternalPlanId string `json:"externalPlanId" dc:"The externalId of plan which code to apply, either planId or externalPlanId is needed"`
}

type CodeApplyPreviewRes struct {
	Valid          bool                       `json:"valid" dc:"The apply preview result, true or false" `
	DiscountAmount int64                      `json:"discountAmount" dc:"The discount amount can apply to plan" `
	DiscountCode   *bean.MerchantDiscountCode `json:"discountCode" dc:"The discount code object" `
	FailureReason  string                     `json:"failureReason" dc:"The apply preview failure reason" `
}
