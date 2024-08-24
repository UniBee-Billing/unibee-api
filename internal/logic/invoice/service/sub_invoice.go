package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	"unibee/api/bean"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	"unibee/internal/controller/link"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/invoice/handler"
	"unibee/internal/logic/operation_log"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func CreateProcessingInvoiceForSub(ctx context.Context, simplify *bean.Invoice, sub *entity.Subscription, gatewayId uint64, paymentMethodId string, isSubLatestInvoice bool, timeNow int64) (*entity.Invoice, error) {
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
