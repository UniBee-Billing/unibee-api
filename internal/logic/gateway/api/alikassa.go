package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"unibee/internal/cmd/config"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/gateway/api/log"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

type AliKassa struct{}

// AliKassa API response structures
type AliKassaPaymentResponse struct {
	Success bool `json:"success"`
	Data    struct {
		InvoiceId string `json:"invoice_id"`
		Status    string `json:"status"`
		Amount    int64  `json:"amount"`
		Currency  string `json:"currency"`
		Link      string `json:"link"`
	} `json:"data"`
	Error string `json:"error,omitempty"`
}

type AliKassaRefundResponse struct {
	Success bool `json:"success"`
	Data    struct {
		RefundId string `json:"refund_id"`
		Status   string `json:"status"`
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
	} `json:"data"`
	Error string `json:"error,omitempty"`
}

type AliKassaPaymentDetailResponse struct {
	Success bool `json:"success"`
	Data    struct {
		InvoiceId string `json:"invoice_id"`
		Status    string `json:"status"`
		Amount    int64  `json:"amount"`
		Currency  string `json:"currency"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	} `json:"data"`
	Error string `json:"error,omitempty"`
}

// https://doc-merchant.alikassa.com/
func (a AliKassa) GatewayInfo(ctx context.Context) *_interface.GatewayInfo {
	return &_interface.GatewayInfo{
		Name:                          "AliKassa",
		Description:                   "AliKassa payment gateway for Russian market",
		DisplayName:                   "AliKassa",
		GatewayWebsiteLink:            "https://alikassa.com/",
		GatewayWebhookIntegrationLink: "https://alikassa.com/",
		GatewayLogo:                   "https://api.unibee.dev/oss/file/dbzbak6zzgw1rnhhjh.svg",
		GatewayIcons:                  []string{"https://api.unibee.dev/oss/file/dbzbak6zzgw1rnhhjh.svg"},
		GatewayType:                   consts.GatewayTypeCard,
		Sort:                          230,
		AutoChargeEnabled:             false,
		PublicKeyName:                 "ShopId",
		PrivateSecretName:             "SecretKey",
		Host:                          "https://api.alikassa.com",
		IsStaging:                     true,
	}
}

func (a AliKassa) GatewayTest(ctx context.Context, req *_interface.GatewayTestReq) (icon string, gatewayType int64, err error) {
	utility.Assert(len(req.Key) > 0, "AliKassa shop ID is required")
	utility.Assert(len(req.Secret) > 0, "AliKassa secret key is required")

	// Test API connection
	_, err = a.makeAPICall(ctx, req.Key, req.Secret, "GET", "/test", nil)
	if err != nil {
		return "", 0, gerror.Newf("AliKassa API test failed: %v", err)
	}

	return "https://api.unibee.top/oss/file/alikassa-icon.png", consts.GatewayTypeCard, nil
}

func (a AliKassa) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	// AliKassa does not need to create users, return user ID as gateway user ID
	return &gateway_bean.GatewayUserCreateResp{
		GatewayUserId: fmt.Sprintf("user_%d", user.Id),
	}, nil
}

func (a AliKassa) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, gatewayUserId string) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	return &gateway_bean.GatewayUserDetailQueryResp{
		GatewayUserId: gatewayUserId,
	}, nil
}

func (a AliKassa) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	// AliKassa does not support balance queries, return empty balance
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

func (a AliKassa) GatewayNewPayment(ctx context.Context, gateway *entity.MerchantGateway, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
	// Build AliKassa payment request
	paymentData := map[string]interface{}{
		"amount":       createPayContext.Pay.PaymentAmount,
		"currency":     createPayContext.Pay.Currency,
		"order_id":     createPayContext.Pay.PaymentId,
		"description":  "Payment for order " + createPayContext.Pay.PaymentId,
		"return_url":   createPayContext.Pay.ReturnUrl,
		"fail_url":     createPayContext.Pay.ReturnUrl,
		"payment_type": "card", // AliKassa supports multiple payment types
	}

	// Call AliKassa API to create payment
	response, err := a.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", "/api/v1/invoice/create", paymentData)
	if err != nil {
		return nil, gerror.Newf("Failed to create AliKassa payment: %v", err)
	}

	// Parse response
	var aliKassaResp AliKassaPaymentResponse
	err = json.Unmarshal([]byte(response), &aliKassaResp)
	if err != nil {
		return nil, gerror.Newf("Failed to parse AliKassa response: %v", err)
	}

	if !aliKassaResp.Success {
		return nil, gerror.Newf("AliKassa payment creation failed: %s", aliKassaResp.Error)
	}

	// Save API log
	log.SaveChannelHttpLog("GatewayNewPayment", paymentData, response, nil, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	return &gateway_bean.GatewayNewPaymentResp{
		Payment:          createPayContext.Pay,
		Status:           a.mapAliKassaStatusToSystemStatus(aliKassaResp.Data.Status),
		GatewayPaymentId: aliKassaResp.Data.InvoiceId,
		Link:             aliKassaResp.Data.Link,
	}, nil
}

func (a AliKassa) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
	// Call AliKassa API to query payment details
	response, err := a.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "GET", fmt.Sprintf("/api/v1/invoice/%s", gatewayPaymentId), nil)
	if err != nil {
		return nil, gerror.Newf("Failed to get AliKassa payment detail: %v", err)
	}

	// Parse response
	var aliKassaResp AliKassaPaymentDetailResponse
	err = json.Unmarshal([]byte(response), &aliKassaResp)
	if err != nil {
		return nil, gerror.Newf("Failed to parse AliKassa response: %v", err)
	}

	if !aliKassaResp.Success {
		return nil, gerror.Newf("AliKassa payment detail failed: %s", aliKassaResp.Error)
	}

	// Save API log
	log.SaveChannelHttpLog("GatewayPaymentDetail", map[string]interface{}{"payment_id": gatewayPaymentId}, response, nil, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	return &gateway_bean.GatewayPaymentRo{
		GatewayPaymentId: gatewayPaymentId,
		Status:           int(a.mapAliKassaStatusToSystemStatus(aliKassaResp.Data.Status)),
		PaymentAmount:    aliKassaResp.Data.Amount,
		Currency:         aliKassaResp.Data.Currency,
		CreateTime:       a.parseAliKassaTime(aliKassaResp.Data.CreatedAt),
	}, nil
}

func (a AliKassa) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	// AliKassa does not support payment list queries, return empty list
	return []*gateway_bean.GatewayPaymentRo{}, nil
}

func (a AliKassa) GatewayCapture(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	// AliKassa payments are instant, no capture needed
	return &gateway_bean.GatewayPaymentCaptureResp{}, nil
}

func (a AliKassa) GatewayCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	// Call AliKassa API to cancel payment
	response, err := a.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", fmt.Sprintf("/api/v1/invoice/%s/cancel", payment.GatewayPaymentId), nil)
	if err != nil {
		return nil, gerror.Newf("Failed to cancel AliKassa payment: %v", err)
	}

	// Save API log
	log.SaveChannelHttpLog("GatewayCancel", map[string]interface{}{"payment_id": payment.GatewayPaymentId}, response, nil, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	return &gateway_bean.GatewayPaymentCancelResp{
		Status: consts.PaymentCancelled,
	}, nil
}

func (a AliKassa) GatewayRefund(ctx context.Context, gateway *entity.MerchantGateway, createPaymentRefundContext *gateway_bean.GatewayNewPaymentRefundReq) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	// Build refund request
	refundData := map[string]interface{}{
		"payment_id": createPaymentRefundContext.Payment.GatewayPaymentId,
		"amount":     createPaymentRefundContext.Refund.RefundAmount,
		"currency":   createPaymentRefundContext.Refund.Currency,
		"reason":     createPaymentRefundContext.Refund.RefundComment,
	}

	// Call AliKassa API to create refund
	response, err := a.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", "/api/v1/refund/create", refundData)
	if err != nil {
		return nil, gerror.Newf("Failed to create AliKassa refund: %v", err)
	}

	// Parse response
	var aliKassaResp AliKassaRefundResponse
	err = json.Unmarshal([]byte(response), &aliKassaResp)
	if err != nil {
		return nil, gerror.Newf("Failed to parse AliKassa response: %v", err)
	}

	if !aliKassaResp.Success {
		return &gateway_bean.GatewayPaymentRefundResp{
			GatewayRefundId: createPaymentRefundContext.Payment.GatewayPaymentId,
			Status:          consts.RefundFailed,
			Type:            consts.RefundTypeGateway,
			Reason:          fmt.Sprintf("AliKassa refund failed: %s", aliKassaResp.Error),
		}, nil
	}

	// Save API log
	log.SaveChannelHttpLog("GatewayRefund", refundData, response, nil, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: aliKassaResp.Data.RefundId,
		Status:          a.mapAliKassaRefundStatusToSystemStatus(aliKassaResp.Data.Status),
		RefundAmount:    aliKassaResp.Data.Amount,
		Currency:        aliKassaResp.Data.Currency,
		Type:            consts.RefundTypeGateway,
	}, nil
}

func (a AliKassa) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	// Call AliKassa API to query refund details
	response, err := a.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "GET", fmt.Sprintf("/api/v1/refund/%s", gatewayRefundId), nil)
	if err != nil {
		return nil, gerror.Newf("Failed to get AliKassa refund detail: %v", err)
	}

	// Parse response
	var aliKassaResp AliKassaRefundResponse
	err = json.Unmarshal([]byte(response), &aliKassaResp)
	if err != nil {
		return nil, gerror.Newf("Failed to parse AliKassa response: %v", err)
	}

	if !aliKassaResp.Success {
		return nil, gerror.Newf("AliKassa refund detail failed: %s", aliKassaResp.Error)
	}

	// Save API log
	log.SaveChannelHttpLog("GatewayRefundDetail", map[string]interface{}{"refund_id": gatewayRefundId}, response, nil, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: gatewayRefundId,
		Status:          a.mapAliKassaRefundStatusToSystemStatus(aliKassaResp.Data.Status),
		RefundAmount:    aliKassaResp.Data.Amount,
		Currency:        aliKassaResp.Data.Currency,
		Type:            consts.RefundTypeGateway,
	}, nil
}

func (a AliKassa) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	// AliKassa does not support refund list query, return empty list
	return []*gateway_bean.GatewayPaymentRefundResp{}, nil
}

func (a AliKassa) GatewayRefundCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	// Call AliKassa API to cancel refund
	response, err := a.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", fmt.Sprintf("/api/v1/refund/%s/cancel", refund.GatewayRefundId), nil)
	if err != nil {
		return nil, gerror.Newf("Failed to cancel AliKassa refund: %v", err)
	}

	// Save API log
	log.SaveChannelHttpLog("GatewayRefundCancel", map[string]interface{}{"refund_id": refund.GatewayRefundId}, response, nil, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: refund.GatewayRefundId,
		Status:          consts.RefundCancelled,
		RefundAmount:    refund.RefundAmount,
		Currency:        refund.Currency,
		Type:            consts.RefundTypeGateway,
	}, nil
}

// The following methods are not supported by AliKassa, return empty implementation
func (a AliKassa) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	return nil, fmt.Errorf("AliKassa does not support crypto transactions")
}

func (a AliKassa) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	return nil, fmt.Errorf("AliKassa does not support payment method management")
}

func (a AliKassa) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	return nil, fmt.Errorf("AliKassa does not support payment method management")
}

func (a AliKassa) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	return nil, fmt.Errorf("AliKassa does not support payment method management")
}

func (a AliKassa) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, currency string, metadata map[string]interface{}) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	return nil, fmt.Errorf("AliKassa does not support payment method management")
}

// Utility methods
func (a AliKassa) makeAPICall(ctx context.Context, apiKey, apiSecret, method, path string, data map[string]interface{}) (string, error) {
	// Determine base URL based on environment
	baseURL := "https://api.alikassa.com"
	if !config.GetConfigInstance().IsProd() {
		baseURL = "https://sandbox.api.alikassa.com"
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
	req.Header.Set("X-API-Secret", apiSecret)

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

	if resp.StatusCode != http.StatusOK {
		return "", gerror.Newf("AliKassa API error: %d - %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}

func (a AliKassa) mapAliKassaStatusToSystemStatus(aliKassaStatus string) consts.PaymentStatusEnum {
	switch strings.ToLower(aliKassaStatus) {
	case "pending":
		return consts.PaymentCreated
	case "succeeded", "completed":
		return consts.PaymentSuccess
	case "failed", "declined":
		return consts.PaymentFailed
	case "cancelled", "canceled":
		return consts.PaymentCancelled
	default:
		return consts.PaymentCreated
	}
}

func (a AliKassa) mapAliKassaRefundStatusToSystemStatus(aliKassaStatus string) consts.RefundStatusEnum {
	switch strings.ToLower(aliKassaStatus) {
	case "pending":
		return consts.RefundCreated
	case "succeeded", "completed":
		return consts.RefundSuccess
	case "failed", "declined":
		return consts.RefundFailed
	case "cancelled", "canceled":
		return consts.RefundCancelled
	default:
		return consts.RefundCreated
	}
}

func (a AliKassa) parseAliKassaTime(timeStr string) *gtime.Time {
	if timeStr == "" {
		return gtime.Now()
	}

	// Try to parse AliKassa time format
	t, err := time.Parse("2006-01-02T15:04:05Z", timeStr)
	if err != nil {
		// If parsing fails, return current time
		return gtime.Now()
	}

	return gtime.NewFromTime(t)
}
