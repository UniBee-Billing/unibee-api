package api

import (
	"context"
	"testing"
	_interface "unibee/internal/interface"
)

func TestForCreateNewFireKassaPayment(t *testing.T) {
	gateway := &FireKassa{}
	_, _, _ = gateway.GatewayTest(context.Background(), &_interface.GatewayTestReq{
		Key:    "",
		Secret: "",
	}) // indigo staging test key
}
