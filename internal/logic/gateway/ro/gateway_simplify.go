package ro

import entity "unibee/internal/model/entity/oversea_pay"

type GatewaySimplify struct {
	Id          uint64 `json:"gatewayId"`
	GatewayName string `json:"gatewayName"`
}

func SimplifyGateway(one *entity.MerchantGateway) *GatewaySimplify {
	if one == nil {
		return nil
	}
	return &GatewaySimplify{
		Id:          one.Id,
		GatewayName: one.GatewayName,
	}
}
