package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq "github.com/jackyang-hk/go-redismq"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	"unibee/internal/logic/gateway/api"
	handler2 "unibee/internal/logic/payment/handler"
	"unibee/internal/query"
	"unibee/utility"
)

type RefundCheckerListener struct {
}

func (t RefundCheckerListener) GetTopic() string {
	return redismq2.TopicRefundChecker.Topic
}

func (t RefundCheckerListener) GetTag() string {
	return redismq2.TopicRefundChecker.Tag
}

func (t RefundCheckerListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Infof(ctx, "RefundCheckerListener Receive Message:%s", utility.MarshalToJsonString(message))
	if message.ReconsumeTimes > 50 {
		g.Log().Infof(ctx, "RefundCheckerListener_Commit by Reach Limit 50 paymentId:%s", message.Body)
		return redismq.CommitMessage
	}
	one := query.GetRefundByRefundId(ctx, message.Body)
	if one != nil {
		if one.Status == consts.RefundCreated {
			gateway := query.GetGatewayById(ctx, one.GatewayId)
			if gateway != nil && gateway.GatewayType != consts.GatewayTypeWireTransfer && len(one.GatewayRefundId) > 0 {
				gatewayRefundRo, err := api.GetGatewayServiceProvider(ctx, one.GatewayId).GatewayRefundDetail(ctx, gateway, one.GatewayRefundId, one)
				if err != nil {
					g.Log().Errorf(ctx, "RefundCheckerListener_Rollback refund:%s error:%s", message.Body, err.Error())
				} else {
					if gatewayRefundRo.Status == consts.RefundSuccess {
						err = handler2.HandleRefundSuccess(ctx, &handler2.HandleRefundReq{
							RefundId:         one.RefundId,
							GatewayRefundId:  gatewayRefundRo.GatewayRefundId,
							RefundAmount:     gatewayRefundRo.RefundAmount,
							RefundStatusEnum: gatewayRefundRo.Status,
							RefundTime:       gatewayRefundRo.RefundTime,
							Reason:           gatewayRefundRo.Reason,
						})
						if err != nil {
							g.Log().Errorf(ctx, "RefundCheckerListener_Rollback paymentId:%s HandleRefundSuccess error:%s", message.Body, err.Error())
						} else {
							g.Log().Infof(ctx, "RefundCheckerListener_Commit by HandleRefundSuccess refundId:%s", message.Body)
							return redismq.CommitMessage
						}
					} else if gatewayRefundRo.Status == consts.RefundFailed {
						err = handler2.HandleRefundFailure(ctx, &handler2.HandleRefundReq{
							RefundId:         one.RefundId,
							GatewayRefundId:  gatewayRefundRo.GatewayRefundId,
							RefundAmount:     gatewayRefundRo.RefundAmount,
							RefundStatusEnum: gatewayRefundRo.Status,
							RefundTime:       gatewayRefundRo.RefundTime,
							Reason:           gatewayRefundRo.Reason,
						})
						if err != nil {
							g.Log().Errorf(ctx, "RefundCheckerListener_Rollback paymentId:%s HandleRefundFailure error:%s", message.Body, err.Error())
						} else {
							g.Log().Infof(ctx, "RefundCheckerListener_Commit by HandleRefundFailure refundId:%s", message.Body)
							return redismq.CommitMessage
						}
					} else if gatewayRefundRo.Status == consts.RefundCancelled {
						err = handler2.HandleRefundCancelled(ctx, &handler2.HandleRefundReq{
							RefundId:         one.RefundId,
							GatewayRefundId:  gatewayRefundRo.GatewayRefundId,
							RefundAmount:     gatewayRefundRo.RefundAmount,
							RefundStatusEnum: gatewayRefundRo.Status,
							RefundTime:       gatewayRefundRo.RefundTime,
							Reason:           gatewayRefundRo.Reason,
						})
						if err != nil {
							g.Log().Errorf(ctx, "RefundCheckerListener_Rollback paymentId:%s HandleRefundCancelled error:%s", message.Body, err.Error())
						} else {
							g.Log().Infof(ctx, "RefundCheckerListener_Commit by HandleRefundCancelled refundId:%s", message.Body)
							return redismq.CommitMessage
						}
					} else if gatewayRefundRo.Status == consts.RefundReverse {
						err = handler2.HandleRefundReversed(ctx, &handler2.HandleRefundReq{
							RefundId:         one.RefundId,
							GatewayRefundId:  gatewayRefundRo.GatewayRefundId,
							RefundAmount:     gatewayRefundRo.RefundAmount,
							RefundStatusEnum: gatewayRefundRo.Status,
							RefundTime:       gatewayRefundRo.RefundTime,
							Reason:           gatewayRefundRo.Reason,
						})
						if err != nil {
							g.Log().Errorf(ctx, "RefundCheckerListener_Rollback paymentId:%s HandleRefundReversed error:%s", message.Body, err.Error())
						} else {
							g.Log().Infof(ctx, "RefundCheckerListener_Commit by HandleRefundReversed refundId:%s", message.Body)
							return redismq.CommitMessage
						}
					}
				}
				return redismq.ReconsumeLater
			} else {
				g.Log().Infof(ctx, "RefundCheckerListener_Commit by gateway not found or wire transfer or gatewayRefundId nil paymentId:%s", message.Body)
			}
		} else {
			g.Log().Infof(ctx, "RefundCheckerListener_Commit by status:%d paymentId:%s", one.Status, message.Body)
		}
	} else {
		g.Log().Infof(ctx, "RefundCheckerListener_Commit by can't find payment paymentId:%s", message.Body)
	}
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewRefundCheckerListener())
	fmt.Println("RefundCheckerListener RegisterListener")
}

func NewRefundCheckerListener() *RefundCheckerListener {
	return &RefundCheckerListener{}
}
