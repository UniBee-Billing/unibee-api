// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package merchant

import (
	"go-oversea-pay/api/merchant"
)

type ControllerPlan struct{}

func NewPlan() merchant.IMerchantPlan {
	return &ControllerPlan{}
}

type ControllerWebhook struct{}

func NewWebhook() merchant.IMerchantWebhook {
	return &ControllerWebhook{}
}

