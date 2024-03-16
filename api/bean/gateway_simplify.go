package bean

import entity "unibee/internal/model/entity/oversea_pay"

type GatewaySimplify struct {
	Id            uint64          `json:"gatewayId"`
	GatewayName   string          `json:"gatewayName"`
	GatewayLogo   string          `json:"gatewayLogo"`
	GatewayType   int64           `json:"gatewayType"           description:"gateway type，1-Default｜ 2-Crypto"` // gateway type，1-Default｜ 2-Crypto
	CountryConfig map[string]bool `json:"countryConfig"`
}

func SimplifyGateway(one *entity.MerchantGateway) *GatewaySimplify {
	if one == nil {
		return nil
	}
	var countryConfig map[string]bool
	_ = UnmarshalFromJsonString(one.CountryConfig, &countryConfig)
	return &GatewaySimplify{
		Id:            one.Id,
		GatewayLogo:   one.Logo,
		GatewayName:   one.GatewayName,
		GatewayType:   one.GatewayType,
		CountryConfig: countryConfig,
	}
}

func SimplifyGatewayList(ones []*entity.MerchantGateway) (list []*GatewaySimplify) {
	if len(ones) == 0 {
		return make([]*GatewaySimplify, 0)
	}
	for _, one := range ones {
		list = append(list, SimplifyGateway(one))
	}
	return list
}
