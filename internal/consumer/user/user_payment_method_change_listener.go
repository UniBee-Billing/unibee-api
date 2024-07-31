package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
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
				subs := query.GetLatestActiveOrIncompleteOrCreateSubscriptionsByUserId(ctx, user.Id, user.MerchantId)
				for _, sub := range subs {
					if sub != nil && (userGatewayId != sub.GatewayId || userPaymentMethod != sub.GatewayDefaultPaymentMethod) {
						_, _ = dao.Subscription.Ctx(ctx).Data(g.Map{
							dao.Subscription.Columns().GmtModify:                   gtime.Now(),
							dao.Subscription.Columns().GatewayId:                   userGatewayId,
							dao.Subscription.Columns().GatewayDefaultPaymentMethod: userPaymentMethod,
						}).Where(dao.Subscription.Columns().SubscriptionId, sub.SubscriptionId).OmitNil().Update()
						sub.GatewayId = userGatewayId
						sub.GatewayDefaultPaymentMethod = userPaymentMethod
					}
					if sub != nil && len(sub.LatestInvoiceId) > 0 {
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
						}
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
