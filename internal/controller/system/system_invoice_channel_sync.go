package system

//
//func (c *ControllerInvoice) ChannelSync(ctx context.Context, req *invoice.ChannelSyncReq) (res *invoice.ChannelSyncRes, err error) {
//	utility.Assert(len(req.MerchantId) > 0, "merchantId invalid")
//	utility.Assert(len(req.InvoiceId) > 0, "invoiceId invalid")
//	go func() {
//		defer func() {
//			if exception := recover(); exception != nil {
//				if v, ok := exception.(error); ok && gerror.HasStack(v) {
//					err = v
//				} else {
//					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
//				}
//				g.Log().Errorf(context.Background(), "BulkChannelSync Background panic error:%s\n", err.Error())
//				return
//			}
//		}()
//		var page = 0
//		var count = 100
//		for {
//			backgroundCtx := context.Background()
//			var mainList []*entity.Invoice
//			err := dao.Invoice.Ctx(backgroundCtx).
//				Where(dao.Invoice.Columns().MerchantId, req.MerchantId).
//				Where(dao.Invoice.Columns().InvoiceId, req.InvoiceId).
//				WhereNotNull(dao.Invoice.Columns().ChannelInvoiceId).
//				OrderDesc("id").
//				Limit(page*count, count).
//				OmitEmpty().Scan(&mainList)
//			if err != nil {
//				fmt.Printf("ChannelSync Background List error%s\n", err.Error())
//				return
//			}
//			for _, one := range mainList {
//				payChannel := query.GetPayChannelById(backgroundCtx, one.ChannelId)
//				utility.Assert(payChannel != nil, "invalid planChannel")
//				details, err := channel.GetPayChannelServiceProvider(backgroundCtx, one.ChannelId).DoRemoteChannelInvoiceDetails(backgroundCtx, payChannel, one.ChannelInvoiceId)
//				if err == nil {
//					err := handler.CreateOrUpdateInvoiceByChannelDetail(backgroundCtx, details, details.ChannelInvoiceId)
//					if err != nil {
//						fmt.Printf("ChannelSync Background CreateOrUpdateInvoiceByChannelDetail InvoiceId:%s error%s\n", one.InvoiceId, err.Error())
//						return
//					}
//					fmt.Printf("ChannelSync Background Fetch InvoiceId:%s success\n", one.InvoiceId)
//				} else {
//					fmt.Printf("ChannelSync Background Fetch InvoiceId:%s error%s\n", one.InvoiceId, err.Error())
//				}
//			}
//			if len(mainList) == 0 {
//				break
//			}
//			clear(mainList)
//			page = page + 1
//		}
//
//	}()
//	return nil, nil
//}
