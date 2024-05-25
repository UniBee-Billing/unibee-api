package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	"unibee/internal/logic/gateway/api"
	handler2 "unibee/internal/logic/payment/handler"
	"unibee/internal/logic/payment/service"
	"unibee/internal/query"
	"unibee/redismq"
	"unibee/utility"
)

type PaymentCheckerListener struct {
}

func (t PaymentCheckerListener) GetTopic() string {
	return redismq2.TopicPaymentChecker.Topic
}

func (t PaymentCheckerListener) GetTag() string {
	return redismq2.TopicPaymentChecker.Tag
}

func (t PaymentCheckerListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Infof(ctx, "PaymentCheckerListener Receive Message:%s", utility.MarshalToJsonString(message))
	if message.ReconsumeTimes > 50 {
		g.Log().Infof(ctx, "PaymentCheckerListener_Commit by Reach Limit 50 paymentId:%s", message.Body)
		return redismq.CommitMessage
	}
	one := query.GetPaymentByPaymentId(ctx, message.Body)
	if one != nil {
		if one.Status == consts.PaymentCreated {
			gateway := query.GetGatewayById(ctx, one.GatewayId)
			if gateway != nil && gateway.GatewayType != consts.GatewayTypeWireTransfer && len(one.GatewayPaymentId) > 0 {
				gatewayPaymentRo, err := api.GetGatewayServiceProvider(ctx, one.GatewayId).GatewayPaymentDetail(ctx, gateway, one.GatewayPaymentId, one)
				if err != nil {
					g.Log().Errorf(ctx, "PaymentCheckerListener_Rollback paymentId:%s error:%s", message.Body, err.Error())
				} else {
					if gatewayPaymentRo.Status == consts.PaymentSuccess {
						err := handler2.HandlePaySuccess(ctx, &handler2.HandlePayReq{
							PaymentId:              one.PaymentId,
							GatewayPaymentIntentId: gatewayPaymentRo.GatewayPaymentId,
							GatewayPaymentId:       gatewayPaymentRo.GatewayPaymentId,
							TotalAmount:            gatewayPaymentRo.TotalAmount,
							PayStatusEnum:          consts.PaymentSuccess,
							PaidTime:               gatewayPaymentRo.PaidTime,
							PaymentAmount:          gatewayPaymentRo.PaymentAmount,
							CaptureAmount:          0,
							Reason:                 gatewayPaymentRo.Reason,
							GatewayPaymentMethod:   gatewayPaymentRo.GatewayPaymentMethod,
							PaymentCode:            gatewayPaymentRo.PaymentCode,
						})
						if err != nil {
							g.Log().Errorf(ctx, "PaymentCheckerListener_Rollback paymentId:%s HandlePaySuccess error:%s", message.Body, err.Error())
						} else {
							g.Log().Infof(ctx, "PaymentCheckerListener_Commit by HandlePaySuccess paymentId:%s", message.Body)
							return redismq.CommitMessage
						}
					} else if gatewayPaymentRo.Status == consts.PaymentFailed {
						err := handler2.HandlePayFailure(ctx, &handler2.HandlePayReq{
							PaymentId:              one.PaymentId,
							GatewayPaymentIntentId: gatewayPaymentRo.GatewayPaymentId,
							GatewayPaymentId:       gatewayPaymentRo.GatewayPaymentId,
							PayStatusEnum:          consts.PaymentFailed,
							Reason:                 gatewayPaymentRo.Reason,
							PaymentCode:            gatewayPaymentRo.PaymentCode,
						})
						if err != nil {
							g.Log().Errorf(ctx, "PaymentCheckerListener_Rollback paymentId:%s HandlePayFailure error:%s", message.Body, err.Error())
						} else {
							g.Log().Infof(ctx, "PaymentCheckerListener_Commit by HandlePayFailure paymentId:%s", message.Body)
							return redismq.CommitMessage
						}
					} else if gatewayPaymentRo.Status == consts.PaymentCancelled {
						err := handler2.HandlePayCancel(ctx, &handler2.HandlePayReq{
							PaymentId:              one.PaymentId,
							GatewayPaymentIntentId: gatewayPaymentRo.GatewayPaymentId,
							GatewayPaymentId:       gatewayPaymentRo.GatewayPaymentId,
							PayStatusEnum:          consts.PaymentCancelled,
							Reason:                 gatewayPaymentRo.Reason,
							PaymentCode:            gatewayPaymentRo.PaymentCode,
						})
						if err != nil {
							g.Log().Errorf(ctx, "PaymentCheckerListener_Rollback paymentId:%s HandlePayCancel error:%s", message.Body, err.Error())
						} else {
							g.Log().Infof(ctx, "PaymentCheckerListener_Commit by HandlePayCancel paymentId:%s", message.Body)
							return redismq.CommitMessage
						}
					} else if one.AuthorizeStatus == consts.Authorized && gateway.GatewayType == consts.GatewayTypePaypal {
						err := service.PaymentGatewayCapture(ctx, one)
						if err != nil {
							g.Log().Errorf(ctx, "PaymentCheckerListener_Rollback PaymentGatewayCapture paymentId:%s error:%s", message.Body, err.Error())
						}
						return redismq.ReconsumeLater
					}
				}
				return redismq.ReconsumeLater
			} else {
				g.Log().Infof(ctx, "PaymentCheckerListener_Commit by gateway not found or wire transfer or gatewayPaymentId nil paymentId:%s", message.Body)
			}
		} else {
			g.Log().Infof(ctx, "PaymentCheckerListener_Commit by status:%d paymentId:%s", one.Status, message.Body)
		}
	} else {
		g.Log().Infof(ctx, "PaymentCheckerListener_Commit by can't find payment paymentId:%s", message.Body)
	}
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewPaymentCheckerListener())
	fmt.Println("PaymentCheckerListener RegisterListener")
}

func NewPaymentCheckerListener() *PaymentCheckerListener {
	return &PaymentCheckerListener{}
}
