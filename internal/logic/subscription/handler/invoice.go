package handler

import (
	"context"
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

	if one == nil {
		//创建
		one := &entity.Invoice{
			MerchantId:         merchantId,
			SubscriptionId:     subscriptionId,
			InvoiceId:          utility.CreateInvoiceOrderNo(),
			TotalAmount:        details.TotalAmount,
			TaxAmount:          details.TaxAmount,
			SubscriptionAmount: details.SubscriptionAmount,
			Currency:           details.Currency,
			Lines:              "",
			ChannelId:          channelId,
			Status:             int(details.Status),
			SendStatus:         0,
			UserId:             userId,
			Data:               utility.FormatToJsonString(details),
			Link:               details.Link,
			ChannelStatus:      details.ChannelStatus,
			ChannelInvoiceId:   details.ChannelInvoiceId,
		}

		result, err := dao.Invoice.Ctx(ctx).Data(one).OmitEmpty().Insert(one)
		if err != nil {
			err = gerror.Newf(`CreateOrUpdateInvoiceByDetail record insert failure %s`, err)
			return err
		}
		id, _ := result.LastInsertId()
		one.Id = uint64(uint(id))

	} else {
		//更新
		update, err := dao.Invoice.Ctx(ctx).Data(g.Map{
			dao.Invoice.Columns().MerchantId:         merchantId,
			dao.Invoice.Columns().SubscriptionId:     subscriptionId,
			dao.Invoice.Columns().ChannelId:          channelId,
			dao.Invoice.Columns().TotalAmount:        details.TotalAmount,
			dao.Invoice.Columns().TaxAmount:          details.TaxAmount,
			dao.Invoice.Columns().SubscriptionAmount: details.SubscriptionAmount,
			dao.Invoice.Columns().Currency:           details.Currency,
			dao.Invoice.Columns().Status:             details.Status,
			dao.Invoice.Columns().Lines:              "",
			dao.Invoice.Columns().ChannelStatus:      details.ChannelStatus,
			dao.Invoice.Columns().ChannelInvoiceId:   details.ChannelInvoiceId,
			dao.Invoice.Columns().SubscriptionId:     subscriptionId,
			dao.Invoice.Columns().Link:               details.Link,
			dao.Invoice.Columns().Data:               utility.FormatToJsonString(details),
			dao.Invoice.Columns().GmtModify:          gtime.Now(),
		}).Where(dao.Invoice.Columns().Id, one.Id).OmitEmpty().Update()
		if err != nil {
			return err
		}
		rowAffected, err := update.RowsAffected()
		if rowAffected != 1 {
			return gerror.Newf("CreateOrUpdateInvoiceByDetail err:%s", update)
		}
	}
	// 异步处理发送邮件事件 todo mark
	return nil
}
