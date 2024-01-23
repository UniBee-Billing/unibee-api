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
	//return CreateOrUpdateInvoiceByChannelDetail(ctx, details, details.ChannelInvoiceId)
	// 不再由 Invoice Hook 触发产生发票，改由 Payment Hook 自行生成发票
	return nil
}

func CreateOrUpdateInvoiceByChannelDetail(ctx context.Context, details *ro.ChannelDetailInvoiceInternalResp, uniqueId string) error {
	utility.Assert(len(details.ChannelInvoiceId) > 0, "invoice id is null")
	utility.Assert(len(details.ChannelSubscriptionId) > 0, "channelSubscriptionId invalid")
	sub := query.GetSubscriptionByChannelSubscriptionId(ctx, details.ChannelSubscriptionId)
	utility.Assert(sub != nil, "subscription of invoice not found ")
	one := query.GetInvoiceByChannelInvoiceId(ctx, details.ChannelInvoiceId)
	var invoiceId string
	var change = false
	if one == nil {
		//创建
		one = &entity.Invoice{
			MerchantId:                     sub.MerchantId,
			SubscriptionId:                 sub.SubscriptionId,
			InvoiceId:                      utility.CreateInvoiceId(),
			TotalAmount:                    details.TotalAmount,
			TotalAmountExcludingTax:        details.TotalAmountExcludingTax,
			TaxAmount:                      details.TaxAmount,
			SubscriptionAmount:             details.SubscriptionAmount,
			SubscriptionAmountExcludingTax: details.SubscriptionAmountExcludingTax,
			PeriodStart:                    details.PeriodStart,
			PeriodEnd:                      details.PeriodEnd,
			PeriodStartTime:                gtime.NewFromTimeStamp(details.PeriodStart),
			PeriodEndTime:                  gtime.NewFromTimeStamp(details.PeriodEnd),
			Currency:                       details.Currency,
			Lines:                          utility.MarshalToJsonString(details.Lines),
			ChannelId:                      sub.ChannelId,
			Status:                         int(details.Status),
			SendStatus:                     0,
			SendEmail:                      sub.CustomerEmail,
			UserId:                         sub.UserId,
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
			err = gerror.Newf(`CreateOrUpdateInvoiceByChannelDetail record insert failure %s`, err.Error())
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
		_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
			dao.Invoice.Columns().MerchantId:                     sub.MerchantId,
			dao.Invoice.Columns().SubscriptionId:                 sub.SubscriptionId,
			dao.Invoice.Columns().ChannelId:                      sub.ChannelId,
			dao.Invoice.Columns().TotalAmount:                    details.TotalAmount,
			dao.Invoice.Columns().TotalAmountExcludingTax:        details.TotalAmountExcludingTax,
			dao.Invoice.Columns().TaxAmount:                      details.TaxAmount,
			dao.Invoice.Columns().SubscriptionAmount:             details.SubscriptionAmount,
			dao.Invoice.Columns().SubscriptionAmountExcludingTax: details.SubscriptionAmountExcludingTax,
			dao.Invoice.Columns().PeriodStart:                    details.PeriodStart,
			dao.Invoice.Columns().PeriodEnd:                      details.PeriodEnd,
			dao.Invoice.Columns().PeriodStartTime:                gtime.NewFromTimeStamp(details.PeriodStart),
			dao.Invoice.Columns().PeriodEndTime:                  gtime.NewFromTimeStamp(details.PeriodEnd),
			dao.Invoice.Columns().Currency:                       details.Currency,
			dao.Invoice.Columns().Status:                         details.Status,
			dao.Invoice.Columns().Lines:                          utility.FormatToJsonString(details.Lines),
			dao.Invoice.Columns().ChannelStatus:                  details.ChannelStatus,
			dao.Invoice.Columns().ChannelInvoiceId:               details.ChannelInvoiceId,
			dao.Invoice.Columns().ChannelUserId:                  details.ChannelUserId,
			dao.Invoice.Columns().ChannelInvoicePdf:              details.ChannelInvoicePdf,
			dao.Invoice.Columns().Link:                           details.Link,
			dao.Invoice.Columns().SendEmail:                      sub.CustomerEmail,
			dao.Invoice.Columns().Data:                           utility.FormatToJsonString(details),
			dao.Invoice.Columns().GmtModify:                      gtime.Now(),
			dao.Invoice.Columns().ChannelPaymentId:               details.ChannelPaymentId,
		}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			return err
		}
		//rowAffected, err := update.RowsAffected()
		//if rowAffected != 1 {
		//	return gerror.Newf("CreateOrUpdateInvoiceByChannelDetail err:%s", update)
		//}
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

func CreateInvoiceFromSubscriptionPaymentFailure(ctx context.Context, subscriptionId string, payment *entity.Payment, channelDetailInvoiceInternalResp *ro.ChannelDetailInvoiceInternalResp) error {
	//https://stripe.com/docs/billing/subscriptions/overview#requires-payment-method
	return nil
}

func CreateInvoiceFromSubscriptionPaymentSuccess(ctx context.Context, subscriptionId string, payment *entity.Payment, channelDetailInvoiceInternalResp *ro.ChannelDetailInvoiceInternalResp) error {
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	one := query.GetInvoiceByChannelUniqueId(ctx, payment.PaymentId)
	invoice := CalculateInvoiceForSubscription(ctx, sub)
	var channelInvoicePdf = ""
	var channelInvoiceStatus = ""
	var channelLink = ""
	if channelDetailInvoiceInternalResp != nil {
		channelInvoicePdf = channelDetailInvoiceInternalResp.ChannelInvoicePdf
		channelInvoiceStatus = channelDetailInvoiceInternalResp.ChannelStatus
		channelLink = channelDetailInvoiceInternalResp.Link
	}
	if one == nil {
		//创建
		one = &entity.Invoice{
			UserId:                         sub.UserId,
			MerchantId:                     sub.MerchantId,
			SubscriptionId:                 sub.SubscriptionId,
			InvoiceId:                      utility.CreateInvoiceId(),
			PeriodStart:                    sub.CurrentPeriodStart,
			PeriodEnd:                      sub.CurrentPeriodEnd,
			PeriodStartTime:                gtime.NewFromTimeStamp(sub.CurrentPeriodStart),
			PeriodEndTime:                  gtime.NewFromTimeStamp(sub.CurrentPeriodEnd),
			Currency:                       sub.Currency,
			ChannelId:                      sub.ChannelId,
			Status:                         consts.InvoiceStatusPaid,
			SendStatus:                     0,
			SendEmail:                      sub.CustomerEmail,
			ChannelUserId:                  sub.ChannelUserId,
			ChannelInvoiceId:               payment.ChannelInvoiceId,
			ChannelPaymentId:               payment.ChannelPaymentId,
			UniqueId:                       payment.PaymentId,
			TotalAmount:                    sub.Amount,
			TotalAmountExcludingTax:        invoice.TotalAmountExcludingTax,
			TaxAmount:                      invoice.TaxAmount,
			SubscriptionAmount:             invoice.SubscriptionAmount,
			SubscriptionAmountExcludingTax: invoice.SubscriptionAmountExcludingTax,
			Lines:                          utility.MarshalToJsonString(invoice.Lines),
			Link:                           channelLink,
			ChannelStatus:                  channelInvoiceStatus,
			ChannelInvoicePdf:              channelInvoicePdf,
		}

		result, err := dao.Invoice.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			err = gerror.Newf(`CreateOrUpdateInvoiceByChannelDetail record insert failure %s`, err.Error())
			return err
		}
		id, _ := result.LastInsertId()
		one.Id = uint64(uint(id))
		//新建 Invoice 发送邮件
		_ = SubscriptionInvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true)
	} else {
		//更新
		_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
			dao.Invoice.Columns().MerchantId:                     sub.MerchantId,
			dao.Invoice.Columns().UserId:                         sub.UserId,
			dao.Invoice.Columns().SubscriptionId:                 sub.SubscriptionId,
			dao.Invoice.Columns().ChannelId:                      sub.ChannelId,
			dao.Invoice.Columns().PeriodStart:                    sub.CurrentPeriodStart,
			dao.Invoice.Columns().PeriodEnd:                      sub.CancelAtPeriodEnd,
			dao.Invoice.Columns().PeriodStartTime:                gtime.NewFromTimeStamp(sub.CurrentPeriodStart),
			dao.Invoice.Columns().PeriodEndTime:                  gtime.NewFromTimeStamp(sub.CurrentPeriodEnd),
			dao.Invoice.Columns().Currency:                       sub.Currency,
			dao.Invoice.Columns().Status:                         consts.InvoiceStatusPaid,
			dao.Invoice.Columns().ChannelInvoiceId:               payment.ChannelInvoiceId,
			dao.Invoice.Columns().ChannelUserId:                  sub.ChannelUserId,
			dao.Invoice.Columns().SendEmail:                      sub.CustomerEmail,
			dao.Invoice.Columns().GmtModify:                      gtime.Now(),
			dao.Invoice.Columns().ChannelPaymentId:               payment.ChannelPaymentId,
			dao.Invoice.Columns().TotalAmount:                    sub.Amount,
			dao.Invoice.Columns().TotalAmountExcludingTax:        invoice.TotalAmountExcludingTax,
			dao.Invoice.Columns().TaxAmount:                      invoice.TaxAmount,
			dao.Invoice.Columns().SubscriptionAmount:             invoice.SubscriptionAmount,
			dao.Invoice.Columns().SubscriptionAmountExcludingTax: invoice.SubscriptionAmountExcludingTax,
			dao.Invoice.Columns().Lines:                          utility.FormatToJsonString(invoice.Lines),
			dao.Invoice.Columns().Link:                           channelLink,
			dao.Invoice.Columns().ChannelStatus:                  channelInvoiceStatus,
			dao.Invoice.Columns().ChannelInvoicePdf:              channelInvoicePdf,
		}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			return err
		}
		if one.Status != consts.InvoiceStatusPaid {
			//更新状态发送邮件
			_ = SubscriptionInvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true)
		}
	}
	return nil
}

func CalculateInvoiceForSubscription(ctx context.Context, sub *entity.Subscription) *ro.ChannelDetailInvoiceRo {
	plan := query.GetPlanById(ctx, sub.PlanId)
	addons := query.GetSubscriptionAddonsByAddonJson(ctx, sub.AddonData)
	var totalAmountExcludingTax = plan.Amount * sub.Quantity
	for _, addon := range addons {
		totalAmountExcludingTax = totalAmountExcludingTax + addon.AddonPlan.Amount*addon.Quantity
	}

	var invoiceItems []*ro.ChannelDetailInvoiceItem
	invoiceItems = append(invoiceItems, &ro.ChannelDetailInvoiceItem{
		Currency:               sub.Currency,
		Amount:                 sub.Quantity*plan.Amount + int64(float64(sub.Quantity*plan.Amount)*utility.ConvertTaxPercentageToInternalFloat(sub.TaxPercentage)),
		AmountExcludingTax:     sub.Quantity * plan.Amount,
		Tax:                    int64(float64(sub.Quantity*plan.Amount) * utility.ConvertTaxPercentageToInternalFloat(sub.TaxPercentage)),
		UnitAmountExcludingTax: plan.Amount,
		Description:            plan.PlanName,
		Quantity:               sub.Quantity,
	})
	for _, addon := range addons {
		invoiceItems = append(invoiceItems, &ro.ChannelDetailInvoiceItem{
			Currency:               sub.Currency,
			Amount:                 addon.Quantity*addon.AddonPlan.Amount + int64(float64(addon.Quantity*addon.AddonPlan.Amount)*utility.ConvertTaxPercentageToInternalFloat(sub.TaxPercentage)),
			Tax:                    int64(float64(addon.Quantity*addon.AddonPlan.Amount) * utility.ConvertTaxPercentageToInternalFloat(sub.TaxPercentage)),
			AmountExcludingTax:     addon.Quantity * addon.AddonPlan.Amount,
			UnitAmountExcludingTax: addon.AddonPlan.Amount,
			Description:            addon.AddonPlan.PlanName,
			Quantity:               addon.Quantity,
		})
	}
	var taxAmount = int64(float64(totalAmountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(sub.TaxPercentage))
	return &ro.ChannelDetailInvoiceRo{
		TotalAmount:                    totalAmountExcludingTax + taxAmount,
		TotalAmountExcludingTax:        totalAmountExcludingTax,
		Currency:                       sub.Currency,
		TaxAmount:                      taxAmount,
		SubscriptionAmount:             totalAmountExcludingTax + taxAmount, // 在没有 discount 之前，保持于 Total 一致
		SubscriptionAmountExcludingTax: totalAmountExcludingTax,             // 在没有 discount 之前，保持于 Total 一致
		Lines:                          invoiceItems,
	}
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
				g.Log().Errorf(context.Background(), "CreateOrUpdateInvoiceByChannelDetail Background Generate PDF panic error:%s\n", err.Error())
				return
			}
		}()
		backgroundCtx := context.Background()
		one := query.GetInvoiceByInvoiceId(backgroundCtx, invoiceId)
		utility.Assert(one != nil, "invoice not found")
		url := GenerateAndUploadInvoicePdf(backgroundCtx, one)
		if len(url) > 0 {
			_, err = dao.Invoice.Ctx(backgroundCtx).Data(g.Map{
				dao.Invoice.Columns().SendPdf:   url,
				dao.Invoice.Columns().GmtModify: gtime.Now(),
			}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
			if err != nil {
				fmt.Printf("GenerateAndUploadInvoicePdf update err:%s", err.Error())
			}
			//rowAffected, err := update.RowsAffected()
			//if rowAffected != 1 {
			//	fmt.Printf("GenerateAndUploadInvoicePdf update err:%s", update)
			//}
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
		_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
			dao.Invoice.Columns().SendStatus: 1,
			dao.Invoice.Columns().GmtModify:  gtime.Now(),
		}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			fmt.Printf("SendInvoiceEmailToUser update err:%s", err.Error())
		}
		//rowAffected, err := update.RowsAffected()
		//if rowAffected != 1 {
		//	fmt.Printf("SendInvoiceEmailToUser update err:%s", update)
		//}
	} else {
		fmt.Printf("SendInvoiceEmailToUser invoice status is pending or init, email not send")
	}
	return nil

}
