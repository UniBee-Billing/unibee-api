package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	"strings"
	"unibee/api/bean"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	"unibee/internal/controller/link"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/credit/payment"
	"unibee/internal/logic/discount"
	discount2 "unibee/internal/logic/invoice/discount"
	"unibee/internal/logic/invoice/handler"
	"unibee/internal/logic/operation_log"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func CreateProcessingInvoiceForSub(ctx context.Context, planId uint64, simplify *bean.Invoice, sub *entity.Subscription, gatewayId uint64, paymentMethodId string, isSubLatestInvoice bool, timeNow int64) (*entity.Invoice, error) {
	utility.Assert(simplify != nil, "invoice data is nil")
	utility.Assert(sub != nil, "sub is nil")
	user := query.GetUserAccountById(ctx, sub.UserId)
	//Try cancel current sub processing invoice
	if isSubLatestInvoice {
		TryCancelSubscriptionLatestInvoice(ctx, sub)
	}
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
	var currentTime = gtime.Now().Timestamp()
	if timeNow > currentTime {
		currentTime = timeNow
	}

	invoiceId := utility.CreateInvoiceId()
	if len(simplify.InvoiceId) > 0 {
		invoiceId = simplify.InvoiceId
	}
	{
		//promo credit
		if simplify.PromoCreditDiscountAmount > 0 && simplify.PromoCreditPayout != nil && simplify.PromoCreditAccount != nil {
			_, err := payment.NewCreditPayment(ctx, &payment.CreditPaymentInternalReq{
				UserId:                  sub.UserId,
				MerchantId:              sub.MerchantId,
				ExternalCreditPaymentId: invoiceId,
				InvoiceId:               invoiceId,
				CurrencyAmount:          simplify.PromoCreditDiscountAmount,
				Currency:                simplify.Currency,
				CreditType:              simplify.PromoCreditAccount.Type,
				Name:                    "InvoicePromoCreditDiscount",
				Description:             "Subscription Invoice Promo Credit Discount",
			})
			if err != nil {
				return nil, err
			}
		}
	}
	if len(simplify.DiscountCode) > 0 {
		_, err := discount.UserDiscountApply(ctx, &discount.UserDiscountApplyReq{
			MerchantId:       sub.MerchantId,
			UserId:           sub.UserId,
			DiscountCode:     simplify.DiscountCode,
			SubscriptionId:   sub.SubscriptionId,
			PLanId:           planId,
			PaymentId:        "",
			InvoiceId:        invoiceId,
			ApplyAmount:      simplify.DiscountAmount,
			Currency:         simplify.Currency,
			IsRecurringApply: strings.Compare(simplify.CreateFrom, consts.InvoiceAutoChargeFlag) == 0,
		})
		if err != nil {
			_ = payment.RollbackCreditPayment(ctx, sub.MerchantId, invoiceId)
			return nil, err
		}
	}

	status := consts.InvoiceStatusProcessing
	st := utility.CreateInvoiceSt()
	one := &entity.Invoice{
		SubscriptionId:                 sub.SubscriptionId,
		BizType:                        consts.BizTypeSubscription,
		UserId:                         sub.UserId,
		MerchantId:                     sub.MerchantId,
		InvoiceName:                    simplify.InvoiceName,
		ProductName:                    simplify.ProductName,
		InvoiceId:                      invoiceId,
		PeriodStart:                    simplify.PeriodStart,
		PeriodEnd:                      simplify.PeriodEnd,
		PeriodStartTime:                gtime.NewFromTimeStamp(simplify.PeriodStart),
		PeriodEndTime:                  gtime.NewFromTimeStamp(simplify.PeriodEnd),
		Currency:                       sub.Currency,
		GatewayId:                      gatewayId,
		GatewayPaymentMethod:           paymentMethodId,
		Status:                         status,
		SendNote:                       simplify.SendNote,
		SendStatus:                     simplify.SendStatus,
		SendEmail:                      sendEmail,
		UniqueId:                       invoiceId,
		SendTerms:                      st,
		TotalAmount:                    simplify.TotalAmount,
		TotalAmountExcludingTax:        simplify.TotalAmountExcludingTax,
		TaxAmount:                      simplify.TaxAmount,
		CountryCode:                    simplify.CountryCode,
		VatNumber:                      simplify.VatNumber,
		TaxPercentage:                  simplify.TaxPercentage,
		SubscriptionAmount:             simplify.SubscriptionAmount,
		SubscriptionAmountExcludingTax: simplify.SubscriptionAmountExcludingTax,
		Lines:                          utility.MarshalToJsonString(simplify.Lines),
		Link:                           link.GetInvoiceLink(invoiceId, st),
		CreateTime:                     gtime.Now().Timestamp(),
		FinishTime:                     currentTime,
		DayUtilDue:                     simplify.DayUtilDue,
		DiscountAmount:                 simplify.DiscountAmount,
		DiscountCode:                   simplify.DiscountCode,
		TrialEnd:                       simplify.TrialEnd,
		BillingCycleAnchor:             simplify.BillingCycleAnchor,
		Data:                           utility.MarshalToJsonString(userSnapshot),
		MetaData:                       utility.MarshalToJsonString(simplify.Metadata),
		CreateFrom:                     simplify.CreateFrom,
		PromoCreditDiscountAmount:      simplify.PromoCreditDiscountAmount,
		PartialCreditPaidAmount:        simplify.PartialCreditPaidAmount,
	}

	result, err := dao.Invoice.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		g.Log().Infof(ctx, "CreateProcessingInvoiceForSub Create Invoice failed subId:%s err:%s", sub.SubscriptionId, err.Error())
		err = gerror.Newf(`CreateProcessingInvoiceForSub record insert failure %s`, err.Error())
		// should roll back discount usage
		rollbackErr := discount2.InvoiceRollbackAllDiscountsFromInvoice(ctx, invoiceId)
		if rollbackErr != nil {
			g.Log().Infof(ctx, "CreateProcessingInvoiceForSub InvoiceRollbackAllDiscountsFromInvoice rollback failed subId:%s err:%s", sub.SubscriptionId, rollbackErr.Error())
		}
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))
	if isSubLatestInvoice {
		_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().LatestInvoiceId: invoiceId,
		}).Where(dao.Subscription.Columns().SubscriptionId, sub.SubscriptionId).OmitNil().Update()
		if err != nil {
			utility.AssertError(err, "CreateProcessingInvoiceForSub")
		}
	}
	if utility.TryLock(ctx, fmt.Sprintf("CreateProcessingInvoiceForSub_%s", one.InvoiceId), 60) {
		_, _ = redismq.Send(&redismq.Message{
			Topic:      redismq2.TopicInvoiceCreated.Topic,
			Tag:        redismq2.TopicInvoiceCreated.Tag,
			Body:       one.InvoiceId,
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
		_, _ = redismq.Send(&redismq.Message{
			Topic:      redismq2.TopicInvoiceProcessed.Topic,
			Tag:        redismq2.TopicInvoiceProcessed.Tag,
			Body:       one.InvoiceId,
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
	}
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Invoice(%s)", one.InvoiceId),
		Content:        "New",
		UserId:         one.UserId,
		SubscriptionId: one.SubscriptionId,
		InvoiceId:      one.InvoiceId,
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	//New Invoice Send Email
	_ = handler.InvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true, false)
	if err != nil {
		return nil, err
	}
	return one, nil
}
