package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	"strconv"
	"time"
	"unibee/api/bean"
	config2 "unibee/internal/cmd/config"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	"unibee/internal/controller/link"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/currency"
	"unibee/internal/logic/email"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/gateway_bean"
	"unibee/internal/logic/merchant_config"
	"unibee/internal/logic/subscription/config"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func CreateProcessInvoiceForNewPayment(ctx context.Context, invoice *bean.Invoice, payment *entity.Payment) (*entity.Invoice, error) {
	utility.Assert(invoice != nil, "invoice data is nil")
	utility.Assert(payment != nil, "payment data is nil")
	utility.Assert(len(payment.PaymentId) > 0, "paymentId is nil")
	utility.Assert(len(payment.InvoiceId) > 0, "payment InvoiceId is nil")
	user := query.GetUserAccountById(ctx, payment.UserId)
	var sendEmail = ""
	var userSnapshot *entity.UserAccount
	if user != nil {
		sendEmail = user.Email
		userSnapshot = &entity.UserAccount{
			Email:         user.Email,
			CountryCode:   user.CountryCode,
			CountryName:   user.CountryName,
			VATNumber:     user.VATNumber,
			TaxPercentage: user.TaxPercentage,
			GatewayId:     user.GatewayId,
			Type:          user.Type,
			UserName:      user.UserName,
			Mobile:        user.Mobile,
			Phone:         user.Phone,
			Address:       user.Address,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			CompanyName:   user.CompanyName,
			City:          user.City,
			ZipCode:       user.ZipCode,
		}
	}
	st := utility.CreateInvoiceSt()
	var name = invoice.InvoiceName
	if len(name) == 0 {
		name = payment.BillingReason
	}
	one := &entity.Invoice{
		SubscriptionId:                 payment.SubscriptionId,
		BizType:                        payment.BizType,
		UserId:                         payment.UserId,
		MerchantId:                     payment.MerchantId,
		InvoiceName:                    name,
		ProductName:                    invoice.ProductName,
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
		CountryCode:                    payment.CountryCode,
		VatNumber:                      invoice.VatNumber,
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
		BillingCycleAnchor:             invoice.BillingCycleAnchor,
		Data:                           utility.MarshalToJsonString(userSnapshot),
		MetaData:                       utility.MarshalToJsonString(invoice.Metadata),
		CreateFrom:                     invoice.CreateFrom,
	}

	result, err := dao.Invoice.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`CreateProcessInvoiceForNewPayment record insert failure %s`, err.Error())
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))
	_, _ = redismq.Send(&redismq.Message{
		Topic: redismq2.TopicInvoiceCreated.Topic,
		Tag:   redismq2.TopicInvoiceCreated.Tag,
		Body:  one.InvoiceId,
	})
	_, _ = redismq.Send(&redismq.Message{
		Topic: redismq2.TopicInvoiceProcessed.Topic,
		Tag:   redismq2.TopicInvoiceProcessed.Tag,
		Body:  one.InvoiceId,
	})
	_ = InvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true, false)
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
	if one.Status == consts.InvoiceStatusFailed || one.Status == consts.InvoiceStatusCancelled {
		if payment.Status == consts.PaymentSuccess {
			_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
				dao.Invoice.Columns().Status:           consts.InvoiceStatusPaid,
				dao.Invoice.Columns().GmtModify:        gtime.Now(),
				dao.Invoice.Columns().GatewayPaymentId: payment.GatewayPaymentId,
				dao.Invoice.Columns().PaymentLink:      payment.Link,
			}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
			if err != nil {
				g.Log().Errorf(ctx, "UpdateInvoiceFromPayment_Reverse invoiceId:%s paymentId:%s error:%s", one.InvoiceId, payment.PaymentId, err.Error())
				return one, gerror.New("invoice reverse failed, invoiceId:" + one.InvoiceId + " paymentId:" + payment.PaymentId + " subId:" + payment.SubscriptionId)
			} else {
				one.Status = consts.InvoiceStatusPaid
				one.GatewayPaymentId = payment.GatewayPaymentId
				one.Link = payment.Link
				g.Log().Infof(ctx, "UpdateInvoiceFromPayment_Reverse invoiceId:%s paymentId:%s", one.InvoiceId, payment.PaymentId)
				_, _ = redismq.Send(&redismq.Message{
					Topic: redismq2.TopicInvoicePaid.Topic,
					Tag:   redismq2.TopicInvoicePaid.Tag,
					Body:  one.InvoiceId,
				})
				return one, nil
			}
		}
		g.Log().Infof(ctx, "UpdateInvoiceFromPayment already failed or cancelled invoiceId:%s paymentId:%s", one.InvoiceId, payment.PaymentId)
		return one, nil
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
		return one, err
	}
	if one.Status != status {
		_, _ = dao.Invoice.Ctx(ctx).Data(g.Map{
			dao.Invoice.Columns().SendPdf: "",
		}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
		_ = InvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true, false)
		if status == consts.InvoiceStatusPaid {
			_, _ = redismq.Send(&redismq.Message{
				Topic: redismq2.TopicInvoicePaid.Topic,
				Tag:   redismq2.TopicInvoicePaid.Tag,
				Body:  one.InvoiceId,
			})
		} else if status == consts.InvoiceStatusCancelled {
			g.Log().Infof(ctx, "CancelProcessingInvoice invoiceId:%s reason:%s", one.InvoiceId, "UpdateInvoiceFromPayment")
			_, _ = redismq.Send(&redismq.Message{
				Topic: redismq2.TopicInvoiceCancelled.Topic,
				Tag:   redismq2.TopicInvoiceCancelled.Tag,
				Body:  one.InvoiceId,
			})
		} else if status == consts.InvoiceStatusFailed {
			g.Log().Infof(ctx, "ProcessingInvoiceFailure invoiceId:%s reason:%s", one.InvoiceId, "UpdateInvoiceFromPayment")
			_, _ = redismq.Send(&redismq.Message{
				Topic: redismq2.TopicInvoiceFailed.Topic,
				Tag:   redismq2.TopicInvoiceFailed.Tag,
				Body:  one.InvoiceId,
			})
		}
	}
	one.Status = status
	one.GatewayPaymentId = payment.GatewayPaymentId
	one.Link = payment.Link
	return one, nil
}

func CreateProcessInvoiceForNewPaymentRefund(ctx context.Context, invoice *bean.Invoice, refund *entity.Refund) (*entity.Invoice, error) {
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
		SubscriptionId:                 payment.SubscriptionId,
		BizType:                        payment.BizType,
		UserId:                         refund.UserId,
		MerchantId:                     refund.MerchantId,
		InvoiceName:                    invoice.InvoiceName,
		ProductName:                    invoice.ProductName,
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
		CountryCode:                    invoice.CountryCode,
		VatNumber:                      invoice.VatNumber,
		TaxPercentage:                  invoice.TaxPercentage,
		SubscriptionAmount:             invoice.SubscriptionAmount,
		SubscriptionAmountExcludingTax: invoice.SubscriptionAmountExcludingTax,
		Lines:                          utility.MarshalToJsonString(invoice.Lines),
		CreateTime:                     gtime.Now().Timestamp(),
		FinishTime:                     gtime.Now().Timestamp(),
		DayUtilDue:                     invoice.DayUtilDue,
		DiscountAmount:                 invoice.DiscountAmount,
		DiscountCode:                   invoice.DiscountCode,
		CreateFrom:                     refund.RefundComment,
	}

	result, err := dao.Invoice.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`CreateProcessInvoiceForNewPaymentRefund record insert failure %s`, err.Error())
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))
	_, _ = redismq.Send(&redismq.Message{
		Topic: redismq2.TopicInvoiceCreated.Topic,
		Tag:   redismq2.TopicInvoiceCreated.Tag,
		Body:  one.InvoiceId,
	})
	_, _ = redismq.Send(&redismq.Message{
		Topic: redismq2.TopicInvoiceProcessed.Topic,
		Tag:   redismq2.TopicInvoiceProcessed.Tag,
		Body:  one.InvoiceId,
	})
	_ = InvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true, false)
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
		return one, gerror.New("invoice has failed, refundId:" + refund.RefundId + " subId:" + payment.SubscriptionId)
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
		return one, err
	}
	if one.Status != status {
		_ = InvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true, false)
		if status == consts.InvoiceStatusPaid {
			_, _ = redismq.Send(&redismq.Message{
				Topic: redismq2.TopicInvoicePaid.Topic,
				Tag:   redismq2.TopicInvoicePaid.Tag,
				Body:  one.InvoiceId,
			})
		} else if status == consts.InvoiceStatusCancelled {
			g.Log().Infof(ctx, "CancelProcessingInvoice invoiceId:%s reason:%s", one.InvoiceId, "UpdateInvoiceFromPaymentRefund")
			_, _ = redismq.Send(&redismq.Message{
				Topic: redismq2.TopicInvoiceCancelled.Topic,
				Tag:   redismq2.TopicInvoiceCancelled.Tag,
				Body:  one.InvoiceId,
			})
		} else if status == consts.InvoiceStatusFailed {
			g.Log().Infof(ctx, "ProcessingInvoiceFailure invoiceId:%s reason:%s", one.InvoiceId, "UpdateInvoiceFromPayment")
			_, _ = redismq.Send(&redismq.Message{
				Topic: redismq2.TopicInvoiceFailed.Topic,
				Tag:   redismq2.TopicInvoiceFailed.Tag,
				Body:  one.InvoiceId,
			})
		}
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
	_ = InvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true, false)
	one.Status = consts.InvoiceStatusPaid
	go func() {
		time.Sleep(1 * time.Second)
		_, _ = redismq.Send(&redismq.Message{
			Topic: redismq2.TopicInvoicePaid.Topic,
			Tag:   redismq2.TopicInvoicePaid.Tag,
			Body:  one.InvoiceId,
		})
	}()

	return one, nil
}

func InvoicePdfGenerateAndEmailSendBackground(invoiceId string, sendUserEmail bool, manualSend bool) (err error) {
	return InvoicePdfGenerateAndEmailSendByTargetTemplateBackground(invoiceId, sendUserEmail, manualSend, "")
}

func InvoicePdfGenerateAndEmailSendByTargetTemplateBackground(invoiceId string, sendUserEmail bool, manualSend bool, targetTemplate string) (err error) {
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
			err := SendInvoiceEmailToUser(backgroundCtx, one.InvoiceId, manualSend, targetTemplate)
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

	exchangeApiKeyConfig := merchant_config.GetMerchantConfig(ctx, user.MerchantId, config.FiatExchangeApiKey)
	var cryptoCurrency string
	var cryptoAmount int64 = -1
	if exchangeApiKeyConfig != nil && len(exchangeApiKeyConfig.ConfigValue) > 0 {
		if one.Currency == "USD" {
			_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
				dao.Invoice.Columns().CryptoCurrency: "USD",
				dao.Invoice.Columns().CryptoAmount:   one.TotalAmount,
			}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
			if err != nil {
				fmt.Printf("ReconvertCryptoDataForInvoice update err:%s", err.Error())
			}
			return err
		} else {
			rate, err := currency.GetExchangeConversionRates(ctx, exchangeApiKeyConfig.ConfigValue, "USD", one.Currency)
			if err != nil {
				return err
			}
			if rate != nil {
				cryptoCurrency = "USD"
				cryptoAmount = utility.RoundUp(float64(one.TotalAmount) / *rate)
			}
		}
	}
	if len(cryptoCurrency) == 0 {
		trans, err := api.GetGatewayServiceProvider(ctx, one.GatewayId).GatewayCryptoFiatTrans(ctx, &gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq{
			Amount:   one.TotalAmount,
			Currency: one.Currency,
			Gateway:  gateway,
		})
		if err != nil {
			return err
		}
		cryptoCurrency = trans.CryptoCurrency
		cryptoAmount = trans.CryptoAmount
	}
	utility.Assert(len(cryptoCurrency) > 0, "transfer to crypto currency error")
	_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
		dao.Invoice.Columns().CryptoCurrency: cryptoCurrency,
		dao.Invoice.Columns().CryptoAmount:   cryptoAmount,
	}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		fmt.Printf("ReconvertCryptoDataForInvoice update err:%s", err.Error())
	}
	return err
}

func SendInvoiceEmailToUser(ctx context.Context, invoiceId string, manualSend bool, sendTemplate string) error {
	one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
	utility.Assert(one != nil, "invoice not found")
	if len(one.RefundId) == 0 && one.TotalAmount <= 0 {
		g.Log().Infof(ctx, "SendInvoiceEmailToUser invoice totalAmount lower then zero, email not send")
		return nil
	}
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
	if !manualSend && !config.GetMerchantSubscriptionConfig(ctx, one.MerchantId).InvoiceEmail {
		fmt.Printf("SendInvoiceEmailToUser merchant configed to stop sending invoice email, email not send")
		return nil
	}
	user := query.GetUserAccountById(ctx, one.UserId)
	merchant := query.GetMerchantById(ctx, one.MerchantId)
	if len(one.RefundId) == 0 {
		if one.Status > consts.InvoiceStatusPending {
			utility.Assert(len(pdfFileName) > 0, "pdfFile download or generate error:"+one.InvoiceId)
			var template = email.TemplateNewProcessingInvoice
			var accountHolder string
			var bic string
			var iban string
			var address string
			gateway := query.GetGatewayById(ctx, one.GatewayId)
			if gateway != nil && gateway.GatewayType == consts.GatewayTypeWireTransfer {
				template = email.TemplateNewProcessingInvoiceForWireTransfer
				gatewaySimplify := bean.SimplifyGateway(gateway)
				if gatewaySimplify != nil {
					accountHolder = gatewaySimplify.Bank.AccountHolder
					bic = gatewaySimplify.Bank.BIC
					iban = gatewaySimplify.Bank.IBAN
					address = gatewaySimplify.Bank.Address
				}
			} else if one.TrialEnd > 0 && one.TrialEnd > one.PeriodStart {
				// paid trial invoice
				template = email.TemplateNewProcessingInvoiceForPaidTrial
			} else if one.TrialEnd == -2 {
				// first cycle invoice after trial
				template = email.TemplateNewProcessingInvoiceAfterTrial
			}
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
			if len(sendTemplate) > 0 {
				template = sendTemplate
			}
			err := email.SendTemplateEmail(ctx, merchant.Id, one.SendEmail, user.TimeZone, user.Language, template, pdfFileName, &email.TemplateVariable{
				InvoiceId:             one.InvoiceId,
				UserName:              user.FirstName + " " + user.LastName,
				MerchantProductName:   one.ProductName,
				MerchantCustomerEmail: merchant.Email,
				MerchantName:          query.GetMerchantCountryConfigName(ctx, one.MerchantId, user.CountryCode),
				DateNow:               gtime.Now(),
				PeriodEnd:             gtime.Now().AddDate(0, 0, 5),
				PaymentAmount:         strconv.FormatInt(one.TotalAmount, 10),
				TokenExpireMinute:     strconv.FormatInt(config2.GetConfigInstance().Auth.Login.Expire/60, 10),
				Link:                  "<a href=\"" + link.GetInvoiceLink(one.InvoiceId, one.SendTerms) + "\">Link</a>",
				HttpLink:              link.GetInvoiceLink(one.InvoiceId, one.SendTerms),
				AccountHolder:         accountHolder,
				BIC:                   bic,
				IBAN:                  iban,
				Address:               address,
			})
			if err != nil {
				g.Log().Errorf(ctx, "SendTemplateEmail SendInvoiceEmailToUser err:%s", err.Error())
			} else {
				//update send status
				_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
					dao.Invoice.Columns().SendStatus: consts.InvoiceSendStatusSend,
					dao.Invoice.Columns().GmtModify:  gtime.Now(),
				}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
				if err != nil {
					fmt.Printf("SendInvoiceEmailToUser update err:%s", err.Error())
				}
			}
		} else {
			g.Log().Errorf(ctx, "SendInvoiceEmailToUser payment invoice status is pending or init, email not send")
		}
	} else {
		refund := query.GetRefundByRefundId(ctx, one.RefundId)
		if refund != nil {
			if one.Status > consts.InvoiceStatusPending {
				var template = email.TemplateInvoiceRefundCreated
				if one.Status == consts.InvoiceStatusProcessing {
					template = email.TemplateInvoiceRefundCreated
				} else if one.Status == consts.InvoiceStatusCancelled || one.Status == consts.InvoiceStatusFailed {
					template = email.TemplateInvoiceCancel
				} else if one.Status == consts.InvoiceStatusPaid {
					template = email.TemplateInvoiceRefundPaid
				} else {
					return nil
				}
				err := email.SendTemplateEmail(ctx, merchant.Id, one.SendEmail, user.TimeZone, user.Language, template, pdfFileName, &email.TemplateVariable{
					InvoiceId:             one.InvoiceId,
					UserName:              user.FirstName + " " + user.LastName,
					MerchantProductName:   one.ProductName,
					MerchantCustomerEmail: merchant.Email,
					MerchantName:          query.GetMerchantCountryConfigName(ctx, one.MerchantId, user.CountryCode),
					DateNow:               gtime.Now(),
					PeriodEnd:             gtime.Now().AddDate(0, 0, 5),
					RefundAmount:          utility.ConvertCentToDollarStr(refund.RefundAmount, refund.Currency),
					TokenExpireMinute:     strconv.FormatInt(config2.GetConfigInstance().Auth.Login.Expire/60, 10),
					Link:                  "<a href=\"" + link.GetInvoiceLink(one.InvoiceId, one.SendTerms) + "\">Link</a>",
					HttpLink:              link.GetInvoiceLink(one.InvoiceId, one.SendTerms),
				})
				if err != nil {
					g.Log().Errorf(ctx, "SendTemplateEmail SendInvoiceEmailToUser err:%s", err.Error())
				} else {
					//update send status
					_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
						dao.Invoice.Columns().SendStatus: consts.InvoiceSendStatusSend,
						dao.Invoice.Columns().GmtModify:  gtime.Now(),
					}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
					if err != nil {
						g.Log().Errorf(ctx, "SendInvoiceEmailToUser update err:%s", err.Error())
					}
				}
			} else {
				g.Log().Errorf(ctx, "SendInvoiceEmailToUser refund invoice status is pending or init, email not send")
			}
		}
	}

	return nil
}
