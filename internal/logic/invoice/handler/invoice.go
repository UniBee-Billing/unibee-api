package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"unibee-api/internal/consts"
	dao "unibee-api/internal/dao/oversea_pay"
	"unibee-api/internal/logic/email"
	"unibee-api/internal/logic/gateway/ro"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

func UpdateInvoiceFromPayment(ctx context.Context, payment *entity.Payment) (*entity.Invoice, error) {
	utility.Assert(payment != nil, "payment data is nil")
	one := query.GetInvoiceByPaymentId(ctx, payment.PaymentId)
	utility.Assert(one != nil, "invoice not found, paymentId:"+payment.PaymentId+" subId:"+payment.SubscriptionId)
	var status = consts.InvoiceStatusProcessing
	if payment.Status == consts.PAY_SUCCESS {
		status = consts.InvoiceStatusPaid
	} else if payment.Status == consts.PAY_FAILED {
		status = consts.InvoiceStatusFailed
	} else if payment.Status == consts.PAY_CANCEL {
		status = consts.InvoiceStatusCancelled
	}
	utility.Assert(one != nil, "invoice not found")
	_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
		dao.Invoice.Columns().Status:           status,
		dao.Invoice.Columns().GmtModify:        gtime.Now(),
		dao.Invoice.Columns().GatewayPaymentId: payment.GatewayPaymentId,
		dao.Invoice.Columns().Link:             payment.Link,
	}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	if one.Status != status {
		//更新状态发送邮件
		_ = SubscriptionInvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true)
	}
	one.Status = status
	one.GatewayPaymentId = payment.GatewayPaymentId
	one.Link = payment.Link
	return one, nil
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

func CreateOrUpdateInvoiceFromPayment(ctx context.Context, simplify *ro.InvoiceDetailSimplify, payment *entity.Payment) (*entity.Invoice, error) {
	utility.Assert(simplify != nil, "invoice data is nil")
	utility.Assert(payment != nil, "payment data is nil")
	one := query.GetInvoiceByPaymentId(ctx, payment.PaymentId)
	if one == nil && len(simplify.InvoiceId) > 0 {
		one = query.GetInvoiceByInvoiceId(ctx, simplify.InvoiceId)
	}
	user := query.GetUserAccountById(ctx, uint64(payment.UserId))
	var sendEmail = ""
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
			InvoiceName:                    payment.BillingReason,
			InvoiceId:                      utility.CreateInvoiceId(),
			PeriodStart:                    simplify.PeriodStart,
			PeriodEnd:                      simplify.PeriodEnd,
			PeriodStartTime:                gtime.NewFromTimeStamp(simplify.PeriodStart),
			PeriodEndTime:                  gtime.NewFromTimeStamp(simplify.PeriodEnd),
			Currency:                       payment.Currency,
			GatewayId:                      payment.GatewayId,
			Status:                         status,
			SendStatus:                     0,
			SendEmail:                      sendEmail,
			GatewayPaymentId:               payment.GatewayPaymentId,
			UniqueId:                       payment.PaymentId,
			PaymentId:                      payment.PaymentId,
			TotalAmount:                    simplify.TotalAmount,
			TotalAmountExcludingTax:        simplify.TotalAmountExcludingTax,
			TaxAmount:                      simplify.TaxAmount,
			SubscriptionAmount:             simplify.SubscriptionAmount,
			SubscriptionAmountExcludingTax: simplify.SubscriptionAmountExcludingTax,
			Lines:                          utility.MarshalToJsonString(simplify.Lines),
			Link:                           payment.Link,
			CreateAt:                       gtime.Now().Timestamp(),
		}

		result, err := dao.Invoice.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			err = gerror.Newf(`CreateOrUpdateInvoiceByChannelDetail record insert failure %s`, err.Error())
			return nil, err
		}
		id, _ := result.LastInsertId()
		one.Id = uint64(uint(id))
		//新建 Invoice 发送邮件
		_ = SubscriptionInvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true)
		err = UpdatePaymentInvoiceId(ctx, payment.PaymentId, one.InvoiceId)
		if err != nil {
			return nil, err
		}
	} else {
		//更新
		_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
			dao.Invoice.Columns().BizType:                        payment.BizType,
			dao.Invoice.Columns().MerchantId:                     payment.MerchantId,
			dao.Invoice.Columns().UserId:                         payment.UserId,
			dao.Invoice.Columns().SubscriptionId:                 payment.SubscriptionId,
			dao.Invoice.Columns().GatewayId:                      payment.GatewayId,
			dao.Invoice.Columns().PaymentId:                      payment.PaymentId,
			dao.Invoice.Columns().UniqueId:                       payment.PaymentId,
			dao.Invoice.Columns().Currency:                       payment.Currency,
			dao.Invoice.Columns().Status:                         status,
			dao.Invoice.Columns().GmtModify:                      gtime.Now(),
			dao.Invoice.Columns().GatewayPaymentId:               payment.GatewayPaymentId,
			dao.Invoice.Columns().TotalAmount:                    simplify.TotalAmount,
			dao.Invoice.Columns().TotalAmountExcludingTax:        simplify.TotalAmountExcludingTax,
			dao.Invoice.Columns().TaxAmount:                      simplify.TaxAmount,
			dao.Invoice.Columns().SubscriptionAmount:             simplify.SubscriptionAmount,
			dao.Invoice.Columns().SubscriptionAmountExcludingTax: simplify.SubscriptionAmountExcludingTax,
			dao.Invoice.Columns().Lines:                          utility.FormatToJsonString(simplify.Lines),
			dao.Invoice.Columns().Link:                           payment.Link,
		}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			return nil, err
		}
		if one.Status != status {
			//更新状态发送邮件
			_ = SubscriptionInvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true)
		}
	}
	one = query.GetInvoiceByPaymentId(ctx, payment.PaymentId)
	return one, nil
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
		if one.BizType == consts.BIZ_TYPE_ONE_TIME || len(one.Lines) == 0 {
			// invoice not one time type and valid lines will send emails
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
			err := SendSubscriptionInvoiceEmailToUser(backgroundCtx, one.InvoiceId)
			utility.Assert(err == nil, "SendInvoiceEmail error")
		}
	}()
	return nil
}

func SendSubscriptionInvoiceEmailToUser(ctx context.Context, invoiceId string) error {
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
	if sub == nil {
		// todo mark invoice not relative to subscription
		sub = query.GetLatestActiveOrCreateSubscriptionByUserId(ctx, one.UserId, merchant.Id)
	}
	if sub != nil {
		plan := query.GetPlanById(ctx, sub.PlanId)
		merchantProductName = plan.PlanName
	}
	if one.Status > consts.InvoiceStatusPending {
		pdfFileName := utility.DownloadFile(one.SendPdf)
		utility.Assert(len(pdfFileName) > 0, "download pdf error:"+one.SendPdf)
		payment := query.GetPaymentByPaymentId(ctx, one.PaymentId)
		var link = one.Link
		if len(link) == 0 {
			link = payment.Link
		}
		var template = email.TemplateNewProcessingInvoice
		if one.Status == consts.InvoiceStatusPaid {
			if payment.Automatic == 0 {
				template = email.TemplateInvoiceManualPaid
			} else {
				template = email.TemplateInvoiceAutomaticPaid
			}
		} else if one.Status == consts.InvoiceStatusCancelled || one.Status == consts.InvoiceStatusFailed {
			template = email.TemplateInvoiceCancel
		}
		//err := email.SendPdfAttachEmailToUser(one.SendEmail, "Invoice", "Invoice", pdfFileName, fmt.Sprintf("%s.pdf", one.InvoiceId))
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
			Link:                "<a href=\"" + link + "\">Link</a>",
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
			fmt.Printf("SendSubscriptionInvoiceEmailToUser update err:%s", err.Error())
		}
		//rowAffected, err := update.RowsAffected()
		//if rowAffected != 1 {
		//	fmt.Printf("SendSubscriptionInvoiceEmailToUser update err:%s", update)
		//}
	} else {
		fmt.Printf("SendSubscriptionInvoiceEmailToUser invoice status is pending or init, email not send")
	}
	return nil
}
