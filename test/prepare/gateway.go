package prepare

import (
	"context"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func CreateTestGateway(ctx context.Context, merchantId uint64) *entity.MerchantGateway {
	one := query.GetGatewayByGatewayName(ctx, merchantId, "autotest")
	if one != nil {
		return one
	}
	//service.SetupGateway(ctx, merchantId, "autotest", "", "")
	//one = query.GetGatewayByGatewayName(ctx, merchantId, "autotest")
	utility.Assert(one != nil, "autotest gateway error")
	return one
}
