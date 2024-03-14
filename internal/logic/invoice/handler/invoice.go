package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"unibee/api/bean"
	"unibee/internal/consts"
	"unibee/internal/controller/link"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/crypto"
	"unibee/internal/logic/email"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func UpdateInvoiceFromPayment(ctx context.Context, payment *entity.Payment) (*entity.Invoice, error) {
	utility.Assert(payment != nil, "payment data is nil")
	one := query.GetInvoiceByInvoiceId(ctx, payment.InvoiceId)
	if one == nil {
		return nil, gerror.New("invoice not found, paymentId:" + payment.PaymentId + " subId:" + payment.SubscriptionId)
	}
	var status = consts.InvoiceStatusProcessing
	if payment.Status == consts.PaymentSuccess {
		status = consts.InvoiceStatusPaid
	} else if payment.Status == consts.PaymentFailed {
		status = consts.InvoiceStatusFailed
	} else if payment.Status == consts.PaymentCancelled {
		status = consts.InvoiceStatusCancelled
	}
	_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
		dao.Invoice.Columns().Status:           status,
		dao.Invoice.Columns().GmtModify:        gtime.Now(),
		dao.Invoice.Columns().GatewayPaymentId: payment.GatewayPaymentId,
		dao.Invoice.Columns().PaymentLink:      payment.Link,
	}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	if one.Status != status && one.BizType == consts.BizTypeSubscription {
		_ = InvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true)
	}
	one.Status = status
	one.GatewayPaymentId = payment.GatewayPaymentId
	one.Link = payment.Link
	return one, nil
}

func CreateProcessingInvoiceForSub(ctx context.Context, simplify *bean.InvoiceSimplify, sub *entity.Subscription) (*entity.Invoice, error) {
	utility.Assert(simplify != nil, "invoice data is nil")
	utility.Assert(sub != nil, "sub is nil")
	user := query.GetUserAccountById(ctx, uint64(sub.UserId))
	var sendEmail = ""
	if user != nil {
		sendEmail = user.Email
	}

	//Create
	invoiceId := utility.CreateInvoiceId()
	one := &entity.Invoice{
		BizType:                        consts.BizTypeSubscription,
		UserId:                         sub.UserId,
		MerchantId:                     sub.MerchantId,
		SubscriptionId:                 sub.SubscriptionId,
		InvoiceName:                    simplify.InvoiceName,
		InvoiceId:                      invoiceId,
		PeriodStart:                    simplify.PeriodStart,
		PeriodEnd:                      simplify.PeriodEnd,
		PeriodStartTime:                gtime.NewFromTimeStamp(simplify.PeriodStart),
		PeriodEndTime:                  gtime.NewFromTimeStamp(simplify.PeriodEnd),
		Currency:                       sub.Currency,
		GatewayId:                      sub.GatewayId,
		Status:                         consts.InvoiceStatusProcessing,
		SendStatus:                     0,
		SendEmail:                      sendEmail,
		UniqueId:                       invoiceId,
		TotalAmount:                    simplify.TotalAmount,
		TotalAmountExcludingTax:        simplify.TotalAmountExcludingTax,
		TaxAmount:                      simplify.TaxAmount,
		SubscriptionAmount:             simplify.SubscriptionAmount,
		SubscriptionAmountExcludingTax: simplify.SubscriptionAmountExcludingTax,
		Lines:                          utility.MarshalToJsonString(simplify.Lines),
		Link:                           link.GetInvoiceLink(invoiceId),
		CreateTime:                     gtime.Now().Timestamp(),
	}

	result, err := dao.Invoice.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`CreateProcessingInvoiceForSub record insert failure %s`, err.Error())
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))
	_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().LatestInvoiceId: invoiceId,
	}).Where(dao.Subscription.Columns().SubscriptionId, sub.SubscriptionId).OmitNil().Update()
	if err != nil {
		utility.AssertError(err, "CreateProcessingInvoiceForSub")
	}
	//todo mark cancel other sub processing invoice
	//New Invoice Send Email
	_ = InvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true)
	if err != nil {
		return nil, err
	}
	return one, nil
}

func CreateOrUpdateInvoiceForNewPayment(ctx context.Context, invoice *bean.InvoiceSimplify, payment *entity.Payment) (*entity.Invoice, error) {
	utility.Assert(invoice != nil, "invoice data is nil")
	utility.Assert(payment != nil, "payment data is nil")
	one := query.GetInvoiceByPaymentId(ctx, payment.PaymentId)
	if one == nil && len(invoice.InvoiceId) > 0 {
		one = query.GetInvoiceByInvoiceId(ctx, invoice.InvoiceId)
	}
	user := query.GetUserAccountById(ctx, uint64(payment.UserId))
	var sendEmail = ""
	if user != nil {
		sendEmail = user.Email
	} else if one != nil && len(one.SendEmail) > 0 {
		sendEmail = one.SendEmail
	}
	var status = consts.InvoiceStatusProcessing
	if payment.Status == consts.PaymentSuccess {
		status = consts.InvoiceStatusPaid
	} else if payment.Status == consts.PaymentFailed {
		status = consts.InvoiceStatusFailed
	} else if payment.Status == consts.PaymentCancelled {
		status = consts.InvoiceStatusCancelled
	}
	if one == nil {
		//创建
		one = &entity.Invoice{
			BizType:                        payment.BizType,
			UserId:                         payment.UserId,
			MerchantId:                     payment.MerchantId,
			SubscriptionId:                 payment.SubscriptionId,
			InvoiceName:                    payment.BillingReason,
			InvoiceId:                      utility.CreateInvoiceId(),
			PeriodStart:                    invoice.PeriodStart,
			PeriodEnd:                      invoice.PeriodEnd,
			PeriodStartTime:                gtime.NewFromTimeStamp(invoice.PeriodStart),
			PeriodEndTime:                  gtime.NewFromTimeStamp(invoice.PeriodEnd),
			Currency:                       payment.Currency,
			GatewayId:                      payment.GatewayId,
			Status:                         status,
			SendStatus:                     0,
			SendEmail:                      sendEmail,
			GatewayPaymentId:               payment.GatewayPaymentId,
			UniqueId:                       payment.PaymentId,
			PaymentId:                      payment.PaymentId,
			TotalAmount:                    invoice.TotalAmount,
			TotalAmountExcludingTax:        invoice.TotalAmountExcludingTax,
			TaxAmount:                      invoice.TaxAmount,
			SubscriptionAmount:             invoice.SubscriptionAmount,
			SubscriptionAmountExcludingTax: invoice.SubscriptionAmountExcludingTax,
			Lines:                          utility.MarshalToJsonString(invoice.Lines),
			PaymentLink:                    payment.Link,
			CreateTime:                     gtime.Now().Timestamp(),
		}

		result, err := dao.Invoice.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			err = gerror.Newf(`CreateOrUpdateInvoiceByChannelDetail record insert failure %s`, err.Error())
			return nil, err
		}
		id, _ := result.LastInsertId()
		one.Id = uint64(uint(id))
		if one.BizType == consts.BizTypeSubscription {
			_ = InvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true)
		}
	} else {
		//Update
		_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
			dao.Invoice.Columns().BizType:                        payment.BizType,
			dao.Invoice.Columns().MerchantId:                     payment.MerchantId,
			dao.Invoice.Columns().UserId:                         payment.UserId,
			dao.Invoice.Columns().SubscriptionId:                 payment.SubscriptionId,
			dao.Invoice.Columns().GatewayId:                      payment.GatewayId,
			dao.Invoice.Columns().PaymentId:                      payment.PaymentId,
			dao.Invoice.Columns().UniqueId:                       payment.PaymentId,
			dao.Invoice.Columns().Status:                         status,
			dao.Invoice.Columns().GmtModify:                      gtime.Now(),
			dao.Invoice.Columns().GatewayPaymentId:               payment.GatewayPaymentId,
			dao.Invoice.Columns().TotalAmount:                    invoice.TotalAmount,
			dao.Invoice.Columns().TotalAmountExcludingTax:        invoice.TotalAmountExcludingTax,
			dao.Invoice.Columns().TaxAmount:                      invoice.TaxAmount,
			dao.Invoice.Columns().SubscriptionAmount:             invoice.SubscriptionAmount,
			dao.Invoice.Columns().SubscriptionAmountExcludingTax: invoice.SubscriptionAmountExcludingTax,
			dao.Invoice.Columns().Lines:                          utility.FormatToJsonString(invoice.Lines),
			dao.Invoice.Columns().PaymentLink:                    payment.Link,
		}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			return nil, err
		}
		if one.Status != status && one.BizType == consts.BizTypeSubscription {
			_ = InvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true)
		}
	}
	one = query.GetInvoiceByPaymentId(ctx, payment.PaymentId)
	return one, nil
}

func InvoicePdfGenerateAndEmailSendBackground(invoiceId string, sendUserEmail bool) (err error) {
	one := query.GetInvoiceByInvoiceId(context.Background(), invoiceId)
	utility.Assert(one != nil, "invoice not found")
	if len(one.Lines) == 0 {
		// invoice with valid lines will send emails
		return gerror.New("invalid lines")
	}
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
			err := SendSubscriptionInvoiceEmailToUser(backgroundCtx, one.InvoiceId)
			utility.Assert(err == nil, "SendInvoiceEmail error")
		}
	}()
	return nil
}

func ReconvertCryptoDataForInvoice(ctx context.Context, invoiceId string) error {
	one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
	utility.Assert(one != nil, "invoice not found")
	utility.Assert(one.UserId > 0, "invoice userId not found")
	utility.Assert(one.MerchantId > 0, "invoice merchantId not found")
	_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
		dao.Invoice.Columns().CryptoCurrency: crypto.GetCryptoCurrency(),
		dao.Invoice.Columns().CryptoAmount:   crypto.GetCryptoAmount(one.TotalAmount, one.TaxAmount),
	}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		fmt.Printf("ReconvertCryptoDataForInvoice update err:%s", err.Error())
	}
	return err
}

func SendSubscriptionInvoiceEmailToUser(ctx context.Context, invoiceId string) error {
	one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
	utility.Assert(one != nil, "invoice not found")
	utility.Assert(one.UserId > 0, "invoice userId not found")
	utility.Assert(one.MerchantId > 0, "invoice merchantId not found")
	utility.Assert(len(one.SendEmail) > 0, "SendEmail Is Nil, InvoiceId:"+one.InvoiceId)
	utility.Assert(len(one.SendPdf) > 0, "pdf not generate is nil")
	user := query.GetUserAccountById(ctx, uint64(one.UserId))
	merchant := query.GetMerchantById(ctx, one.MerchantId)
	var merchantProductName = ""
	sub := query.GetSubscriptionBySubscriptionId(ctx, one.SubscriptionId)
	if sub == nil {
		// todo mark invoice not relative to subscription
		sub = query.GetLatestActiveOrIncompleteOrCreateSubscriptionByUserId(ctx, one.UserId, merchant.Id)
	}
	if sub != nil {
		plan := query.GetPlanById(ctx, sub.PlanId)
		merchantProductName = plan.PlanName
	}
	if one.Status > consts.InvoiceStatusPending {
		pdfFileName := utility.DownloadFile(one.SendPdf)
		utility.Assert(len(pdfFileName) > 0, "download pdf error:"+one.SendPdf)
		var template = email.TemplateNewProcessingInvoice
		if one.Status == consts.InvoiceStatusPaid {
			payment := query.GetPaymentByPaymentId(ctx, one.PaymentId)
			if payment.Automatic == 0 {
				template = email.TemplateInvoiceManualPaid
			} else {
				template = email.TemplateInvoiceAutomaticPaid
			}
		} else if one.Status == consts.InvoiceStatusCancelled || one.Status == consts.InvoiceStatusFailed {
			template = email.TemplateInvoiceCancel
		}
		err := email.SendTemplateEmail(ctx, merchant.Id, one.SendEmail, user.TimeZone, template, pdfFileName, &email.TemplateVariable{
			InvoiceId:           one.InvoiceId,
			UserName:            user.FirstName + " " + user.LastName,
			MerchantProductName: merchantProductName,
			MerchantCustomEmail: merchant.Email,
			MerchantName:        merchant.Name,
			DateNow:             gtime.Now(),
			PeriodEnd:           gtime.Now().AddDate(0, 0, 5),
			PaymentAmount:       strconv.FormatInt(one.TotalAmount, 10),
			TokenExpireMinute:   strconv.FormatInt(consts.GetConfigInstance().Auth.Login.Expire/60, 10),
			Link:                "<a href=\"" + one.Link + "\">Link</a>",
		})
		utility.AssertError(err, "send email error")
		//update send status
		_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
			dao.Invoice.Columns().SendStatus: 1,
			dao.Invoice.Columns().GmtModify:  gtime.Now(),
		}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			fmt.Printf("SendSubscriptionInvoiceEmailToUser update err:%s", err.Error())
		}
	} else {
		fmt.Printf("SendSubscriptionInvoiceEmailToUser invoice status is pending or init, email not send")
	}
	return nil
}
