package ro

import entity "unibee-api/internal/model/entity/oversea_pay"

type PlanSimplify struct {
	Id              uint64 `json:"id"                        description:""`
	MerchantId      uint64 `json:"merchantId"                description:"merchant id"`                                               // merchant id
	PlanName        string `json:"planName"                  description:"PlanName"`                                                  // PlanName
	Amount          int64  `json:"amount"                    description:"amount, cent, without tax"`                                 // amount, cent, without tax
	Currency        string `json:"currency"                  description:"currency"`                                                  // currency
	IntervalUnit    string `json:"intervalUnit"              description:"period unit,day|month|year|week"`                           // period unit,day|month|year|week
	IntervalCount   int    `json:"intervalCount"             description:"period unit count"`                                         // period unit count
	Description     string `json:"description"               description:"description"`                                               // description
	ImageUrl        string `json:"imageUrl"                  description:"image_url"`                                                 // image_url
	HomeUrl         string `json:"homeUrl"                   description:"home_url"`                                                  // home_url
	TaxScale        int    `json:"taxScale"                  description:"tax scale 1000 = 10%"`                                      // tax scale 1000 = 10%
	Type            int    `json:"type"                      description:"type，1-main plan，2-addon plan"`                             // type，1-main plan，2-addon plan
	Status          int    `json:"status"                    description:"status，1-editing，2-active，3-inactive，4-expired"`            // status，1-editing，2-active，3-inactive，4-expired
	BindingAddonIds string `json:"bindingAddonIds"           description:"binded addon planIds，split with ,"`                         // binded addon planIds，split with ,
	PublishStatus   int    `json:"publishStatus"             description:"1-UnPublish,2-Publish, Use For Display Plan At UserPortal"` // 1-UnPublish,2-Publish, Use For Display Plan At UserPortal
	CreateTime      int64  `json:"createTime"                description:"create utc time"`                                           // create utc time
}

func SimplifyPlan(one *entity.SubscriptionPlan) *PlanSimplify {
	if one == nil {
		return nil
	}
	return &PlanSimplify{
		Id:              one.Id,
		MerchantId:      one.MerchantId,
		PlanName:        one.PlanName,
		Amount:          one.Amount,
		Currency:        one.Currency,
		IntervalUnit:    one.IntervalUnit,
		IntervalCount:   one.IntervalCount,
		Description:     one.Description,
		ImageUrl:        one.ImageUrl,
		HomeUrl:         one.HomeUrl,
		TaxScale:        one.TaxScale,
		Type:            one.Type,
		Status:          one.Status,
		BindingAddonIds: one.BindingAddonIds,
		PublishStatus:   one.PublishStatus,
		CreateTime:      one.CreateTime,
	}
}

func SimplifyPlanList(ones []*entity.SubscriptionPlan) (list []*PlanSimplify) {
	if len(ones) == 0 {
		return make([]*PlanSimplify, 0)
	}
	for _, one := range ones {
		list = append(list, SimplifyPlan(one))
	}
	return list
}
