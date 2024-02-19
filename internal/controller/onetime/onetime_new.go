// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package onetime

import (
	"unibee-api/api/onetime"
)

type ControllerMock struct{}

func NewMock() onetime.IOnetimeMock {
	return &ControllerMock{}
}

type ControllerPayment struct{}

func NewPayment() onetime.IOnetimePayment {
	return &ControllerPayment{}
}
