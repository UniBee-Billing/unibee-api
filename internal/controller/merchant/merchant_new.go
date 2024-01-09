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


type ControllerAuth struct{}

func NewAuth() merchant.IMerchantAuth {
	return &ControllerAuth{}
}

type ControllerProfile struct{}

func NewProfile() merchant.IMerchantProfile {
	return &ControllerProfile{}
}
type ControllerSubscription struct{}

func NewSubscription() merchant.IMerchantSubscription {
	return &ControllerSubscription{}
}

type ControllerOss struct{}

func NewOss() merchant.IMerchantOss {
	return &ControllerOss{}
}

