package bean

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type Plan struct {
	Id                     uint64                          `json:"id"                        description:""`
	MerchantId             uint64                          `json:"merchantId"                description:"merchant id"`                     // merchant id
	PlanName               string                          `json:"planName"                  description:"PlanName"`                        // PlanName
	Amount                 int64                           `json:"amount"                    description:"amount, cent, without tax"`       // amount, cent, without tax
	Currency               string                          `json:"currency"                  description:"currency"`                        // currency
	IntervalUnit           string                          `json:"intervalUnit"              description:"period unit,day|month|year|week"` // period unit,day|month|year|week
	IntervalCount          int                             `json:"intervalCount"             description:"period unit count"`               // period unit count
	Description            string                          `json:"description"               description:"description"`                     // description
	ImageUrl               string                          `json:"imageUrl"                  description:"image_url"`                       // image_url
	HomeUrl                string                          `json:"homeUrl"                   description:"home_url"`                        // home_url
	TaxPercentage          int                             `json:"taxPercentage"                  description:"TaxPercentage 1000 = 10%"`   // tax scale 1000 = 10%
	Type                   int                             `json:"type"                      description:"type，1-main plan，2-addon plan"`   // type，1-main plan，2-addon plan
	Status                 int                             `json:"status"                    description:"status，1-editing，2-active，3-inactive，4-soft archive, 5-hard archive"`
	BindingAddonIds        string                          `json:"bindingAddonIds"           description:"binded recurring addon planIds，split with ,"`               // binded addon planIds，split with ,
	BindingOnetimeAddonIds string                          `json:"bindingOnetimeAddonIds"    description:"binded onetime addon planIds，split with ,"`                 // binded onetime addon planIds，split with ,
	PublishStatus          int                             `json:"publishStatus"             description:"1-UnPublish,2-Publish, Use For Display Plan At UserPortal"` // 1-UnPublish,2-Publish, Use For Display Plan At UserPortal
	CreateTime             int64                           `json:"createTime"                description:"create utc time"`                                           // create utc time 	// product description
	ExtraMetricData        string                          `json:"extraMetricData"           description:""`                                                          //
	Metadata               map[string]interface{}          `json:"metadata"                  description:""`
	GasPayer               string                          `json:"gasPayer"                  description:"who pay the gas, merchant|user"` // who pay the gas, merchant|user
	TrialAmount            int64                           `json:"trialAmount"                description:"price of trial period"`         // price of trial period
	TrialDurationTime      int64                           `json:"trialDurationTime"         description:"duration of trial"`              // duration of trial
	TrialDemand            string                          `json:"trialDemand"               description:""`
	CancelAtTrialEnd       int                             `json:"cancelAtTrialEnd"          description:"whether cancel at subscripiton first trial end，0-false | 1-true, will pass to cancelAtPeriodEnd of subscription"` // whether cancel at subscripiton first trial end，0-false | 1-true, will pass to cancelAtPeriodEnd of subscription
	ExternalPlanId         string                          `json:"externalPlanId"            description:"external_user_id"`                                                                                                // external_user_id
	ProductId              int64                           `json:"productId"                 description:"product id"`                                                                                                      // product id
	MetricLimits           []*PlanMetricLimitParam         `json:"metricLimits"  dc:"Plan's MetricLimit List" `
	MetricMeteredCharge    []*PlanMetricMeteredChargeParam `json:"metricMeteredCharge"  dc:"Plan's MetricMeteredCharge" `
	MetricRecurringCharge  []*PlanMetricMeteredChargeParam `json:"metricRecurringCharge"  dc:"Plan's MetricRecurringCharge" `
}

func SimplifyPlan(one *entity.Plan) *Plan {
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
	var metricPlanCharge = &MetricPlanBindingEntity{}
	if len(one.MetricCharge) > 0 {
		_ = utility.UnmarshalFromJsonString(one.MetricCharge, &metricPlanCharge)
	}
	return &Plan{
		Id:                     one.Id,
		MerchantId:             one.MerchantId,
		PlanName:               one.PlanName,
		Amount:                 one.Amount,
		Currency:               one.Currency,
		IntervalUnit:           one.IntervalUnit,
		IntervalCount:          one.IntervalCount,
		Description:            one.Description,
		ImageUrl:               one.ImageUrl,
		HomeUrl:                one.HomeUrl,
		TaxPercentage:          one.TaxPercentage,
		Type:                   one.Type,
		Status:                 one.Status,
		BindingAddonIds:        one.BindingAddonIds,
		BindingOnetimeAddonIds: one.BindingOnetimeAddonIds,
		PublishStatus:          one.PublishStatus,
		CreateTime:             one.CreateTime,
		ExtraMetricData:        one.ExtraMetricData,
		Metadata:               metadata,
		GasPayer:               one.GasPayer,
		TrialDemand:            one.TrialDemand,
		TrialDurationTime:      one.TrialDurationTime,
		TrialAmount:            one.TrialAmount,
		CancelAtTrialEnd:       one.CancelAtTrialEnd,
		ExternalPlanId:         one.ExternalPlanId,
		ProductId:              one.ProductId,
		MetricLimits:           metricPlanCharge.MetricLimits,
		MetricMeteredCharge:    metricPlanCharge.MetricMeteredCharge,
		MetricRecurringCharge:  metricPlanCharge.MetricRecurringCharge,
	}
}

func SimplifyPlanList(ones []*entity.Plan) (list []*Plan) {
	if len(ones) == 0 {
		return make([]*Plan, 0)
	}
	for _, one := range ones {
		list = append(list, SimplifyPlan(one))
	}
	return list
}
