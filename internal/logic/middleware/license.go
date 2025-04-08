package middleware

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	go_redismq "github.com/jackyang-hk/go-redismq"
	"strconv"
	"unibee/api/bean"
	"unibee/internal/cmd/config"
	_interface "unibee/internal/interface"
	context2 "unibee/internal/interface/context"
	"unibee/utility"
)

func GetMerchantAPIRateLimit(ctx context.Context, merchantId uint64) int {
	var defaultRateLimit = 50
	if !config.GetConfigInstance().IsProd() {
		defaultRateLimit = 10
	}
	cache, err := g.Redis().Get(ctx, fmt.Sprintf("UniBee#Cloud#MerchantAPIRateLimit#%d", merchantId))
	if err == nil && cache != nil {
		if !cache.IsEmpty() {
			g.Log().Infof(ctx, "GetMerchantAPIRateLimit Cache hint:%d data:%s", merchantId, cache)
			qps := cache.Int()
			if qps > defaultRateLimit {
				return qps
			}
		}
	}
	license := GetMerchantLicense(ctx, merchantId)
	if license != nil && license.IsPaid && license.Plan != nil && license.Plan.Metadata != nil {
		if license.Plan.Metadata["QPSLimit"] != nil {
			if v, ok := license.Plan.Metadata["QPSLimit"].(float64); ok {
				qps := int(v)
				if qps > defaultRateLimit {
					_, _ = g.Redis().Set(ctx, fmt.Sprintf("UniBee#Cloud#MerchantAPIRateLimit#%d", merchantId), qps)
					_, _ = g.Redis().Expire(ctx, fmt.Sprintf("UniBee#Cloud#MerchantAPIRateLimit#%d", merchantId), 60)
					return qps
				}
			}
		}
	}
	return defaultRateLimit
}

func GetMerchantMemberLimit(ctx context.Context, merchantId uint64) int {
	var defaultMemberLimit = -1
	license := GetMerchantLicense(ctx, merchantId)
	if license != nil && license.IsPaid && license.Plan != nil && license.Plan.Metadata != nil {
		if license.Plan.Metadata["MemberCount"] != nil {
			memberCount, err := strconv.Atoi(fmt.Sprintf("%s", license.Plan.Metadata["MemberCount"]))
			if err == nil && memberCount != 0 {
				return memberCount
			}
		}
	}
	return defaultMemberLimit
}

func GetMerchantLicense(ctx context.Context, merchantId uint64) (one *License) {
	if merchantId <= 0 {
		return nil
	}
	license := go_redismq.Invoke(ctx, &go_redismq.InvoiceRequest{
		Group:   "GID_UniBee_License",
		Method:  "GetLicenseByMerchantId",
		Request: merchantId,
	}, 0)
	if license == nil || !license.Status {
		return nil
	}
	err := utility.UnmarshalFromJsonString(utility.MarshalToJsonString(license.Response), &one)
	if err != nil {
		return nil
	}
	return one
}

func IsPremiumVersion(ctx context.Context, merchantId uint64) bool {
	if config.GetConfigInstance().IsLocal() {
		return true
	}
	if merchantId <= 0 {
		return false
	}
	license := go_redismq.Invoke(ctx, &go_redismq.InvoiceRequest{
		Group:   "GID_UniBee_License",
		Method:  "GetLicenseByMerchantId",
		Request: merchantId,
	}, 0)
	if license == nil || !license.Status {
		g.Log().Errorf(ctx, "Get IsPremiumVersion error ,license:%s", utility.MarshalToJsonString(license))
		return false
	}
	var one *License
	err := utility.UnmarshalFromJsonString(utility.MarshalToJsonString(license.Response), &one)
	if err != nil {
		return false
	}
	if one == nil {
		return false
	}
	if one.Version == nil || !one.Version.IsPaid || one.Version.Expired {
		return false
	}
	return true
}

func PremiumLicenseHandler(r *ghttp.Request) {
	if config.GetConfigInstance().IsLocal() {
		r.Middleware.Next()
		return
	}
	uniBeeContext := context2.Context().Get(r.Context())
	if uniBeeContext == nil || uniBeeContext.MerchantId <= 0 {
		//merchant not found
		_interface.OpenApiJsonExit(r, 61, "Merchant Not found")
		r.Exit()
		return
	}
	license := go_redismq.Invoke(r.Context(), &go_redismq.InvoiceRequest{
		Group:   "GID_UniBee_License",
		Method:  "GetLicenseByMerchantId",
		Request: uniBeeContext.MerchantId,
	}, 0)
	if license == nil || !license.Status {
		_interface.OpenApiJsonExit(r, 61, fmt.Sprintf("Get License failed:%s", utility.MarshalToJsonString(license)))
		r.Exit()
		return
	}
	var one *License
	err := utility.UnmarshalFromJsonString(utility.MarshalToJsonString(license.Response), &one)
	if err != nil {
		_interface.OpenApiJsonExit(r, 61, fmt.Sprintf("Get License UnmarshalFromJsonString failed:%s", err.Error()))
		r.Exit()
		return
	}
	if one == nil {
		_interface.OpenApiJsonExit(r, 61, fmt.Sprintf("Get License UnmarshalFromJsonString failed"))
		r.Exit()
		return
	}
	if one.Version == nil || !one.Version.IsPaid || one.Version.Expired {
		_interface.OpenApiJsonExit(r, 61, fmt.Sprintf("Feature analytics need premium license, please contact out support team"))
		r.Exit()
		return
	}
	r.Middleware.Next()
}

type Plan struct {
	// amount, cent, without tax
	Amount *int64 `json:"amount,omitempty"`
	// binded recurring addon planIds，split with ,
	BindingAddonIds *string `json:"bindingAddonIds,omitempty"`
	// binded onetime addon planIds，split with ,
	BindingOnetimeAddonIds *string `json:"bindingOnetimeAddonIds,omitempty"`
	// whether cancel at subscripiton first trial end，0-false | 1-true, will pass to cancelAtPeriodEnd of subscription
	CancelAtTrialEnd *int32 `json:"cancelAtTrialEnd,omitempty"`
	// create utc time
	CreateTime *int64 `json:"createTime,omitempty"`
	// currency
	Currency *string `json:"currency,omitempty"`
	// description
	Description *string `json:"description,omitempty"`
	// external_user_id
	ExternalPlanId  *string `json:"externalPlanId,omitempty"`
	ExtraMetricData *string `json:"extraMetricData,omitempty"`
	// who pay the gas, merchant|user
	GasPayer *string `json:"gasPayer,omitempty"`
	// home_url
	HomeUrl *string `json:"homeUrl,omitempty"`
	Id      *int64  `json:"id,omitempty"`
	// image_url
	ImageUrl *string `json:"imageUrl,omitempty"`
	// period unit count
	IntervalCount *int32 `json:"intervalCount,omitempty"`
	// period unit,day|month|year|week
	IntervalUnit *string `json:"intervalUnit,omitempty"`
	// merchant id
	MerchantId *int64                  `json:"merchantId,omitempty"`
	Metadata   *map[string]interface{} `json:"metadata,omitempty"`
	// PlanName
	PlanName *string `json:"planName,omitempty"`
	// product id
	ProductId *int64 `json:"productId,omitempty"`
	// 1-UnPublish,2-Publish, Use For Display Plan At UserPortal
	PublishStatus *int32 `json:"publishStatus,omitempty"`
	// status，1-editing，2-active，3-inactive，4-expired
	Status *int32 `json:"status,omitempty"`
	// TaxPercentage 1000 = 10%
	TaxPercentage *int32 `json:"taxPercentage,omitempty"`
	// price of trial period
	TrialAmount *int64  `json:"trialAmount,omitempty"`
	TrialDemand *string `json:"trialDemand,omitempty"`
	// duration of trial
	TrialDurationTime *int64 `json:"trialDurationTime,omitempty"`
	// type，1-main plan，2-addon plan
	Type *int32 `json:"type,omitempty"`
}

type MerchantVersion struct {
	Name      string `json:"name" dc:"Name"`
	IsPaid    bool   `json:"isPaid" dc:"IsPaid"`
	Expired   bool   `json:"expired" dc:"Expired"`
	Plan      *Plan  `json:"plan" dc:"Plan"`
	StartTime int64  `json:"startTime" dc:"StartTime,UTC, The Start Time Of Plan,0 for free"`
	EndTime   int64  `json:"endTime" dc:"EndTime,UTC, The End Time Of Plan,0 for free"`
}

type License struct {
	IsPaid       bool               `json:"isPaid,omitempty"`
	OwnerEmail   string             `json:"ownerEmail" dc:"OwnerEmail"`
	Version      *MerchantVersion   `json:"version" dc:"Version Info"`
	License      string             `json:"license" dc:"License, Premium Version will contain License"`
	Plan         *bean.Plan         `json:"plan,omitempty"`
	Subscription *bean.Subscription `json:"subscription,omitempty"`
}
