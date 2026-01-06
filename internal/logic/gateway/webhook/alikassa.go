package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type AliKassaWebhook struct{}

func (a AliKassaWebhook) GatewayCheckAndSetupWebhook(ctx context.Context, gateway *entity.MerchantGateway) (err error) {
	// AliKassa webhook setup is usually done on the frontend, here we only need to validate the configuration
	utility.Assert(len(gateway.WebhookSecret) > 0, "AliKassa webhook secret is required")
	return nil
}

func (a AliKassaWebhook) GatewayWebhook(r *ghttp.Request, gateway *entity.MerchantGateway) {
	ctx := r.Context()

	// 1. Read webhook data
	body, err := io.ReadAll(r.Body)
	if err != nil {
		g.Log().Errorf(ctx, "AliKassa webhook: failed to read body: %v", err)
		r.Response.WriteJsonExit(g.Map{"error": "failed to read body"})
		return
	}

	if len(body) == 0 {
		g.Log().Errorf(ctx, "AliKassa webhook: empty body")
		r.Response.WriteJsonExit(g.Map{"error": "empty body"})
		return
	}

	// 2. Verify signature
	if !a.verifyWebhookSignature(r, gateway, body) {
		g.Log().Errorf(ctx, "AliKassa webhook: invalid signature")
		r.Response.WriteJsonExit(g.Map{"error": "invalid signature"})
		return
	}

	// 3. Parse webhook data
	webhookData := gjson.New(body)
	eventType := webhookData.Get("event").String()

	g.Log().Infof(ctx, "AliKassa webhook received: %s", eventType)

	// 4. Handle different types of webhook events
	switch eventType {
	case "payment.finished":
		a.handlePaymentFinished(ctx, webhookData, gateway)
	case "refund.finished":
		a.handleRefundFinished(ctx, webhookData, gateway)
	default:
		g.Log().Infof(ctx, "AliKassa webhook: unhandled event type: %s", eventType)
	}

	// 5. Return success response
	r.Response.WriteJsonExit(g.Map{"status": "ok"})
}

func (a AliKassaWebhook) GatewayRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayRedirectResp, err error) {
	ctx := r.Context()

	// Handle AliKassa payment completion redirect
	paymentId := r.Get("payment_id").String()
	status := r.Get("status").String()

	g.Log().Infof(ctx, "AliKassa redirect: payment_id=%s, status=%s", paymentId, status)

	// Find corresponding payment record
	payment := query.GetPaymentByGatewayPaymentId(ctx, paymentId)
	if payment == nil {
		return nil, gerror.Newf("payment not found: %s", paymentId)
	}

	success := status == "succeeded"

	return &gateway_bean.GatewayRedirectResp{
		Payment:   payment,
		Status:    success,
		Success:   success,
		Message:   fmt.Sprintf("Payment %s", status),
		ReturnUrl: payment.ReturnUrl,
	}, nil
}

func (a AliKassaWebhook) GatewayNewPaymentMethodRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (err error) {
	// AliKassa does not support payment method management
	return gerror.New("AliKassa does not support payment method management")
}

// Utility functions
func (a AliKassaWebhook) verifyWebhookSignature(r *ghttp.Request, gateway *entity.MerchantGateway, body []byte) bool {
	// Get signature header
	signature := r.Header.Get("X-AliKassa-Signature")
	if signature == "" {
		return false
	}

	// Use webhook secret to verify signature
	h := hmac.New(sha256.New, []byte(gateway.WebhookSecret))
	h.Write(body)
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return signature == expectedSignature
}

func (a AliKassaWebhook) handlePaymentFinished(ctx context.Context, webhookData *gjson.Json, gateway *entity.MerchantGateway) {
	paymentId := webhookData.Get("payment.id").String()
	status := webhookData.Get("payment.status").String()

	g.Log().Infof(ctx, "AliKassa payment finished: payment_id=%s, status=%s", paymentId, status)

	// Find payment record
	payment := query.GetPaymentByGatewayPaymentId(ctx, paymentId)
	if payment == nil {
		g.Log().Errorf(ctx, "AliKassa webhook: payment not found: %s", paymentId)
		return
	}

	// Use existing webhook processing mechanism
	err := ProcessPaymentWebhook(ctx, payment.PaymentId, paymentId, gateway)
	if err != nil {
		g.Log().Errorf(ctx, "AliKassa webhook: failed to process payment: %v", err)
		return
	}
}

func (a AliKassaWebhook) handleRefundFinished(ctx context.Context, webhookData *gjson.Json, gateway *entity.MerchantGateway) {
	refundId := webhookData.Get("refund.id").String()
	status := webhookData.Get("refund.status").String()

	g.Log().Infof(ctx, "AliKassa refund finished: refund_id=%s, status=%s", refundId, status)

	// Find refund record
	refund := query.GetRefundByGatewayRefundId(ctx, refundId)
	if refund == nil {
		g.Log().Errorf(ctx, "AliKassa webhook: refund not found: %s", refundId)
		return
	}

	// Use existing webhook processing mechanism
	err := ProcessRefundWebhook(ctx, "refund.finished", refundId, gateway)
	if err != nil {
		g.Log().Errorf(ctx, "AliKassa webhook: failed to process refund: %v", err)
		return
	}
}
