package system

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/system/refund"
)

func (c *ControllerRefund) GatewayDetail(ctx context.Context, req *refund.GatewayDetailReq) (res *refund.GatewayDetailRes, err error) {
	one := query.GetRefundByRefundId(ctx, req.RefundId)
	utility.Assert(one != nil, "refund not found")
	gateway := query.GetGatewayById(ctx, one.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	detail, err := api.GetGatewayServiceProvider(ctx, one.GatewayId).GatewayRefundDetail(ctx, gateway, one.GatewayRefundId, one)
	if err != nil {
		return nil, err
	}
	return &refund.GatewayDetailRes{RefundDetail: gjson.New(detail)}, nil
}
