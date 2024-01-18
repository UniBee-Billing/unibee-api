package system

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/gateway"
	"go-oversea-pay/internal/logic/subscription/handler"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/system/invoice"
)

func (c *ControllerInvoice) BulkChannelSync(ctx context.Context, req *invoice.BulkChannelSyncReq) (res *invoice.BulkChannelSyncRes, err error) {
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
			err := dao.Invoice.Ctx(backgroundCtx).
				Where(dao.Invoice.Columns().MerchantId, req.MerchantId).
				WhereNotNull(dao.Invoice.Columns().ChannelInvoiceId).
				OrderDesc("id").
				Limit(page*count, count).
				OmitEmpty().Scan(&mainList)
			if err != nil {
				fmt.Printf("BulkChannelSync Background List error%s\n", err.Error())
				return
			}
			for _, one := range mainList {
				payChannel := query.GetPayChannelById(backgroundCtx, one.ChannelId)
				utility.Assert(payChannel != nil, "invalid planChannel")
				details, err := gateway.GetPayChannelServiceProvider(backgroundCtx, one.ChannelId).DoRemoteChannelInvoiceDetails(backgroundCtx, payChannel, one.ChannelInvoiceId)
				if err == nil {
					err := handler.CreateOrUpdateInvoiceByDetail(backgroundCtx, details)
					if err != nil {
						fmt.Printf("BulkChannelSync Background CreateOrUpdateInvoiceByDetail InvoiceId:%s error%s\n", one.InvoiceId, err.Error())
						return
					}
					fmt.Printf("BulkChannelSync Background Fetch InvoiceId:%s success\n", one.InvoiceId)
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
