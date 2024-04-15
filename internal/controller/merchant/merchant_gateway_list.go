package merchant

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee/api/bean"
	_interface "unibee/internal/interface"
	"unibee/internal/query"

	"unibee/api/merchant/gateway"
)

func UnmarshalFromJsonString(target string, one interface{}) error {
	if len(target) > 0 {
		return gjson.Unmarshal([]byte(target), &one)
	} else {
		return gerror.New("target is nil")
	}
}

func (c *ControllerGateway) List(ctx context.Context, req *gateway.ListReq) (res *gateway.ListRes, err error) {
	data := query.GetMerchantGatewayList(ctx, _interface.GetMerchantId(ctx))
	var list = make([]*bean.GatewaySimplify, 0)
	for _, one := range data {
		var countryConfig map[string]bool
		_ = UnmarshalFromJsonString(one.CountryConfig, &countryConfig)
		list = append(list, &bean.GatewaySimplify{
			Id:            one.Id,
			GatewayLogo:   one.Logo,
			GatewayName:   one.GatewayName,
			GatewayKey:    one.GatewayKey,
			GatewayType:   one.GatewayType,
			CountryConfig: countryConfig,
			CreateTime:    one.CreateTime,
		})
	}
	return &gateway.ListRes{
		Gateways: list,
	}, nil
}
