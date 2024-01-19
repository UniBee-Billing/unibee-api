package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/email"
	"go-oversea-pay/internal/logic/gateway/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func HandleInvoiceWebhookEvent(ctx context.Context, eventType string, details *ro.ChannelDetailInvoiceInternalResp) error {
	return CreateOrUpdateInvoiceByDetail(ctx, details, details.ChannelInvoiceId)
}

func CreateOrUpdateInvoiceByDetail(ctx context.Context, details *ro.ChannelDetailInvoiceInternalResp, uniqueId string) error {
	utility.Assert(len(details.ChannelInvoiceId) > 0, "invoice id is null")
	var subscriptionId string
	var merchantId int64
	var channelId int64
	var userId int64
	var sendEmail string
	if len(details.ChannelSubscriptionId) > 0 {
		sub := query.GetSubscriptionByChannelSubscriptionId(ctx, details.ChannelSubscriptionId)
		if sub != nil {
			subscriptionId = sub.SubscriptionId
			merchantId = sub.MerchantId
			channelId = sub.ChannelId
			userId = sub.UserId
			merchantInfo := query.GetMerchantInfoById(ctx, sub.MerchantId)
			if merchantInfo != nil {
				sendEmail = merchantInfo.Email
			}
		}
	}
	one := query.GetInvoiceByChannelInvoiceId(ctx, details.ChannelInvoiceId)

	var invoiceId string
	var change = false
	if one == nil {
		//创建
		one := &entity.Invoice{
			MerchantId:                     merchantId,
			SubscriptionId:                 subscriptionId,
			InvoiceId:                      utility.CreateInvoiceId(),
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
			SendEmail:                      sendEmail,
			UserId:                         userId,
			Data:                           utility.MarshalToJsonString(details),
			Link:                           details.Link,
			ChannelUserId:                  details.ChannelUserId,
			ChannelStatus:                  details.ChannelStatus,
			ChannelInvoiceId:               details.ChannelInvoiceId,
			ChannelInvoicePdf:              details.ChannelInvoicePdf,
			ChannelPaymentId:               details.ChannelPaymentId,
			UniqueId:                       uniqueId,
		}

		result, err := dao.Invoice.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			err = gerror.Newf(`CreateOrUpdateInvoiceByDetail record insert failure %s`, err)
			return err
		}
		id, _ := result.LastInsertId()
		one.Id = uint64(uint(id))
		invoiceId = one.InvoiceId
		change = true
	} else {
		//更新
		if one.Status != int(details.Status) {
			change = true
		}
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
			dao.Invoice.Columns().ChannelUserId:                  details.ChannelUserId,
			dao.Invoice.Columns().ChannelInvoicePdf:              details.ChannelInvoicePdf,
			dao.Invoice.Columns().Link:                           details.Link,
			dao.Invoice.Columns().SendEmail:                      sendEmail,
			dao.Invoice.Columns().Data:                           utility.FormatToJsonString(details),
			dao.Invoice.Columns().GmtModify:                      gtime.Now(),
			dao.Invoice.Columns().ChannelPaymentId:               details.ChannelPaymentId,
		}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			return err
		}
		rowAffected, err := update.RowsAffected()
		if rowAffected != 1 {
			return gerror.Newf("CreateOrUpdateInvoiceByDetail err:%s", update)
		}
		invoiceId = one.InvoiceId
	}

	if change {
		//脱离草稿状态每次变化都生成并发送邮件
		//  1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled
		// 发送 Invoice  Email 节点
		// 1、Invoice 状态从 Pending->Processing 等待用户支付
		// 2、Invoice 状态从 Pending->Paid (自动支付）
		// 3、Invoice 状态从 Processing->Paid（手动支付）
		// 4、Invoice 状态从 Processing->Cancelled（手动取消）
		// 5、Invoice 状态从 Processing->Failed (支付超时->支付失败）
		_ = SubscriptionInvoicePdfGenerateAndEmailSendBackground(invoiceId, true)
	}

	return nil
}

func SubscriptionInvoicePdfGenerateAndEmailSendBackground(invoiceId string, sendUserEmail bool) (err error) {
	go func() {
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				g.Log().Errorf(context.Background(), "CreateOrUpdateInvoiceByDetail Background Generate PDF panic error:%s\n", err.Error())
				return
			}
		}()
		backgroundCtx := context.Background()
		one := query.GetInvoiceByInvoiceId(backgroundCtx, invoiceId)
		utility.Assert(one != nil, "invoice not found")
		url := GenerateAndUploadInvoicePdf(backgroundCtx, one)
		if len(url) > 0 {
			update, err := dao.Invoice.Ctx(backgroundCtx).Data(g.Map{
				dao.Invoice.Columns().SendPdf:   url,
				dao.Invoice.Columns().GmtModify: gtime.Now(),
			}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
			if err != nil {
				fmt.Printf("GenerateAndUploadInvoicePdf update err:%s", update)
			}
			rowAffected, err := update.RowsAffected()
			if rowAffected != 1 {
				fmt.Printf("GenerateAndUploadInvoicePdf update err:%s", update)
			}
		}
		if sendUserEmail {
			err := SendInvoiceEmailToUser(backgroundCtx, one.InvoiceId)
			utility.Assert(err == nil, "SendInvoiceEmail error")
		}
	}()
	return nil
}

func SendInvoiceEmailToUser(ctx context.Context, invoiceId string) error {
	one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
	utility.Assert(one != nil, "invoice not found")
	utility.Assert(len(one.SendEmail) > 0, "SendEmail is nil")
	utility.Assert(len(one.SendPdf) > 0, "pdf not generate is nil")
	if one.Status > consts.InvoiceStatusPending {
		pdfFileName := utility.DownloadFile(one.SendPdf)
		utility.Assert(len(pdfFileName) > 0, "download pdf error:"+one.SendPdf)
		err := email.SendPdfAttachEmailToUser(one.SendEmail, "Invoice", "Invoice", pdfFileName)
		if err != nil {
			return err
		}
		//修改发送状态
		update, err := dao.Invoice.Ctx(ctx).Data(g.Map{
			dao.Invoice.Columns().SendStatus: 1,
			dao.Invoice.Columns().GmtModify:  gtime.Now(),
		}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			fmt.Printf("SendInvoiceEmailToUser update err:%s", update)
		}
		rowAffected, err := update.RowsAffected()
		if rowAffected != 1 {
			fmt.Printf("SendInvoiceEmailToUser update err:%s", update)
		}
	} else {
		fmt.Printf("SendInvoiceEmailToUser invoice status is pending or init, email not send")
	}
	return nil

}
