package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/payment/service"
	"unibee/internal/query"
	"unibee/redismq"
	"unibee/utility"
)

type UserPaymentMethodChangeListener struct {
}

func (t UserPaymentMethodChangeListener) GetTopic() string {
	return redismq2.TopicUserPaymentMethodChanged.Topic
}

func (t UserPaymentMethodChangeListener) GetTag() string {
	return redismq2.TopicUserPaymentMethodChanged.Tag
}

func (t UserPaymentMethodChangeListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Infof(ctx, "UserPaymentMethodChangeListener Receive Message:%s", utility.MarshalToJsonString(message))
	if len(message.Body) > 0 {
		userId, _ := strconv.ParseUint(message.Body, 10, 64)
		if userId > 0 {
			user := query.GetUserAccountById(ctx, userId)
			userGatewayId, _ := strconv.ParseUint(user.GatewayId, 10, 64)
			userPaymentMethod := user.PaymentMethod
			if user != nil && userGatewayId > 0 {
				sub := query.GetLatestActiveOrIncompleteOrCreateSubscriptionByUserId(ctx, user.Id, user.MerchantId)
				if sub != nil && (userGatewayId != sub.GatewayId || userPaymentMethod != sub.GatewayDefaultPaymentMethod) {
					_, _ = dao.Subscription.Ctx(ctx).Data(g.Map{
						dao.Subscription.Columns().GmtModify:                   gtime.Now(),
						dao.Subscription.Columns().GatewayId:                   userGatewayId,
						dao.Subscription.Columns().GatewayDefaultPaymentMethod: userPaymentMethod,
					}).Where(dao.Subscription.Columns().SubscriptionId, sub.SubscriptionId).OmitNil().Update()
					sub.GatewayId = userGatewayId
					sub.GatewayDefaultPaymentMethod = userPaymentMethod
				}
				if len(sub.LatestInvoiceId) > 0 {
					invoice := query.GetInvoiceByInvoiceId(ctx, sub.LatestInvoiceId)
					if invoice != nil && invoice.Status == consts.InvoiceStatusProcessing && (userGatewayId != invoice.GatewayId || userPaymentMethod != invoice.GatewayPaymentMethod) {
						if len(invoice.PaymentId) == 0 && len(invoice.PaymentLink) == 0 {
							_, _ = dao.Invoice.Ctx(ctx).Data(g.Map{
								dao.Invoice.Columns().GmtModify:            gtime.Now(),
								dao.Invoice.Columns().GatewayId:            userGatewayId,
								dao.Invoice.Columns().GatewayPaymentMethod: user.PaymentMethod,
							}).Where(dao.Invoice.Columns().InvoiceId, invoice.InvoiceId).OmitNil().Update()
						} else {
							gateway := query.GetGatewayById(ctx, invoice.GatewayId)
							if gateway != nil && gateway.GatewayType != consts.GatewayTypeCrypto {
								// try cancel old payment
								_, _ = dao.Invoice.Ctx(ctx).Data(g.Map{
									dao.Invoice.Columns().PaymentId:   "",
									dao.Invoice.Columns().PaymentLink: "",
								}).Where(dao.Invoice.Columns().InvoiceId, invoice.InvoiceId).OmitNil().Update()
								lastPayment := query.GetPaymentByPaymentId(ctx, invoice.PaymentId)
								if lastPayment != nil {
									err := service.PaymentGatewayCancel(ctx, lastPayment)
									if err != nil {
										g.Log().Print(ctx, "UserPaymentMethodChangeListener CancelLastPayment PaymentGatewayCancel:%s err:", lastPayment.PaymentId, err.Error())
									}
								}
								_, _ = dao.Invoice.Ctx(ctx).Data(g.Map{
									dao.Invoice.Columns().GmtModify:            gtime.Now(),
									dao.Invoice.Columns().GatewayId:            userGatewayId,
									dao.Invoice.Columns().GatewayPaymentMethod: user.PaymentMethod,
								}).Where(dao.Invoice.Columns().InvoiceId, invoice.InvoiceId).OmitNil().Update()
							}
						}
						//} else {
						//	// recreate invoice
						//	var lines []*bean.InvoiceItemSimplify
						//	err := utility.UnmarshalFromJsonString(invoice.Lines, &lines)
						//	if err != nil {
						//		g.Log().Errorf(ctx, "UserPaymentMethodChangeListener UnmarshalFromJsonString err:", err.Error())
						//		return redismq.ReconsumeLater
						//	}
						//	_, err = service.CreateProcessingInvoiceForSub(ctx, &bean.InvoiceSimplify{
						//		InvoiceName:                    invoice.InvoiceName,
						//		ProductName:                    invoice.ProductName,
						//		DiscountCode:                   invoice.DiscountCode,
						//		TotalAmount:                    invoice.TotalAmount,
						//		DiscountAmount:                 invoice.DiscountAmount,
						//		TotalAmountExcludingTax:        invoice.TotalAmountExcludingTax,
						//		Currency:                       invoice.Currency,
						//		TaxAmount:                      invoice.TaxAmount,
						//		TaxPercentage:                  invoice.TaxPercentage,
						//		SubscriptionAmount:             invoice.SubscriptionAmount,
						//		SubscriptionAmountExcludingTax: invoice.SubscriptionAmountExcludingTax,
						//		Lines:                          lines,
						//		PeriodEnd:                      invoice.PeriodEnd,
						//		PeriodStart:                    invoice.PeriodStart,
						//		FinishTime:                     invoice.FinishTime,
						//		SendNote:                       invoice.SendNote,
						//		SubscriptionId:                 invoice.SubscriptionId,
						//		BizType:                        invoice.BizType,
						//		CryptoAmount:                   invoice.CryptoAmount,
						//		CryptoCurrency:                 invoice.CryptoCurrency,
						//		SendStatus:                     invoice.SendStatus,
						//		DayUtilDue:                     invoice.DayUtilDue,
						//		TrialEnd:                       invoice.TrialEnd,
						//	}, sub)
						//	if err != nil {
						//		g.Log().Errorf(ctx, "UserPaymentMethodChangeListener CreateProcessingInvoiceForSub err:", err.Error())
						//		return redismq.ReconsumeLater
						//	}
						//}
					}
				}
			}
		}
	}

	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewUserPaymentMethodChangeListener())
	fmt.Println("UserPaymentMethodChangeListener RegisterListener")
}

func NewUserPaymentMethodChangeListener() *UserPaymentMethodChangeListener {
	return &UserPaymentMethodChangeListener{}
}
