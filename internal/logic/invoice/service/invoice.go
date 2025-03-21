package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
	redismq "github.com/jackyang-hk/go-redismq"
	"math"
	"strconv"
	"strings"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/api/merchant/invoice"
	"unibee/internal/cmd/config"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	"unibee/internal/controller/link"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/invoice/handler"
	"unibee/internal/logic/operation_log"
	handler2 "unibee/internal/logic/payment/handler"
	"unibee/internal/logic/payment/service"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func TryCancelSubscriptionLatestInvoice(ctx context.Context, subscription *entity.Subscription) {
	one := query.GetInvoiceByInvoiceId(ctx, subscription.LatestInvoiceId)
	if one != nil && one.Status == consts.InvoiceStatusProcessing {
		err := CancelProcessingInvoice(ctx, one.InvoiceId, "TryCancelSubscriptionLatestInvoice")
		if err != nil {
			g.Log().Errorf(ctx, `TryCancelSubscriptionLatestInvoice failure error:%s`, err.Error())
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
	if req.GatewayId <= 0 {
		gatewayId, _ := strconv.ParseUint(user.GatewayId, 10, 64)
		if gatewayId > 0 {
			req.GatewayId = gatewayId
		}
	}
	utility.Assert(req.GatewayId > 0, "invalid gatewayId")
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
			OriginAmount:           amountExcludingTax + tax,
			Amount:                 amountExcludingTax + tax,
			DiscountAmount:         0,
			Tax:                    tax,
			TaxPercentage:          req.TaxPercentage,
			AmountExcludingTax:     amountExcludingTax,
			UnitAmountExcludingTax: line.UnitAmountExcludingTax,
			Quantity:               line.Quantity,
			Name:                   line.Name,
			Description:            line.Description,
		})
		totalTax = totalTax + tax
		totalAmountExcludingTax = totalAmountExcludingTax + amountExcludingTax
	}
	totalTax = int64(math.Round(float64(totalAmountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(req.TaxPercentage)))
	var totalAmount = totalTax + totalAmountExcludingTax

	invoiceId := utility.CreateInvoiceId()
	one := &entity.Invoice{
		BizType:                        consts.BizTypeInvoice,
		MerchantId:                     merchantId,
		InvoiceId:                      invoiceId,
		InvoiceName:                    req.Name,
		ProductName:                    req.Name,
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
		CountryCode:                    user.CountryCode,
		CreateFrom:                     "Admin",
	}

	result, err := dao.Invoice.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		return nil, gerror.Newf(`CreateInvoice record insert failure %s`, err)
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))

	one.Lines = utility.MarshalToJsonString(invoiceItems)
	_, _ = redismq.Send(&redismq.Message{
		Topic:      redismq2.TopicInvoiceCreated.Topic,
		Tag:        redismq2.TopicInvoiceCreated.Tag,
		Body:       one.InvoiceId,
		CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
	})
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Invoice(%s)", one.InvoiceId),
		Content:        "New",
		UserId:         one.UserId,
		SubscriptionId: one.SubscriptionId,
		InvoiceId:      one.InvoiceId,
		PlanId:         0,
		DiscountCode:   "",
	}, err)
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
			OriginAmount:           amountExcludingTax + tax,
			Amount:                 amountExcludingTax + tax,
			DiscountAmount:         0,
			Tax:                    tax,
			TaxPercentage:          req.TaxPercentage,
			AmountExcludingTax:     amountExcludingTax,
			UnitAmountExcludingTax: line.UnitAmountExcludingTax,
			Quantity:               line.Quantity,
			Name:                   line.Name,
			Description:            line.Description,
		})
		totalTax = totalTax + tax
		totalAmountExcludingTax = totalAmountExcludingTax + amountExcludingTax
	}
	totalTax = int64(math.Round(float64(totalAmountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(req.TaxPercentage)))
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
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Invoice(%s)", one.InvoiceId),
		Content:        "Edit",
		UserId:         one.UserId,
		SubscriptionId: one.SubscriptionId,
		InvoiceId:      one.InvoiceId,
		PlanId:         0,
		DiscountCode:   "",
	}, err)
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

		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     one.MerchantId,
			Target:         fmt.Sprintf("Invoice(%s)", one.InvoiceId),
			Content:        "Delete",
			UserId:         one.UserId,
			SubscriptionId: one.SubscriptionId,
			InvoiceId:      one.InvoiceId,
			PlanId:         0,
			DiscountCode:   "",
		}, err)

		if err != nil {
			return err
		}
		return nil
	}
}

func CancelProcessingInvoice(ctx context.Context, invoiceId string, reason string) error {
	one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
	utility.Assert(one != nil, fmt.Sprintf("invoice not found:%s", invoiceId))
	if one.Status == consts.InvoiceStatusCancelled || one.Status == consts.InvoiceStatusFailed {
		return nil
	}
	utility.Assert(one.Status == consts.InvoiceStatusProcessing, "invoice not in processing status")
	utility.Assert(one.IsDeleted == 0, "invoice is deleted")
	g.Log().Infof(ctx, "CancelProcessingInvoice invoiceId:%s reason:%s", invoiceId, reason)
	invoiceStatus := consts.InvoiceStatusCancelled
	_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
		dao.Invoice.Columns().Status:    invoiceStatus,
		dao.Invoice.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return err
	}
	one.Status = invoiceStatus
	_ = handler.InvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true, false)
	_, _ = redismq.Send(&redismq.Message{
		Topic:      redismq2.TopicInvoiceCancelled.Topic,
		Tag:        redismq2.TopicInvoiceCancelled.Tag,
		Body:       one.InvoiceId,
		CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
	})

	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Invoice(%s)", one.InvoiceId),
		Content:        "Cancel",
		UserId:         one.UserId,
		SubscriptionId: one.SubscriptionId,
		InvoiceId:      one.InvoiceId,
		PlanId:         0,
		DiscountCode:   "",
	}, err)

	if len(one.RefundId) > 0 {
		refund := query.GetRefundByRefundId(ctx, one.RefundId)
		if refund != nil {
			err = service.PaymentRefundGatewayCancel(ctx, refund)
			if err != nil {
				g.Log().Errorf(ctx, `PaymentRefundGatewayCancel failure %s`, err.Error())
			}
			return err
		}
	} else if len(one.PaymentId) > 0 {
		payment := query.GetPaymentByPaymentId(ctx, one.PaymentId)
		if payment != nil {
			err = service.PaymentGatewayCancel(ctx, payment)
			if err != nil {
				g.Log().Errorf(ctx, `PaymentGatewayCancel failure %s`, err.Error())
			}
			return err
		}
	}
	return nil
}

func ProcessingInvoiceFailure(ctx context.Context, invoiceId string, reason string) error {
	one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
	utility.Assert(one != nil, fmt.Sprintf("invoice not found:%s", invoiceId))
	if one.Status == consts.InvoiceStatusCancelled || one.Status == consts.InvoiceStatusFailed {
		return nil
	}
	utility.Assert(one.Status == consts.InvoiceStatusProcessing, "invoice not in processing status")
	utility.Assert(one.IsDeleted == 0, "invoice is deleted")
	g.Log().Infof(ctx, "ProcessingInvoiceFailure invoiceId:%s reason:%s", invoiceId, reason)
	_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
		dao.Invoice.Columns().Status:     consts.InvoiceStatusFailed,
		dao.Invoice.Columns().SendStatus: consts.InvoiceSendStatusUnnecessary,
		dao.Invoice.Columns().GmtModify:  gtime.Now(),
	}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()

	if err == nil {
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     one.MerchantId,
			Target:         fmt.Sprintf("Invoice(%s)", one.InvoiceId),
			Content:        "Expired",
			UserId:         one.UserId,
			SubscriptionId: one.SubscriptionId,
			InvoiceId:      one.InvoiceId,
			PlanId:         0,
			DiscountCode:   "",
		}, err)
	}

	_, _ = redismq.Send(&redismq.Message{
		Topic:      redismq2.TopicInvoiceFailed.Topic,
		Tag:        redismq2.TopicInvoiceFailed.Tag,
		Body:       one.InvoiceId,
		CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
	})

	if len(one.RefundId) > 0 {
		refund := query.GetRefundByRefundId(ctx, one.RefundId)
		if refund != nil {
			err = service.PaymentRefundGatewayCancel(ctx, refund)
			if err != nil {
				g.Log().Errorf(ctx, `PaymentRefundGatewayCancel failure %s`, err.Error())
			}
			return err
		}
	} else if len(one.PaymentId) > 0 {
		payment := query.GetPaymentByPaymentId(ctx, one.PaymentId)
		if payment != nil {
			err = service.PaymentGatewayCancel(ctx, payment)
			if err != nil {
				g.Log().Errorf(ctx, `PaymentGatewayCancel failure %s`, err.Error())
			}
			return err
		}
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
	st := utility.CreateInvoiceSt()
	invoiceLink := link.GetInvoiceLink(one.InvoiceId, st)
	_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
		dao.Invoice.Columns().SendStatus: consts.InvoiceSendStatusUnSend,
		dao.Invoice.Columns().Status:     invoiceStatus,
		dao.Invoice.Columns().SendTerms:  st,
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
	one.SendTerms = st
	_ = handler.InvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true, false)
	_, _ = redismq.Send(&redismq.Message{
		Topic:      redismq2.TopicInvoiceProcessed.Topic,
		Tag:        redismq2.TopicInvoiceProcessed.Tag,
		Body:       one.InvoiceId,
		CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
	})
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Invoice(%s)", one.InvoiceId),
		Content:        "Finish",
		UserId:         one.UserId,
		SubscriptionId: one.SubscriptionId,
		InvoiceId:      one.InvoiceId,
		PlanId:         0,
		DiscountCode:   "",
	}, err)
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
	gateway := query.GetGatewayById(ctx, payment.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	if _interface.Context().Get(ctx).IsOpenApiCall {
		utility.Assert(gateway.GatewayType != consts.GatewayTypeCrypto, "crypto payment refund not available, refund manual is need and then mark a payment refund")
		if config.GetConfigInstance().IsProd() {
			utility.Assert(gateway.GatewayType != consts.GatewayTypeWireTransfer, "wire transfer payment refund not available, should refund manual and then mark a payment refund in admin portal")
		}
	} else if gateway.GatewayType == consts.GatewayTypeWireTransfer || gateway.GatewayType == consts.GatewayTypeCrypto {
		utility.Assert(len(req.Reason) > 0, "reason is need for crypto|wire transfer refund")
	}
	var reason = "Refund Requested"
	if len(req.Reason) > 0 {
		reason = fmt.Sprintf("%s: %s", reason, req.Reason)
	}
	refund, err := service.GatewayPaymentRefundCreate(ctx, &service.NewPaymentRefundInternalReq{
		PaymentId:        one.PaymentId,
		ExternalRefundId: fmt.Sprintf("%s-%s", one.PaymentId, req.RefundNo),
		Reason:           reason,
		RefundAmount:     req.RefundAmount,
		Currency:         one.Currency,
	})
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Invoice(%s)", one.InvoiceId),
		Content:        "NewRefund",
		UserId:         one.UserId,
		SubscriptionId: one.SubscriptionId,
		InvoiceId:      one.InvoiceId,
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	if err != nil {
		return nil, err
	}

	return refund, nil
}

func MarkInvoiceRefundSuccess(ctx context.Context, merchantId uint64, invoiceId string, reason string) {
	utility.Assert(len(invoiceId) > 0, "invoiceId invalid")
	one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
	utility.Assert(one != nil, "invoice not found")
	utility.Assert(one.MerchantId == merchantId, "wrong merchant")
	//gateway := query.GetGatewayById(ctx, one.GatewayId)
	//utility.Assert(gateway != nil && (gateway.GatewayType == consts.GatewayTypeWireTransfer || gateway.GatewayType == consts.GatewayTypeCrypto), "gateway not wire transfer or changelly type")
	utility.Assert(len(one.RefundId) > 0, "invoiceId not refund invoice")
	refund := query.GetRefundByRefundId(ctx, one.RefundId)
	utility.Assert(refund != nil, "refund not found")
	//gateway = query.GetGatewayById(ctx, refund.GatewayId)
	//utility.Assert(gateway != nil && (gateway.GatewayType == consts.GatewayTypeWireTransfer || gateway.GatewayType == consts.GatewayTypeCrypto), "gateway not wire transfer or changelly type")

	err := handler2.HandleRefundSuccess(ctx, &handler2.HandleRefundReq{
		RefundId:         one.RefundId,
		GatewayRefundId:  refund.GatewayRefundId,
		RefundAmount:     refund.RefundAmount,
		RefundStatusEnum: consts.RefundSuccess,
		RefundTime:       gtime.Now(),
		Reason:           reason,
	})
	if err != nil {
		g.Log().Errorf(ctx, "MarkInvoiceRefundSuccess invoiceId:%s error:%s", invoiceId, err.Error())
	}
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Invoice(%s)", one.InvoiceId),
		Content:        "MarkRefundSuccess",
		UserId:         one.UserId,
		SubscriptionId: one.SubscriptionId,
		InvoiceId:      one.InvoiceId,
		PlanId:         0,
		DiscountCode:   "",
	}, err)
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
	utility.Assert(gateway.GatewayType == consts.GatewayTypeCrypto || gateway.GatewayType == consts.GatewayTypeWireTransfer, "mark refund only support crypto or wire transfer invoice")
	refund, err := service.MarkPaymentRefundCreate(ctx, &service.NewPaymentRefundInternalReq{
		PaymentId:        one.PaymentId,
		ExternalRefundId: fmt.Sprintf("%s-%s", one.PaymentId, req.RefundNo),
		Reason:           req.Reason,
		RefundAmount:     req.RefundAmount,
		Currency:         one.Currency,
	})
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Invoice(%s)", one.InvoiceId),
		Content:        "MarkRefund",
		UserId:         one.UserId,
		SubscriptionId: one.SubscriptionId,
		InvoiceId:      one.InvoiceId,
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	if err != nil {
		return nil, err
	}

	return refund, nil
}

func HardDeleteInvoice(ctx context.Context, merchantId uint64, invoiceId string) error {
	utility.Assert(merchantId > 0, "invalid merchantId")
	utility.Assert(len(invoiceId) > 0, "invalid invoiceId")
	_, err := dao.Invoice.Ctx(ctx).Where(dao.Invoice.Columns().InvoiceId, invoiceId).Delete()
	return err
}

func MarkWireTransferInvoiceAsPaid(ctx context.Context, invoiceId string, transferNumber string, reason string) (*entity.Invoice, error) {
	utility.Assert(len(invoiceId) > 0, "invalid invoiceId")
	utility.Assert(len(transferNumber) > 0, "invalid transferNumber")
	one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
	utility.Assert(one != nil, "invoice not found, InvoiceId:"+invoiceId)
	utility.Assert(one.Status == consts.InvoiceStatusProcessing, "invoice not process status, InvoiceId:"+invoiceId)
	utility.Assert(one.TotalAmount != 0, "invoice totalAmount not zero, InvoiceId:"+invoiceId)
	gateway := query.GetGatewayById(ctx, one.GatewayId)
	utility.Assert(gateway != nil, "invoice gateway not found")
	utility.Assert(gateway.GatewayType == consts.GatewayTypeWireTransfer, "invoice not wire transfer type")
	utility.Assert(one.TotalAmount >= gateway.MinimumAmount, "Total Amount not reach the gateway's minimum amount")
	utility.Assert(strings.ToUpper(one.Currency) == strings.ToUpper(gateway.Currency), "Invoice currency not reach the gateway's currency")
	payment := query.GetPaymentByPaymentId(ctx, one.PaymentId)
	if payment == nil {
		res, err := service.CreateSubInvoicePaymentDefaultAutomatic(ctx, &service.CreateSubInvoicePaymentDefaultAutomaticReq{
			Invoice:       one,
			ManualPayment: true,
			ReturnUrl:     "",
			CancelUrl:     "",
			Source:        "MarkWireTransferInvoiceAsPaid",
			TimeNow:       0,
		})
		utility.AssertError(err, "Mark as success error")
		payment = res.Payment
		utility.Assert(payment != nil, "payment not found")
	}
	err := handler2.HandlePaySuccess(ctx, &handler2.HandlePayReq{
		PaymentId:              payment.PaymentId,
		GatewayPaymentIntentId: transferNumber,
		GatewayPaymentId:       transferNumber,
		TotalAmount:            payment.TotalAmount,
		PayStatusEnum:          consts.PaymentSuccess,
		PaidTime:               gtime.Now(),
		PaymentAmount:          payment.TotalAmount,
		Reason:                 reason,
	})
	utility.AssertError(err, "MarkWireTransferInvoiceAsPaid")
	one = query.GetInvoiceByInvoiceId(ctx, invoiceId)
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Invoice(%s)", one.InvoiceId),
		Content:        "MarkInvoiceAsPaid(WireTransfer)",
		UserId:         one.UserId,
		SubscriptionId: one.SubscriptionId,
		InvoiceId:      one.InvoiceId,
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return one, nil
}
