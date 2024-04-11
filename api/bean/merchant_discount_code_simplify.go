package bean

import (
	entity "unibee/internal/model/entity/oversea_pay"
)

type MerchantDiscountCodeSimplify struct {
	Id                 int64  `json:"id"             description:"userId"`                                                                         // userId
	MerchantId         uint64 `json:"merchantId"         description:"merchantId"`                                                                 // merchantId
	Name               string `json:"name"               description:"name"`                                                                       // name
	Code               string `json:"code"               description:"code"`                                                                       // code
	Status             int    `json:"status"             description:"status, 1-editable, 2-active, 3-deactive, 4-expire"`                         // status, 1-editable, 2-active, 3-deactive, 4-expire
	BillingType        int    `json:"billingType"        description:"billing_type, 1-one-time, 2-recurring"`                                      // billing_type, 1-one-time, 2-recurring
	DiscountType       int    `json:"discountType"       description:"discount_type, 1-percentage, 2-fixed_amount"`                                // discount_type, 1-percentage, 2-fixed_amount
	DiscountAmount     int64  `json:"discountAmount"     description:"amount of discount, available when discount_type is fixed_amount"`           // amount of discount, available when discount_type is fixed_amount
	DiscountPercentage int64  `json:"discountPercentage" description:"percentage of discount, 100=1%, available when discount_type is percentage"` // percentage of discount, 100=1%, available when discount_type is percentage
	Currency           string `json:"currency"           description:"currency of discount, available when discount_type is fixed_amount"`         // currency of discount, available when discount_type is fixed_amount
	UserLimit          int    `json:"userLimit"          description:"the limit of every user apply, 0-unlimited"`                                 // the limit of every user apply, 0-unlimited
	SubscriptionLimit  int    `json:"subscriptionLimit"  description:"the limit of every subscription apply, 0-unlimited"`                         // the limit of every subscription apply, 0-unlimited
	StartTime          int64  `json:"startTime"          description:"start of discount available utc time"`                                       // start of discount available utc time
	EndTime            int64  `json:"endTime"            description:"end of discount available utc time"`                                         // end of discount available utc time
	CreateTime         int64  `json:"createTime"     description:"create utc time"`                                                                // create utc time
}

func SimplifyMerchantDiscountCode(one *entity.MerchantDiscountCode) *MerchantDiscountCodeSimplify {
	if one == nil {
		return nil
	}
	return &MerchantDiscountCodeSimplify{
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
		UserLimit:          one.UserLimit,
		SubscriptionLimit:  one.SubscriptionLimit,
		StartTime:          one.StartTime,
		EndTime:            one.EndTime,
		CreateTime:         one.CreateTime,
	}
}
