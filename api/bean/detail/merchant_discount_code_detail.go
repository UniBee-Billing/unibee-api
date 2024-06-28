package detail

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"strconv"
	"strings"
	"unibee/api/bean"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
)

type MerchantDiscountCodeDetail struct {
	Id                 uint64                 `json:"id"                 description:"Id"`                                                                         // Id
	MerchantId         uint64                 `json:"merchantId"         description:"merchantId"`                                                                 // merchantId
	Name               string                 `json:"name"               description:"name"`                                                                       // name
	Code               string                 `json:"code"               description:"code"`                                                                       // code
	Status             int                    `json:"status"             description:"status, 1-editable, 2-active, 3-deactive, 4-expire"`                         // status, 1-editable, 2-active, 3-deactive, 4-expire
	BillingType        int                    `json:"billingType"        description:"billing_type, 1-one-time, 2-recurring"`                                      // billing_type, 1-one-time, 2-recurring
	DiscountType       int                    `json:"discountType"       description:"discount_type, 1-percentage, 2-fixed_amount"`                                // discount_type, 1-percentage, 2-fixed_amount
	DiscountAmount     int64                  `json:"discountAmount"     description:"amount of discount, available when discount_type is fixed_amount"`           // amount of discount, available when discount_type is fixed_amount
	DiscountPercentage int64                  `json:"discountPercentage" description:"percentage of discount, 100=1%, available when discount_type is percentage"` // percentage of discount, 100=1%, available when discount_type is percentage
	Currency           string                 `json:"currency"           description:"currency of discount, available when discount_type is fixed_amount"`         // currency of discount, available when discount_type is fixed_amount
	CycleLimit         int                    `json:"cycleLimit"         description:"the count limitation of subscription cycle , 0-no limit"`                    // the count limitation of subscription cycle , 0-no limit
	StartTime          int64                  `json:"startTime"          description:"start of discount available utc time"`                                       // start of discount available utc time
	EndTime            int64                  `json:"endTime"            description:"end of discount available utc time, 0-invalid"`                              // end of discount available utc time
	CreateTime         int64                  `json:"createTime"         description:"create utc time"`                                                            // create utc time
	PlanIds            []int64                `json:"planIds"  description:"Ids of plan which discount code can effect, default effect all plans if not set" `
	Plans              []*bean.PlanSimplify   `json:"plans"         description:"plans which discount code can effect, default effect all plans if not set"` // create utc time
	Metadata           map[string]interface{} `json:"metadata"           description:""`
}

func ConvertMerchantDiscountCodeDetail(ctx context.Context, one *entity.MerchantDiscountCode) *MerchantDiscountCodeDetail {
	if one == nil {
		return nil
	}
	var metadata = make(map[string]interface{})
	if len(one.MetaData) > 0 {
		err := gjson.Unmarshal([]byte(one.MetaData), &metadata)
		if err != nil {
			fmt.Printf("SimplifyPlan Unmarshal Metadata error:%s", err.Error())
		}
	}
	var planIds = make([]int64, 0)
	if len(one.PlanIds) > 0 {
		strList := strings.Split(one.PlanIds, ",")

		for _, s := range strList {
			num, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				fmt.Println("Internal Error converting string to int:", err)
			} else {
				planIds = append(planIds, num)
			}
		}
	}

	return &MerchantDiscountCodeDetail{
		Id:                 one.Id,
		MerchantId:         one.MerchantId,
		Name:               one.Name,
		Code:               one.Code,
		Status:             one.Status,
		BillingType:        one.BillingType,
		DiscountType:       one.DiscountType,
		DiscountAmount:     one.DiscountAmount,
		DiscountPercentage: one.DiscountPercentage,
		Currency:           one.Currency,
		CycleLimit:         one.CycleLimit,
		StartTime:          one.StartTime,
		EndTime:            one.EndTime,
		CreateTime:         one.CreateTime,
		PlanIds:            planIds,
		Plans:              bean.SimplifyPlanList(query.GetPlansByIds(ctx, planIds)),
		Metadata:           metadata,
	}
}

type MerchantUserDiscountCodeDetail struct {
	Id             int64                     `json:"id"             description:"ID"`         // ID
	MerchantId     uint64                    `json:"merchantId"     description:"merchantId"` // merchantId
	User           *bean.UserAccountSimplify `json:"user"     description:"User"`
	Code           string                    `json:"code"           description:"code"` // code
	Plan           *bean.PlanSimplify        `json:"plan"     description:"Plan"`
	SubscriptionId string                    `json:"subscriptionId" description:"subscription_id"` // subscription_id
	PaymentId      string                    `json:"paymentId"      description:"payment_id"`      // payment_id
	InvoiceId      string                    `json:"invoiceId"      description:"invoice_id"`      // invoice_id
	CreateTime     int64                     `json:"createTime"     description:"create utc time"` // create utc time
	ApplyAmount    int64                     `json:"applyAmount"    description:"apply_amount"`    // apply_amount
	Currency       string                    `json:"currency"       description:"currency"`        // currency
}

func ConvertMerchantUserDiscountCodeDetail(ctx context.Context, one *entity.MerchantUserDiscountCode) *MerchantUserDiscountCodeDetail {
	if one == nil {
		return nil
	}
	planId, _ := strconv.ParseInt(one.PlanId, 10, 64)
	if planId <= 0 {
		sub := query.GetSubscriptionBySubscriptionId(ctx, one.SubscriptionId)
		if sub != nil {
			planId = int64(sub.PlanId)
		}
	}
	var plan *bean.PlanSimplify
	if planId > 0 {
		plan = bean.SimplifyPlan(query.GetPlanById(ctx, uint64(planId)))
	}
	return &MerchantUserDiscountCodeDetail{
		Id:             one.Id,
		MerchantId:     one.MerchantId,
		User:           bean.SimplifyUserAccount(query.GetUserAccountById(ctx, one.UserId)),
		Code:           one.Code,
		Plan:           plan,
		SubscriptionId: one.SubscriptionId,
		PaymentId:      one.PaymentId,
		InvoiceId:      one.InvoiceId,
		CreateTime:     one.CreateTime,
		ApplyAmount:    one.ApplyAmount,
		Currency:       one.Currency,
	}
}
