package service

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"strconv"
	"unibee/api/bean"
	"unibee/internal/consts"
	"unibee/internal/controller/link"
	"unibee/internal/logic/gateway/gateway_bean"
	"unibee/internal/logic/invoice/handler"
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

			merchantInfo := query.GetMerchantById(ctx, one.MerchantId)
			user := query.GetUserAccountById(ctx, one.UserId)

			createRes, err := service.GatewayPaymentCreate(ctx, &gateway_bean.GatewayNewPaymentReq{
				CheckoutMode: true,
				Gateway:      gateway,
				Pay: &entity.Payment{
					SubscriptionId:    one.SubscriptionId,
					BizType:           one.BizType,
					ExternalPaymentId: one.SubscriptionId,
					UserId:            one.UserId,
					GatewayId:         gateway.Id,
					TotalAmount:       one.TotalAmount,
					Currency:          one.Currency,
					CryptoAmount:      one.CryptoAmount,
					CryptoCurrency:    one.CryptoCurrency,
					CountryCode:       user.CountryCode,
					MerchantId:        one.MerchantId,
					CompanyId:         merchantInfo.CompanyId,
					BillingReason:     one.InvoiceName,
					ReturnUrl:         "",
					//GasPayer:          one.GasPayer, // todo mark
				},
				ExternalUserId: strconv.FormatUint(one.UserId, 10),
				Email:          user.Email,
				Invoice:        bean.SimplifyInvoice(one),
				Metadata:       map[string]interface{}{"BillingReason": one.InvoiceName},
			})
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
		if len(one.SendPdf) > 0 {
			res.Link = link.GetInvoicePdfLink(one.InvoiceId, one.SendTerms)
		} else {
			res.FileName = handler.GenerateInvoicePdf(ctx, one)
		}
	}
	return res
}
