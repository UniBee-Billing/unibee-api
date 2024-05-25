package webhook

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"net/http"
	"strings"
	_gateway "unibee/internal/logic/gateway"
	"unibee/internal/logic/gateway/api/log"
	"unibee/internal/logic/gateway/api/paypal"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

type PaypalWebhook struct {
}

func init() {
	//注册 gateway_webhook_entry
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

// GatewayCheckAndSetupWebhook https://developer.paypal.com/docs/subscriptions/webhooks/
func (p PaypalWebhook) GatewayCheckAndSetupWebhook(ctx context.Context, gateway *entity.MerchantGateway) (err error) {
	utility.Assert(gateway != nil, "gateway is nil")
	client, _ := NewClient(gateway.GatewayKey, gateway.GatewaySecret, gateway.Host)
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
			URL: _gateway.GetPaymentWebhookEntranceUrl(gateway.Id),
			EventTypes: []paypal.WebhookEventType{
				{Name: "PAYMENT.SALE.COMPLETED"},
				{Name: "PAYMENT.SALE.REFUNDED"},
				{Name: "PAYMENT.SALE.REVERSED"},
			},
		}
		response, err := client.CreateWebhook(ctx, param)
		log.SaveChannelHttpLog("GatewayCheckAndSetupWebhook", param, response, err, "", nil, gateway)
		if err != nil {
			return err
		}
		if err != nil {
			return nil
		}
		//更新 secret
		//utility.Assert(len(result.Secret) > 0, "secret is nil")
		//err = query.UpdateGatewayWebhookSecret(ctx, gateway.Id, result.Secret)
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
					{Name: "PAYMENT.SALE.COMPLETED"},
					{Name: "PAYMENT.SALE.REFUNDED"},
					{Name: "PAYMENT.SALE.REVERSED"},
				},
			},
			{
				Operation: "replace",
				Path:      "/url",
				Value:     strings.Replace(_gateway.GetPaymentWebhookEntranceUrl(gateway.Id), "http://", "https://", 1), //paypal 只支持 https
			},
		}
		response, err := client.UpdateWebhook(ctx, webhook.ID, param)
		log.SaveChannelHttpLog("GatewayCheckAndSetupWebhook", param, response, err, webhook.ID, nil, gateway)
		if err != nil && strings.Compare(err.(*paypal.ErrorResponse).Name, "WEBHOOK_PATCH_REQUEST_NO_CHANGE") != 0 {
			//WEBHOOK_PATCH_REQUEST_NO_CHANGE 忽略没有更改的错误
			return err
		}
	}

	return nil
}

func (p PaypalWebhook) GatewayRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayRedirectResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p PaypalWebhook) GatewayWebhook(r *ghttp.Request, gateway *entity.MerchantGateway) {
	jsonData, err := r.GetJson()
	if err != nil {
		g.Log().Errorf(r.Context(), "⚠️  Webhook Gateway:%s, Webhook Get PortalJson failed. %v\n", gateway.GatewayName, err.Error())
		r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}
	client, _ := NewClient(gateway.GatewayKey, gateway.GatewaySecret, gateway.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}
	signature, err := client.VerifyWebhookSignature(r.Context(), r.Request, jsonData.Get("id").String())
	if err != nil {
		g.Log().Errorf(r.Context(), "⚠️  Webhook Gateway:%s, Webhook signature verification success\n", gateway.GatewayName)
		r.Response.WriteHeader(http.StatusBadRequest)
		return
	}
	if strings.Compare(signature.VerificationStatus, "SUCCESS") == 0 {
		g.Log().Info(r.Context(), "Receive_Webhook_Channel:", gateway.GatewayName, " hook:", jsonData.String())
		eventType := jsonData.Get("event_type").String()
		var responseBack = http.StatusOK
		switch eventType {
		default:
			g.Log().Errorf(r.Context(), "Webhook Gateway:%s, Unhandled event type: %s\n", gateway.GatewayName, eventType)
		}
		r.Response.WriteHeader(http.StatusOK)
		log.SaveChannelHttpLog("GatewayWebhook", jsonData, responseBack, err, "", nil, gateway)
		return
	} else {
		g.Log().Errorf(r.Context(), "⚠️  Webhook Gateway:%s, Webhook signature verification failed.\n", gateway.GatewayName)
		r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}
}
