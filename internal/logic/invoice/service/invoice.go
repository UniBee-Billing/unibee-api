package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"unibee/api/merchant/invoice"
	v1 "unibee/api/onetime/payment"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/email"
	"unibee/internal/logic/gateway/ro"
	"unibee/internal/logic/invoice/handler"
	"unibee/internal/logic/invoice/invoice_compute"
	"unibee/internal/logic/payment/service"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
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
	gateway := query.GetGatewayById(ctx, req.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")

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
		BizType:                        consts.BizTypeSubscription,
		MerchantId:                     _interface.GetMerchantId(ctx),
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
		GatewayId:                      req.GatewayId,
		Status:                         consts.InvoiceStatusPending,
		SendStatus:                     0,
		SendEmail:                      user.Email,
		UserId:                         req.UserId,
		CreateTime:                     gtime.Now().Timestamp(),
	}

	result, err := dao.Invoice.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		return nil, gerror.Newf(`CreateInvoice record insert failure %s`, err)
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))

	one.Lines = utility.MarshalToJsonString(invoiceItems)
	if req.Finish {
		finishRes, err := FinishInvoice(ctx, &invoice.FinishInvoiceForPayReq{
			InvoiceId:   one.InvoiceId,
			PayMethod:   2,
			DaysUtilDue: 3,
		})
		if err != nil {
			return nil, err
		}
		one = finishRes.Invoice
	}
	return &invoice.NewInvoiceCreateRes{Invoice: invoice_compute.ConvertInvoiceToRo(ctx, one)}, nil
}

func EditInvoice(ctx context.Context, req *invoice.NewInvoiceEditReq) (res *invoice.NewInvoiceEditRes, err error) {
	one := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	utility.Assert(one != nil, fmt.Sprintf("invoice not found:%s", req.InvoiceId))
	utility.Assert(one.Status == consts.InvoiceStatusPending, "invoice not in pending status")
	utility.Assert(one.IsDeleted == 0, "invoice is deleted")
	if req.GatewayId > 0 {
		gateway := query.GetGatewayById(ctx, req.GatewayId)
		utility.Assert(gateway != nil, "gateway not found")
	} else {
		req.GatewayId = one.GatewayId
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
		dao.Invoice.Columns().BizType:                        consts.BizTypeSubscription,
		dao.Invoice.Columns().InvoiceName:                    req.Name,
		dao.Invoice.Columns().TotalAmount:                    totalAmount,
		dao.Invoice.Columns().TotalAmountExcludingTax:        totalAmountExcludingTax,
		dao.Invoice.Columns().TaxAmount:                      totalTax,
		dao.Invoice.Columns().SubscriptionAmount:             totalAmount,
		dao.Invoice.Columns().SubscriptionAmountExcludingTax: totalAmountExcludingTax,
		dao.Invoice.Columns().Currency:                       strings.ToUpper(req.Currency),
		dao.Invoice.Columns().Currency:                       req.Currency,
		dao.Invoice.Columns().TaxScale:                       req.TaxScale,
		dao.Invoice.Columns().GatewayId:                      req.GatewayId,
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
	one.GatewayId = req.GatewayId
	one.Lines = utility.MarshalToJsonString(invoiceItems)
	if req.Finish {
		finishRes, err := FinishInvoice(ctx, &invoice.FinishInvoiceForPayReq{
			InvoiceId:   one.InvoiceId,
			PayMethod:   2,
			DaysUtilDue: 3,
		})
		if err != nil {
			return nil, err
		}
		one = finishRes.Invoice
	}
	return &invoice.NewInvoiceEditRes{Invoice: invoice_compute.ConvertInvoiceToRo(ctx, one)}, nil
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
			dao.Invoice.Columns().IsDeleted: 1,
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
	utility.Assert(len(one.PaymentId) > 0, "invoice payment not found")
	if one.IsDeleted == 1 {
		return nil
	} else {
		payment := query.GetPaymentByPaymentId(ctx, one.PaymentId)
		utility.Assert(payment != nil, "payment not found")
		gateway := query.GetGatewayById(ctx, one.GatewayId)
		utility.Assert(gateway != nil, "gateway not found")
		err := service.PaymentGatewayCancel(ctx, payment)
		if err != nil {
			g.Log().Errorf(ctx, `PaymentGatewayCancel failure %s`, err.Error())
		}
		return err
	}
}

func FinishInvoice(ctx context.Context, req *invoice.FinishInvoiceForPayReq) (*invoice.FinishInvoiceForPayRes, error) {
	one := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	utility.Assert(one != nil, fmt.Sprintf("invoice not found:%s", req.InvoiceId))
	utility.Assert(one.Status == consts.InvoiceStatusPending, "invoice not in pending status")
	utility.Assert(one.IsDeleted == 0, "invoice is deleted")
	checkInvoice(invoice_compute.ConvertInvoiceToRo(ctx, one))
	invoiceStatus := consts.InvoiceStatusProcessing
	invoiceLink := invoice_compute.GetInvoiceLink(one.InvoiceId)
	_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
		dao.Invoice.Columns().Status:    invoiceStatus,
		dao.Invoice.Columns().Link:      invoiceLink,
		dao.Invoice.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	one.Status = invoiceStatus
	one.Link = invoiceLink
	_ = handler.InvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true)

	return &invoice.FinishInvoiceForPayRes{Invoice: one}, nil
}

func CreateInvoiceRefund(ctx context.Context, req *invoice.NewInvoiceRefundReq) (*entity.Refund, error) {
	utility.Assert(req.RefundAmount > 0, "refundFee should > 0")
	utility.Assert(len(req.InvoiceId) > 0, "invoiceId invalid")
	utility.Assert(len(req.Reason) > 0, "reason should not be blank")
	utility.Assert(len(req.RefundNo) > 0, "refundNo should not be blank")
	one := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	utility.Assert(one != nil, "invoice not found")
	utility.Assert(one.TotalAmount > req.RefundAmount, "not enough amount to refund")
	utility.Assert(len(one.PaymentId) > 0, "paymentId not found")
	payment := query.GetPaymentByPaymentId(ctx, one.PaymentId)
	utility.Assert(payment != nil, "payment not found")
	refund, err := service.GatewayPaymentRefundCreate(ctx, payment.BizType, &v1.NewPaymentRefundReq{
		PaymentId:        one.PaymentId,
		MerchantRefundId: fmt.Sprintf("%s-%s", one.PaymentId, req.RefundNo),
		Reason:           req.Reason,
		Amount: &v1.AmountVo{
			Currency: one.Currency,
			Amount:   req.RefundAmount,
		},
	}, 0)
	if err != nil {
		return nil, err
	}
	user := query.GetUserAccountById(ctx, uint64(payment.UserId))
	if user != nil {
		merchant := query.GetMerchantInfoById(ctx, payment.MerchantId)
		if merchant != nil {
			err = email.SendTemplateEmail(ctx, merchant.Id, user.Email, user.TimeZone, email.TemplateInvoiceRefundCreated, "", &email.TemplateVariable{
				UserName:            user.FirstName + " " + user.LastName,
				MerchantCustomEmail: merchant.Email,
				MerchantName:        merchant.Name,
				PaymentAmount:       utility.ConvertCentToDollarStr(refund.RefundAmount, refund.Currency),
				Currency:            strings.ToUpper(refund.Currency),
			})
			if err != nil {
				fmt.Printf("CreateInvoiceRefund SendTemplateEmail err:%s", err.Error())
			}
		}
	}

	return refund, nil
}
