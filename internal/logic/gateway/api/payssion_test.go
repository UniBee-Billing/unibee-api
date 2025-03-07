package api

import (
	"context"
	"fmt"
	"testing"
	_interface "unibee/internal/interface"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

func TestForPayssion(t *testing.T) {
	pay := &Payssion{}
	_, _, _ = pay.GatewayTest(context.Background(), &_interface.GatewayTestReq{
		Key:                 "sandbox_6340c0569ae5339c",
		Secret:              "hdvh5MkJMCQ5ZhtgatLzukbJXwbRMra4",
		SubGateway:          "",
		GatewayPaymentTypes: nil,
	})
}

func TestForPayssionGetPaymentDetail(t *testing.T) {
	pay := &Payssion{}
	res, err := pay.GatewayPaymentDetail(context.Background(), &entity.MerchantGateway{
		GatewayKey:    "sandbox_6340c0569ae5339c",
		GatewaySecret: "hdvh5MkJMCQ5ZhtgatLzukbJXwbRMra4",
	}, "T120926320315372", &entity.Payment{
		PaymentId: "pay20250120lYdJKhnvfoBi7W9",
	})
	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("%s", utility.MarshalToJsonString(res))
	}
}

func TestForGetSubGatewayData(t *testing.T) {
	list := fetchPayssionPaymentTypes(context.Background())
	for _, i := range list {
		fmt.Println(utility.MarshalToJsonString(i))
	}
}
