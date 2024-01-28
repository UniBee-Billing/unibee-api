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
	"go-oversea-pay/internal/logic/channel/ro"
	"go-oversea-pay/internal/logic/email"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strconv"
)

//
//func HandleInvoiceWebhookEvent(ctx context.Context, eventType string, details *ro.ChannelDetailInvoiceInternalResp) error {
//	//return CreateOrUpdateInvoiceByChannelDetail(ctx, details, details.ChannelInvoiceId)
//	// Not Generate Invoice Hook From Invoice Hook，Do it UniBee Self
//	return nil
//}

//func CreateOrUpdateInvoiceByChannelDetail(ctx context.Context, details *ro.ChannelDetailInvoiceInternalResp, uniqueId string) error {
//	utility.Assert(len(details.ChannelInvoiceId) > 0, "invoice id is null")
//	utility.Assert(len(details.ChannelSubscriptionId) > 0, "channelSubscriptionId invalid")
//	sub := query.GetSubscriptionByChannelSubscriptionId(ctx, details.ChannelSubscriptionId)
//	utility.Assert(sub != nil, "subscription of invoice not found ")
//	one := query.GetInvoiceByChannelInvoiceId(ctx, details.ChannelInvoiceId)
//	var invoiceId string
//	var change = false
//	if one == nil {
//		//Create
//		one = &entity.Invoice{
//			MerchantId:                     sub.MerchantId,
//			SubscriptionId:                 sub.SubscriptionId,
//			InvoiceId:                      utility.CreateInvoiceId(),
//			TotalAmount:                    details.TotalAmount,
//			TotalAmountExcludingTax:        details.TotalAmountExcludingTax,
//			TaxAmount:                      details.TaxAmount,
//			SubscriptionAmount:             details.SubscriptionAmount,
//			SubscriptionAmountExcludingTax: details.SubscriptionAmountExcludingTax,
//			PeriodStart:                    details.PeriodStart,
//			PeriodEnd:                      details.PeriodEnd,
//			PeriodStartTime:                gtime.NewFromTimeStamp(details.PeriodStart),
//			PeriodEndTime:                  gtime.NewFromTimeStamp(details.PeriodEnd),
//			Currency:                       details.Currency,
//			Lines:                          utility.MarshalToJsonString(details.Lines),
//			ChannelId:                      sub.ChannelId,
//			Status:                         int(details.Status),
//			SendStatus:                     0,
//			SendEmail:                      sub.CustomerEmail,
//			UserId:                         sub.UserId,
//			Data:                           utility.MarshalToJsonString(details),
//			Link:                           details.Link,
//			ChannelUserId:                  details.ChannelUserId,
//			ChannelStatus:                  details.ChannelStatus,
//			ChannelInvoiceId:               details.ChannelInvoiceId,
//			ChannelInvoicePdf:              details.ChannelInvoicePdf,
//			ChannelPaymentIntentId:               details.ChannelPaymentIntentId,
//			UniqueId:                       uniqueId,
//		}
//
//		result, err := dao.Invoice.Ctx(ctx).Data(one).OmitNil().Insert(one)
//		if err != nil {
//			err = gerror.Newf(`CreateOrUpdateInvoiceByChannelDetail record insert failure %s`, err.Error())
//			return err
//		}
//		id, _ := result.LastInsertId()
//		one.Id = uint64(uint(id))
//		invoiceId = one.InvoiceId
//		change = true
//	} else {
//		//Update
//		if one.Status != int(details.Status) {
//			change = true
//		}
//		_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
//			dao.Invoice.Columns().MerchantId:                     sub.MerchantId,
//			dao.Invoice.Columns().SubscriptionId:                 sub.SubscriptionId,
//			dao.Invoice.Columns().ChannelId:                      sub.ChannelId,
//			dao.Invoice.Columns().TotalAmount:                    details.TotalAmount,
//			dao.Invoice.Columns().TotalAmountExcludingTax:        details.TotalAmountExcludingTax,
//			dao.Invoice.Columns().TaxAmount:                      details.TaxAmount,
//			dao.Invoice.Columns().SubscriptionAmount:             details.SubscriptionAmount,
//			dao.Invoice.Columns().SubscriptionAmountExcludingTax: details.SubscriptionAmountExcludingTax,
//			dao.Invoice.Columns().PeriodStart:                    details.PeriodStart,
//			dao.Invoice.Columns().PeriodEnd:                      details.PeriodEnd,
//			dao.Invoice.Columns().PeriodStartTime:                gtime.NewFromTimeStamp(details.PeriodStart),
//			dao.Invoice.Columns().PeriodEndTime:                  gtime.NewFromTimeStamp(details.PeriodEnd),
//			dao.Invoice.Columns().Currency:                       details.Currency,
//			dao.Invoice.Columns().Status:                         details.Status,
//			dao.Invoice.Columns().Lines:                          utility.FormatToJsonString(details.Lines),
//			dao.Invoice.Columns().ChannelStatus:                  details.ChannelStatus,
//			dao.Invoice.Columns().ChannelInvoiceId:               details.ChannelInvoiceId,
//			dao.Invoice.Columns().ChannelUserId:                  details.ChannelUserId,
//			dao.Invoice.Columns().ChannelInvoicePdf:              details.ChannelInvoicePdf,
//			dao.Invoice.Columns().Link:                           details.Link,
//			dao.Invoice.Columns().SendEmail:                      sub.CustomerEmail,
//			dao.Invoice.Columns().Data:                           utility.FormatToJsonString(details),
//			dao.Invoice.Columns().GmtModify:                      gtime.Now(),
//			dao.Invoice.Columns().ChannelPaymentIntentId:               details.ChannelPaymentIntentId,
//		}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
//		if err != nil {
//			return err
//		}
//		//rowAffected, err := update.RowsAffected()
//		//if rowAffected != 1 {
//		//	return gerror.Newf("CreateOrUpdateInvoiceByChannelDetail err:%s", update)
//		//}
//		invoiceId = one.InvoiceId
//	}
//
//	if change {
//		// Send Email Any State Change After Stage
//		//  1-pending｜2-processing｜3-paid | 4-failed | 5-cancelled
//		// Send Invoice Email
//		// 1、Invoice Change Pending->Processing WaitingForPay
//		// 2、Invoice Change Pending->Paid (Automatic）
//		// 3、Invoice Change Processing->Paid（Manual）
//		// 4、Invoice Change Processing->Cancelled（Cancel Manual）
//		// 5、Invoice Change Processing->Failed (Payment Failure）
//		_ = SubscriptionInvoicePdfGenerateAndEmailSendBackground(invoiceId, true)
//	}
//
//	return nil
//}

func CreateInvoiceFromSubscriptionPaymentFailure(ctx context.Context, subscriptionId string, payment *entity.Payment, channelDetailInvoiceInternalResp *ro.ChannelDetailInvoiceInternalResp) error {
	//https://stripe.com/docs/billing/subscriptions/overview#requires-payment-method
	return nil
}

type CreateInvoiceInternalReq struct {
	Payment                          *entity.Payment                      `json:"payment"`
	ChannelInvoiceId                 string                               `json:"channelInvoiceId"`
	Currency                         string                               `json:"currency"`
	PlanId                           int64                                `json:"planId"`
	Quantity                         int64                                `json:"quantity"`
	AddonJsonData                    string                               `json:"addonJsonData"`
	TaxScale                         int64                                `json:"taxScale"`
	UserId                           int64                                `json:"userId"`
	MerchantId                       int64                                `json:"merchantId"`
	SubscriptionId                   string                               `json:"subscriptionId"`
	ChannelId                        int64                                `json:"channelId"`
	InvoiceStatus                    int                                  `json:"invoiceStatus"`
	ChannelDetailInvoiceInternalResp *ro.ChannelDetailInvoiceInternalResp `json:"channelDetailInvoiceInternalResp"`
	PeriodStart                      int64                                `json:"periodStart"                    description:"period_start，发票项目被添加到此发票的使用期限开始。，并非发票对应 sub 的周期"` // period_start，发票项目被添加到此发票的使用期限开始。，并非发票对应 sub 的周期
	PeriodEnd                        int64                                `json:"periodEnd"                      description:"period_end"`                                      // period_end
}

func UpdateInvoiceFromPayment(ctx context.Context, payment *entity.Payment, channelDetailInvoiceInternalResp *ro.ChannelDetailInvoiceInternalResp) error {
	utility.Assert(payment != nil, "payment data is nil")
	one := query.GetInvoiceByPaymentId(ctx, payment.PaymentId)
	utility.Assert(one != nil, "invoice not found, paymentId:"+payment.PaymentId)
	var status = consts.InvoiceStatusProcessing
	if payment.Status == consts.PAY_SUCCESS {
		status = consts.InvoiceStatusPaid
	} else if payment.Status == consts.PAY_FAILED {
		status = consts.InvoiceStatusFailed
	}
	var channelInvoicePdf = ""
	var channelInvoiceStatus = ""
	var channelLink = ""
	var channelInvoiceId = ""
	if channelDetailInvoiceInternalResp != nil {
		channelInvoiceId = channelDetailInvoiceInternalResp.ChannelInvoiceId
		channelInvoicePdf = channelDetailInvoiceInternalResp.ChannelInvoicePdf
		channelInvoiceStatus = channelDetailInvoiceInternalResp.ChannelStatus
		channelLink = channelDetailInvoiceInternalResp.Link
	} else {
		channelInvoiceId = one.ChannelInvoiceId
		channelInvoicePdf = one.ChannelInvoicePdf
		channelInvoiceStatus = one.ChannelStatus
		channelLink = one.Link
	}
	utility.Assert(one != nil, "invoice not found")
	_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
		dao.Invoice.Columns().Status:            status,
		dao.Invoice.Columns().ChannelInvoiceId:  channelInvoiceId,
		dao.Invoice.Columns().GmtModify:         gtime.Now(),
		dao.Invoice.Columns().ChannelPaymentId:  payment.ChannelPaymentId,
		dao.Invoice.Columns().Link:              channelLink,
		dao.Invoice.Columns().ChannelStatus:     channelInvoiceStatus,
		dao.Invoice.Columns().ChannelInvoicePdf: channelInvoicePdf,
	}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return err
	}
	if one.Status != status {
		//更新状态发送邮件
		_ = SubscriptionInvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true)
	}
	return nil
}

func UpdatePaymentInvoiceId(ctx context.Context, paymentId string, invoiceId string) error {
	_, err := dao.Payment.Ctx(ctx).Data(g.Map{
		dao.Payment.Columns().InvoiceId: invoiceId,
		dao.Invoice.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Payment.Columns().PaymentId, paymentId).OmitNil().Update()
	if err != nil {
		return err
	}
	return nil
}

func CreateOrUpdateInvoiceFromPayment(ctx context.Context, simplify *ro.InvoiceDetailSimplify, payment *entity.Payment, channelDetailInvoiceInternalResp *ro.ChannelDetailInvoiceInternalResp) error {
	utility.Assert(simplify != nil, "invoice data is nil")
	utility.Assert(payment != nil, "payment data is nil")
	one := query.GetInvoiceByPaymentId(ctx, payment.PaymentId)
	user := query.GetUserAccountById(ctx, uint64(payment.UserId))
	var channelInvoicePdf = ""
	var channelInvoiceId = ""
	var channelInvoiceStatus = ""
	var channelLink = ""
	var sendEmail = ""
	if channelDetailInvoiceInternalResp != nil {
		channelInvoiceId = channelDetailInvoiceInternalResp.ChannelInvoiceId
		channelInvoicePdf = channelDetailInvoiceInternalResp.ChannelInvoicePdf
		channelInvoiceStatus = channelDetailInvoiceInternalResp.ChannelStatus
		channelLink = channelDetailInvoiceInternalResp.Link
	}
	if user != nil {
		sendEmail = user.Email
	} else if one != nil && len(one.SendEmail) > 0 {
		sendEmail = one.SendEmail
	}
	var status = consts.InvoiceStatusProcessing
	if payment.Status == consts.PAY_SUCCESS {
		status = consts.InvoiceStatusPaid
	} else if payment.Status == consts.PAY_FAILED {
		status = consts.InvoiceStatusFailed
	}
	if one == nil {
		//创建
		one = &entity.Invoice{
			BizType:                        payment.BizType,
			UserId:                         payment.UserId,
			MerchantId:                     payment.MerchantId,
			SubscriptionId:                 payment.SubscriptionId,
			InvoiceId:                      utility.CreateInvoiceId(),
			PeriodStart:                    simplify.PeriodStart,
			PeriodEnd:                      simplify.PeriodEnd,
			PeriodStartTime:                gtime.NewFromTimeStamp(simplify.PeriodStart),
			PeriodEndTime:                  gtime.NewFromTimeStamp(simplify.PeriodEnd),
			Currency:                       payment.Currency,
			ChannelId:                      payment.ChannelId,
			Status:                         status,
			SendStatus:                     0,
			SendEmail:                      sendEmail,
			ChannelInvoiceId:               channelInvoiceId,
			ChannelPaymentId:               payment.ChannelPaymentId,
			UniqueId:                       payment.PaymentId,
			PaymentId:                      payment.PaymentId,
			TotalAmount:                    simplify.TotalAmount,
			TotalAmountExcludingTax:        simplify.TotalAmountExcludingTax,
			TaxAmount:                      simplify.TaxAmount,
			SubscriptionAmount:             simplify.SubscriptionAmount,
			SubscriptionAmountExcludingTax: simplify.SubscriptionAmountExcludingTax,
			Lines:                          utility.MarshalToJsonString(simplify.Lines),
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
		err = UpdatePaymentInvoiceId(ctx, payment.PaymentId, one.InvoiceId)
		if err != nil {
			return err
		}
	} else {
		//更新
		_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
			dao.Invoice.Columns().BizType:                        payment.BizType,
			dao.Invoice.Columns().MerchantId:                     payment.MerchantId,
			dao.Invoice.Columns().UserId:                         payment.UserId,
			dao.Invoice.Columns().SubscriptionId:                 payment.SubscriptionId,
			dao.Invoice.Columns().ChannelId:                      payment.ChannelId,
			dao.Invoice.Columns().PaymentId:                      payment.PaymentId,
			dao.Invoice.Columns().UniqueId:                       payment.PaymentId,
			dao.Invoice.Columns().Currency:                       payment.Currency,
			dao.Invoice.Columns().Status:                         status,
			dao.Invoice.Columns().ChannelInvoiceId:               channelInvoiceId,
			dao.Invoice.Columns().SendEmail:                      sendEmail,
			dao.Invoice.Columns().GmtModify:                      gtime.Now(),
			dao.Invoice.Columns().ChannelPaymentId:               payment.ChannelPaymentId,
			dao.Invoice.Columns().TotalAmount:                    simplify.TotalAmount,
			dao.Invoice.Columns().TotalAmountExcludingTax:        simplify.TotalAmountExcludingTax,
			dao.Invoice.Columns().TaxAmount:                      simplify.TaxAmount,
			dao.Invoice.Columns().SubscriptionAmount:             simplify.SubscriptionAmount,
			dao.Invoice.Columns().SubscriptionAmountExcludingTax: simplify.SubscriptionAmountExcludingTax,
			dao.Invoice.Columns().Lines:                          utility.FormatToJsonString(simplify.Lines),
			dao.Invoice.Columns().Link:                           channelLink,
			dao.Invoice.Columns().ChannelStatus:                  channelInvoiceStatus,
			dao.Invoice.Columns().ChannelInvoicePdf:              channelInvoicePdf,
		}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			return err
		}
		if one.Status != status {
			//更新状态发送邮件
			_ = SubscriptionInvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true)
		}
	}
	return nil
}

//func CreateNewSubscriptionBillingCycleInvoiceForPayment(ctx context.Context, req *CreateInvoiceInternalReq) error {
//	one := query.GetInvoiceByPaymentId(ctx, req.Payment.PaymentId)
//	invoice := invoice_compute.ConvertInvoiceDetailSimplifyFromPlanAddons(ctx, &invoice_compute.CalculateInvoiceReq{
//		Currency:      req.Currency,
//		PlanId:        req.PlanId,
//		Quantity:      req.Quantity,
//		AddonJsonData: req.AddonJsonData,
//		TaxScale:      req.TaxScale,
//	})
//	var channelInvoicePdf = ""
//	var channelInvoiceStatus = ""
//	var channelLink = ""
//	var sendEmail = ""
//	if req.ChannelDetailInvoiceInternalResp != nil {
//		channelInvoicePdf = req.ChannelDetailInvoiceInternalResp.ChannelInvoicePdf
//		channelInvoiceStatus = req.ChannelDetailInvoiceInternalResp.ChannelStatus
//		channelLink = req.ChannelDetailInvoiceInternalResp.Link
//	}
//	user := query.GetUserAccountById(ctx, uint64(req.UserId))
//	if user != nil {
//		sendEmail = user.Email
//	} else if one != nil && len(one.SendEmail) > 0 {
//		sendEmail = one.SendEmail
//	}
//	if one == nil {
//		//创建
//		one = &entity.Invoice{
//			BizType:                        req.Payment.BizType,
//			UserId:                         req.UserId,
//			MerchantId:                     req.MerchantId,
//			SubscriptionId:                 req.SubscriptionId,
//			InvoiceId:                      utility.CreateInvoiceId(),
//			PeriodStart:                    req.PeriodStart,
//			PeriodEnd:                      req.PeriodEnd,
//			PeriodStartTime:                gtime.NewFromTimeStamp(req.PeriodStart),
//			PeriodEndTime:                  gtime.NewFromTimeStamp(req.PeriodEnd),
//			Currency:                       req.Currency,
//			ChannelId:                      req.ChannelId,
//			Status:                         req.InvoiceStatus,
//			SendStatus:                     0,
//			SendEmail:                      sendEmail,
//			ChannelInvoiceId:               req.ChannelInvoiceId,
//			UniqueId:                       req.Payment.PaymentId,
//			PaymentId:                      req.Payment.PaymentId,
//			TotalAmount:                    invoice.TotalAmount,
//			TotalAmountExcludingTax:        invoice.TotalAmountExcludingTax,
//			TaxAmount:                      invoice.TaxAmount,
//			SubscriptionAmount:             invoice.SubscriptionAmount,
//			SubscriptionAmountExcludingTax: invoice.SubscriptionAmountExcludingTax,
//			Lines:                          utility.MarshalToJsonString(invoice.Lines),
//			Link:                           channelLink,
//			ChannelStatus:                  channelInvoiceStatus,
//			ChannelInvoicePdf:              channelInvoicePdf,
//		}
//
//		result, err := dao.Invoice.Ctx(ctx).Data(one).OmitNil().Insert(one)
//		if err != nil {
//			err = gerror.Newf(`CreateOrUpdateInvoiceByChannelDetail record insert failure %s`, err.Error())
//			return err
//		}
//		id, _ := result.LastInsertId()
//		one.Id = uint64(uint(id))
//		//新建 Invoice 发送邮件
//		_ = SubscriptionInvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true)
//	}
//	return nil
//}

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
		if one.BizType != consts.BIZ_TYPE_SUBSCRIPTION && len(one.Lines) == 0 {
			// invoice with subscription type and valid lines will send emails
			return
		}
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
	utility.Assert(one.UserId > 0, "invoice userId not found")
	utility.Assert(one.MerchantId > 0, "invoice merchantId not found")
	utility.Assert(len(one.SendEmail) > 0, "SendEmail Is Nil, InvoiceId:"+one.InvoiceId)
	utility.Assert(len(one.SendPdf) > 0, "pdf not generate is nil")
	user := query.GetUserAccountById(ctx, uint64(one.UserId))
	merchant := query.GetMerchantInfoById(ctx, one.MerchantId)
	var merchantProductName = ""
	sub := query.GetSubscriptionBySubscriptionId(ctx, one.SubscriptionId)
	if sub != nil {
		plan := query.GetPlanById(ctx, sub.PlanId)
		merchantProductName = plan.PlanName
	}
	if one.Status > consts.InvoiceStatusPending {
		pdfFileName := utility.DownloadFile(one.SendPdf)
		utility.Assert(len(pdfFileName) > 0, "download pdf error:"+one.SendPdf)
		//err := email.SendPdfAttachEmailToUser(one.SendEmail, "Invoice", "Invoice", pdfFileName, fmt.Sprintf("%s.pdf", one.InvoiceId))
		err := email.SendTemplateEmail(ctx, merchant.Id, one.SendEmail, email.TemplateInvoiceAutomaticPaid, pdfFileName, &email.TemplateVariable{
			InvoiceId:           one.InvoiceId,
			UserName:            user.UserName,
			MerchantProductName: merchantProductName,
			MerchantCustomEmail: merchant.Email,
			MerchantName:        merchant.Name,
			DateNow:             gtime.Now().Layout(`2006-01-02`),
			PaymentAmount:       strconv.FormatInt(one.TotalAmount, 10),
			TokenExpireMinute:   strconv.FormatInt(consts.GetConfigInstance().Auth.Login.Expire/60, 10),
		})
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
