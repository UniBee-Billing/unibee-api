package bean

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"strconv"
	"strings"
	entity "unibee/internal/model/entity/default"
)

type MerchantDiscountCode struct {
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
	PlanIds            []int64                `json:"planIds"  dc:"Ids of plan which discount code can effect, default effect all plans if not set" `
	Metadata           map[string]interface{} `json:"metadata"           description:""`
}

func SimplifyMerchantDiscountCode(one *entity.MerchantDiscountCode) *MerchantDiscountCode {
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
	return &MerchantDiscountCode{
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
		Metadata:           metadata,
	}
}
