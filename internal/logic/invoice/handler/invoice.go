package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"time"
	"unibee/api/bean"
	config2 "unibee/internal/cmd/config"
	"unibee/internal/consts"
	"unibee/internal/controller/link"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/email"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/gateway_bean"
	"unibee/internal/logic/subscription/config"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func CreateProcessingInvoiceForSub(ctx context.Context, simplify *bean.InvoiceSimplify, sub *entity.Subscription) (*entity.Invoice, error) {
	utility.Assert(simplify != nil, "invoice data is nil")
	utility.Assert(sub != nil, "sub is nil")
	user := query.GetUserAccountById(ctx, sub.UserId)
	var sendEmail = ""
	if user != nil {
		sendEmail = user.Email
	}
	status := consts.InvoiceStatusProcessing
	invoiceId := utility.CreateInvoiceId()
	st := utility.CreateInvoiceSt()
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
		Status:                         status,
		SendStatus:                     simplify.SendStatus,
		SendEmail:                      sendEmail,
		UniqueId:                       invoiceId,
		SendTerms:                      st,
		TotalAmount:                    simplify.TotalAmount,
		TotalAmountExcludingTax:        simplify.TotalAmountExcludingTax,
		TaxAmount:                      simplify.TaxAmount,
		TaxPercentage:                  simplify.TaxPercentage,
		SubscriptionAmount:             simplify.SubscriptionAmount,
		SubscriptionAmountExcludingTax: simplify.SubscriptionAmountExcludingTax,
		Lines:                          utility.MarshalToJsonString(simplify.Lines),
		Link:                           link.GetInvoiceLink(invoiceId, st),
		CreateTime:                     gtime.Now().Timestamp(),
		FinishTime:                     gtime.Now().Timestamp(),
		DayUtilDue:                     simplify.DayUtilDue,
		DiscountAmount:                 simplify.DiscountAmount,
		DiscountCode:                   simplify.DiscountCode,
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

func CreateProcessInvoiceForNewPayment(ctx context.Context, invoice *bean.InvoiceSimplify, payment *entity.Payment) (*entity.Invoice, error) {
	utility.Assert(invoice != nil, "invoice data is nil")
	utility.Assert(payment != nil, "payment data is nil")
	utility.Assert(len(payment.PaymentId) > 0, "paymentId is nil")
	utility.Assert(len(payment.InvoiceId) > 0, "payment InvoiceId is nil")
	user := query.GetUserAccountById(ctx, payment.UserId)
	var sendEmail = ""
	if user != nil {
		sendEmail = user.Email
	}
	st := utility.CreateInvoiceSt()
	one := &entity.Invoice{
		BizType:                        payment.BizType,
		UserId:                         payment.UserId,
		MerchantId:                     payment.MerchantId,
		SubscriptionId:                 payment.SubscriptionId,
		InvoiceName:                    payment.BillingReason,
		InvoiceId:                      payment.InvoiceId,
		UniqueId:                       payment.PaymentId,
		PaymentId:                      payment.PaymentId,
		Link:                           link.GetInvoiceLink(payment.InvoiceId, st),
		SendTerms:                      st,
		PeriodStart:                    invoice.PeriodStart,
		PeriodEnd:                      invoice.PeriodEnd,
		PeriodStartTime:                gtime.NewFromTimeStamp(invoice.PeriodStart),
		PeriodEndTime:                  gtime.NewFromTimeStamp(invoice.PeriodEnd),
		Currency:                       payment.Currency,
		CryptoCurrency:                 payment.CryptoCurrency,
		GatewayId:                      payment.GatewayId,
		Status:                         consts.InvoiceStatusProcessing,
		SendStatus:                     invoice.SendStatus,
		SendEmail:                      sendEmail,
		GatewayPaymentId:               payment.GatewayPaymentId,
		TotalAmount:                    invoice.TotalAmount,
		CryptoAmount:                   payment.CryptoAmount,
		TotalAmountExcludingTax:        invoice.TotalAmountExcludingTax,
		TaxAmount:                      invoice.TaxAmount,
		TaxPercentage:                  invoice.TaxPercentage,
		SubscriptionAmount:             invoice.SubscriptionAmount,
		SubscriptionAmountExcludingTax: invoice.SubscriptionAmountExcludingTax,
		Lines:                          utility.MarshalToJsonString(invoice.Lines),
		PaymentLink:                    payment.Link,
		CreateTime:                     gtime.Now().Timestamp(),
		FinishTime:                     gtime.Now().Timestamp(),
		DayUtilDue:                     invoice.DayUtilDue,
		DiscountAmount:                 invoice.DiscountAmount,
		DiscountCode:                   invoice.DiscountCode,
	}

	result, err := dao.Invoice.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`CreateProcessInvoiceForNewPayment record insert failure %s`, err.Error())
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))
	_ = InvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true)
	if err != nil {
		return nil, err
	}
	return one, nil
}

func UpdateInvoiceFromPayment(ctx context.Context, payment *entity.Payment) (*entity.Invoice, error) {
	utility.Assert(payment != nil, "payment data is nil")
	utility.Assert(len(payment.PaymentId) > 0, "paymentId is nil")
	one := query.GetInvoiceByPaymentId(ctx, payment.PaymentId)
	if one == nil {
		return nil, gerror.New("invoice not found, paymentId:" + payment.PaymentId + " subId:" + payment.SubscriptionId)
	}
	if one.Status == consts.InvoiceStatusFailed {
		return nil, gerror.New("invoice has failed, payment:" + payment.PaymentId + " subId:" + payment.SubscriptionId)
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

func CreateProcessInvoiceForNewPaymentRefund(ctx context.Context, invoice *bean.InvoiceSimplify, refund *entity.Refund) (*entity.Invoice, error) {
	utility.Assert(invoice != nil, "invoice data is nil")
	utility.Assert(refund != nil, "refund data is nil")
	utility.Assert(len(refund.RefundId) > 0, "refundId is nil")
	utility.Assert(len(refund.PaymentId) > 0, "paymentId is nil")
	utility.Assert(len(refund.InvoiceId) > 0, "refund InvoiceId is nil")
	payment := query.GetPaymentByPaymentId(ctx, refund.PaymentId)
	utility.Assert(payment != nil, "payment data is nil")
	user := query.GetUserAccountById(ctx, refund.UserId)
	var sendEmail = ""
	if user != nil {
		sendEmail = user.Email
	}
	st := utility.CreateInvoiceSt()
	one := &entity.Invoice{
		BizType:                        refund.BizType,
		UserId:                         refund.UserId,
		MerchantId:                     refund.MerchantId,
		SubscriptionId:                 refund.SubscriptionId,
		InvoiceName:                    payment.BillingReason,
		InvoiceId:                      refund.InvoiceId,
		UniqueId:                       refund.RefundId,
		PaymentId:                      refund.PaymentId,
		RefundId:                       refund.RefundId,
		Link:                           link.GetInvoiceLink(refund.InvoiceId, st),
		SendNote:                       invoice.SendNote,
		SendTerms:                      st,
		PeriodStart:                    invoice.PeriodStart,
		PeriodEnd:                      invoice.PeriodEnd,
		PeriodStartTime:                gtime.NewFromTimeStamp(invoice.PeriodStart),
		PeriodEndTime:                  gtime.NewFromTimeStamp(invoice.PeriodEnd),
		Currency:                       refund.Currency,
		CryptoCurrency:                 payment.CryptoCurrency,
		GatewayId:                      refund.GatewayId,
		Status:                         consts.InvoiceStatusProcessing,
		SendStatus:                     invoice.SendStatus,
		SendEmail:                      sendEmail,
		TotalAmount:                    invoice.TotalAmount,
		CryptoAmount:                   payment.CryptoAmount,
		TotalAmountExcludingTax:        invoice.TotalAmountExcludingTax,
		TaxAmount:                      invoice.TaxAmount,
		TaxPercentage:                  invoice.TaxPercentage,
		SubscriptionAmount:             invoice.SubscriptionAmount,
		SubscriptionAmountExcludingTax: invoice.SubscriptionAmountExcludingTax,
		Lines:                          utility.MarshalToJsonString(invoice.Lines),
		CreateTime:                     gtime.Now().Timestamp(),
		FinishTime:                     gtime.Now().Timestamp(),
		DayUtilDue:                     invoice.DayUtilDue,
		DiscountAmount:                 invoice.DiscountAmount,
		DiscountCode:                   invoice.DiscountCode,
	}

	result, err := dao.Invoice.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`CreateProcessInvoiceForNewPaymentRefund record insert failure %s`, err.Error())
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))
	_ = InvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true)
	if err != nil {
		return nil, err
	}
	return one, nil
}

func UpdateInvoiceFromPaymentRefund(ctx context.Context, refund *entity.Refund) (*entity.Invoice, error) {
	utility.Assert(refund != nil, "refund data is nil")
	utility.Assert(len(refund.RefundId) > 0, "refundId is nil")
	utility.Assert(len(refund.PaymentId) > 0, "paymentId is nil")
	payment := query.GetPaymentByPaymentId(ctx, refund.PaymentId)
	utility.Assert(payment != nil, "payment data is nil")
	one := query.GetInvoiceByRefundId(ctx, refund.RefundId)
	if one == nil {
		return nil, gerror.New("invoice not found, refundId:" + refund.RefundId + " subId:" + payment.SubscriptionId)
	}
	if one.Status == consts.InvoiceStatusFailed {
		return nil, gerror.New("invoice has failed, refundId:" + refund.RefundId + " subId:" + payment.SubscriptionId)
	}
	var status = consts.InvoiceStatusProcessing
	if refund.Status == consts.RefundSuccess {
		status = consts.InvoiceStatusPaid
	} else if refund.Status == consts.RefundFailed {
		status = consts.InvoiceStatusFailed
	} else if refund.Status == consts.RefundCancelled {
		status = consts.InvoiceStatusCancelled
	}
	_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
		dao.Invoice.Columns().Status:    status,
		dao.Invoice.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	if one.Status != status && one.BizType == consts.BizTypeSubscription {
		_ = InvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true)
	}
	one.Status = status
	return one, nil
}

func MarkInvoiceAsPaidForZeroPayment(ctx context.Context, invoiceId string) (*entity.Invoice, error) {
	one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
	if one == nil {
		return nil, gerror.New("invoice not found, InvoiceId:" + invoiceId)
	}
	if one.TotalAmount != 0 {
		return nil, gerror.New("invoice totalAmount not zero, InvoiceId:" + invoiceId)
	}
	_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
		dao.Invoice.Columns().Status:    consts.InvoiceStatusPaid,
		dao.Invoice.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	if one.BizType == consts.BizTypeSubscription {
		_ = InvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true)
	}
	one.Status = consts.InvoiceStatusPaid
	return one, nil
}

func InvoicePdfGenerateAndEmailSendBackground(invoiceId string, sendUserEmail bool) (err error) {
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
		time.Sleep(2 * time.Second)
		one := query.GetInvoiceByInvoiceId(backgroundCtx, invoiceId)
		if one == nil {
			g.Log().Errorf(backgroundCtx, "InvoicePdfGenerateAndEmailSendBackground Error one is null")
			return
		}
		if len(one.Lines) == 0 {
			// invoice with valid lines will send emails
			g.Log().Errorf(backgroundCtx, "InvoicePdfGenerateAndEmailSendBackground Error one.lines is null")
			return
		}

		filepath := GenerateInvoicePdf(backgroundCtx, one)
		if len(filepath) > 0 {
			url, _ := UploadInvoicePdf(backgroundCtx, one.InvoiceId, filepath)
			if len(url) > 0 {
				_, err = dao.Invoice.Ctx(backgroundCtx).Data(g.Map{
					dao.Invoice.Columns().SendPdf:   url,
					dao.Invoice.Columns().GmtModify: gtime.Now(),
				}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
				if err != nil {
					fmt.Printf("UploadInvoice SendPdf err:%s", err.Error())
				}
			}
		}
		if sendUserEmail && one.SendStatus != consts.InvoiceSendStatusUnnecessary {
			err := SendInvoiceEmailToUser(backgroundCtx, one.InvoiceId)
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
	gateway := query.GetGatewayById(ctx, one.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	user := query.GetUserAccountById(ctx, uint64(one.UserId))
	utility.Assert(user != nil, "user not found")
	trans, err := api.GetGatewayServiceProvider(ctx, one.GatewayId).GatewayCryptoFiatTrans(ctx, &gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq{
		Amount:   one.TotalAmount,
		Currency: one.Currency,
		Gateway:  gateway,
	})
	if err != nil {
		return err
	}

	_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
		dao.Invoice.Columns().CryptoCurrency: trans.CryptoCurrency,
		dao.Invoice.Columns().CryptoAmount:   trans.CryptoAmount,
	}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		fmt.Printf("ReconvertCryptoDataForInvoice update err:%s", err.Error())
	}
	//todo mark cancel the currency payment and regeneration
	return err
}

func SendInvoiceEmailToUser(ctx context.Context, invoiceId string) error {
	one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
	utility.Assert(one != nil, "invoice not found")
	utility.Assert(one.UserId > 0, "invoice userId not found")
	utility.Assert(one.MerchantId > 0, "invoice merchantId not found")
	utility.Assert(len(one.SendEmail) > 0, "SendEmail Is Nil, InvoiceId:"+one.InvoiceId)
	_, emailKey := email.GetDefaultMerchantEmailConfig(ctx, one.MerchantId)
	if len(emailKey) == 0 {
		return gerror.New("Email gateway not setup")
	}
	var pdfFileName string
	if len(one.SendPdf) > 0 {
		pdfFileName = utility.DownloadFile(one.SendPdf)
	} else {
		pdfFileName = GenerateInvoicePdf(ctx, one)
	}
	if len(pdfFileName) == 0 {
		return gerror.New("pdfFile download or generate error")
	}
	if !config.GetMerchantSubscriptionConfig(ctx, one.MerchantId).InvoiceEmail {
		fmt.Printf("SendInvoiceEmailToUser merchant configed to stop sending invoice email, email not send")
		return nil
	}
	user := query.GetUserAccountById(ctx, one.UserId)
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
		utility.Assert(len(pdfFileName) > 0, "pdfFile download or generate error:"+one.InvoiceId)
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
			MerchantName:        query.GetMerchantCountryConfigName(ctx, one.MerchantId, user.CountryCode),
			DateNow:             gtime.Now(),
			PeriodEnd:           gtime.Now().AddDate(0, 0, 5),
			PaymentAmount:       strconv.FormatInt(one.TotalAmount, 10),
			TokenExpireMinute:   strconv.FormatInt(config2.GetConfigInstance().Auth.Login.Expire/60, 10),
			Link:                "<a href=\"" + link.GetInvoiceLink(one.InvoiceId, one.SendTerms) + "\">Link</a>",
		})
		utility.AssertError(err, "send email error")
		//update send status
		_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
			dao.Invoice.Columns().SendStatus: consts.InvoiceSendStatusSend,
			dao.Invoice.Columns().GmtModify:  gtime.Now(),
		}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			fmt.Printf("SendInvoiceEmailToUser update err:%s", err.Error())
		}
	} else {
		fmt.Printf("SendInvoiceEmailToUser invoice status is pending or init, email not send")
	}
	return nil
}
