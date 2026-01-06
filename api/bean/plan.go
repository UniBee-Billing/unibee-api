package bean

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"strings"
	"unibee/internal/cmd/config"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type Plan struct {
	Id                     uint64                          `json:"id"                        description:""`
	MerchantId             uint64                          `json:"merchantId"                description:"merchant id"`                     // merchant id
	PlanName               string                          `json:"planName"                  description:"PlanName"`                        // PlanName
	InternalName           string                          `json:"internalName"              description:"PlanInternalName"`                //
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
	CancelAtTrialEnd       int                             `json:"cancelAtTrialEnd"          description:"whether cancel at subscription first trial end，0-false | 1-true, will pass to cancelAtPeriodEnd of subscription"` // whether cancel at subscripiton first trial end，0-false | 1-true, will pass to cancelAtPeriodEnd of subscription
	ExternalPlanId         string                          `json:"externalPlanId"            description:"external_user_id"`                                                                                                // external_user_id
	ProductId              int64                           `json:"productId"                 description:"product id"`                                                                                                      // product id
	DisableAutoCharge      int                             `json:"disableAutoCharge"         description:"disable auto-charge, 0-false,1-true"`                                                                             // disable auto-charge, 0-false,1-true
	MetricLimits           []*PlanMetricLimitParam         `json:"metricLimits"  dc:"Plan's MetricLimit List" `
	MetricMeteredCharge    []*PlanMetricMeteredChargeParam `json:"metricMeteredCharge"  dc:"Plan's MetricMeteredCharge" `
	MetricRecurringCharge  []*PlanMetricMeteredChargeParam `json:"metricRecurringCharge"  dc:"Plan's MetricRecurringCharge" `
	CheckoutUrl            string                          `json:"checkoutUrl"                 description:"CheckoutUrl"`
	MultiCurrencies        []*PlanMultiCurrency            `json:"multiCurrencies"  dc:"Plan's MultiCurrencies" `
}

const MerchantMultiCurrenciesConfig = "MerchantMultiCurrenciesConfig"

func GetMerchantMultiCurrencyConfig(ctx context.Context, merchantId uint64) []*MerchantMultiCurrencyConfig {
	utility.Assert(merchantId > 0, "invalid merchantId")
	var one *entity.MerchantConfig
	err := dao.MerchantConfig.Ctx(ctx).
		Where(dao.MerchantConfig.Columns().MerchantId, merchantId).
		Where(dao.MerchantConfig.Columns().ConfigKey, MerchantMultiCurrenciesConfig).
		Scan(&one)
	if err != nil {
		one = nil
	}
	var configs = make([]*MerchantMultiCurrencyConfig, 0)
	if one != nil {
		_ = utility.UnmarshalFromJsonString(one.ConfigValue, &configs)
	}
	return configs
}

func SimplifyPlanWithContext(ctx context.Context, one *entity.Plan) *Plan {
	plan := SimplifyPlan(one)
	var multiCurrencies = make([]*PlanMultiCurrency, 0)
	if len(one.GatewayProductDescription) > 0 {
		_ = utility.UnmarshalFromJsonString(one.GatewayProductDescription, &multiCurrencies)
	}

	exchangeCurrencyMap := make(map[string]*PlanMultiCurrency)
	for _, multiCurrency := range multiCurrencies {
		//if multiCurrency.AutoExchange {
		//	multiCurrencies[i].ExchangeRate = currency_exchange.GetMerchantExchangeCurrencyRate(ctx, one.MerchantId, one.Currency, multiCurrency.Currency)
		//}
		//multiCurrencies[i].Amount = utility.ExchangeCurrencyConvert(one.Amount, one.Currency, multiCurrency.Currency, multiCurrency.ExchangeRate)
		exchangeCurrencyMap[multiCurrency.Currency] = multiCurrency
	}
	multiCurrencyConfigs := GetMerchantMultiCurrencyConfig(ctx, one.MerchantId)
	var targetMultiCurrencies = make([]*PlanMultiCurrency, 0)
	for _, multiCurrencyConfig := range multiCurrencyConfigs {
		if strings.ToUpper(multiCurrencyConfig.DefaultCurrency) == one.Currency {
			for _, multiCurrency := range multiCurrencyConfig.MultiCurrencies {
				if multiCurrency.Currency == one.Currency {
					continue
				}
				target := &PlanMultiCurrency{
					Currency:     multiCurrency.Currency,
					AutoExchange: multiCurrency.AutoExchange,
					ExchangeRate: multiCurrency.ExchangeRate,
					Amount:       utility.ExchangeCurrencyConvert(one.Amount, one.Currency, multiCurrency.Currency, multiCurrency.ExchangeRate),
				}
				if _, ok := exchangeCurrencyMap[multiCurrency.Currency]; ok {
					target.Disable = exchangeCurrencyMap[multiCurrency.Currency].Disable
				}
				targetMultiCurrencies = append(targetMultiCurrencies, target)
			}
			break
		}
	}
	plan.MultiCurrencies = targetMultiCurrencies
	return plan
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
		InternalName:           one.InternalName,
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
		DisableAutoCharge:      one.DisableAutoCharge,
		MetricLimits:           metricPlanCharge.MetricLimits,
		MetricMeteredCharge:    metricPlanCharge.MetricMeteredCharge,
		MetricRecurringCharge:  metricPlanCharge.MetricRecurringCharge,
		CheckoutUrl:            fmt.Sprintf("%s/checkout?planId=%d&env=%s", config.GetConfigInstance().Server.GetHostedPath(), one.Id, config.GetConfigInstance().Env),
	}
}

func SimplifyPlanList(ctx context.Context, ones []*entity.Plan) (list []*Plan) {
	if len(ones) == 0 {
		return make([]*Plan, 0)
	}
	for _, one := range ones {
		list = append(list, SimplifyPlanWithContext(ctx, one))
	}
	return list
}

type PlanMultiCurrency struct {
	Currency     string  `json:"currency" description:"target currency"`
	AutoExchange bool    `json:"autoExchange" description:"using https://app.exchangerate-api.com/ to update exchange rate if true, the exchange APIKey need setup first"`
	ExchangeRate float64 `json:"exchangeRate" description:"exchange rate, no setup required if AutoExchange is true"`
	Amount       int64   `json:"amount" description:"the amount of exchange rate"`
	Disable      bool    `json:"disable" description:"disable currency exchange"`
}

func GetPlanCurrencyAmount(ctx context.Context, p *entity.Plan, currency string) int64 {
	if currency == "" || strings.ToUpper(currency) == strings.ToUpper(p.Currency) {
		return p.Amount
	}
	multiCurrencyConfigs := GetMerchantMultiCurrencyConfig(ctx, p.MerchantId)
	for _, multiCurrency := range multiCurrencyConfigs {
		if strings.ToUpper(multiCurrency.DefaultCurrency) == strings.ToUpper(p.Currency) {
			for _, one := range multiCurrency.MultiCurrencies {
				if strings.ToUpper(one.Currency) == strings.ToUpper(currency) {
					return utility.ExchangeCurrencyConvert(p.Amount, p.Currency, one.Currency, one.ExchangeRate)
				}
			}
		}
	}
	utility.Assert(false, fmt.Sprintf("No currency(%s) config found for plan", currency))
	return p.Amount
}

func (p *Plan) CurrencyAmount(ctx context.Context, targetCurrency string) int64 {
	if targetCurrency == "" || strings.ToUpper(targetCurrency) == strings.ToUpper(p.Currency) {
		return p.Amount
	}

	for _, one := range p.MultiCurrencies {
		if strings.ToUpper(one.Currency) == strings.ToUpper(targetCurrency) {
			if one.Disable {
				utility.Assert(false, fmt.Sprintf("Exchange: currency(%s) disabled for plan(%d)", targetCurrency, p.Id))
			}
		}
	}

	multiCurrencyConfigs := GetMerchantMultiCurrencyConfig(ctx, p.MerchantId)
	for _, multiCurrency := range multiCurrencyConfigs {
		if strings.ToUpper(multiCurrency.DefaultCurrency) == strings.ToUpper(p.Currency) {
			for _, one := range multiCurrency.MultiCurrencies {
				if strings.ToUpper(one.Currency) == strings.ToUpper(targetCurrency) {
					return utility.ExchangeCurrencyConvert(p.Amount, p.Currency, one.Currency, one.ExchangeRate)
				}
			}
		}
	}
	utility.Assert(false, fmt.Sprintf("No currency(%s) config found for plan(%d)", targetCurrency, p.Id))
	return p.Amount
}
