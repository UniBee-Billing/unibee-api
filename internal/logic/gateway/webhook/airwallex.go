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

type AirwallexWebhook struct{}

func (a AirwallexWebhook) GatewayCheckAndSetupWebhook(ctx context.Context, gateway *entity.MerchantGateway) (err error) {
	// Airwallex webhook setup is usually done on the frontend, here we only need to validate the configuration
	utility.Assert(len(gateway.WebhookSecret) > 0, "Airwallex webhook secret is required")

	// Log webhook setup attempt
	g.Log().Infof(ctx, "Airwallex webhook setup validated for gateway %d", gateway.Id)

	return nil
}

func (a AirwallexWebhook) GatewayWebhook(r *ghttp.Request, gateway *entity.MerchantGateway) {
	ctx := r.Context()

	// 1. Read webhook data
	body, err := io.ReadAll(r.Body)
	if err != nil {
		g.Log().Errorf(ctx, "Airwallex webhook: failed to read body: %v", err)
		r.Response.WriteJsonExit(g.Map{"error": "failed to read body"})
		return
	}

	if len(body) == 0 {
		g.Log().Errorf(ctx, "Airwallex webhook: empty body")
		r.Response.WriteJsonExit(g.Map{"error": "empty body"})
		return
	}

	// 2. Verify signature
	if !a.verifyWebhookSignature(r, gateway, body) {
		g.Log().Errorf(ctx, "Airwallex webhook: invalid signature")
		r.Response.WriteJsonExit(g.Map{"error": "invalid signature"})
		return
	}

	// 3. Parse webhook data
	webhookData := gjson.New(body)
	eventType := webhookData.Get("name").String()

	g.Log().Infof(ctx, "Airwallex webhook received: %s", eventType)

	// 4. Handle different types of webhook events
	switch eventType {
	case "payment_intent.succeeded":
		a.handlePaymentSucceeded(ctx, webhookData, gateway)
	case "payment_intent.failed":
		a.handlePaymentFailed(ctx, webhookData, gateway)
	case "payment_intent.canceled":
		a.handlePaymentCanceled(ctx, webhookData, gateway)
	case "refund.succeeded":
		a.handleRefundSucceeded(ctx, webhookData, gateway)
	case "refund.failed":
		a.handleRefundFailed(ctx, webhookData, gateway)
	case "payment_method.attached":
		a.handlePaymentMethodAttached(ctx, webhookData, gateway)
	case "payment_method.detached":
		a.handlePaymentMethodDetached(ctx, webhookData, gateway)
	default:
		g.Log().Infof(ctx, "Airwallex webhook: unhandled event type: %s", eventType)
	}

	// 5. Return success response
	r.Response.WriteJsonExit(g.Map{"status": "ok"})
}

func (a AirwallexWebhook) GatewayRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayRedirectResp, err error) {
	ctx := r.Context()

	// Handle Airwallex payment completion redirect
	paymentIntentId := r.Get("payment_intent_id").String()
	status := r.Get("status").String()

	g.Log().Infof(ctx, "Airwallex redirect: payment_intent_id=%s, status=%s", paymentIntentId, status)

	// Find corresponding payment record
	payment := query.GetPaymentByGatewayPaymentId(ctx, paymentIntentId)
	if payment == nil {
		return nil, gerror.Newf("payment not found: %s", paymentIntentId)
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

func (a AirwallexWebhook) GatewayNewPaymentMethodRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (err error) {
	ctx := r.Context()

	// Handle Airwallex payment method setup redirect
	setupIntentId := r.Get("setup_intent_id").String()
	status := r.Get("status").String()

	g.Log().Infof(ctx, "Airwallex payment method setup: setup_intent_id=%s, status=%s", setupIntentId, status)

	if status == "succeeded" {
		// Payment method setup successful
		return nil
	} else {
		return gerror.Newf("Payment method setup failed: %s", status)
	}
}

// Utility functions
func (a AirwallexWebhook) verifyWebhookSignature(r *ghttp.Request, gateway *entity.MerchantGateway, body []byte) bool {
	// Get signature header (Airwallex uses X-Airwallex-Signature)
	signature := r.Header.Get("X-Airwallex-Signature")
	if signature == "" {
		// Try alternative header names
		signature = r.Header.Get("X-Signature")
	}

	if signature == "" {
		return false
	}

	// Use webhook secret to verify signature
	// Airwallex uses HMAC-SHA256 with the webhook secret
	h := hmac.New(sha256.New, []byte(gateway.WebhookSecret))
	h.Write(body)
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return signature == expectedSignature
}

func (a AirwallexWebhook) handlePaymentSucceeded(ctx context.Context, webhookData *gjson.Json, gateway *entity.MerchantGateway) {
	paymentIntentId := webhookData.Get("data.object.id").String()
	amount := webhookData.Get("data.object.amount").Int64()
	currency := webhookData.Get("data.object.currency").String()
	status := webhookData.Get("data.object.status").String()

	g.Log().Infof(ctx, "Airwallex payment succeeded: payment_intent_id=%s, amount=%d, currency=%s, status=%s",
		paymentIntentId, amount, currency, status)

	// Find payment record
	payment := query.GetPaymentByGatewayPaymentId(ctx, paymentIntentId)
	if payment == nil {
		g.Log().Errorf(ctx, "Airwallex webhook: payment not found: %s", paymentIntentId)
		return
	}

	// Use existing webhook processing mechanism
	err := ProcessPaymentWebhook(ctx, payment.PaymentId, paymentIntentId, gateway)
	if err != nil {
		g.Log().Errorf(ctx, "Airwallex webhook: failed to process payment: %v", err)
		return
	}
}

func (a AirwallexWebhook) handlePaymentFailed(ctx context.Context, webhookData *gjson.Json, gateway *entity.MerchantGateway) {
	paymentIntentId := webhookData.Get("data.object.id").String()
	status := webhookData.Get("data.object.status").String()
	lastPaymentError := webhookData.Get("data.object.last_payment_error.message").String()

	g.Log().Infof(ctx, "Airwallex payment failed: payment_intent_id=%s, status=%s, error=%s",
		paymentIntentId, status, lastPaymentError)

	// Find payment record
	payment := query.GetPaymentByGatewayPaymentId(ctx, paymentIntentId)
	if payment == nil {
		g.Log().Errorf(ctx, "Airwallex webhook: payment not found: %s", paymentIntentId)
		return
	}

	// Use existing webhook processing mechanism
	err := ProcessPaymentWebhook(ctx, payment.PaymentId, paymentIntentId, gateway)
	if err != nil {
		g.Log().Errorf(ctx, "Airwallex webhook: failed to process payment: %v", err)
		return
	}
}

func (a AirwallexWebhook) handlePaymentCanceled(ctx context.Context, webhookData *gjson.Json, gateway *entity.MerchantGateway) {
	paymentIntentId := webhookData.Get("data.object.id").String()
	status := webhookData.Get("data.object.status").String()

	g.Log().Infof(ctx, "Airwallex payment canceled: payment_intent_id=%s, status=%s",
		paymentIntentId, status)

	// Find payment record
	payment := query.GetPaymentByGatewayPaymentId(ctx, paymentIntentId)
	if payment == nil {
		g.Log().Errorf(ctx, "Airwallex webhook: payment not found: %s", paymentIntentId)
		return
	}

	// Use existing webhook processing mechanism
	err := ProcessPaymentWebhook(ctx, payment.PaymentId, paymentIntentId, gateway)
	if err != nil {
		g.Log().Errorf(ctx, "Airwallex webhook: failed to process payment: %v", err)
		return
	}
}

func (a AirwallexWebhook) handleRefundSucceeded(ctx context.Context, webhookData *gjson.Json, gateway *entity.MerchantGateway) {
	refundId := webhookData.Get("data.object.id").String()
	status := webhookData.Get("data.object.status").String()
	amount := webhookData.Get("data.object.amount").Int64()
	currency := webhookData.Get("data.object.currency").String()

	g.Log().Infof(ctx, "Airwallex refund succeeded: refund_id=%s, status=%s, amount=%d, currency=%s",
		refundId, status, amount, currency)

	// Find refund record
	refund := query.GetRefundByGatewayRefundId(ctx, refundId)
	if refund == nil {
		g.Log().Errorf(ctx, "Airwallex webhook: refund not found: %s", refundId)
		return
	}

	// Use existing webhook processing mechanism
	err := ProcessRefundWebhook(ctx, "refund.succeeded", refundId, gateway)
	if err != nil {
		g.Log().Errorf(ctx, "Airwallex webhook: failed to process refund: %v", err)
		return
	}
}

func (a AirwallexWebhook) handleRefundFailed(ctx context.Context, webhookData *gjson.Json, gateway *entity.MerchantGateway) {
	refundId := webhookData.Get("data.object.id").String()
	status := webhookData.Get("data.object.status").String()
	reason := webhookData.Get("data.object.reason").String()

	g.Log().Infof(ctx, "Airwallex refund failed: refund_id=%s, status=%s, reason=%s",
		refundId, status, reason)

	// Find refund record
	refund := query.GetRefundByGatewayRefundId(ctx, refundId)
	if refund == nil {
		g.Log().Errorf(ctx, "Airwallex webhook: refund not found: %s", refundId)
		return
	}

	// Use existing webhook processing mechanism
	err := ProcessRefundWebhook(ctx, "refund.failed", refundId, gateway)
	if err != nil {
		g.Log().Errorf(ctx, "Airwallex webhook: failed to process refund: %v", err)
		return
	}
}

func (a AirwallexWebhook) handlePaymentMethodAttached(ctx context.Context, webhookData *gjson.Json, gateway *entity.MerchantGateway) {
	paymentMethodId := webhookData.Get("data.object.id").String()
	customerId := webhookData.Get("data.object.customer_id").String()

	g.Log().Infof(ctx, "Airwallex payment method attached: payment_method_id=%s, customer_id=%s",
		paymentMethodId, customerId)

	// Log the event for audit purposes
	// In a real implementation, you might want to update user payment method records
}

func (a AirwallexWebhook) handlePaymentMethodDetached(ctx context.Context, webhookData *gjson.Json, gateway *entity.MerchantGateway) {
	paymentMethodId := webhookData.Get("data.object.id").String()
	customerId := webhookData.Get("data.object.customer_id").String()

	g.Log().Infof(ctx, "Airwallex payment method detached: payment_method_id=%s, customer_id=%s",
		paymentMethodId, customerId)

	// Log the event for audit purposes
	// In a real implementation, you might want to update user payment method records
}
