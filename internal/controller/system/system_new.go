// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package system

import (
	"go-oversea-pay/api/system"
)

type ControllerSubscription struct{}

func NewSubscription() system.ISystemSubscription {
	return &ControllerSubscription{}
}

type ControllerInvoice struct{}

func NewInvoice() system.ISystemInvoice {
	return &ControllerInvoice{}
}

