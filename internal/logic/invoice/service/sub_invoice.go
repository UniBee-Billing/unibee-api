package service

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	"unibee/internal/controller/link"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/invoice/handler"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/redismq"
	"unibee/utility"
)

func CreateProcessingInvoiceForSub(ctx context.Context, simplify *bean.InvoiceSimplify, sub *entity.Subscription, gatewayId uint64, paymentMethodId string, isSubLatestInvoice bool, timeNow int64) (*entity.Invoice, error) {
	utility.Assert(simplify != nil, "invoice data is nil")
	utility.Assert(sub != nil, "sub is nil")
	user := query.GetUserAccountById(ctx, sub.UserId)
	//Try cancel current sub processing invoice
	if isSubLatestInvoice {
		TryCancelSubscriptionLatestInvoice(ctx, sub)
	}
	var sendEmail = ""
	if user != nil {
		sendEmail = user.Email
	}
	var currentTime = gtime.Now().Timestamp()
	if timeNow > currentTime {
		currentTime = timeNow
	}
	status := consts.InvoiceStatusProcessing
	invoiceId := utility.CreateInvoiceId()
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
		CountryCode:                    sub.CountryCode,
		BillingCycleAnchor:             simplify.BillingCycleAnchor,
		CreateFrom:                     simplify.CreateFrom,
		MetaData:                       utility.MarshalToJsonString(simplify.Metadata),
	}

	result, err := dao.Invoice.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`CreateProcessingInvoiceForSub record insert failure %s`, err.Error())
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
	//New Invoice Send Email
	_ = handler.InvoicePdfGenerateAndEmailSendBackground(one.InvoiceId, true, false)
	if err != nil {
		return nil, err
	}
	return one, nil
}
