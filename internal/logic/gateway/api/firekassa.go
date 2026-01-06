package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"unibee/internal/cmd/config"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	webhook2 "unibee/internal/logic/gateway"
	"unibee/internal/logic/gateway/api/log"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type FireKassa struct{}

// FireKassa API response structures
type FireKassaPaymentResponse struct {
	Id         string `json:"id"`
	Amount     string `json:"amount"`
	Commission string `json:"commission"`
	PaymentUrl string `json:"payment_url"`
	Error      string `json:"error"`
}

type FireKassaPaymentDetailResponse struct {
	Id               string `json:"id"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	Action           string `json:"action"`
	Method           string `json:"method"`
	Type             string `json:"type"`
	PaymentId        string `json:"payment_id"`
	OrderId          string `json:"order_id"`
	Account          string `json:"account"`
	Amount           string `json:"amount"`
	PaymentAmount    string `json:"payment_amount"`
	Commission       string `json:"commission"`
	Comment          string `json:"comment"`
	Status           string `json:"status"`
	PaymentUrl       string `json:"payment_url"`
	PaymentCode      string `json:"payment_code"`
	PaymentErrorCode string `json:"payment_error_code"`
	PaymentError     string `json:"payment_error"`
	Currency         string `json:"currency"`
	WalletNumber     string `json:"wallet_number"`
	CardNumber       string `json:"card_number"`
	CryptoAddress    string `json:"crypto_address"`
	CryptoBlockchain string `json:"crypto_blockchain"`
	Error            string `json:"error"`
}

// https://fkassa.gitbook.io/firekassa-api-v2/
func (f FireKassa) GatewayInfo(ctx context.Context) *_interface.GatewayInfo {
	return &_interface.GatewayInfo{
		Name:                          "FireKassa",
		Description:                   "FireKassa payment gateway for Russian market",
		DisplayName:                   "FireKassa",
		GatewayWebsiteLink:            "https://firekassa.com/",
		GatewayWebhookIntegrationLink: "https://fkassa.gitbook.io/firekassa-api-v2/overview/webhooks",
		GatewayLogo:                   "https://api.unibee.dev/oss/file/dbyr8bvtxg0r1qegb5.png",
		GatewayIcons:                  []string{"https://api.unibee.dev/oss/file/dbyr8bvtxg0r1qegb5.png"},
		GatewayType:                   consts.GatewayTypeCard,
		Sort:                          210,
		AutoChargeEnabled:             false,
		PublicKeyName:                 "API Bearer Token",
		PrivateSecretName:             "API Sign Token",
		Host:                          "https://api.firekassa.com",
		IsStaging:                     false,
	}
}

func (f FireKassa) GatewayTest(ctx context.Context, req *_interface.GatewayTestReq) (icon string, gatewayType int64, err error) {
	utility.Assert(len(req.Key) > 0, "FireKassa API Bearer Token is required")
	//utility.Assert(len(req.Secret) > 0, "FireKassa secret key is required")

	// Build FireKassa payment request
	paymentData := map[string]interface{}{
		"order_id": uuid.New().String(),
		"amount":   "100.00",
	}

	// Call FireKassa API to create payment
	response, err := f.makeAPICall(ctx, req.Key, req.Secret, "POST", "/api/v2/invoices", paymentData)

	if err != nil {
		return "", consts.GatewayTypeCard, gerror.Newf("Failed to create FireKassa payment: %v", err)
	}

	// Parse response
	var fireKassaResp FireKassaPaymentResponse
	err = json.Unmarshal([]byte(response), &fireKassaResp)
	if err != nil {
		return "", consts.GatewayTypeCard, gerror.Newf("Failed to parse FireKassa response: %v", err)
	}

	// Check for error in response
	if fireKassaResp.Error != "" {
		return "", consts.GatewayTypeCard, gerror.Newf("FireKassa payment creation failed: %s", fireKassaResp.Error)
	}

	// Check if ID is valid
	if fireKassaResp.Id == "" {
		return "", consts.GatewayTypeCard, gerror.New("invalid request, payment id is empty")
	}

	return "https://api.unibee.top/oss/file/firekassa-icon.png", consts.GatewayTypeCard, nil
}

func (f FireKassa) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	// FireKassa does not need to create users, return user ID as gateway user ID
	return &gateway_bean.GatewayUserCreateResp{
		GatewayUserId: fmt.Sprintf("user_%d", user.Id),
	}, nil
}

func (f FireKassa) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, gatewayUserId string) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	return &gateway_bean.GatewayUserDetailQueryResp{
		GatewayUserId: gatewayUserId,
	}, nil
}

func (f FireKassa) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	// FireKassa does not support balance queries, return empty balance
	balances := []*gateway_bean.GatewayBalance{
		{
			Amount:   0,
			Currency: "RUB",
		},
	}

	return &gateway_bean.GatewayMerchantBalanceQueryResp{
		AvailableBalance: balances,
		PendingBalance:   balances,
	}, nil
}

func (f FireKassa) GatewayNewPayment(ctx context.Context, gateway *entity.MerchantGateway, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
	// Get product name and description like Stripe does
	//productName, productDescription := createPayContext.GetInvoiceSingleProductNameAndDescription()

	// Build comment with meaningful product information
	//comment := productName
	//if productDescription != "" && productDescription != productName {
	//	comment = fmt.Sprintf("%s - %s", productName, productDescription)
	//}

	// Build FireKassa payment request
	paymentData := map[string]interface{}{
		"order_id": createPayContext.Pay.PaymentId,
		"amount":   utility.ConvertCentToDollarStr(createPayContext.Pay.TotalAmount, createPayContext.Pay.Currency),
		"currency": createPayContext.Pay.Currency,
		//"comment":          comment,
		"success_url":      strings.ReplaceAll(webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, true), "&session_id={CHECKOUT_SESSION_ID}", ""),
		"fail_url":         strings.ReplaceAll(webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, false), "&session_id={CHECKOUT_SESSION_ID}", ""),
		"notification_url": webhook2.GetPaymentWebhookEntranceUrl(gateway.Id),
	}

	// Call FireKassa API to create payment
	response, err := f.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", "/api/v2/invoices", paymentData)

	// Always save API log first, regardless of success or failure
	log.SaveChannelHttpLog("GatewayNewPayment", paymentData, response, err, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	if err != nil {
		return nil, gerror.Newf("Failed to create FireKassa payment: %v", err)
	}

	// Parse response
	var fireKassaResp FireKassaPaymentResponse
	err = json.Unmarshal([]byte(response), &fireKassaResp)
	if err != nil {
		return nil, gerror.Newf("Failed to parse FireKassa response: %v", err)
	}

	// Check for error in response
	if fireKassaResp.Error != "" {
		return nil, gerror.Newf("FireKassa payment creation failed: %s", fireKassaResp.Error)
	}

	// Check if ID is valid
	if fireKassaResp.Id == "" {
		return nil, gerror.New("invalid request, payment id is empty")
	}

	return &gateway_bean.GatewayNewPaymentResp{
		Payment:          createPayContext.Pay,
		Status:           f.mapFireKassaStatusToSystemStatus("process"), // 新创建的发票状态为 process
		GatewayPaymentId: fireKassaResp.Id,
		Link:             fireKassaResp.PaymentUrl,
	}, nil
}

func (f FireKassa) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
	// Call FireKassa API to query payment details
	response, err := f.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "GET", fmt.Sprintf("/api/v2/transactions/%s", gatewayPaymentId), nil)

	// Always save API log first, regardless of success or failure
	log.SaveChannelHttpLog("GatewayPaymentDetail", map[string]interface{}{"payment_id": gatewayPaymentId}, response, err, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	if err != nil {
		return nil, gerror.Newf("Failed to get FireKassa payment detail: %v", err)
	}

	// Parse response
	var fireKassaResp FireKassaPaymentDetailResponse
	err = json.Unmarshal([]byte(response), &fireKassaResp)
	if err != nil {
		return nil, gerror.Newf("Failed to parse FireKassa response: %v", err)
	}

	// Check for error in response
	if fireKassaResp.Error != "" {
		return nil, gerror.Newf("FireKassa payment detail failed: %s", fireKassaResp.Error)
	}

	return &gateway_bean.GatewayPaymentRo{
		GatewayPaymentId: gatewayPaymentId,
		Status:           int(f.mapFireKassaStatusToSystemStatus(fireKassaResp.Status)),
		PaymentAmount:    f.parseAmount(fireKassaResp.Amount, fireKassaResp.Currency),
		Currency:         fireKassaResp.Currency,
		CreateTime:       f.parseFireKassaTime(fireKassaResp.CreatedAt),
	}, nil
}

func (f FireKassa) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	// FireKassa does not support payment list queries, return empty list
	return []*gateway_bean.GatewayPaymentRo{}, nil
}

func (f FireKassa) GatewayCapture(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	// FireKassa payments are instant, no capture needed
	return &gateway_bean.GatewayPaymentCaptureResp{}, nil
}

func (f FireKassa) GatewayCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	// Call FireKassa API to cancel payment
	response, err := f.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", fmt.Sprintf("/api/v2/invoices/%s/cancel", payment.GatewayPaymentId), nil)

	// Always save API log first, regardless of success or failure
	log.SaveChannelHttpLog("GatewayCancel", map[string]interface{}{"payment_id": payment.GatewayPaymentId}, response, err, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	if err != nil {
		return nil, gerror.Newf("Failed to cancel FireKassa payment: %v", err)
	}

	return &gateway_bean.GatewayPaymentCancelResp{
		Status: consts.PaymentCancelled,
	}, nil
}

func (f FireKassa) GatewayRefund(ctx context.Context, gateway *entity.MerchantGateway, createPaymentRefundContext *gateway_bean.GatewayNewPaymentRefundReq) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	// FireKassa does not support refunds through API
	// Return marked refund status (manual refund required)
	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: createPaymentRefundContext.Refund.RefundId,
		Status:          consts.RefundCreated,
		Type:            consts.RefundTypeMarked,
	}, nil
}

func (f FireKassa) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	// FireKassa does not support refunds through API
	// Return marked refund status (manual refund required)
	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: gatewayRefundId,
		Status:          consts.RefundCreated,
		Type:            consts.RefundTypeMarked,
	}, nil
}

func (f FireKassa) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	// FireKassa does not support refund list query, return empty list
	return []*gateway_bean.GatewayPaymentRefundResp{}, nil
}

func (f FireKassa) GatewayRefundCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	// FireKassa does not support refunds through API
	// Return marked refund cancelled status
	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: refund.GatewayRefundId,
		Status:          consts.RefundCancelled,
		RefundAmount:    refund.RefundAmount,
		Currency:        refund.Currency,
		Type:            consts.RefundTypeMarked,
	}, nil
}

// The following methods are not supported by FireKassa, return empty implementation
func (f FireKassa) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	return nil, fmt.Errorf("FireKassa does not support crypto transactions")
}

func (f FireKassa) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	return nil, fmt.Errorf("FireKassa does not support payment method management")
}

func (f FireKassa) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	return nil, fmt.Errorf("FireKassa does not support payment method management")
}

func (f FireKassa) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	return nil, fmt.Errorf("FireKassa does not support payment method management")
}

func (f FireKassa) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, currency string, metadata map[string]interface{}) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	return nil, fmt.Errorf("FireKassa does not support payment method management")
}

// Utility methods
func (f FireKassa) makeAPICall(ctx context.Context, apiKey, apiSecret, method, path string, data map[string]interface{}) (string, error) {
	// Determine base URL based on environment
	baseURL := "https://admin.vanilapay.com"
	if !config.GetConfigInstance().IsProd() {
		baseURL = "https://admin.vanilapay.com"
	}
	url := baseURL + path

	var req *http.Request
	var err error

	if method == "GET" {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	} else {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		req, err = http.NewRequestWithContext(ctx, method, url, strings.NewReader(string(jsonData)))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	if err != nil {
		return "", err
	}

	// Add authentication headers
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Add signature header if apiSecret is provided (optional security feature)
	if apiSecret != "" {
		signature := f.generateSignature(method, path, data, apiSecret)
		req.Header.Set("Signature", signature)
	}

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response body correctly
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", gerror.Newf("Failed to read response body: %v", err)
	}

	// Accept both 200 (OK) and 201 (Created) as successful responses
	// 200: OK (GET requests, queries)
	// 201: Created (POST requests, resource creation)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", gerror.Newf("FireKassa API error: %d - %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}

func (f FireKassa) mapFireKassaStatusToSystemStatus(fireKassaStatus string) consts.PaymentStatusEnum {
	switch strings.ToLower(fireKassaStatus) {
	case "process":
		return consts.PaymentCreated // 处理中 -> 已创建
	case "paid":
		return consts.PaymentSuccess // 已支付 -> 支付成功
	case "expired":
		return consts.PaymentFailed // 已过期 -> 支付失败
	case "cancelled", "canceled":
		return consts.PaymentCancelled // 已取消 -> 已取消
	case "failed", "declined":
		return consts.PaymentFailed // 失败 -> 支付失败
	default:
		// 对于未知状态，记录日志并返回已创建状态
		g.Log().Warningf(context.Background(), "Unknown FireKassa status: %s, defaulting to PaymentCreated", fireKassaStatus)
		return consts.PaymentCreated
	}
}

func (f FireKassa) mapFireKassaRefundStatusToSystemStatus(fireKassaStatus string) consts.RefundStatusEnum {
	switch strings.ToLower(fireKassaStatus) {
	case "process":
		return consts.RefundCreated // 处理中 -> 已创建
	case "paid":
		return consts.RefundSuccess // 已支付 -> 退款成功
	case "expired":
		return consts.RefundFailed // 已过期 -> 退款失败
	case "cancelled", "canceled":
		return consts.RefundCancelled // 已取消 -> 已取消
	case "failed", "declined":
		return consts.RefundFailed // 失败 -> 退款失败
	default:
		// 对于未知状态，记录日志并返回已创建状态
		g.Log().Warningf(context.Background(), "Unknown FireKassa refund status: %s, defaulting to RefundCreated", fireKassaStatus)
		return consts.RefundCreated
	}
}

func (f FireKassa) parseFireKassaTime(timeStr string) *gtime.Time {
	if timeStr == "" {
		return gtime.Now()
	}

	// Try to parse FireKassa time format
	t, err := time.Parse("2006-01-02T15:04:05Z", timeStr)
	if err != nil {
		// If parsing fails, return current time
		return gtime.Now()
	}

	return gtime.NewFromTime(t)
}

func (f FireKassa) parseAmount(amountStr string, currency string) int64 {
	// FireKassa returns amount as a string (in dollars/rubles), convert to int64 cents
	return utility.ConvertDollarStrToCent(amountStr, currency)
}

// generateSignature generates HMAC-SHA512 signature for request security
// According to FireKassa API v2 security documentation
func (f FireKassa) generateSignature(method, path string, data map[string]interface{}, secret string) string {
	// Create string from request path and body (if exists)
	var payload string

	if data != nil && len(data) > 0 {
		// Convert data to JSON string
		jsonData, err := json.Marshal(data)
		if err == nil {
			payload = path + string(jsonData)
		} else {
			// If JSON marshaling fails, use path only
			payload = path
		}
	} else {
		// No body data, use path only
		payload = path
	}

	// Generate HMAC-SHA512 signature
	h := hmac.New(sha512.New, []byte(secret))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}
