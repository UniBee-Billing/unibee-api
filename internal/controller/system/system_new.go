// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package system

import (
	"unibee-api/api/system"
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

