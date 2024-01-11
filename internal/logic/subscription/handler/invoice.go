package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/payment/outchannel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func HandleInvoiceWebhookEvent(ctx context.Context, eventType string, details *ro.ChannelDetailInvoiceInternalResp) error {
	return CreateOrUpdateInvoiceByDetail(ctx, details)
}

func CreateOrUpdateInvoiceByDetail(ctx context.Context, details *ro.ChannelDetailInvoiceInternalResp) error {
	utility.Assert(len(details.ChannelInvoiceId) > 0, "invoice id is null")
	var subscriptionId string
	var merchantId int64
	var channelId int64
	var userId int64
	if len(details.ChannelSubscriptionId) > 0 {
		sub := query.GetSubscriptionByChannelSubscriptionId(ctx, details.ChannelSubscriptionId)
		if sub != nil {
			subscriptionId = sub.SubscriptionId
			merchantId = sub.MerchantId
			channelId = sub.ChannelId
			userId = sub.UserId
		}
	}
	one := query.GetInvoiceByChannelInvoiceId(ctx, details.ChannelInvoiceId)

	var invoiceId string
	if one == nil {
		//创建
		one := &entity.Invoice{
			MerchantId:                     merchantId,
			SubscriptionId:                 subscriptionId,
			InvoiceId:                      utility.CreateInvoiceOrderNo(),
			TotalAmount:                    details.TotalAmount,
			TotalAmountExcludingTax:        details.TotalAmountExcludingTax,
			TaxAmount:                      details.TaxAmount,
			SubscriptionAmount:             details.SubscriptionAmount,
			SubscriptionAmountExcludingTax: details.SubscriptionAmountExcludingTax,
			PeriodStart:                    details.PeriodStart,
			PeriodEnd:                      details.PeriodEnd,
			Currency:                       details.Currency,
			Lines:                          utility.MarshalToJsonString(details.Lines),
			ChannelId:                      channelId,
			Status:                         int(details.Status),
			SendStatus:                     0,
			UserId:                         userId,
			Data:                           utility.MarshalToJsonString(details),
			Link:                           details.Link,
			ChannelStatus:                  details.ChannelStatus,
			ChannelInvoiceId:               details.ChannelInvoiceId,
			ChannelInvoicePdf:              details.ChannelInvoicePdf,
		}

		result, err := dao.Invoice.Ctx(ctx).Data(one).OmitEmpty().Insert(one)
		if err != nil {
			err = gerror.Newf(`CreateOrUpdateInvoiceByDetail record insert failure %s`, err)
			return err
		}
		id, _ := result.LastInsertId()
		one.Id = uint64(uint(id))
		invoiceId = one.InvoiceId

	} else {
		//更新
		update, err := dao.Invoice.Ctx(ctx).Data(g.Map{
			dao.Invoice.Columns().MerchantId:                     merchantId,
			dao.Invoice.Columns().SubscriptionId:                 subscriptionId,
			dao.Invoice.Columns().ChannelId:                      channelId,
			dao.Invoice.Columns().TotalAmount:                    details.TotalAmount,
			dao.Invoice.Columns().TotalAmountExcludingTax:        details.TotalAmountExcludingTax,
			dao.Invoice.Columns().TaxAmount:                      details.TaxAmount,
			dao.Invoice.Columns().SubscriptionAmount:             details.SubscriptionAmount,
			dao.Invoice.Columns().SubscriptionAmountExcludingTax: details.SubscriptionAmountExcludingTax,
			dao.Invoice.Columns().PeriodStart:                    details.PeriodStart,
			dao.Invoice.Columns().PeriodEnd:                      details.PeriodEnd,
			dao.Invoice.Columns().Currency:                       details.Currency,
			dao.Invoice.Columns().Status:                         details.Status,
			dao.Invoice.Columns().Lines:                          utility.FormatToJsonString(details.Lines),
			dao.Invoice.Columns().ChannelStatus:                  details.ChannelStatus,
			dao.Invoice.Columns().ChannelInvoiceId:               details.ChannelInvoiceId,
			dao.Invoice.Columns().SubscriptionId:                 subscriptionId,
			dao.Invoice.Columns().Link:                           details.Link,
			dao.Invoice.Columns().Data:                           utility.FormatToJsonString(details),
			dao.Invoice.Columns().GmtModify:                      gtime.Now(),
		}).Where(dao.Invoice.Columns().Id, one.Id).OmitEmpty().Update()
		if err != nil {
			return err
		}
		rowAffected, err := update.RowsAffected()
		if rowAffected != 1 {
			return gerror.Newf("CreateOrUpdateInvoiceByDetail err:%s", update)
		}
		invoiceId = one.InvoiceId
	}

	_ = SubscriptionInvoicePdfGenerateBackground(invoiceId)

	return nil
}

func SubscriptionInvoicePdfGenerateBackground(invoiceId string) (err error) {
	go func() {
		defer func() {
			if exception := recover(); exception != nil {
				fmt.Printf("CreateOrUpdateInvoiceByDetail Background Generate PDF panic error:%s\n", exception)
				return
			}
		}()
		backgroundCtx := context.Background()
		one := query.GetInvoiceByInvoiceId(backgroundCtx, invoiceId)
		url := GenerateAndUploadInvoicePdf(backgroundCtx, one)
		if len(url) > 0 {
			update, err := dao.Invoice.Ctx(backgroundCtx).Data(g.Map{
				dao.Invoice.Columns().SendPdf:   url,
				dao.Invoice.Columns().GmtModify: gtime.Now(),
			}).Where(dao.Invoice.Columns().Id, one.Id).OmitEmpty().Update()
			if err != nil {
				fmt.Printf("GenerateAndUploadInvoicePdf update err:%s", update)
			}
			rowAffected, err := update.RowsAffected()
			if rowAffected != 1 {
				fmt.Printf("GenerateAndUploadInvoicePdf update err:%s", update)
			}
		}
		// 异步处理发送邮件事件 todo mark
	}()
	return nil
}
