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

type SberPay struct{}

// SberPay API response structure
type SberPayPaymentResponse struct {
	Success bool `json:"success"`
	Data    struct {
		PaymentId string `json:"payment_id"`
		Status    string `json:"status"`
		Amount    int64  `json:"amount"`
		Currency  string `json:"currency"`
		Link      string `json:"link"`
	} `json:"data"`
	Error string `json:"error,omitempty"`
}

type SberPayRefundResponse struct {
	Success bool `json:"success"`
	Data    struct {
		RefundId string `json:"refund_id"`
		Status   string `json:"status"`
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
	} `json:"data"`
	Error string `json:"error,omitempty"`
}

type SberPayPaymentDetailResponse struct {
	Success bool `json:"success"`
	Data    struct {
		PaymentId string `json:"payment_id"`
		Status    string `json:"status"`
		Amount    int64  `json:"amount"`
		Currency  string `json:"currency"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	} `json:"data"`
	Error string `json:"error,omitempty"`
}

// https://ecomtest.sberbank.ru/en/doc， https://www.sberbank.ru/en/individualclients
func (s SberPay) GatewayInfo(ctx context.Context) *_interface.GatewayInfo {
	return &_interface.GatewayInfo{
		Name:                          "SberPay",
		Description:                   "SberPay mobile payment gateway for Russian market",
		DisplayName:                   "SberPay",
		GatewayWebsiteLink:            "https://www.sberbank.ru/",
		GatewayWebhookIntegrationLink: "https://www.sberbank.ru/",
		GatewayLogo:                   "https://api.unibee.dev/oss/file/dbyr5ts5ykzlbbn5ly.png",
		GatewayIcons:                  []string{"https://api.unibee.dev/oss/file/dbyr5ts5ykzlbbn5ly.png"},
		GatewayType:                   consts.GatewayTypeCard,
		Sort:                          200,
		AutoChargeEnabled:             false,
		PublicKeyName:                 "GatewayKey",
		PrivateSecretName:             "GatewaySecret",
		Host:                          "https://api.sber.ru",
		IsStaging:                     true,
	}
}

func (s SberPay) GatewayTest(ctx context.Context, req *_interface.GatewayTestReq) (icon string, gatewayType int64, err error) {
	utility.Assert(len(req.Key) > 0, "SberPay project name is required")
	utility.Assert(len(req.Secret) > 0, "SberPay API secret is required")

	// Test API connection
	_, err = s.makeAPICall(ctx, req.Key, req.Secret, "GET", "/payments/test", nil)
	if err != nil {
		return "", 0, gerror.Newf("SberPay API test failed: %v", err)
	}

	return "https://api.unibee.top/oss/file/sberpay-icon.png", consts.GatewayTypeCard, nil
}

func (s SberPay) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	// SberPay doesn't need to create users, return user ID as gateway user ID
	return &gateway_bean.GatewayUserCreateResp{
		GatewayUserId: fmt.Sprintf("user_%d", user.Id),
	}, nil
}

func (s SberPay) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, gatewayUserId string) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	return &gateway_bean.GatewayUserDetailQueryResp{
		GatewayUserId: gatewayUserId,
	}, nil
}

func (s SberPay) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	// SberPay doesn't support balance queries, return empty balance
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

func (s SberPay) GatewayNewPayment(ctx context.Context, gateway *entity.MerchantGateway, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
	// Build SberPay payment request
	paymentData := map[string]interface{}{
		"amount":       createPayContext.Pay.PaymentAmount,
		"currency":     createPayContext.Pay.Currency,
		"order_id":     createPayContext.Pay.PaymentId,
		"description":  "Payment for order " + createPayContext.Pay.PaymentId,
		"return_url":   createPayContext.Pay.ReturnUrl,
		"fail_url":     createPayContext.Pay.ReturnUrl,
		"payment_type": "sberpay", // Always use SberPay for this gateway
	}

	// Call SberPay API to create payment
	response, err := s.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", "/payments/create", paymentData)
	if err != nil {
		return nil, gerror.Newf("Failed to create SberPay payment: %v", err)
	}

	// Parse response
	var sberPayResp SberPayPaymentResponse
	err = json.Unmarshal([]byte(response), &sberPayResp)
	if err != nil {
		return nil, gerror.Newf("Failed to parse SberPay response: %v", err)
	}

	if !sberPayResp.Success {
		return nil, gerror.Newf("SberPay payment creation failed: %s", sberPayResp.Error)
	}

	// Save API log
	log.SaveChannelHttpLog("GatewayNewPayment", paymentData, response, nil, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	return &gateway_bean.GatewayNewPaymentResp{
		Payment:          createPayContext.Pay,
		Status:           s.mapSberPayStatusToSystemStatus(sberPayResp.Data.Status),
		GatewayPaymentId: sberPayResp.Data.PaymentId,
		Link:             sberPayResp.Data.Link,
	}, nil
}

func (s SberPay) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
	// Call SberPay API to query payment details
	response, err := s.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "GET", fmt.Sprintf("/payments/%s", gatewayPaymentId), nil)
	if err != nil {
		return nil, gerror.Newf("Failed to get SberPay payment detail: %v", err)
	}

	// Parse response
	var sberPayResp SberPayPaymentDetailResponse
	err = json.Unmarshal([]byte(response), &sberPayResp)
	if err != nil {
		return nil, gerror.Newf("Failed to parse SberPay response: %v", err)
	}

	if !sberPayResp.Success {
		return nil, gerror.Newf("SberPay payment detail failed: %s", sberPayResp.Error)
	}

	// Save API log
	log.SaveChannelHttpLog("GatewayPaymentDetail", map[string]interface{}{"payment_id": gatewayPaymentId}, response, nil, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	return &gateway_bean.GatewayPaymentRo{
		GatewayPaymentId: gatewayPaymentId,
		Status:           int(s.mapSberPayStatusToSystemStatus(sberPayResp.Data.Status)),
		PaymentAmount:    sberPayResp.Data.Amount,
		Currency:         sberPayResp.Data.Currency,
		CreateTime:       s.parseSberPayTime(sberPayResp.Data.CreatedAt),
	}, nil
}

func (s SberPay) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	// SberPay doesn't support payment list queries, return empty list
	return []*gateway_bean.GatewayPaymentRo{}, nil
}

func (s SberPay) GatewayCapture(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	// SberPay payments are instant, no capture needed
	return &gateway_bean.GatewayPaymentCaptureResp{}, nil
}

func (s SberPay) GatewayCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	// Call SberPay API to cancel payment
	response, err := s.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", fmt.Sprintf("/payments/%s/cancel", payment.GatewayPaymentId), nil)
	if err != nil {
		return nil, gerror.Newf("Failed to cancel SberPay payment: %v", err)
	}

	// Save API log
	log.SaveChannelHttpLog("GatewayCancel", map[string]interface{}{"payment_id": payment.GatewayPaymentId}, response, nil, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	return &gateway_bean.GatewayPaymentCancelResp{
		Status: consts.PaymentCancelled,
	}, nil
}

func (s SberPay) GatewayRefund(ctx context.Context, gateway *entity.MerchantGateway, createPaymentRefundContext *gateway_bean.GatewayNewPaymentRefundReq) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	// Build refund request
	refundData := map[string]interface{}{
		"payment_id": createPaymentRefundContext.Payment.GatewayPaymentId,
		"amount":     createPaymentRefundContext.Refund.RefundAmount,
		"currency":   createPaymentRefundContext.Refund.Currency,
		"reason":     createPaymentRefundContext.Refund.RefundComment,
	}

	// Call SberPay API to create refund
	response, err := s.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", "/refunds/create", refundData)
	if err != nil {
		return nil, gerror.Newf("Failed to create SberPay refund: %v", err)
	}

	// Parse response
	var sberPayResp SberPayRefundResponse
	err = json.Unmarshal([]byte(response), &sberPayResp)
	if err != nil {
		return nil, gerror.Newf("Failed to parse SberPay response: %v", err)
	}

	if !sberPayResp.Success {
		return &gateway_bean.GatewayPaymentRefundResp{
			GatewayRefundId: createPaymentRefundContext.Payment.GatewayPaymentId,
			Status:          consts.RefundFailed,
			Type:            consts.RefundTypeGateway,
			Reason:          fmt.Sprintf("SberPay refund failed: %s", sberPayResp.Error),
		}, nil
	}

	// Save API log
	log.SaveChannelHttpLog("GatewayRefund", refundData, response, nil, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: sberPayResp.Data.RefundId,
		Status:          s.mapSberPayRefundStatusToSystemStatus(sberPayResp.Data.Status),
		RefundAmount:    sberPayResp.Data.Amount,
		Currency:        sberPayResp.Data.Currency,
		Type:            consts.RefundTypeGateway,
	}, nil
}

func (s SberPay) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	// Call SberPay API to query refund details
	response, err := s.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "GET", fmt.Sprintf("/refunds/%s", gatewayRefundId), nil)
	if err != nil {
		return nil, gerror.Newf("Failed to get SberPay refund detail: %v", err)
	}

	// Parse response
	var sberPayResp SberPayRefundResponse
	err = json.Unmarshal([]byte(response), &sberPayResp)
	if err != nil {
		return nil, gerror.Newf("Failed to parse SberPay response: %v", err)
	}

	if !sberPayResp.Success {
		return nil, gerror.Newf("SberPay refund detail failed: %s", sberPayResp.Error)
	}

	// Save API log
	log.SaveChannelHttpLog("GatewayRefundDetail", map[string]interface{}{"refund_id": gatewayRefundId}, response, nil, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: gatewayRefundId,
		Status:          s.mapSberPayRefundStatusToSystemStatus(sberPayResp.Data.Status),
		RefundAmount:    sberPayResp.Data.Amount,
		Currency:        sberPayResp.Data.Currency,
		Type:            consts.RefundTypeGateway,
	}, nil
}

func (s SberPay) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	// SberPay does not support refund list query, return empty list
	return []*gateway_bean.GatewayPaymentRefundResp{}, nil
}

func (s SberPay) GatewayRefundCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	// Call SberPay API to cancel refund
	response, err := s.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", fmt.Sprintf("/refunds/%s/cancel", refund.GatewayRefundId), nil)
	if err != nil {
		return nil, gerror.Newf("Failed to cancel SberPay refund: %v", err)
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

// The following methods are not supported by SberPay, return empty implementation
func (s SberPay) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	return nil, fmt.Errorf("SberPay does not support crypto transactions")
}

func (s SberPay) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	return nil, fmt.Errorf("SberPay does not support payment method management")
}

func (s SberPay) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	return nil, fmt.Errorf("SberPay does not support payment method management")
}

func (s SberPay) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	return nil, fmt.Errorf("SberPay does not support payment method management")
}

func (s SberPay) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, currency string, metadata map[string]interface{}) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	return nil, fmt.Errorf("SberPay does not support payment method management")
}

// Utility methods
func (s SberPay) makeAPICall(ctx context.Context, apiKey, apiSecret, method, path string, data map[string]interface{}) (string, error) {
	// Determine base URL based on environment
	baseURL := "https://api.sber.ru"
	if !config.GetConfigInstance().IsProd() {
		baseURL = "https://sandbox.api.sber.ru"
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
		return "", gerror.Newf("SberPay API error: %d - %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}

func (s SberPay) mapSberPayStatusToSystemStatus(sberPayStatus string) consts.PaymentStatusEnum {
	switch strings.ToLower(sberPayStatus) {
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

func (s SberPay) mapSberPayRefundStatusToSystemStatus(sberPayStatus string) consts.RefundStatusEnum {
	switch strings.ToLower(sberPayStatus) {
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

func (s SberPay) parseSberPayTime(timeStr string) *gtime.Time {
	if timeStr == "" {
		return gtime.Now()
	}

	// Try to parse SberPay time format
	t, err := time.Parse("2006-01-02T15:04:05Z", timeStr)
	if err != nil {
		// If parsing fails, return current time
		return gtime.Now()
	}

	return gtime.NewFromTime(t)
}
