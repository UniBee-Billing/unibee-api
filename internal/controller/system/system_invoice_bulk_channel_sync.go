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
	//			WhereNotNull(dao.Invoice.Columns().ChannelInvoiceId).
	//			OrderDesc("id").
	//			Limit(page*count, count).
	//			OmitEmpty().Scan(&mainList)
	//		if err != nil {
	//			fmt.Printf("BulkChannelSync Background List error%s\n", err.Error())
	//			return
	//		}
	//		for _, one := range mainList {
	//			payChannel := query.GetPayChannelById(backgroundCtx, one.ChannelId)
	//			utility.Assert(payChannel != nil, "invalid planChannel")
	//			details, err := channel.GetPayChannelServiceProvider(backgroundCtx, one.ChannelId).DoRemoteChannelInvoiceDetails(backgroundCtx, payChannel, one.ChannelInvoiceId)
	//			if err == nil {
	//				err := handler.CreateOrUpdateInvoiceByChannelDetail(backgroundCtx, details, details.ChannelInvoiceId)
	//				if err != nil {
	//					fmt.Printf("BulkChannelSync Background CreateOrUpdateInvoiceByChannelDetail ChannelInvoiceId:%s error%s\n", one.ChannelInvoiceId, err.Error())
	//					return
	//				}
	//				fmt.Printf("BulkChannelSync Background Fetch ChannelInvoiceId:%s success\n", one.ChannelInvoiceId)
	//			} else {
	//				fmt.Printf("BulkChannelSync Background Fetch ChannelInvoiceId:%s error%s\n", one.ChannelInvoiceId, err.Error())
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
