package webhook

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/gateway_bean"
	handler2 "unibee/internal/logic/payment/handler"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
)

func ProcessPaymentWebhook(ctx context.Context, metaPaymentId string, gatewayPaymentId string, gateway *entity.MerchantGateway) error {
	if len(metaPaymentId) > 0 {
		// PaymentIntent Under UniBee Control
		payment := query.GetPaymentByPaymentId(ctx, metaPaymentId)
		if payment != nil {
			paymentIntentDetail, err := api.GetGatewayServiceProvider(ctx, gateway.Id).GatewayPaymentDetail(ctx, gateway, gatewayPaymentId, payment)
			if err != nil {
				return err
			}
			err = handler2.HandlePaymentWebhookEvent(ctx, payment.PaymentId, &gateway_bean.GatewayPaymentRo{
				Status:               paymentIntentDetail.Status,
				AuthorizeStatus:      paymentIntentDetail.AuthorizeStatus,
				AuthorizeReason:      paymentIntentDetail.AuthorizeReason,
				Currency:             paymentIntentDetail.Currency,
				TotalAmount:          paymentIntentDetail.TotalAmount,
				PaymentAmount:        paymentIntentDetail.PaymentAmount,
				BalanceAmount:        paymentIntentDetail.BalanceAmount,
				BalanceStart:         paymentIntentDetail.BalanceStart,
				BalanceEnd:           paymentIntentDetail.BalanceEnd,
				Reason:               paymentIntentDetail.Reason,
				CancelReason:         paymentIntentDetail.CancelReason,
				PaymentData:          paymentIntentDetail.PaymentData,
				LastError:            paymentIntentDetail.LastError,
				PaymentCode:          paymentIntentDetail.PaymentCode,
				PaidTime:             paymentIntentDetail.PaidTime,
				CreateTime:           paymentIntentDetail.CreateTime,
				CancelTime:           paymentIntentDetail.CancelTime,
				GatewayPaymentId:     paymentIntentDetail.GatewayPaymentId,
				GatewayPaymentMethod: paymentIntentDetail.GatewayPaymentMethod,
				GatewayUserId:        paymentIntentDetail.GatewayUserId,
			})
			if err != nil {
				return err
			}
		} else {
			return gerror.New("Payment Not Found")
		}
	} else {
		//Maybe PaymentIntent Create By Invoice
		g.Log().Errorf(ctx, "No PaymentId Metadata PaymentIntentId:%s", gatewayPaymentId)
		return nil
	}
	return nil
}

func ProcessRefundWebhook(ctx context.Context, eventType string, gatewayRefundId string, gateway *entity.MerchantGateway) error {
	refundDetail, err := api.GetGatewayServiceProvider(ctx, gateway.Id).GatewayRefundDetail(ctx, gateway, gatewayRefundId, nil)
	if err != nil {
		return err
	}
	err = handler2.HandleRefundWebhookEvent(ctx, refundDetail)
	if err != nil {
		return err
	}

	return nil
}

func GetPaymentRedirectUrl(ctx context.Context, payment *entity.Payment, success string) string {
	if success == "false" {
		var metadata = make(map[string]string)
		if len(payment.MetaData) > 0 {
			err := gjson.Unmarshal([]byte(payment.MetaData), &metadata)
			if err != nil {
				fmt.Printf("SimplifyPayment Unmarshal Metadata error:%s", err.Error())
			}
		}
		cancelUrl := metadata["CancelUrl"]
		if cancelUrl != "" && len(cancelUrl) > 0 {
			return cancelUrl
		} else {
			return payment.ReturnUrl
		}
	} else {
		return payment.ReturnUrl
	}
}
