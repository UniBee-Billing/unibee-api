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

type MulenPay struct{}

// MulenPay API response structures
type MulenPayPaymentResponse struct {
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

type MulenPayRefundResponse struct {
	Success bool `json:"success"`
	Data    struct {
		RefundId string `json:"refund_id"`
		Status   string `json:"status"`
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
	} `json:"data"`
	Error string `json:"error,omitempty"`
}

type MulenPayPaymentDetailResponse struct {
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

// https://mulenpay.ru/docs/api
func (m MulenPay) GatewayInfo(ctx context.Context) *_interface.GatewayInfo {
	return &_interface.GatewayInfo{
		Name:                          "MulenPay",
		Description:                   "MulenPay payment gateway for Russian market",
		DisplayName:                   "MulenPay",
		GatewayWebsiteLink:            "https://mulenpay.ru/",
		GatewayWebhookIntegrationLink: "https://mulenpay.ru/",
		GatewayLogo:                   "https://api.unibee.dev/oss/file/dbzb7irqjw39mublde.svg",
		GatewayIcons:                  []string{"https://api.unibee.dev/oss/file/dbzb7irqjw39mublde.svg"},
		GatewayType:                   consts.GatewayTypeCard,
		Sort:                          220,
		AutoChargeEnabled:             false,
		PublicKeyName:                 "merchant_id",
		PrivateSecretName:             "secret_key",
		Host:                          "https://api.mulenpay.ru",
		IsStaging:                     true,
	}
}

func (m MulenPay) GatewayTest(ctx context.Context, req *_interface.GatewayTestReq) (icon string, gatewayType int64, err error) {
	utility.Assert(len(req.Key) > 0, "MulenPay merchant ID is required")
	utility.Assert(len(req.Secret) > 0, "MulenPay secret key is required")

	// Test API connection
	_, err = m.makeAPICall(ctx, req.Key, req.Secret, "GET", "/test", nil)
	if err != nil {
		return "", 0, gerror.Newf("MulenPay API test failed: %v", err)
	}

	return "https://api.unibee.top/oss/file/mulenpay-icon.png", consts.GatewayTypeCard, nil
}

func (m MulenPay) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	// MulenPay does not need to create users, return user ID as gateway user ID
	return &gateway_bean.GatewayUserCreateResp{
		GatewayUserId: fmt.Sprintf("user_%d", user.Id),
	}, nil
}

func (m MulenPay) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, gatewayUserId string) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	return &gateway_bean.GatewayUserDetailQueryResp{
		GatewayUserId: gatewayUserId,
	}, nil
}

func (m MulenPay) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	// MulenPay does not support balance queries, return empty balance
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

func (m MulenPay) GatewayNewPayment(ctx context.Context, gateway *entity.MerchantGateway, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
	_, description := createPayContext.GetInvoiceSingleProductNameAndDescription()
	if description == "" {
		description = "Payment for order " + createPayContext.Pay.PaymentId
	}
	paymentData := map[string]interface{}{
		"amount":       createPayContext.Pay.PaymentAmount, // 单位为分
		"currency":     createPayContext.Pay.Currency,
		"order_id":     createPayContext.Pay.PaymentId,
		"description":  description,
		"return_url":   createPayContext.Pay.ReturnUrl,
		"fail_url":     createPayContext.Pay.ReturnUrl,
		"payment_type": "card",
	}
	response, err := m.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", "/api/v2/payments", paymentData)
	if err != nil {
		return nil, gerror.Newf("Failed to create MulenPay payment: %v", err)
	}
	var mulenPayResp MulenPayPaymentResponse
	err = json.Unmarshal([]byte(response), &mulenPayResp)
	if err != nil {
		return nil, gerror.Newf("Failed to parse MulenPay response: %v", err)
	}
	if !mulenPayResp.Success {
		return nil, gerror.Newf("MulenPay payment creation failed: %s", mulenPayResp.Error)
	}
	log.SaveChannelHttpLog("GatewayNewPayment", paymentData, response, nil, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)
	return &gateway_bean.GatewayNewPaymentResp{
		Payment:          createPayContext.Pay,
		Status:           m.mapMulenPayStatusToSystemStatus(mulenPayResp.Data.Status),
		GatewayPaymentId: mulenPayResp.Data.PaymentId,
		Link:             mulenPayResp.Data.Link,
	}, nil
}

func (m MulenPay) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
	// Call MulenPay API to query payment details
	response, err := m.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "GET", fmt.Sprintf("/api/v2/payments/%s", gatewayPaymentId), nil)
	if err != nil {
		return nil, gerror.Newf("Failed to get MulenPay payment detail: %v", err)
	}

	// Parse response
	var mulenPayResp MulenPayPaymentDetailResponse
	err = json.Unmarshal([]byte(response), &mulenPayResp)
	if err != nil {
		return nil, gerror.Newf("Failed to parse MulenPay response: %v", err)
	}

	if !mulenPayResp.Success {
		return nil, gerror.Newf("MulenPay payment detail failed: %s", mulenPayResp.Error)
	}

	// Save API log
	log.SaveChannelHttpLog("GatewayPaymentDetail", map[string]interface{}{"payment_id": gatewayPaymentId}, response, nil, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	return &gateway_bean.GatewayPaymentRo{
		GatewayPaymentId: gatewayPaymentId,
		Status:           int(m.mapMulenPayStatusToSystemStatus(mulenPayResp.Data.Status)),
		PaymentAmount:    mulenPayResp.Data.Amount,
		Currency:         mulenPayResp.Data.Currency,
		CreateTime:       m.parseMulenPayTime(mulenPayResp.Data.CreatedAt),
	}, nil
}

func (m MulenPay) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	// MulenPay does not support payment list queries, return empty list
	return []*gateway_bean.GatewayPaymentRo{}, nil
}

func (m MulenPay) GatewayCapture(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	// MulenPay payments are instant, no capture needed
	return &gateway_bean.GatewayPaymentCaptureResp{}, nil
}

func (m MulenPay) GatewayCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	// Call MulenPay API to cancel payment
	response, err := m.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", fmt.Sprintf("/api/v2/payments/%s/cancel", payment.GatewayPaymentId), nil)
	if err != nil {
		return nil, gerror.Newf("Failed to cancel MulenPay payment: %v", err)
	}

	// Save API log
	log.SaveChannelHttpLog("GatewayCancel", map[string]interface{}{"payment_id": payment.GatewayPaymentId}, response, nil, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	return &gateway_bean.GatewayPaymentCancelResp{
		Status: consts.PaymentCancelled,
	}, nil
}

func (m MulenPay) GatewayRefund(ctx context.Context, gateway *entity.MerchantGateway, createPaymentRefundContext *gateway_bean.GatewayNewPaymentRefundReq) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	// Build refund request
	refundData := map[string]interface{}{
		"payment_id": createPaymentRefundContext.Payment.GatewayPaymentId,
		"amount":     createPaymentRefundContext.Refund.RefundAmount,
		"currency":   createPaymentRefundContext.Refund.Currency,
		"reason":     createPaymentRefundContext.Refund.RefundComment,
	}

	// Call MulenPay API to create refund
	response, err := m.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", "/api/v2/refunds", refundData)
	if err != nil {
		return nil, gerror.Newf("Failed to create MulenPay refund: %v", err)
	}

	// Parse response
	var mulenPayResp MulenPayRefundResponse
	err = json.Unmarshal([]byte(response), &mulenPayResp)
	if err != nil {
		return nil, gerror.Newf("Failed to parse MulenPay response: %v", err)
	}

	if !mulenPayResp.Success {
		return &gateway_bean.GatewayPaymentRefundResp{
			GatewayRefundId: createPaymentRefundContext.Payment.GatewayPaymentId,
			Status:          consts.RefundFailed,
			Type:            consts.RefundTypeGateway,
			Reason:          fmt.Sprintf("MulenPay refund failed: %s", mulenPayResp.Error),
		}, nil
	}

	// Save API log
	log.SaveChannelHttpLog("GatewayRefund", refundData, response, nil, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: mulenPayResp.Data.RefundId,
		Status:          m.mapMulenPayRefundStatusToSystemStatus(mulenPayResp.Data.Status),
		RefundAmount:    mulenPayResp.Data.Amount,
		Currency:        mulenPayResp.Data.Currency,
		Type:            consts.RefundTypeGateway,
	}, nil
}

func (m MulenPay) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	// Call MulenPay API to query refund details
	response, err := m.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "GET", fmt.Sprintf("/api/v2/refunds/%s", gatewayRefundId), nil)
	if err != nil {
		return nil, gerror.Newf("Failed to get MulenPay refund detail: %v", err)
	}

	// Parse response
	var mulenPayResp MulenPayRefundResponse
	err = json.Unmarshal([]byte(response), &mulenPayResp)
	if err != nil {
		return nil, gerror.Newf("Failed to parse MulenPay response: %v", err)
	}

	if !mulenPayResp.Success {
		return nil, gerror.Newf("MulenPay refund detail failed: %s", mulenPayResp.Error)
	}

	// Save API log
	log.SaveChannelHttpLog("GatewayRefundDetail", map[string]interface{}{"refund_id": gatewayRefundId}, response, nil, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: gatewayRefundId,
		Status:          m.mapMulenPayRefundStatusToSystemStatus(mulenPayResp.Data.Status),
		RefundAmount:    mulenPayResp.Data.Amount,
		Currency:        mulenPayResp.Data.Currency,
		Type:            consts.RefundTypeGateway,
	}, nil
}

func (m MulenPay) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	// MulenPay does not support refund list query, return empty list
	return []*gateway_bean.GatewayPaymentRefundResp{}, nil
}

func (m MulenPay) GatewayRefundCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	// Call MulenPay API to cancel refund
	response, err := m.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", fmt.Sprintf("/api/v2/refunds/%s/cancel", refund.GatewayRefundId), nil)
	if err != nil {
		return nil, gerror.Newf("Failed to cancel MulenPay refund: %v", err)
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

// The following methods are not supported by MulenPay, return empty implementation
func (m MulenPay) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	return nil, fmt.Errorf("MulenPay does not support crypto transactions")
}

func (m MulenPay) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	return nil, fmt.Errorf("MulenPay does not support payment method management")
}

func (m MulenPay) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	return nil, fmt.Errorf("MulenPay does not support payment method management")
}

func (m MulenPay) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	return nil, fmt.Errorf("MulenPay does not support payment method management")
}

func (m MulenPay) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, currency string, metadata map[string]interface{}) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	return nil, fmt.Errorf("MulenPay does not support payment method management")
}

// Utility methods
func (m MulenPay) makeAPICall(ctx context.Context, apiKey, apiSecret, method, path string, data map[string]interface{}) (string, error) {
	baseURL := "https://api.mulenpay.ru"
	if !config.GetConfigInstance().IsProd() {
		baseURL = "https://sandbox.api.mulenpay.ru"
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
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", gerror.Newf("Failed to read response body: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", gerror.Newf("MulenPay API error: %d - %s", resp.StatusCode, string(body))
	}
	return string(body), nil
}

func (m MulenPay) mapMulenPayStatusToSystemStatus(mulenPayStatus string) consts.PaymentStatusEnum {
	switch strings.ToLower(mulenPayStatus) {
	case "pending", "created", "waiting":
		return consts.PaymentCreated
	case "succeeded", "completed", "success":
		return consts.PaymentSuccess
	case "failed", "declined", "error":
		return consts.PaymentFailed
	case "cancelled", "canceled", "void":
		return consts.PaymentCancelled
	default:
		return consts.PaymentCreated
	}
}

func (m MulenPay) mapMulenPayRefundStatusToSystemStatus(mulenPayStatus string) consts.RefundStatusEnum {
	switch strings.ToLower(mulenPayStatus) {
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

func (m MulenPay) parseMulenPayTime(timeStr string) *gtime.Time {
	if timeStr == "" {
		return gtime.Now()
	}

	// Try to parse MulenPay time format
	t, err := time.Parse("2006-01-02T15:04:05Z", timeStr)
	if err != nil {
		// If parsing fails, return current time
		return gtime.Now()
	}

	return gtime.NewFromTime(t)
}
