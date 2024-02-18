// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package onetime

import (
	"unibee-api/api/onetime"
)

type ControllerMock struct{}

func NewMock() onetime.IOpenMock {
	return &ControllerMock{}
}

type ControllerPayment struct{}

func NewPayment() onetime.IOpenPayment {
	return &ControllerPayment{}
}
