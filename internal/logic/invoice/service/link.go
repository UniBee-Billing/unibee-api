package service

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/internal/consts"
	"unibee/internal/controller/link"
	"unibee/internal/logic/payment/service"
	entity "unibee/internal/model/entity/oversea_pay"
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
			createRes, err := service.CreateSubInvoicePaymentDefaultAutomatic(ctx, one, true, "", "", "InvoiceLink", 0)
			if err != nil {
				g.Log().Infof(ctx, "GatewayPaymentCreate Error:%s", err.Error())
				res.Message = "Server Error"
				return res
			}
			res.Link = createRes.Link
		} else {
			res.Link = one.PaymentLink
		}
	} else if one.Status == consts.InvoiceStatusPaid || one.Status == consts.InvoiceStatusReversed {
		res.Link = link.GetInvoicePdfLink(one.InvoiceId, one.SendTerms)
	}
	return res
}
