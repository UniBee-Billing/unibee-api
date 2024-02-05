package system

import (
	"context"
	"go-oversea-pay/api/system/invoice"
)

func (c *ControllerInvoice) BulkChannelSync(ctx context.Context, req *invoice.BulkChannelSyncReq) (res *invoice.BulkChannelSyncRes, err error) {
	//utility.Assert(len(req.MerchantId) > 0, "merchantId invalid")
	//go func() {
	//	defer func() {
	//		if exception := recover(); exception != nil {
	//			if v, ok := exception.(error); ok && gerror.HasStack(v) {
	//				err = v
	//			} else {
	//				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
	//			}
	//			g.Log().Errorf(context.Background(), "BulkChannelSync Background panic error:%s\n", err.Error())
	//			return
	//		}
	//	}()
	//	var page = 0
	//	var count = 10
	//	for {t sta
	//		backgroundCtx := context.Background()
	//		var mainList []*entity.Invoice
	//		err := dao.Invoice.Ctx(backgroundCtx).
	//			Where(dao.Invoice.Columns().MerchantId, req.MerchantId).
	//			WhereNotNull(dao.Invoice.Columns().GatewayInvoiceId).
	//			OrderDesc("id").
	//			Limit(page*count, count).
	//			OmitEmpty().Scan(&mainList)
	//		if err != nil {
	//			fmt.Printf("BulkChannelSync Background List error%s\n", err.Error())
	//			return
	//		}
	//		for _, one := range mainList {
	//			gateway := query.GetPayChannelById(backgroundCtx, one.GatewayId)
	//			utility.Assert(gateway != nil, "invalid planChannel")
	//			details, err := channel.GetPayChannelServiceProvider(backgroundCtx, one.GatewayId).GatewayInvoiceDetails(backgroundCtx, gateway, one.GatewayInvoiceId)
	//			if err == nil {
	//				err := handler.CreateOrUpdateInvoiceByChannelDetail(backgroundCtx, details, details.GatewayInvoiceId)
	//				if err != nil {
	//					fmt.Printf("BulkChannelSync Background CreateOrUpdateInvoiceByChannelDetail GatewayInvoiceId:%s error%s\n", one.GatewayInvoiceId, err.Error())
	//					return
	//				}
	//				fmt.Printf("BulkChannelSync Background Fetch GatewayInvoiceId:%s success\n", one.GatewayInvoiceId)
	//			} else {
	//				fmt.Printf("BulkChannelSync Background Fetch GatewayInvoiceId:%s error%s\n", one.GatewayInvoiceId, err.Error())
	//			}
	//		}
	//		if len(mainList) == 0 {
	//			break
	//		}
	//		clear(mainList)
	//		page = page + 1
	//	}
	//}()
	return nil, nil
}
