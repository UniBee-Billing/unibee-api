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

type CreateProcessingInvoiceForSubReq struct {
	PlanId             uint64
	Simplify           *bean.Invoice
	Sub                *entity.Subscription
	GatewayId          uint64
	GatewayPaymentType string
	PaymentMethodId    string
	IsSubLatestInvoice bool
	TimeNow            int64
}

func CreateProcessingInvoiceForSub(ctx context.Context, req *CreateProcessingInvoiceForSubReq) (*entity.Invoice, error) {
	utility.Assert(req.Simplify != nil, "invoice data is nil")
	utility.Assert(req.Sub != nil, "sub is nil")
	user := query.GetUserAccountById(ctx, req.Sub.UserId)
	//Try cancel current sub processing invoice
	if req.IsSubLatestInvoice {
		TryCancelSubscriptionLatestInvoice(ctx, req.Sub)
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
	if req.TimeNow > currentTime {
		currentTime = req.TimeNow
	}

	invoiceId := utility.CreateInvoiceId()
	if len(req.Simplify.InvoiceId) > 0 {
		invoiceId = req.Simplify.InvoiceId
	}
	{
		//promo credit
		if req.Simplify.PromoCreditDiscountAmount > 0 && req.Simplify.PromoCreditPayout != nil && req.Simplify.PromoCreditAccount != nil {
			_, err := payment.NewCreditPayment(ctx, &payment.CreditPaymentInternalReq{
				UserId:                  req.Sub.UserId,
				MerchantId:              req.Sub.MerchantId,
				ExternalCreditPaymentId: invoiceId,
				InvoiceId:               invoiceId,
				CurrencyAmount:          req.Simplify.PromoCreditDiscountAmount,
				Currency:                req.Simplify.Currency,
				CreditType:              req.Simplify.PromoCreditAccount.Type,
				Name:                    "InvoicePromoCreditDiscount",
				Description:             "Subscription Invoice Promo Credit Discount",
			})
			if err != nil {
				return nil, err
			}
		}
	}
	if len(req.Simplify.DiscountCode) > 0 {
		_, err := discount.UserDiscountApply(ctx, &discount.UserDiscountApplyReq{
			MerchantId:       req.Sub.MerchantId,
			UserId:           req.Sub.UserId,
			DiscountCode:     req.Simplify.DiscountCode,
			SubscriptionId:   req.Sub.SubscriptionId,
			PLanId:           req.PlanId,
			PaymentId:        "",
			InvoiceId:        invoiceId,
			ApplyAmount:      req.Simplify.DiscountAmount,
			Currency:         req.Simplify.Currency,
			IsRecurringApply: strings.Compare(req.Simplify.CreateFrom, consts.InvoiceAutoChargeFlag) == 0,
		})
		if err != nil {
			_ = payment.RollbackCreditPayment(ctx, req.Sub.MerchantId, invoiceId)
			return nil, err
		}
	}

	status := consts.InvoiceStatusProcessing
	st := utility.CreateInvoiceSt()
	one := &entity.Invoice{
		SubscriptionId:                 req.Sub.SubscriptionId,
		BizType:                        consts.BizTypeSubscription,
		UserId:                         req.Sub.UserId,
		MerchantId:                     req.Sub.MerchantId,
		InvoiceName:                    req.Simplify.InvoiceName,
		ProductName:                    req.Simplify.ProductName,
		InvoiceId:                      invoiceId,
		PeriodStart:                    req.Simplify.PeriodStart,
		PeriodEnd:                      req.Simplify.PeriodEnd,
		PeriodStartTime:                gtime.NewFromTimeStamp(req.Simplify.PeriodStart),
		PeriodEndTime:                  gtime.NewFromTimeStamp(req.Simplify.PeriodEnd),
		Currency:                       req.Sub.Currency,
		GatewayId:                      req.GatewayId,
		GatewayInvoiceId:               req.GatewayPaymentType,
		GatewayPaymentMethod:           req.PaymentMethodId,
		Status:                         status,
		SendNote:                       req.Simplify.SendNote,
		SendStatus:                     req.Simplify.SendStatus,
		SendEmail:                      sendEmail,
		UniqueId:                       invoiceId,
		SendTerms:                      st,
		TotalAmount:                    req.Simplify.TotalAmount,
		TotalAmountExcludingTax:        req.Simplify.TotalAmountExcludingTax,
		TaxAmount:                      req.Simplify.TaxAmount,
		CountryCode:                    req.Simplify.CountryCode,
		VatNumber:                      req.Simplify.VatNumber,
		TaxPercentage:                  req.Simplify.TaxPercentage,
		SubscriptionAmount:             req.Simplify.SubscriptionAmount,
		SubscriptionAmountExcludingTax: req.Simplify.SubscriptionAmountExcludingTax,
		Lines:                          utility.MarshalToJsonString(req.Simplify.Lines),
		Link:                           link.GetInvoiceLink(invoiceId, st),
		CreateTime:                     gtime.Now().Timestamp(),
		FinishTime:                     currentTime,
		DayUtilDue:                     req.Simplify.DayUtilDue,
		DiscountAmount:                 req.Simplify.DiscountAmount,
		DiscountCode:                   req.Simplify.DiscountCode,
		TrialEnd:                       req.Simplify.TrialEnd,
		BillingCycleAnchor:             req.Simplify.BillingCycleAnchor,
		Data:                           utility.MarshalToJsonString(userSnapshot),
		MetaData:                       utility.MarshalToJsonString(req.Simplify.Metadata),
		CreateFrom:                     req.Simplify.CreateFrom,
		PromoCreditDiscountAmount:      req.Simplify.PromoCreditDiscountAmount,
		PartialCreditPaidAmount:        req.Simplify.PartialCreditPaidAmount,
		MetricCharge:                   utility.MarshalToJsonString(req.Simplify.UserMetricChargeForInvoice),
	}

	result, err := dao.Invoice.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		g.Log().Infof(ctx, "CreateProcessingInvoiceForSub Create Invoice failed subId:%s err:%s", req.Sub.SubscriptionId, err.Error())
		err = gerror.Newf(`CreateProcessingInvoiceForSub record insert failure %s`, err.Error())
		// should roll back discount usage
		rollbackErr := discount2.InvoiceRollbackAllDiscountsFromInvoice(ctx, invoiceId)
		if rollbackErr != nil {
			g.Log().Infof(ctx, "CreateProcessingInvoiceForSub InvoiceRollbackAllDiscountsFromInvoice rollback failed subId:%s err:%s", req.Sub.SubscriptionId, rollbackErr.Error())
		}
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))
	if req.IsSubLatestInvoice {
		_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().LatestInvoiceId: invoiceId,
		}).Where(dao.Subscription.Columns().SubscriptionId, req.Sub.SubscriptionId).OmitNil().Update()
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
