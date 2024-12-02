package system

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/system/payment"
)

func (c *ControllerPayment) PaymentGatewayDetail(ctx context.Context, req *payment.PaymentGatewayDetailReq) (res *payment.PaymentGatewayDetailRes, err error) {
	one := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	utility.Assert(one != nil, "payment not found")
	invoice := query.GetInvoiceByInvoiceId(ctx, one.InvoiceId)
	utility.Assert(invoice != nil, "invoice not found")
	gateway := query.GetGatewayById(ctx, one.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	detail, err := api.GetGatewayServiceProvider(ctx, one.GatewayId).GatewayPaymentDetail(ctx, gateway, one.GatewayPaymentId, one)
	if err != nil {
		return nil, err
	}
	return &payment.PaymentGatewayDetailRes{PaymentDetail: gjson.New(detail)}, nil
}
