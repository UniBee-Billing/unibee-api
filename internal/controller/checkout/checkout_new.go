// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package checkout

import (
	"unibee/api/checkout"
)

type ControllerGateway struct{}

func NewGateway() checkout.ICheckoutGateway {
	return &ControllerGateway{}
}

type ControllerIp struct{}

func NewIp() checkout.ICheckoutIp {
	return &ControllerIp{}
}

type ControllerSubscription struct{}

func NewSubscription() checkout.ICheckoutSubscription {
	return &ControllerSubscription{}
}

type ControllerVat struct{}

func NewVat() checkout.ICheckoutVat {
	return &ControllerVat{}
}
