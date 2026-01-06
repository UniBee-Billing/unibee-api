package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"unibee/internal/consts"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/api/log"
	"unibee/internal/logic/gateway/gateway_bean"
	"unibee/internal/logic/gateway/util"
	handler2 "unibee/internal/logic/payment/handler"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type FireKassaWebhook struct{}

func (f FireKassaWebhook) GatewayCheckAndSetupWebhook(ctx context.Context, gateway *entity.MerchantGateway) (err error) {
	// FireKassa webhook setup is done on the frontend
	// We only need to validate that the gateway key is configured
	utility.Assert(len(gateway.GatewayKey) > 0, "FireKassa site key is required for webhook signature verification")
	return nil
}

func (f FireKassaWebhook) GatewayWebhook(r *ghttp.Request, gateway *entity.MerchantGateway) {
	ctx := r.Context()

	// 1. Read webhook data
	body, err := io.ReadAll(r.Body)
	if err != nil {
		g.Log().Errorf(ctx, "FireKassa webhook: failed to read body: %v", err)
		r.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		g.Log().Error(ctx, "FireKassa webhook: empty body")
		r.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	// 2. Verify signature
	if !f.verifyWebhookSignature(r, gateway, body) {
		g.Log().Error(ctx, "FireKassa webhook: invalid signature")
		r.Response.WriteHeader(http.StatusUnauthorized)
		return
	}

	// 3. Parse webhook data
	webhookData := gjson.New(body)

	// 4. Handle different types of webhook events
	transactionType := webhookData.Get("type").String()

	g.Log().Infof(ctx, "FireKassa webhook received: type=%s", transactionType)

	switch transactionType {
	case "deposit":
		f.handleDepositFinished(ctx, webhookData, gateway)
	case "withdrawal":
		f.handleWithdrawalFinished(ctx, webhookData, gateway)
	default:
		g.Log().Warningf(ctx, "FireKassa webhook: unhandled transaction type: %s", transactionType)
	}

	// 5. Return success response as required by FireKassa
	// Must return exactly "OK" with HTTP 200 status
	r.Response.WriteHeader(http.StatusOK)
	r.Response.Write("OK")
}

func (f FireKassaWebhook) GatewayRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayRedirectResp, err error) {
	ctx := r.Context()

	payIdStr := r.Get("paymentId").String()
	var response string
	var status = false
	var returnUrl = ""
	var isSuccess = false
	var payment *entity.Payment

	if len(payIdStr) > 0 {
		response = ""
		// Payment Redirect
		payment = query.GetPaymentByPaymentId(ctx, payIdStr)
		if payment != nil {
			success := r.Get("success")
			if success != nil {
				if success.String() == "true" {
					isSuccess = true
				}
				returnUrl = util.GetPaymentRedirectUrl(ctx, payment, success.String())
			} else {
				returnUrl = util.GetPaymentRedirectUrl(ctx, payment, "")
			}
		}

		if r.Get("success").Bool() {
			if payment == nil || len(payment.GatewayPaymentIntentId) == 0 {
				response = "paymentId invalid"
			} else if len(payment.GatewayPaymentId) > 0 && payment.Status == consts.PaymentSuccess {
				response = "success"
				status = true
			} else {
				// Query payment status from FireKassa
				paymentIntentDetail, err := api.GetGatewayServiceProvider(ctx, gateway.Id).GatewayPaymentDetail(ctx, gateway, payment.GatewayPaymentId, payment)
				if err != nil {
					response = fmt.Sprintf("%v", err)
				} else {
					if paymentIntentDetail.Status == consts.PaymentSuccess {
						err := handler2.HandlePaySuccess(ctx, &handler2.HandlePayReq{
							PaymentId:              payment.PaymentId,
							GatewayPaymentIntentId: payment.GatewayPaymentIntentId,
							GatewayPaymentId:       paymentIntentDetail.GatewayPaymentId,
							GatewayUserId:          paymentIntentDetail.GatewayUserId,
							TotalAmount:            paymentIntentDetail.TotalAmount,
							PayStatusEnum:          consts.PaymentSuccess,
							PaidTime:               paymentIntentDetail.PaidTime,
							PaymentAmount:          paymentIntentDetail.PaymentAmount,
							Reason:                 paymentIntentDetail.Reason,
							GatewayPaymentMethod:   paymentIntentDetail.GatewayPaymentMethod,
							PaymentCode:            paymentIntentDetail.PaymentCode,
						})
						if err != nil {
							response = fmt.Sprintf("%v", err)
						} else {
							response = "payment success"
							status = true
						}
					} else if paymentIntentDetail.Status == consts.PaymentFailed {
						err := handler2.HandlePayFailure(ctx, &handler2.HandlePayReq{
							PaymentId:              payment.PaymentId,
							GatewayPaymentIntentId: payment.GatewayPaymentIntentId,
							GatewayPaymentId:       paymentIntentDetail.GatewayPaymentId,
							PayStatusEnum:          consts.PaymentFailed,
							Reason:                 paymentIntentDetail.Reason,
							PaymentCode:            paymentIntentDetail.PaymentCode,
						})
						if err != nil {
							response = fmt.Sprintf("%v", err)
						}
					}
				}
			}
		} else {
			response = "user cancelled"
		}
	}

	// Log the redirect event
	log.SaveChannelHttpLog("GatewayRedirect", r.URL, response, err, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	return &gateway_bean.GatewayRedirectResp{
		Payment:   payment,
		Status:    status,
		Message:   response,
		Success:   isSuccess,
		ReturnUrl: returnUrl,
		QueryPath: r.URL.RawQuery,
	}, nil
}

func (f FireKassaWebhook) GatewayNewPaymentMethodRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (err error) {
	// FireKassa does not support payment method management
	return gerror.New("FireKassa does not support payment method management")
}

// Utility functions

// verifyWebhookSignature verifies the webhook signature according to FireKassa documentation
// Algorithm: HMAC-SHA512 with sorted parameters + timestamp
func (f FireKassaWebhook) verifyWebhookSignature(r *ghttp.Request, gateway *entity.MerchantGateway, body []byte) bool {
	// Get signature header
	signature := r.Header.Get("X-Sign")
	if signature == "" {
		g.Log().Error(r.Context(), "FireKassa webhook: missing X-Sign header")
		return false
	}

	// Get timestamp header
	xTime := r.Header.Get("X-Time")
	if xTime == "" {
		g.Log().Error(r.Context(), "FireKassa webhook: missing X-Time header")
		return false
	}

	// Build parameters map using r.Get() for form data
	params := make(map[string]string)

	// Get all the required parameters according to FireKassa documentation
	requiredParams := []string{"id", "order_id", "type", "site_id", "amount", "currency", "commission", "account", "status", "error_code", "error"}

	for _, param := range requiredParams {
		value := r.Get(param).String()
		if value != "" {
			params[param] = value
		}
	}

	// Sort parameters by key (as required by FireKassa)
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build signature string according to FireKassa algorithm
	var signString string
	for _, key := range keys {
		value := params[key]
		if value != "" {
			// Convert to lowercase as per official documentation
			signString += strings.ToLower(value)
		}
	}
	// Add timestamp at the end
	signString += xTime

	// Generate HMAC-SHA512 signature using site key (GatewayKey)
	h := hmac.New(sha512.New, []byte(gateway.GatewayKey))
	h.Write([]byte(strings.ToLower(signString)))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	// Compare signatures
	isValid := signature == expectedSignature

	if !isValid {
		g.Log().Warningf(r.Context(), "FireKassa webhook signature verification failed. Expected: %s, Received: %s", expectedSignature, signature)
		g.Log().Debugf(r.Context(), "FireKassa webhook signature string: %s", signString)
	}

	return isValid
}

// handleDepositFinished handles deposit (payment) webhook events
func (f FireKassaWebhook) handleDepositFinished(ctx context.Context, webhookData *gjson.Json, gateway *entity.MerchantGateway) {
	// Extract webhook data according to official documentation
	transactionId := webhookData.Get("id").String()
	orderId := webhookData.Get("order_id").String()
	status := webhookData.Get("status").String()
	amount := webhookData.Get("amount").String()
	currency := webhookData.Get("currency").String()
	errorCode := webhookData.Get("error_code").String()
	errorMsg := webhookData.Get("error").String()

	g.Log().Infof(ctx, "FireKassa deposit webhook: id=%s, order_id=%s, status=%s, amount=%s, currency=%s",
		transactionId, orderId, status, amount, currency)

	if transactionId == "" {
		g.Log().Error(ctx, "FireKassa webhook: transaction ID not found")
		return
	}

	// Find payment record by gateway payment ID
	payment := query.GetPaymentByGatewayPaymentId(ctx, transactionId)
	if payment == nil {
		g.Log().Errorf(ctx, "FireKassa webhook: payment not found for transaction ID: %s", transactionId)
		return
	}

	// Log additional information
	if errorCode != "" || errorMsg != "" {
		g.Log().Warningf(ctx, "FireKassa webhook: transaction has errors - code: %s, message: %s", errorCode, errorMsg)
	}

	// Process the webhook using existing mechanism
	err := ProcessPaymentWebhook(ctx, payment.PaymentId, transactionId, gateway)
	if err != nil {
		g.Log().Errorf(ctx, "FireKassa webhook: failed to process payment: %v", err)
		return
	}

	g.Log().Infof(ctx, "FireKassa webhook: successfully processed deposit for payment ID: %s", payment.PaymentId)
}

// handleWithdrawalFinished handles withdrawal webhook events
func (f FireKassaWebhook) handleWithdrawalFinished(ctx context.Context, webhookData *gjson.Json, gateway *entity.MerchantGateway) {
	// Extract webhook data according to official documentation
	transactionId := webhookData.Get("id").String()
	orderId := webhookData.Get("order_id").String()
	status := webhookData.Get("status").String()
	amount := webhookData.Get("amount").String()
	currency := webhookData.Get("currency").String()
	errorCode := webhookData.Get("error_code").String()
	errorMsg := webhookData.Get("error").String()

	g.Log().Infof(ctx, "FireKassa withdrawal webhook: id=%s, order_id=%s, status=%s, amount=%s, currency=%s",
		transactionId, orderId, status, amount, currency)

	// For now, we just log withdrawal events
	// If you need to handle withdrawals, implement the logic here
	g.Log().Infof(ctx, "FireKassa webhook: withdrawal event received (not yet implemented)")

	// Log any errors
	if errorCode != "" || errorMsg != "" {
		g.Log().Warningf(ctx, "FireKassa webhook: withdrawal has errors - code: %s, message: %s", errorCode, errorMsg)
	}
}
