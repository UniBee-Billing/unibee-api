// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package system

import (
	"unibee/api/system"
)

type ControllerSubscription struct{}

func NewSubscription() system.ISystemSubscription {
	return &ControllerSubscription{}
}

type ControllerInvoice struct{}

func NewInvoice() system.ISystemInvoice {
	return &ControllerInvoice{}
}

type ControllerPayment struct{}

func NewPayment() system.ISystemPayment {
	return &ControllerPayment{}
}

type ControllerRefund struct{}

func NewRefund() system.ISystemRefund {
	return &ControllerRefund{}
}

type ControllerInformation struct{}

func NewInformation() system.ISystemInformation {
	return &ControllerInformation{}
}

type ControllerPlan struct{}

func NewPlan() system.ISystemPlan {
	return &ControllerPlan{}
}

type ControllerAuth struct{}

func NewAuth() system.ISystemAuth {
	return &ControllerAuth{}
}

type ControllerUser struct{}

func NewUser() system.ISystemUser {
	return &ControllerUser{}
}
