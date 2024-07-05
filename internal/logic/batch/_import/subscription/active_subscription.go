package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
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
		Gateway:                fmt.Sprintf("%s", row["Gateway(stripe|paypal|wire_transfer|changelly)"]),
		CurrentPeriodStart:     fmt.Sprintf("%s", row["CurrentPeriodStart"]),
		CurrentPeriodEnd:       fmt.Sprintf("%s", row["CurrentPeriodEnd"]),
		BillingCycleAnchor:     fmt.Sprintf("%s", row["BillingCycleAnchor"]),
		FirstPaidTime:          fmt.Sprintf("%s", row["FirstPaidTime"]),
		CreateTime:             fmt.Sprintf("%s", row["CreateTime"]),
		StripeUserId:           fmt.Sprintf("%s", row["StripeUserId(Auto-Charge Required)"]),
		StripePaymentMethod:    fmt.Sprintf("%s", row["StripePaymentMethod(Auto-Charge Required)"]),
		PaypalVaultId:          fmt.Sprintf("%s", row["PaypalVaultId(Auto-Charge Required)"]),
	}
	if len(target.ExternalUserId) == 0 {
		return target, gerror.New("Error, ExternalUserId is blank")
	}
	if len(target.ExternalPlanId) == 0 {
		return target, gerror.New("Error, ExternalPlanId is blank")
	}
	if len(target.Amount) == 0 {
		return target, gerror.New("Error, Amount is blank")
	}
	if len(target.Currency) == 0 {
		return target, gerror.New("Error, Currency is blank")
	}
	if len(target.Gateway) == 0 {
		return target, gerror.New("Error, Gateway is blank")
	}
	if len(target.CurrentPeriodStart) == 0 {
		return target, gerror.New("Error, CurrentPeriodStart is blank")
	}
	if len(target.CurrentPeriodEnd) == 0 {
		return target, gerror.New("Error, CurrentPeriodEnd is blank")
	}
	if len(target.BillingCycleAnchor) == 0 {
		return target, gerror.New("Error, BillingCycleAnchor is blank")
	}
	if len(target.FirstPaidTime) == 0 {
		return target, gerror.New("Error, FirstPaidTime is blank")
	}
	if len(target.CreateTime) == 0 {
		return target, gerror.New("Error, CreateTime is blank")
	}
	plan := query.GetPlanByExternalPlanId(ctx, task.MerchantId, target.ExternalPlanId)
	if plan == nil {
		return target, gerror.New("Error, can't find plan by ExternalPlanId")
	}
	user := query.GetUserAccountByExternalUserId(ctx, task.MerchantId, target.ExternalUserId)
	if user == nil {
		return target, gerror.New("Error, can't find user by ExternalUserId")
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
	Gateway                string `json:"Gateway(stripe|paypal|wire_transfer|changelly)"            `
	CurrentPeriodStart     string `json:"CurrentPeriodStart" `
	CurrentPeriodEnd       string `json:"CurrentPeriodEnd"   `
	BillingCycleAnchor     string `json:"BillingCycleAnchor" `
	FirstPaidTime          string `json:"FirstPaidTime"      `
	CreateTime             string `json:"CreateTime"         `
	StripeUserId           string `json:"StripeUserId(Auto-Charge Required)"             `
	StripePaymentMethod    string `json:"StripePaymentMethod(Auto-Charge Required)"      `
	PaypalVaultId          string `json:"PaypalVaultId(Auto-Charge Required)"      `
	Features               string `json:"Features"         `
}
