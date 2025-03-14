package service

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/internal/consts"
	"unibee/internal/controller/link"
	"unibee/internal/logic/gateway/api"
	handler2 "unibee/internal/logic/payment/handler"
	"unibee/internal/logic/payment/service"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type LinkCheckRes struct {
	Message  string
	Link     string
	FileName string
	Invoice  *entity.Invoice
}

func LinkCheck(ctx context.Context, invoiceId string, time int64) *LinkCheckRes {
	var res = &LinkCheckRes{
		Message: "",
		Link:    "",
		Invoice: nil,
	}
	one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
	if one == nil {
		g.Log().Errorf(ctx, "LinkEntry invoice not found invoiceId: %s", invoiceId)
		res.Message = "Invoice Not Found"
		return res
	}
	res.Invoice = one
	if one.IsDeleted > 0 {
		res.Message = "Invoice Deleted"
	} else if one.Status == consts.InvoiceStatusCancelled {
		res.Message = "Invoice Cancelled"
	} else if one.Status == consts.InvoiceStatusFailed {
		res.Message = "Invoice Failure"
	} else if one.Status < consts.InvoiceStatusProcessing {
		res.Message = "Invoice Not Ready"
	} else if one.Status == consts.InvoiceStatusProcessing {
		dayUtilDue := one.DayUtilDue
		if dayUtilDue <= 0 {
			dayUtilDue = consts.DEFAULT_DAY_UTIL_DUE
		}
		if one.FinishTime > 0 && one.FinishTime+(dayUtilDue*86400) < time {
			res.Message = "Invoice Expire"
			return res
		}
		payment := query.GetPaymentByPaymentId(ctx, one.PaymentId)
		gateway := query.GetGatewayById(ctx, one.GatewayId)
		if gateway == nil {
			res.Message = "Gateway Error"
			return res
		}
		if payment != nil && payment.Status == consts.PaymentCreated && len(one.GatewayPaymentId) > 0 {
			gatewayPaymentRo, _ := api.GetGatewayServiceProvider(ctx, one.GatewayId).GatewayPaymentDetail(ctx, gateway, one.GatewayPaymentId, payment)
			if gatewayPaymentRo != nil {
				if gatewayPaymentRo.Status == consts.PaymentSuccess {
					_ = handler2.HandlePaySuccess(ctx, &handler2.HandlePayReq{
						PaymentId:              one.PaymentId,
						GatewayPaymentIntentId: gatewayPaymentRo.GatewayPaymentId,
						GatewayPaymentId:       gatewayPaymentRo.GatewayPaymentId,
						GatewayUserId:          gatewayPaymentRo.GatewayUserId,
						TotalAmount:            gatewayPaymentRo.TotalAmount,
						PayStatusEnum:          consts.PaymentSuccess,
						PaidTime:               gatewayPaymentRo.PaidTime,
						PaymentAmount:          gatewayPaymentRo.PaymentAmount,
						CaptureAmount:          0,
						Reason:                 gatewayPaymentRo.Reason,
						GatewayPaymentMethod:   gatewayPaymentRo.GatewayPaymentMethod,
						PaymentCode:            gatewayPaymentRo.PaymentCode,
					})
				} else if gatewayPaymentRo.Status == consts.PaymentFailed {
					_ = handler2.HandlePayFailure(ctx, &handler2.HandlePayReq{
						PaymentId:              one.PaymentId,
						GatewayPaymentIntentId: gatewayPaymentRo.GatewayPaymentId,
						GatewayPaymentId:       gatewayPaymentRo.GatewayPaymentId,
						PayStatusEnum:          consts.PaymentFailed,
						Reason:                 gatewayPaymentRo.Reason,
						PaymentCode:            gatewayPaymentRo.PaymentCode,
					})
				} else if gatewayPaymentRo.Status == consts.PaymentCancelled {
					_ = handler2.HandlePayCancel(ctx, &handler2.HandlePayReq{
						PaymentId:              one.PaymentId,
						GatewayPaymentIntentId: gatewayPaymentRo.GatewayPaymentId,
						GatewayPaymentId:       gatewayPaymentRo.GatewayPaymentId,
						PayStatusEnum:          consts.PaymentCancelled,
						Reason:                 gatewayPaymentRo.Reason,
						PaymentCode:            gatewayPaymentRo.PaymentCode,
					})
				}
				payment = query.GetPaymentByPaymentId(ctx, one.PaymentId)
			}
		}
		if payment != nil && payment.Status == consts.PaymentSuccess {
			// status haven't sync completely
			res.Link = link.GetInvoicePdfLink(one.InvoiceId, one.SendTerms)
		} else if gateway.GatewayType != consts.GatewayTypeCrypto || payment == nil || payment.Status == consts.PaymentCancelled || payment.Status == consts.PaymentFailed {
			// create new payment for this invoice
			var lines []*bean.InvoiceItemSimplify
			err := utility.UnmarshalFromJsonString(one.Lines, &lines)
			if err != nil {
				res.Message = "Server Error"
				return res
			}
			createRes, err := service.CreateSubInvoicePaymentDefaultAutomatic(ctx, &service.CreateSubInvoicePaymentDefaultAutomaticReq{
				Invoice:       one,
				ManualPayment: true,
				ReturnUrl:     "",
				CancelUrl:     "",
				Source:        "InvoiceLink",
				TimeNow:       0,
			})
			if err != nil {
				g.Log().Errorf(ctx, "GatewayPaymentCreate Error:%s\n", err.Error())
				res.Message = "Server Error"
				return res
			}
			res.Link = createRes.Link
		} else {
			//res.Link = one.PaymentLink
			res.Link = link.GetPaymentLink(payment.PaymentId)
		}
	} else if one.Status == consts.InvoiceStatusPaid || one.Status == consts.InvoiceStatusReversed {
		res.Link = link.GetInvoicePdfLink(one.InvoiceId, one.SendTerms)
	}
	return res
}
