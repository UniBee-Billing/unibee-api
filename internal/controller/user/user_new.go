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

