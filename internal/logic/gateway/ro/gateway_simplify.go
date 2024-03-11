package ro

import entity "unibee/internal/model/entity/oversea_pay"

type GatewaySimplify struct {
	Id          uint64 `json:"gatewayId"`
	GatewayName string `json:"gatewayName"`
	GatewayLogo string `json:"gatewayLogo"`
	GatewayType int64  `json:"gatewayType"           description:"gateway type，1-Default｜ 2-Crypto"` // gateway type，1-Default｜ 2-Crypto
}

func SimplifyGateway(one *entity.MerchantGateway) *GatewaySimplify {
	if one == nil {
		return nil
	}
	return &GatewaySimplify{
		Id:          one.Id,
		GatewayLogo: one.Logo,
		GatewayName: one.GatewayName,
		GatewayType: one.GatewayType,
	}
}
