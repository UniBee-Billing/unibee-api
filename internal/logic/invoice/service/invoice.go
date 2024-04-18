package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
	"strings"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/api/merchant/invoice"
	"unibee/internal/consts"
	"unibee/internal/controller/link"
	dao "unibee/internal/dao/oversea_pay"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/email"
	"unibee/internal/logic/invoice/handler"
	"unibee/internal/logic/payment/service"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func CancelInvoiceForSubscription(ctx context.Context, subscription *entity.Subscription) {
	var mainList = make([]*entity.Invoice, 0)
	m := dao.Invoice.Ctx(ctx)
	_ = m.Where(dao.Invoice.Columns().IsDeleted, 0).
		Where(dao.Invoice.Columns().MerchantId, subscription.MerchantId).
		Where(dao.Invoice.Columns().SubscriptionId, subscription.SubscriptionId).
		Where(dao.Invoice.Columns().Status, consts.InvoiceStatusProcessing).
		OmitEmpty().Scan(&mainList)
	for _, one := range mainList {
		if len(one.PaymentId) > 0 {
			payment := query.GetPaymentByPaymentId(ctx, one.PaymentId)
			if payment != nil {
				gateway := query.GetGatewayById(ctx, one.GatewayId)
				if gateway != nil {
					err := service.PaymentGatewayCancel(ctx, payment)
					if err != nil {
						g.Log().Errorf(ctx, `PaymentGatewayCancel failure for CancelInvoiceForSubscription %s`, err.Error())
					}
				}
			}
		}
	}
}

func checkInvoice(one *detail.InvoiceDetail) {
	var totalAmountExcludingTax int64 = 0
	var totalTax int64 = 0
	for _, line := range one.Lines {
		amountExcludingTax := line.UnitAmountExcludingTax * line.Quantity
		tax := int64(float64(amountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(one.TaxPercentage))
		utility.Assert(line.AmountExcludingTax == amountExcludingTax, "line amountExcludingTax mistake")
		utility.Assert(strings.Compare(line.Currency, one.Currency) == 0, "line currency not match invoice currency")
		utility.Assert(line.Amount == amountExcludingTax+tax, "line amount mistake")
		//utility.Assert(line.TaxPercentage == one.TaxPercentage, "line TaxPercentage mistake")
		totalTax = totalTax + tax
		totalAmountExcludingTax = totalAmountExcludingTax + amountExcludingTax
	}
	var totalAmount = totalTax + totalAmountExcludingTax
	utility.Assert(one.TaxAmount == totalTax, "invoice taxAmount mistake")
	utility.Assert(one.TotalAmountExcludingTax == totalAmountExcludingTax, "invoice totalAmountExcludingTax mistake")
	utility.Assert(one.TotalAmount == totalAmount, "line totalAmount mistake")
}

func CreateInvoice(ctx context.Context, merchantId uint64, req *invoice.NewReq) (res *invoice.NewRes, err error) {
	user := query.GetUserAccountById(ctx, req.UserId)
	utility.Assert(user != nil, fmt.Sprintf("send user not found:%d", req.UserId))
	utility.Assert(len(user.Email) > 0, fmt.Sprintf("send user email not found:%d", req.UserId))
	gateway := query.GetGatewayById(ctx, req.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")

	var invoiceItems []*bean.InvoiceItemSimplify
	var totalAmountExcludingTax int64 = 0
	var totalTax int64 = 0
	for _, line := range req.Lines {
		amountExcludingTax := line.UnitAmountExcludingTax * line.Quantity
		tax := int64(float64(amountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(req.TaxPercentage))
		invoiceItems = append(invoiceItems, &bean.InvoiceItemSimplify{
			Currency:               req.Currency,
			TaxPercentage:          req.TaxPercentage,
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

	invoiceId := utility.CreateInvoiceId()
	one := &entity.Invoice{
		BizType:                        consts.BizTypeInvoice,
		MerchantId:                     merchantId,
		InvoiceId:                      invoiceId,
		InvoiceName:                    req.Name,
		UniqueId:                       invoiceId,
		TotalAmount:                    totalAmount,
		TotalAmountExcludingTax:        totalAmountExcludingTax,
		TaxAmount:                      totalTax,
		TaxPercentage:                  req.TaxPercentage,
		SubscriptionAmount:             totalAmount,
		SubscriptionAmountExcludingTax: totalAmountExcludingTax,
		Currency:                       strings.ToUpper(req.Currency),
		Lines:                          utility.MarshalToJsonString(invoiceItems),
		GatewayId:                      req.GatewayId,
		Status:                         consts.InvoiceStatusPending,
		SendStatus:                     consts.InvoiceSendStatusUnSend,
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
		finishRes, err := FinishInvoice(ctx, &invoice.FinishReq{
			InvoiceId: one.InvoiceId,
			//PayMethod:   2,
			DaysUtilDue: 3,
		})
		if err != nil {
			return nil, err
		}
		one.Link = finishRes.Invoice.Link
		one.PaymentLink = finishRes.Invoice.PaymentLink
		one.Status = finishRes.Invoice.Status
		one.PaymentId = finishRes.Invoice.PaymentId
	}
	return &invoice.NewRes{Invoice: detail.ConvertInvoiceToDetail(ctx, one)}, nil
}

func EditInvoice(ctx context.Context, req *invoice.EditReq) (res *invoice.EditRes, err error) {
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

	var invoiceItems []*bean.InvoiceItemSimplify
	var totalAmountExcludingTax int64 = 0
	var totalTax int64 = 0
	for _, line := range req.Lines {
		amountExcludingTax := line.UnitAmountExcludingTax * line.Quantity
		tax := int64(float64(amountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(req.TaxPercentage))
		invoiceItems = append(invoiceItems, &bean.InvoiceItemSimplify{
			Currency:               req.Currency,
			TaxPercentage:          req.TaxPercentage,
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
		dao.Invoice.Columns().TaxPercentage:                  req.TaxPercentage,
		dao.Invoice.Columns().GatewayId:                      req.GatewayId,
		dao.Invoice.Columns().Lines:                          utility.MarshalToJsonString(invoiceItems),
		dao.Invoice.Columns().GmtModify:                      gtime.Now(),
	}).Where(dao.Subscription.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	one.Currency = req.Currency
	one.TaxPercentage = req.TaxPercentage
	one.GatewayId = req.GatewayId
	one.Lines = utility.MarshalToJsonString(invoiceItems)
	if req.Finish {
		finishRes, err := FinishInvoice(ctx, &invoice.FinishReq{
			InvoiceId: one.InvoiceId,
			//PayMethod:   2,
			DaysUtilDue: 3,
		})
		if err != nil {
			return nil, err
		}
		one.Link = finishRes.Invoice.Link
		one.PaymentLink = finishRes.Invoice.PaymentLink
		one.Status = finishRes.Invoice.Status
		one.PaymentId = finishRes.Invoice.PaymentId
	}
	return &invoice.EditRes{Invoice: detail.ConvertInvoiceToDetail(ctx, one)}, nil
}

func DeletePendingInvoice(ctx context.Context, invoiceId string) error {
	one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
	utility.Assert(one != nil, fmt.Sprintf("invoice not found:%s", invoiceId))
	utility.Assert(one.Status == consts.InvoiceStatusPending, "invoice not in pending status")
	if one.IsDeleted == 1 {
		return nil
	} else {
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
	invoiceStatus := consts.InvoiceStatusCancelled
	_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
		dao.Invoice.Columns().Status:    invoiceStatus,
		dao.Invoice.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return err
	}
	one.Status = invoiceStatus
	_ = handler.InvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true)

	payment := query.GetPaymentByPaymentId(ctx, one.PaymentId)
	if payment != nil {
		err = service.PaymentGatewayCancel(ctx, payment)
		if err != nil {
			g.Log().Errorf(ctx, `PaymentGatewayCancel failure %s`, err.Error())
		}
		return err
	}
	return nil

}

func FinishInvoice(ctx context.Context, req *invoice.FinishReq) (*invoice.FinishRes, error) {
	one := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	utility.Assert(one != nil, fmt.Sprintf("invoice not found:%s", req.InvoiceId))
	utility.Assert(one.Status == consts.InvoiceStatusPending, "invoice not in pending status")
	utility.Assert(one.IsDeleted == 0, "invoice is deleted")
	checkInvoice(detail.ConvertInvoiceToDetail(ctx, one))
	if req.DaysUtilDue <= 0 {
		req.DaysUtilDue = consts.DEFAULT_DAY_UTIL_DUE
	}
	invoiceStatus := consts.InvoiceStatusProcessing
	invoiceLink := link.GetInvoiceLink(one.InvoiceId)
	_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
		dao.Invoice.Columns().SendStatus: consts.InvoiceSendStatusUnSend,
		dao.Invoice.Columns().Status:     invoiceStatus,
		dao.Invoice.Columns().Link:       invoiceLink,
		dao.Invoice.Columns().DayUtilDue: req.DaysUtilDue,
		dao.Invoice.Columns().FinishTime: gtime.Now().Timestamp(),
		dao.Invoice.Columns().GmtModify:  gtime.Now(),
	}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	one.Status = invoiceStatus
	one.Link = invoiceLink
	_ = handler.InvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true)

	return &invoice.FinishRes{Invoice: bean.SimplifyInvoice(one)}, nil
}

func CreateInvoiceRefund(ctx context.Context, req *invoice.RefundReq) (*entity.Refund, error) {
	utility.Assert(req.RefundAmount > 0, "refundFee should > 0")
	utility.Assert(len(req.InvoiceId) > 0, "invoiceId invalid")
	utility.Assert(len(req.Reason) > 0, "reason should not be blank")
	if _interface.Context().Get(ctx).IsOpenApiCall {
		utility.Assert(len(req.RefundNo) > 0, "refundNo should not be blank")
	} else if len(req.RefundNo) == 0 {
		req.RefundNo = uuid.New().String()
	}
	one := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	utility.Assert(one != nil, "invoice not found")
	utility.Assert(one.TotalAmount >= req.RefundAmount, "not enough amount to refund")
	utility.Assert(len(one.PaymentId) > 0, "paymentId not found")
	payment := query.GetPaymentByPaymentId(ctx, one.PaymentId)
	utility.Assert(payment != nil, "payment not found")
	refund, err := service.GatewayPaymentRefundCreate(ctx, &service.NewPaymentRefundInternalReq{
		PaymentId:        one.PaymentId,
		ExternalRefundId: fmt.Sprintf("%s-%s", one.PaymentId, req.RefundNo),
		Reason:           req.Reason,
		RefundAmount:     req.RefundAmount,
		Currency:         one.Currency,
	})
	if err != nil {
		return nil, err
	}
	user := query.GetUserAccountById(ctx, payment.UserId)
	if user != nil {
		merchant := query.GetMerchantById(ctx, payment.MerchantId)
		if merchant != nil {
			err = email.SendTemplateEmail(ctx, merchant.Id, user.Email, user.TimeZone, email.TemplateInvoiceRefundCreated, "", &email.TemplateVariable{
				UserName:            user.FirstName + " " + user.LastName,
				MerchantCustomEmail: merchant.Email,
				MerchantName:        query.GetMerchantCountryConfigName(ctx, payment.MerchantId, user.CountryCode),
				PaymentAmount:       utility.ConvertCentToDollarStr(refund.RefundAmount, refund.Currency),
				Currency:            strings.ToUpper(refund.Currency),
			})
			if err != nil {
				fmt.Printf("CreateInvoiceRefund SendTemplateEmail err:%s", err.Error())
			}
		}
	}
	_ = handler.InvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true)

	return refund, nil
}

func MarkInvoiceRefund(ctx context.Context, req *invoice.MarkRefundReq) (*entity.Refund, error) {
	utility.Assert(req.RefundAmount > 0, "refundFee should > 0")
	utility.Assert(len(req.InvoiceId) > 0, "invoiceId invalid")
	utility.Assert(len(req.Reason) > 0, "reason should not be blank")
	if _interface.Context().Get(ctx).IsOpenApiCall {
		utility.Assert(len(req.RefundNo) > 0, "refundNo should not be blank")
	} else if len(req.RefundNo) == 0 {
		req.RefundNo = uuid.New().String()
	}
	one := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	utility.Assert(one != nil, "invoice not found")
	utility.Assert(one.TotalAmount >= req.RefundAmount, "not enough amount to refund")
	utility.Assert(len(one.PaymentId) > 0, "paymentId not found")
	payment := query.GetPaymentByPaymentId(ctx, one.PaymentId)
	utility.Assert(payment != nil, "payment not found")
	gateway := query.GetGatewayById(ctx, payment.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	utility.Assert(gateway.GatewayType == consts.GatewayTypeCrypto, "mark refund only support crypto invoice")
	refund, err := service.MarkPaymentRefundCreate(ctx, &service.NewPaymentRefundInternalReq{
		PaymentId:        one.PaymentId,
		ExternalRefundId: fmt.Sprintf("%s-%s", one.PaymentId, req.RefundNo),
		Reason:           req.Reason,
		RefundAmount:     req.RefundAmount,
		Currency:         one.Currency,
	})
	if err != nil {
		return nil, err
	}
	user := query.GetUserAccountById(ctx, payment.UserId)
	if user != nil {
		merchant := query.GetMerchantById(ctx, payment.MerchantId)
		if merchant != nil {
			err = email.SendTemplateEmail(ctx, merchant.Id, user.Email, user.TimeZone, email.TemplateInvoiceRefundCreated, "", &email.TemplateVariable{
				UserName:            user.FirstName + " " + user.LastName,
				MerchantCustomEmail: merchant.Email,
				MerchantName:        query.GetMerchantCountryConfigName(ctx, payment.MerchantId, user.CountryCode),
				PaymentAmount:       utility.ConvertCentToDollarStr(refund.RefundAmount, refund.Currency),
				Currency:            strings.ToUpper(refund.Currency),
			})
			if err != nil {
				fmt.Printf("CreateInvoiceRefund SendTemplateEmail err:%s", err.Error())
			}
		}
	}
	_ = handler.InvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true)

	return refund, nil
}

func HardDeleteInvoice(ctx context.Context, merchantId uint64, invoiceId string) error {
	utility.Assert(merchantId > 0, "invalid merchantId")
	utility.Assert(len(invoiceId) > 0, "invalid invoiceId")
	_, err := dao.Invoice.Ctx(ctx).Where(dao.Invoice.Columns().InvoiceId, invoiceId).Delete()
	return err
}
