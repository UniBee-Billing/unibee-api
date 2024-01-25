// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package user

import (
	"go-oversea-pay/api/user"
)

type ControllerAuth struct{}

func NewAuth() user.IUserAuth {
	return &ControllerAuth{}
}

type ControllerProfile struct{}

func NewProfile() user.IUserProfile {
	return &ControllerProfile{}
}

type ControllerSubscription struct{}

func NewSubscription() user.IUserSubscription {
	return &ControllerSubscription{}
}

type ControllerPlan struct{}

func NewPlan() user.IUserPlan {
	return &ControllerPlan{}
}

type ControllerVat struct{}

func NewVat() user.IUserVat {
	return &ControllerVat{}
}

type ControllerInvoice struct{}

func NewInvoice() user.IUserInvoice {
	return &ControllerInvoice{}
}

type ControllerPayment struct{}

func NewPayment() user.IUserPayment {
	return &ControllerPayment{}
}

type ControllerUser struct{}

func NewUser() user.IUserUser {
	return &ControllerUser{}
}

