package service

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"strconv"
	"unibee/api/bean"
	"unibee/internal/consts"
	"unibee/internal/logic/gateway/gateway_bean"
	"unibee/internal/logic/payment/service"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type LinkCheckRes struct {
	Message string
	Link    string
	Invoice *entity.Invoice
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
		if len(one.PaymentLink) == 0 {
			// create payment link for this invoice
			gateway := query.GetGatewayById(ctx, one.GatewayId)
			if gateway == nil {
				res.Message = "Gateway Error"
				return res
			}
			var lines []*bean.InvoiceItemSimplify
			err := utility.UnmarshalFromJsonString(one.Lines, &lines)
			if err != nil {
				res.Message = "Server Error"
				return res
			}

			merchantInfo := query.GetMerchantById(ctx, one.MerchantId)
			user := query.GetUserAccountById(ctx, one.UserId)
			createPayContext := &gateway_bean.GatewayNewPaymentReq{
				Gateway: gateway,
				Pay: &entity.Payment{
					ExternalPaymentId: one.InvoiceId,
					BizType:           consts.BizTypeInvoice,
					AuthorizeStatus:   consts.Authorized,
					UserId:            one.UserId,
					GatewayId:         gateway.Id,
					TotalAmount:       one.TotalAmount,
					Currency:          one.Currency,
					CountryCode:       user.CountryCode,
					MerchantId:        one.MerchantId,
					CompanyId:         merchantInfo.CompanyId,
					BillingReason:     one.InvoiceName,
				},
				ExternalUserId: strconv.FormatUint(one.UserId, 10),
				Email:          user.Email,
				Invoice:        bean.SimplifyInvoice(one),
				Metadata:       map[string]string{"BillingReason": one.InvoiceName},
			}

			createRes, err := service.GatewayPaymentCreate(ctx, createPayContext)
			if err != nil {
				g.Log().Infof(ctx, "GatewayPaymentCreate Error:%s", err.Error())
				res.Message = "Server Error"
				return res
			}
			res.Link = createRes.Link
		} else {
			res.Link = one.PaymentLink
		}
	} else if one.Status == consts.InvoiceStatusPaid {
		res.Link = one.SendPdf
	}
	return res
}