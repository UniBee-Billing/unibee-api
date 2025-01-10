package api

import (
	"context"
	"testing"
)

func TestForGetUnitPay(t *testing.T) {
	unitpay := &UnitPay{}
	_, _, _ = unitpay.GatewayTest(context.Background(), "423641", "41AD9A1AA16-ECB837EFDC2-6C3D77F2F7")

}
