package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/api/merchant/invoice"
	v1 "go-oversea-pay/api/open/payment"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/channel/out"
	"go-oversea-pay/internal/logic/channel/ro"
	"go-oversea-pay/internal/logic/invoice/invoice_compute"
	"go-oversea-pay/internal/logic/payment/service"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strconv"
	"strings"
)

func checkInvoice(one *ro.InvoiceDetailRo) {
	var totalAmountExcludingTax int64 = 0
	var totalTax int64 = 0
	for _, line := range one.Lines {
		amountExcludingTax := line.UnitAmountExcludingTax * line.Quantity
		tax := int64(float64(amountExcludingTax) * utility.ConvertTaxScaleToInternalFloat(one.TaxScale))
		utility.Assert(line.AmountExcludingTax == amountExcludingTax, "line amountExcludingTax mistake")
		utility.Assert(strings.Compare(line.Currency, one.Currency) == 0, "line currency not match invoice currency")
		utility.Assert(line.Amount == amountExcludingTax+tax, "line amount mistake")
		//utility.Assert(line.TaxScale == one.TaxScale, "line taxScale mistake")
		totalTax = totalTax + tax
		totalAmountExcludingTax = totalAmountExcludingTax + amountExcludingTax
	}
	var totalAmount = totalTax + totalAmountExcludingTax
	utility.Assert(one.TaxAmount == totalTax, "invoice taxAmount mistake")
	utility.Assert(one.TotalAmountExcludingTax == totalAmountExcludingTax, "invoice totalAmountExcludingTax mistake")
	utility.Assert(one.TotalAmount == totalAmount, "line totalAmount mistake")
}

func CreateInvoice(ctx context.Context, req *invoice.NewInvoiceCreateReq) (res *invoice.NewInvoiceCreateRes, err error) {
	user := query.GetUserAccountById(ctx, uint64(req.UserId))
	utility.Assert(user != nil, fmt.Sprintf("send user not found:%d", req.UserId))
	utility.Assert(len(user.Email) > 0, fmt.Sprintf("send user email not found:%d", req.UserId))
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, req.ChannelId)
	utility.Assert(payChannel != nil, "payChannel not found")

	var invoiceItems []*ro.InvoiceItemDetailRo
	var totalAmountExcludingTax int64 = 0
	var totalTax int64 = 0
	for _, line := range req.Lines {
		amountExcludingTax := line.UnitAmountExcludingTax * line.Quantity
		tax := int64(float64(amountExcludingTax) * utility.ConvertTaxScaleToInternalFloat(req.TaxScale))
		invoiceItems = append(invoiceItems, &ro.InvoiceItemDetailRo{
			Currency:               req.Currency,
			TaxScale:               req.TaxScale,
			Tax:                    tax,
			Quantity:               line.Quantity,
			Amount:                 amountExcludingTax + tax,
			AmountExcludingTax:     amountExcludingTax,
			UnitAmountExcludingTax: line.UnitAmountExcludingTax,
			Description:            line.Description,
		})
		totalTax = totalTax + tax
		totalAmountExcludingTax = totalAmountExcludingTax + amountExcludingTax
	}
	var totalAmount = totalTax + totalAmountExcludingTax

	//创建
	invoiceId := utility.CreateInvoiceId()
	one := &entity.Invoice{
		BizType:                        consts.BIZ_TYPE_SUBSCRIPTION,
		MerchantId:                     req.MerchantId,
		InvoiceId:                      invoiceId,
		InvoiceName:                    req.Name,
		UniqueId:                       invoiceId,
		TotalAmount:                    totalAmount,
		TotalAmountExcludingTax:        totalAmountExcludingTax,
		TaxAmount:                      totalTax,
		TaxScale:                       req.TaxScale,
		SubscriptionAmount:             totalAmount,
		SubscriptionAmountExcludingTax: totalAmountExcludingTax,
		Currency:                       strings.ToUpper(req.Currency),
		Lines:                          utility.MarshalToJsonString(invoiceItems),
		ChannelId:                      req.ChannelId,
		Status:                         consts.InvoiceStatusPending,
		SendStatus:                     0,
		SendEmail:                      user.Email,
		UserId:                         req.UserId,
	}

	result, err := dao.Invoice.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		return nil, gerror.Newf(`CreateInvoice record insert failure %s`, err)
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))

	one.Lines = utility.MarshalToJsonString(invoiceItems)

	return &invoice.NewInvoiceCreateRes{Invoice: invoice_compute.ConvertInvoiceToRo(one)}, nil
}

func EditInvoice(ctx context.Context, req *invoice.NewInvoiceEditReq) (res *invoice.NewInvoiceEditRes, err error) {
	one := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	utility.Assert(one != nil, fmt.Sprintf("invoice not found:%s", req.InvoiceId))
	utility.Assert(one.Status == consts.InvoiceStatusPending, "invoice not in pending status")
	utility.Assert(one.IsDeleted == 0, "invoice is deleted")
	if req.ChannelId > 0 {
		payChannel := query.GetSubscriptionTypePayChannelById(ctx, req.ChannelId)
		utility.Assert(payChannel != nil, "payChannel not found")
	} else {
		req.ChannelId = one.ChannelId
	}
	if len(req.Currency) == 0 {
		req.Currency = one.Currency
	}

	var invoiceItems []*ro.InvoiceItemDetailRo
	var totalAmountExcludingTax int64 = 0
	var totalTax int64 = 0
	for _, line := range req.Lines {
		amountExcludingTax := line.UnitAmountExcludingTax * line.Quantity
		tax := int64(float64(amountExcludingTax) * utility.ConvertTaxScaleToInternalFloat(req.TaxScale))
		invoiceItems = append(invoiceItems, &ro.InvoiceItemDetailRo{
			Currency:               req.Currency,
			TaxScale:               req.TaxScale,
			Tax:                    tax,
			Amount:                 amountExcludingTax + tax,
			AmountExcludingTax:     amountExcludingTax,
			UnitAmountExcludingTax: line.UnitAmountExcludingTax,
			Description:            line.Description,
			Quantity:               line.Quantity,
		})
		totalTax = totalTax + tax
		totalAmountExcludingTax = totalAmountExcludingTax + amountExcludingTax
	}
	var totalAmount = totalTax + totalAmountExcludingTax

	//更新 Subscription
	_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
		dao.Invoice.Columns().BizType:                        consts.BIZ_TYPE_SUBSCRIPTION,
		dao.Invoice.Columns().InvoiceName:                    req.Name,
		dao.Invoice.Columns().TotalAmount:                    totalAmount,
		dao.Invoice.Columns().TotalAmountExcludingTax:        totalAmountExcludingTax,
		dao.Invoice.Columns().TaxAmount:                      totalTax,
		dao.Invoice.Columns().SubscriptionAmount:             totalAmount,
		dao.Invoice.Columns().SubscriptionAmountExcludingTax: totalAmountExcludingTax,
		dao.Invoice.Columns().Currency:                       strings.ToUpper(req.Currency),
		dao.Invoice.Columns().Currency:                       req.Currency,
		dao.Invoice.Columns().TaxScale:                       req.TaxScale,
		dao.Invoice.Columns().ChannelId:                      req.ChannelId,
		dao.Invoice.Columns().Lines:                          utility.MarshalToJsonString(invoiceItems),
		dao.Invoice.Columns().GmtModify:                      gtime.Now(),
	}).Where(dao.Subscription.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	//rowAffected, err := update.RowsAffected()
	//if rowAffected != 1 {
	//	return nil, gerror.Newf("EditInvoice update err:%s", update)
	//}
	one.Currency = req.Currency
	one.TaxScale = req.TaxScale
	one.ChannelId = req.ChannelId
	one.Lines = utility.MarshalToJsonString(invoiceItems)
	return &invoice.NewInvoiceEditRes{Invoice: invoice_compute.ConvertInvoiceToRo(one)}, nil
}

func DeletePendingInvoice(ctx context.Context, invoiceId string) error {
	one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
	utility.Assert(one != nil, fmt.Sprintf("invoice not found:%s", invoiceId))
	utility.Assert(one.Status == consts.InvoiceStatusPending, "invoice not in pending status")
	if one.IsDeleted == 1 {
		return nil
	} else {
		//更新 Subscription
		_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
			dao.Invoice.Columns().IsDeleted: 0,
			dao.Invoice.Columns().GmtModify: gtime.Now(),
		}).Where(dao.Subscription.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			return err
		}
		return nil
	}
}

func CancelProcessingInvoice(ctx context.Context, invoiceId string) error {
	one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
	utility.Assert(one != nil, fmt.Sprintf("invoice not found:%s", invoiceId))
	if one.Status == consts.InvoiceStatusCancelled {
		return nil
	}
	utility.Assert(one.Status == consts.InvoiceStatusProcessing, "invoice not in processing status")
	utility.Assert(one.IsDeleted == 0, "invoice is deleted")
	if one.IsDeleted == 1 {
		return nil
	} else {
		payChannel := query.GetSubscriptionTypePayChannelById(ctx, one.ChannelId)
		utility.Assert(payChannel != nil, "payChannel not found")
		_, err := out.GetPayChannelServiceProvider(ctx, one.ChannelId).DoRemoteChannelInvoiceCancel(ctx, payChannel, &ro.ChannelCancelInvoiceInternalReq{
			ChannelInvoiceId: one.ChannelInvoiceId,
		})
		if err != nil {
			return gerror.Newf(`FinishInvoice failure %v`, err)
		}
		// todo mark 重新生成 cancel 状态的 pdf 并发送邮件
		//更新 Subscription
		_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
			dao.Invoice.Columns().Status:    consts.InvoiceStatusCancelled,
			dao.Invoice.Columns().GmtModify: gtime.Now(),
		}).Where(dao.Subscription.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			return err
		}
		//rowAffected, err := update.RowsAffected()
		//if rowAffected != 1 {
		//	return gerror.Newf("EditInvoice update err:%s", update)
		//}
		return nil
	}
}

func FinishInvoice(ctx context.Context, req *invoice.ProcessInvoiceForPayReq) (*invoice.ProcessInvoiceForPayRes, error) {
	one := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	utility.Assert(one != nil, fmt.Sprintf("invoice not found:%s", req.InvoiceId))
	utility.Assert(one.Status == consts.InvoiceStatusPending, "invoice not in pending status")
	utility.Assert(one.IsDeleted == 0, "invoice is deleted")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, one.ChannelId)
	utility.Assert(payChannel != nil, "payChannel not found")
	var lines []*ro.InvoiceItemDetailRo
	err := utility.UnmarshalFromJsonString(one.Lines, &lines)
	if err != nil {
		return nil, err
	}
	checkInvoice(invoice_compute.ConvertInvoiceToRo(one))

	merchantInfo := query.GetMerchantInfoById(ctx, one.MerchantId)
	user := query.GetUserAccountById(ctx, uint64(one.UserId))

	createPayContext := &ro.CreatePayContext{
		PayChannel: payChannel,
		Pay: &entity.Payment{
			BizId:       one.InvoiceId,
			BizType:     consts.BIZ_TYPE_INVOICE,
			UserId:      one.UserId,
			ChannelId:   int64(payChannel.Id),
			TotalAmount: one.TotalAmount,
			Currency:    one.Currency,
			CountryCode: user.CountryCode,
			MerchantId:  one.MerchantId,
			CompanyId:   merchantInfo.CompanyId,
		},
		Platform:      "WEB",
		DeviceType:    "Web",
		ShopperUserId: strconv.FormatInt(one.UserId, 10),
		ShopperEmail:  user.Email,
		ShopperLocale: "en",
		Mobile:        user.Mobile,
		Invoice: &ro.InvoiceDetailSimplify{
			InvoiceId:                      one.InvoiceId,
			TotalAmount:                    one.TotalAmount,
			TotalAmountExcludingTax:        one.TotalAmountExcludingTax,
			Currency:                       one.Currency,
			TaxAmount:                      one.TaxAmount,
			SubscriptionAmount:             one.SubscriptionAmount,
			SubscriptionAmountExcludingTax: one.SubscriptionAmountExcludingTax,
			Lines:                          lines,
			PeriodEnd:                      one.PeriodEnd,
			PeriodStart:                    one.PeriodStart,
		},
		ShopperName: &v1.OutShopperName{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Gender:    user.Gender,
		},
		MerchantOrderReference: one.InvoiceId,
		PayMethod:              1, //automatic
		DaysUtilDue:            5, //one day expire
	}

	createRes, err := service.DoChannelPay(ctx, createPayContext)

	//createRes, err := out.GetPayChannelServiceProvider(ctx, one.ChannelId).DoRemoteChannelInvoiceCreateAndPay(ctx, payChannel, &ro.ChannelCreateInvoiceInternalReq{
	//	Invoice:      one,
	//	InvoiceLines: lines,
	//	PayMethod:    req.PayMethod,
	//	DaysUtilDue:  req.DaysUtilDue,
	//})

	if err != nil {
		return nil, gerror.Newf(`FinishInvoice failure %v`, err)
	}
	//更新 Subscription
	_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
		dao.Invoice.Columns().ChannelInvoiceId:  createRes.ChannelInvoiceId,
		dao.Invoice.Columns().ChannelInvoicePdf: createRes.ChannelInvoicePdf,
		dao.Invoice.Columns().Status:            createRes.InvoiceStatus,
		dao.Invoice.Columns().Link:              createRes.Link,
		dao.Invoice.Columns().GmtModify:         gtime.Now(),
	}).Where(dao.Subscription.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	//rowAffected, err := update.RowsAffected()
	//if rowAffected != 1 {
	//	return nil, gerror.Newf("FinishInvoice update err:%s", update)
	//}
	one.Status = int(createRes.InvoiceStatus)
	one.Link = createRes.Link

	return &invoice.ProcessInvoiceForPayRes{Invoice: one}, nil
}

func CreateInvoiceRefund(ctx context.Context, req *invoice.NewInvoiceRefundReq) (*entity.Refund, error) {
	utility.Assert(req.RefundAmount > 0, "refundFee should > 0")
	utility.Assert(len(req.InvoiceId) > 0, "invoiceId invalid")
	utility.Assert(len(req.Reason) > 0, "reason should not be blank")
	one := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	utility.Assert(one != nil, "invoice not found")
	utility.Assert(one.TotalAmount > req.RefundAmount, "not enough amount to refund")
	utility.Assert(len(one.PaymentId) > 0, "paymentId not found")
	payment := query.GetPaymentByPaymentId(ctx, one.PaymentId)
	refund, err := service.DoChannelRefund(ctx, payment.BizType, &v1.RefundsReq{
		PaymentId:        one.PaymentId,
		MerchantId:       one.MerchantId,
		MerchantRefundId: one.InvoiceId,
		Reason:           req.Reason,
		Amount: &v1.PayAmountVo{
			Currency: one.Currency,
			Value:    req.RefundAmount,
		},
	}, 0)
	if err != nil {
		return nil, err
	}
	return refund, nil
}
