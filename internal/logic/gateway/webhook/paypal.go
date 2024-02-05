package webhook

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/plutov/paypal/v4"
	"go-oversea-pay/internal/logic/gateway"
	"go-oversea-pay/internal/logic/gateway/api"
	"go-oversea-pay/internal/logic/gateway/api/log"
	"go-oversea-pay/internal/logic/gateway/ro"
	"go-oversea-pay/internal/logic/subscription/handler"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"net/http"
	"strings"
)

type PaypalWebhook struct {
}

func init() {
	//注册 channel_webhook_entry
}

// todo mark 确认改造成单例是否可行，不用每次都去获取 accessToken
func NewClient(clientID string, secret string, APIBase string) (*paypal.Client, error) {
	if clientID == "" || secret == "" || APIBase == "" {
		return nil, errors.New("ClientID, Secret and APIBase are required to create a Client")
	}

	return &paypal.Client{
		Client:   &http.Client{},
		ClientID: clientID,
		Secret:   secret,
		APIBase:  APIBase,
	}, nil
}

// DoRemoteChannelCheckAndSetupWebhook https://developer.paypal.com/docs/subscriptions/webhooks/
func (p PaypalWebhook) DoRemoteChannelCheckAndSetupWebhook(ctx context.Context, payChannel *entity.MerchantChannelConfig) (err error) {
	utility.Assert(payChannel != nil, "payChannel is nil")
	client, _ := NewClient(payChannel.ChannelKey, payChannel.ChannelSecret, payChannel.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return err
	}
	result, err := client.ListWebhooks(ctx, paypal.AncorTypeApplication)
	if err != nil {
		return err
	}
	if len(result.Webhooks) > 1 {
		return gerror.New("webhook endpoints count > 1")
	}
	//过滤不可用
	if len(result.Webhooks) == 0 {
		//创建
		param := &paypal.CreateWebhookRequest{
			URL: gateway.GetPaymentWebhookEntranceUrl(int64(payChannel.Id)),
			EventTypes: []paypal.WebhookEventType{
				{Name: "BILLING.SUBSCRIPTION.CREATED"},
				{Name: "BILLING.SUBSCRIPTION.ACTIVATED"},
				{Name: "BILLING.SUBSCRIPTION.UPDATED"},
				{Name: "BILLING.SUBSCRIPTION.EXPIRED"},
				{Name: "BILLING.SUBSCRIPTION.CANCELLED"},
				{Name: "BILLING.SUBSCRIPTION.SUSPENDED"},
				{Name: "BILLING.SUBSCRIPTION.PAYMENT.FAILED"},
				{Name: "PAYMENT.SALE.COMPLETED"},
				{Name: "PAYMENT.SALE.REFUNDED"},
				{Name: "PAYMENT.SALE.REVERSED"},
			},
		}
		response, err := client.CreateWebhook(ctx, param)
		log.SaveChannelHttpLog("DoRemoteChannelCheckAndSetupWebhook", param, response, err, "", nil, payChannel)
		if err != nil {
			return err
		}
		if err != nil {
			return nil
		}
		//更新 secret
		//utility.Assert(len(result.Secret) > 0, "secret is nil")
		//err = query.UpdatePayChannelWebhookSecret(ctx, int64(payChannel.Id), result.Secret)
		//if err != nil {
		//	return err
		//}
	} else {
		utility.Assert(len(result.Webhooks) == 1, "internal webhook update, count is not 1")
		//检查并更新, todo mark 优化检查逻辑，如果 evert 一致不用发起更新
		webhook := result.Webhooks[0]
		//utility.Assert(strings.Compare(result.Status, "enabled") == 0, "webhook not status enabled after updated")// todo mark 需要检查里面的每一项
		param := []paypal.WebhookField{
			{
				Operation: "replace",
				Path:      "/event_types",
				Value: []paypal.WebhookEventType{
					{Name: "BILLING.SUBSCRIPTION.CREATED"},
					{Name: "BILLING.SUBSCRIPTION.ACTIVATED"},
					{Name: "BILLING.SUBSCRIPTION.UPDATED"},
					{Name: "BILLING.SUBSCRIPTION.EXPIRED"},
					{Name: "BILLING.SUBSCRIPTION.CANCELLED"},
					{Name: "BILLING.SUBSCRIPTION.SUSPENDED"},
					{Name: "BILLING.SUBSCRIPTION.PAYMENT.FAILED"},
					{Name: "PAYMENT.SALE.COMPLETED"},
					{Name: "PAYMENT.SALE.REFUNDED"},
					{Name: "PAYMENT.SALE.REVERSED"},
				},
			},
			{
				Operation: "replace",
				Path:      "/url",
				Value:     strings.Replace(gateway.GetPaymentWebhookEntranceUrl(int64(payChannel.Id)), "http://", "https://", 1), //paypal 只支持 https
			},
		}
		response, err := client.UpdateWebhook(ctx, webhook.ID, param)
		log.SaveChannelHttpLog("DoRemoteChannelCheckAndSetupWebhook", param, response, err, webhook.ID, nil, payChannel)
		if err != nil && strings.Compare(err.(*paypal.ErrorResponse).Name, "WEBHOOK_PATCH_REQUEST_NO_CHANGE") != 0 {
			//WEBHOOK_PATCH_REQUEST_NO_CHANGE 忽略没有更改的错误
			return err
		}
	}

	return nil
}

func (p PaypalWebhook) DoRemoteChannelRedirect(r *ghttp.Request, payChannel *entity.MerchantChannelConfig) (res *ro.ChannelRedirectInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p PaypalWebhook) processWebhook(ctx context.Context, eventType string, resource *gjson.Json) error {
	unibSub := query.GetSubscriptionByChannelSubscriptionId(ctx, resource.Get("id").String())
	if unibSub != nil {
		plan := query.GetPlanById(ctx, unibSub.PlanId)
		planChannel := query.GetPlanChannel(ctx, unibSub.PlanId, unibSub.ChannelId)
		details, err := api.GetPayChannelServiceProvider(ctx, unibSub.ChannelId).DoRemoteChannelSubscriptionDetails(ctx, plan, planChannel, unibSub)
		if err != nil {
			return err
		}

		err = handler.HandleSubscriptionWebhookEvent(ctx, unibSub, eventType, details)
		if err != nil {
			return err
		}
		return nil
	} else {
		return gerror.New("subscription not found on channelSubId:" + resource.Get("id").String())
	}
}

func (p PaypalWebhook) DoRemoteChannelWebhook(r *ghttp.Request, payChannel *entity.MerchantChannelConfig) {
	jsonData, err := r.GetJson()
	if err != nil {
		g.Log().Errorf(r.Context(), "⚠️  Webhook Channel:%s, Webhook Get Json failed. %v\n", payChannel.Channel, err.Error())
		r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}
	client, _ := NewClient(payChannel.ChannelKey, payChannel.ChannelSecret, payChannel.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}
	signature, err := client.VerifyWebhookSignature(r.Context(), r.Request, jsonData.Get("id").String())
	if err != nil {
		g.Log().Errorf(r.Context(), "⚠️  Webhook Channel:%s, Webhook signature verification success\n", payChannel.Channel)
		r.Response.WriteHeader(http.StatusBadRequest)
		return
	}
	if strings.Compare(signature.VerificationStatus, "SUCCESS") == 0 {
		g.Log().Info(r.Context(), "Receive_Webhook_Channel:", payChannel.Channel, " hook:", jsonData.String())
		eventType := jsonData.Get("event_type").String()
		var responseBack = http.StatusOK
		switch eventType {
		case "BILLING.SUBSCRIPTION.EXPIRED":
			resource := jsonData.GetJson("resource")
			if resource == nil || !resource.Contains("id") {
				g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error parsing webhook resource is nil\n", payChannel.Channel)
				r.Response.WriteHeader(http.StatusBadRequest)
				responseBack = http.StatusBadRequest
			} else {
				g.Log().Infof(r.Context(), "Webhook Channel:%s, Subscription deleted for %d.", payChannel.Channel, resource.Get("id").String())
				// Then define and call a func to handle the deleted subscription.
				// handleSubscriptionCanceled(subscription)
				err := p.processWebhook(r.Context(), eventType, resource)
				if err != nil {
					g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error HandleSubscriptionWebhookEvent: %v\n", payChannel.Channel, err.Error())
					r.Response.WriteHeader(http.StatusBadRequest)
					responseBack = http.StatusBadRequest
				}
			}
		case "BILLING.SUBSCRIPTION.UPDATED":
			resource := jsonData.GetJson("resource")
			if resource == nil || !resource.Contains("id") {
				g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error parsing webhook resource is nil\n", payChannel.Channel)
				r.Response.WriteHeader(http.StatusBadRequest)
				responseBack = http.StatusBadRequest
			} else {
				g.Log().Infof(r.Context(), "Webhook Channel:%s, Subscription updated for %d.", payChannel.Channel, resource.Get("id").String())
				// Then define and call a func to handle the successful attachment of a ChannelDefaultPaymentMethod.
				// handleSubscriptionUpdated(subscription)
				err := p.processWebhook(r.Context(), eventType, resource)
				if err != nil {
					g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error HandleSubscriptionWebhookEvent: %v\n", payChannel.Channel, err.Error())
					r.Response.WriteHeader(http.StatusBadRequest)
					responseBack = http.StatusBadRequest
				}
			}
		case "BILLING.SUBSCRIPTION.CREATED":
			resource := jsonData.GetJson("resource")
			if resource == nil || !resource.Contains("id") {
				g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error parsing webhook resource is nil\n", payChannel.Channel)
				r.Response.WriteHeader(http.StatusBadRequest)
				responseBack = http.StatusBadRequest
			} else {
				g.Log().Infof(r.Context(), "Webhook Channel:%s, Subscription created for %d.", payChannel.Channel, resource.Get("id").String())
				// Then define and call a func to handle the successful attachment of a ChannelDefaultPaymentMethod.
				// handleSubscriptionCreated(subscription)
				err := p.processWebhook(r.Context(), eventType, resource)
				if err != nil {
					g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error HandleSubscriptionWebhookEvent: %v\n", payChannel.Channel, err.Error())
					r.Response.WriteHeader(http.StatusBadRequest)
					responseBack = http.StatusBadRequest
				}
			}
		default:
			g.Log().Errorf(r.Context(), "Webhook Channel:%s, Unhandled event type: %s\n", payChannel.Channel, eventType)
		}
		r.Response.WriteHeader(http.StatusOK)
		log.SaveChannelHttpLog("DoRemoteChannelWebhook", jsonData, responseBack, err, "", nil, payChannel)
		return
	} else {
		g.Log().Errorf(r.Context(), "⚠️  Webhook Channel:%s, Webhook signature verification failed. %v\n", payChannel.Channel)
		r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}
}
