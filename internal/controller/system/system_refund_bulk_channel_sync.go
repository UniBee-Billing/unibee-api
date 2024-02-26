package system

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/gateway/api"
	handler2 "unibee/internal/logic/payment/handler"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/system/refund"
)

func (c *ControllerRefund) BulkChannelSync(ctx context.Context, req *refund.BulkChannelSyncReq) (res *refund.BulkChannelSyncRes, err error) {
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
			var mainList []*entity.Payment
			err := dao.Payment.Ctx(backgroundCtx).
				Where(dao.Payment.Columns().MerchantId, req.MerchantId).
				WhereNotNull(dao.Payment.Columns().GatewayPaymentId).
				OrderDesc("id").
				Limit(page*count, count).
				OmitEmpty().Scan(&mainList)
			if err != nil {
				fmt.Printf("BulkChannelSync Background List error%s\n", err.Error())
				return
			}
			for _, one := range mainList {
				gateway := query.GetGatewayById(backgroundCtx, one.GatewayId)
				utility.Assert(gateway != nil, "invalid gatewayPlan")
				details, err := api.GetGatewayServiceProvider(backgroundCtx, one.GatewayId).GatewayRefundList(backgroundCtx, gateway, one.GatewayPaymentId)
				if err == nil {
					for _, detail := range details {
						err := handler2.CreateOrUpdateRefundByDetail(backgroundCtx, one, detail, detail.GatewayRefundId)
						if err != nil {
							fmt.Printf("BulkChannelSync Background CreateOrUpdateRefundByDetail GatewayRefundId:%s error%s\n", detail.GatewayRefundId, err.Error())
							return
						}
						fmt.Printf("BulkChannelSync Background Fetch GatewayRefundId:%s success\n", detail.GatewayRefundId)
					}
				} else {
					fmt.Printf("BulkChannelSync Background Fetch GatewayPaymentId:%s error%s\n", one.GatewayPaymentId, err.Error())
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
