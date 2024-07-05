package subscription

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	entity "unibee/internal/model/entity/oversea_pay"
)

type TaskActiveSubscriptionImport struct {
}

func (t TaskActiveSubscriptionImport) TaskName() string {
	return "ActiveSubscriptionImport"
}

func (t TaskActiveSubscriptionImport) TemplateHeader() interface{} {
	return ImportActiveSubscriptionEntity{}
}

func (t TaskActiveSubscriptionImport) ImportRow(ctx context.Context, task *entity.MerchantBatchTask, row map[string]string) (interface{}, error) {
	var err error
	target := &ImportActiveSubscriptionEntity{}
	return target, err
}

type ImportActiveSubscriptionEntity struct {
	ExternalSubscriptionId string      `json:"ExternalSubscriptionId"     `
	ExternalUserId         string      `json:"ExternalUserId"     `
	ExternalPlanId         string      `json:"ExternalPlanId"     `
	Amount                 string      `json:"Amount"             `
	Currency               string      `json:"Currency"           `
	Quantity               string      `json:"Quantity"           `
	Gateway                string      `json:"Gateway"            `
	Status                 string      `json:"Status"             `
	CurrentPeriodStart     *gtime.Time `json:"CurrentPeriodStart" layout:"2006-01-02 15:04:05"`
	CurrentPeriodEnd       *gtime.Time `json:"CurrentPeriodEnd"  layout:"2006-01-02 15:04:05" `
	BillingCycleAnchor     *gtime.Time `json:"BillingCycleAnchor" layout:"2006-01-02 15:04:05"`
	DunningTime            *gtime.Time `json:"DunningTime"        layout:"2006-01-02 15:04:05"`
	TrialEnd               *gtime.Time `json:"TrialEnd"          layout:"2006-01-02 15:04:05" `
	FirstPaidTime          *gtime.Time `json:"FirstPaidTime"    layout:"2006-01-02 15:04:05"  `
	CreateTime             *gtime.Time `json:"CreateTime"      layout:"2006-01-02 15:04:05"   `
	StripeUserId           string      `json:"StripeUserId(Auto-Charge Required)"             `
	StripePaymentMethod    string      `json:"StripePaymentMethod(Auto-Charge Required)"      `
	PaypalVaultId          string      `json:"PaypalVaultId(Auto-Charge Required)"      `
}
