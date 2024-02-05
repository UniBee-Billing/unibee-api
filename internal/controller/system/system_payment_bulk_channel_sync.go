package system

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/gateway/api"
	handler2 "go-oversea-pay/internal/logic/payment/handler"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/system/payment"
)

func (c *ControllerPayment) BulkChannelSync(ctx context.Context, req *payment.BulkChannelSyncReq) (res *payment.BulkChannelSyncRes, err error) {
	utility.Assert(len(req.MerchantId) > 0, "merchantId invalid")
	go func() {
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				g.Log().Errorf(context.Background(), "BulkChannelSync Background panic error:%s\n", err.Error())
				return
			}
		}()
		var page = 0
		var count = 10
		for {
			backgroundCtx := context.Background()
			var mainList []*entity.Invoice
			err = dao.Invoice.Ctx(backgroundCtx).
				Where(dao.Invoice.Columns().MerchantId, req.MerchantId).
				WhereNotNull(dao.Invoice.Columns().GatewayPaymentId).
				OrderDesc("id").
				Limit(page*count, count).
				OmitEmpty().Scan(&mainList)
			if err != nil {
				fmt.Printf("BulkChannelSync Background List error%s\n", err.Error())
				return
			}
			for _, one := range mainList {
				gateway := query.GetGatewayById(backgroundCtx, one.GatewayId)
				utility.Assert(gateway != nil, "invalid planChannel")
				details, err := api.GetGatewayServiceProvider(backgroundCtx, one.GatewayId).GatewayPaymentDetail(backgroundCtx, gateway, one.GatewayPaymentId)
				details.UniqueId = details.GatewayPaymentId
				if err == nil {
					pay, err := handler2.CreateOrUpdateSubscriptionPaymentFromChannel(backgroundCtx, details)
					if err != nil {
						fmt.Printf("BulkChannelSync Background CreateOrUpdateSubscriptionPaymentFromChannel GatewayPaymentIntentId:%s error%s\n", details.GatewayPaymentId, err.Error())
						return
					}
					_, _ = dao.Invoice.Ctx(backgroundCtx).Data(g.Map{
						dao.Invoice.Columns().PaymentId: pay.PaymentId,
					}).Where(dao.Invoice.Columns().InvoiceId, one.InvoiceId).OmitNil().Update()
					fmt.Printf("BulkChannelSync Background Fetch GatewayPaymentIntentId:%s success\n", details.GatewayPaymentId)
				} else {
					fmt.Printf("BulkChannelSync Background Fetch InvoiceId:%s error%s\n", one.InvoiceId, err.Error())
				}
			}
			if len(mainList) == 0 {
				break
			}
			clear(mainList)
			page = page + 1
		}
	}()
	return nil, nil
}
