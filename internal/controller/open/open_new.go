// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package open

import (
	"go-oversea-pay/api/open"
)

type ControllerMock struct{}

func NewMock() open.IOpenMock {
	return &ControllerMock{}
}

type ControllerPayment struct{}

func NewPayment() open.IOpenPayment {
	return &ControllerPayment{}
}

