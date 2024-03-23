package invoice

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"strconv"
	"unibee/api/bean"
	"unibee/internal/consts"
	"unibee/internal/logic/gateway/gateway_bean"
	"unibee/internal/logic/payment/service"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func InvoiceLinkEntry(r *ghttp.Request) {
	invoiceId := r.Get("invoiceId").String()
	one := query.GetInvoiceByInvoiceId(r.Context(), invoiceId)
	if one == nil {
		g.Log().Errorf(r.Context(), "InvoiceLinkEntry invoice not found url: %s invoiceId: %s", r.GetUrl(), invoiceId)
		return
	}
	if one.IsDeleted > 0 {
		r.Response.Writeln("Invoice Deleted")
	} else if one.Status == consts.InvoiceStatusCancelled {
		r.Response.Writeln("Invoice Cancelled")
	} else if one.Status == consts.InvoiceStatusFailed {
		r.Response.Writeln("Invoice Failure")
	} else if one.Status < consts.InvoiceStatusProcessing {
		r.Response.Writeln("Invoice Not Ready")
	} else if one.Status == consts.InvoiceStatusProcessing {
		if len(one.PaymentLink) == 0 {
			// create payment link for this invoice
			gateway := query.GetGatewayById(r.Context(), one.GatewayId)
			if gateway == nil {
				r.Response.Writeln("Gateway Error")
				return
			}
			var lines []*bean.InvoiceItemSimplify
			err := utility.UnmarshalFromJsonString(one.Lines, &lines)
			if err != nil {
				r.Response.Writeln("Server Error")
				return
			}

			merchantInfo := query.GetMerchantById(r.Context(), one.MerchantId)
			user := query.GetUserAccountById(r.Context(), uint64(one.UserId))
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

			createRes, err := service.GatewayPaymentCreate(r.Context(), createPayContext)
			if err != nil {
				g.Log().Infof(r.Context(), "GatewayPaymentCreate Error:%s", err.Error())
				r.Response.Writeln("Server Error")
				return
			}
			r.Response.RedirectTo(createRes.Link)
		} else {
			r.Response.RedirectTo(one.PaymentLink)
		}
	} else if one.Status == consts.InvoiceStatusPaid {
		r.Response.RedirectTo(one.SendPdf)
	}
}
