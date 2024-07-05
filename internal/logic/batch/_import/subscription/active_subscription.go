package subscription

import (
	"context"
	"fmt"
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
	target := &ImportActiveSubscriptionEntity{
		ExternalSubscriptionId: fmt.Sprintf("%s", row["ExternalSubscriptionId"]),
		ExternalUserId:         fmt.Sprintf("%s", row["ExternalUserId"]),
		ExternalPlanId:         fmt.Sprintf("%s", row["ExternalPlanId"]),
		Amount:                 fmt.Sprintf("%s", row["Amount"]),
		Currency:               fmt.Sprintf("%s", row["Currency"]),
		Quantity:               fmt.Sprintf("%s", row["Quantity"]),
		Gateway:                fmt.Sprintf("%s", row["Gateway"]),
		Status:                 fmt.Sprintf("%s", row["Status"]),
		CurrentPeriodStart:     fmt.Sprintf("%s", row["CurrentPeriodStart"]),
		CurrentPeriodEnd:       fmt.Sprintf("%s", row["CurrentPeriodEnd"]),
		BillingCycleAnchor:     fmt.Sprintf("%s", row["BillingCycleAnchor"]),
		DunningTime:            fmt.Sprintf("%s", row["DunningTime"]),
		TrialEnd:               fmt.Sprintf("%s", row["TrialEnd"]),
		FirstPaidTime:          fmt.Sprintf("%s", row["FirstPaidTime"]),
		CreateTime:             fmt.Sprintf("%s", row["CreateTime"]),
		StripeUserId:           fmt.Sprintf("%s", row["StripeUserId(Auto-Charge Required)"]),
		StripePaymentMethod:    fmt.Sprintf("%s", row["StripePaymentMethod(Auto-Charge Required)"]),
		PaypalVaultId:          fmt.Sprintf("%s", row["PaypalVaultId(Auto-Charge Required)"]),
	}
	return target, err
}

type ImportActiveSubscriptionEntity struct {
	ExternalSubscriptionId string `json:"ExternalSubscriptionId"     `
	ExternalUserId         string `json:"ExternalUserId"     `
	ExternalPlanId         string `json:"ExternalPlanId"     `
	Amount                 string `json:"Amount"             `
	Currency               string `json:"Currency"           `
	Quantity               string `json:"Quantity"           `
	Gateway                string `json:"Gateway"            `
	Status                 string `json:"Status"             `
	CurrentPeriodStart     string `json:"CurrentPeriodStart" `
	CurrentPeriodEnd       string `json:"CurrentPeriodEnd"   `
	BillingCycleAnchor     string `json:"BillingCycleAnchor" `
	DunningTime            string `json:"DunningTime"        `
	TrialEnd               string `json:"TrialEnd"           `
	FirstPaidTime          string `json:"FirstPaidTime"      `
	CreateTime             string `json:"CreateTime"         `
	StripeUserId           string `json:"StripeUserId(Auto-Charge Required)"             `
	StripePaymentMethod    string `json:"StripePaymentMethod(Auto-Charge Required)"      `
	PaypalVaultId          string `json:"PaypalVaultId(Auto-Charge Required)"      `
}
