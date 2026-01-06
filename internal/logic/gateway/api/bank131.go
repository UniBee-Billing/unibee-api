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

type Bank131 struct{}

// Bank131 API response structure
type Bank131PaymentResponse struct {
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

type Bank131RefundResponse struct {
	Success bool `json:"success"`
	Data    struct {
		RefundId string `json:"refund_id"`
		Status   string `json:"status"`
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
	} `json:"data"`
	Error string `json:"error,omitempty"`
}

type Bank131PaymentDetailResponse struct {
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

func (b Bank131) GatewayInfo(ctx context.Context) *_interface.GatewayInfo {
	return &_interface.GatewayInfo{
		Name:                          "Bank131",
		Description:                   "Bank131 payment gateway for Russian market",
		DisplayName:                   "Bank131",
		GatewayWebsiteLink:            "https://developer.131.ru",
		GatewayWebhookIntegrationLink: "https://developer.131.ru",
		GatewayLogo:                   "https://api.unibee.dev/oss/file/dbyr430sbbm3hmbqtf.png",
		GatewayIcons:                  []string{"https://api.unibee.dev/oss/file/dbyr430sbbm3hmbqtf.png"},
		GatewayType:                   consts.GatewayTypeCard,
		Sort:                          200,
		AutoChargeEnabled:             false,
		GatewayPaymentTypes: []*_interface.GatewayPaymentType{
			{
				Name:        "SberPay",
				PaymentType: "sberpay",
				CountryName: "Russia",
				AutoCharge:  false,
				Category:    "mobile_payment",
			},
			{
				Name:        "Yandex.Money",
				PaymentType: "yamoney",
				CountryName: "Russia",
				AutoCharge:  false,
				Category:    "digital_wallet",
			},
			{
				Name:        "Bank Card (Yandex.Money)",
				PaymentType: "yamoneyac",
				CountryName: "Russia",
				AutoCharge:  false,
				Category:    "card_payment",
			},
			{
				Name:        "Cash (Yandex.Money)",
				PaymentType: "yamoneygp",
				CountryName: "Russia",
				AutoCharge:  false,
				Category:    "cash_payment",
			},
			{
				Name:        "Moneta",
				PaymentType: "moneta_ru",
				CountryName: "Russia",
				AutoCharge:  false,
				Category:    "digital_wallet",
			},
			{
				Name:        "Alfa-Click",
				PaymentType: "alfaclick_ru",
				CountryName: "Russia",
				AutoCharge:  false,
				Category:    "online_banking",
			},
			{
				Name:        "Promsvyazbank",
				PaymentType: "promsvyazbank_ru",
				CountryName: "Russia",
				AutoCharge:  false,
				Category:    "online_banking",
			},
			{
				Name:        "Faktura",
				PaymentType: "faktura_ru",
				CountryName: "Russia",
				AutoCharge:  false,
				Category:    "invoice_payment",
			},
			{
				Name:        "Russia Bank Transfer",
				PaymentType: "banktransfer_ru",
				CountryName: "Russia",
				AutoCharge:  false,
				Category:    "bank_transfer",
			},
			{
				Name:        "QIWI Wallet",
				PaymentType: "qiwi",
				CountryName: "Russia",
				AutoCharge:  false,
				Category:    "digital_wallet",
			},
			{
				Name:        "WebMoney",
				PaymentType: "webmoney",
				CountryName: "Russia",
				AutoCharge:  false,
				Category:    "digital_wallet",
			},
		},
		PublicKeyName:     "GatewayKey",
		PrivateSecretName: "GatewaySecret",
		Host:              "https://proxy.bank131.ru/api/v1",
		IsStaging:         true,
	}
}

func (b Bank131) GatewayTest(ctx context.Context, req *_interface.GatewayTestReq) (icon string, gatewayType int64, err error) {
	utility.Assert(len(req.Key) > 0, "Bank131 project name is required")
	utility.Assert(len(req.Secret) > 0, "Bank131 API secret is required")

	// Test API connection
	_, err = b.makeAPICall(ctx, req.Key, req.Secret, "GET", "/payments/test", nil)
	if err != nil {
		return "", 0, gerror.Newf("Bank131 API test failed: %v", err)
	}

	return "https://api.unibee.top/oss/file/bank131-icon.png", consts.GatewayTypeCard, nil
}

func (b Bank131) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	// Bank131 doesn't need to create users, return user ID as gateway user ID
	return &gateway_bean.GatewayUserCreateResp{
		GatewayUserId: fmt.Sprintf("user_%d", user.Id),
	}, nil
}

func (b Bank131) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, gatewayUserId string) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	return &gateway_bean.GatewayUserDetailQueryResp{
		GatewayUserId: gatewayUserId,
	}, nil
}

func (b Bank131) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	// Bank131 doesn't support balance queries, return empty balance
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

func (b Bank131) GatewayNewPayment(ctx context.Context, gateway *entity.MerchantGateway, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
	// Build Bank131 payment request
	paymentData := map[string]interface{}{
		"amount":       createPayContext.Pay.PaymentAmount,
		"currency":     createPayContext.Pay.Currency,
		"order_id":     createPayContext.Pay.PaymentId,
		"description":  "Payment for order " + createPayContext.Pay.PaymentId,
		"return_url":   createPayContext.Pay.ReturnUrl,
		"fail_url":     createPayContext.Pay.ReturnUrl,
		"payment_type": createPayContext.GatewayPaymentType,
	}

	// Call Bank131 API to create payment
	response, err := b.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", "/payments/create", paymentData)
	if err != nil {
		return nil, gerror.Newf("Failed to create Bank131 payment: %v", err)
	}

	// Parse response
	var bank131Resp Bank131PaymentResponse
	err = json.Unmarshal([]byte(response), &bank131Resp)
	if err != nil {
		return nil, gerror.Newf("Failed to parse Bank131 response: %v", err)
	}

	if !bank131Resp.Success {
		return nil, gerror.Newf("Bank131 payment creation failed: %s", bank131Resp.Error)
	}

	// Save API log
	log.SaveChannelHttpLog("GatewayNewPayment", paymentData, response, nil, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	return &gateway_bean.GatewayNewPaymentResp{
		Payment:          createPayContext.Pay,
		Status:           b.mapBank131StatusToSystemStatus(bank131Resp.Data.Status),
		GatewayPaymentId: bank131Resp.Data.PaymentId,
		Link:             bank131Resp.Data.Link,
	}, nil
}

func (b Bank131) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
	// Call Bank131 API to query payment details
	response, err := b.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "GET", fmt.Sprintf("/payments/%s", gatewayPaymentId), nil)
	if err != nil {
		return nil, gerror.Newf("Failed to get Bank131 payment detail: %v", err)
	}

	// Parse response
	var bank131Resp Bank131PaymentDetailResponse
	err = json.Unmarshal([]byte(response), &bank131Resp)
	if err != nil {
		return nil, gerror.Newf("Failed to parse Bank131 response: %v", err)
	}

	if !bank131Resp.Success {
		return nil, gerror.Newf("Bank131 payment detail failed: %s", bank131Resp.Error)
	}

	// Save API log
	log.SaveChannelHttpLog("GatewayPaymentDetail", map[string]interface{}{"payment_id": gatewayPaymentId}, response, nil, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	return &gateway_bean.GatewayPaymentRo{
		GatewayPaymentId: gatewayPaymentId,
		Status:           int(b.mapBank131StatusToSystemStatus(bank131Resp.Data.Status)),
		PaymentAmount:    bank131Resp.Data.Amount,
		Currency:         bank131Resp.Data.Currency,
		CreateTime:       b.parseBank131Time(bank131Resp.Data.CreatedAt),
	}, nil
}

func (b Bank131) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	// Bank131 doesn't support payment list queries, return empty list
	return []*gateway_bean.GatewayPaymentRo{}, nil
}

func (b Bank131) GatewayCapture(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	// Bank131 payments are instant, no capture needed
	return &gateway_bean.GatewayPaymentCaptureResp{}, nil
}

func (b Bank131) GatewayCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	// Call Bank131 API to cancel payment
	response, err := b.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", fmt.Sprintf("/payments/%s/cancel", payment.GatewayPaymentId), nil)
	if err != nil {
		return nil, gerror.Newf("Failed to cancel Bank131 payment: %v", err)
	}

	// Save API log
	log.SaveChannelHttpLog("GatewayCancel", map[string]interface{}{"payment_id": payment.GatewayPaymentId}, response, nil, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	return &gateway_bean.GatewayPaymentCancelResp{
		Status: consts.PaymentCancelled,
	}, nil
}

func (b Bank131) GatewayRefund(ctx context.Context, gateway *entity.MerchantGateway, createPaymentRefundContext *gateway_bean.GatewayNewPaymentRefundReq) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	// Build refund request
	refundData := map[string]interface{}{
		"payment_id": createPaymentRefundContext.Payment.GatewayPaymentId,
		"amount":     createPaymentRefundContext.Refund.RefundAmount,
		"currency":   createPaymentRefundContext.Refund.Currency,
		"reason":     createPaymentRefundContext.Refund.RefundComment,
	}

	// Call Bank131 API to create refund
	response, err := b.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", "/refunds/create", refundData)
	if err != nil {
		return nil, gerror.Newf("Failed to create Bank131 refund: %v", err)
	}

	// Parse response
	var bank131Resp Bank131RefundResponse
	err = json.Unmarshal([]byte(response), &bank131Resp)
	if err != nil {
		return nil, gerror.Newf("Failed to parse Bank131 response: %v", err)
	}

	if !bank131Resp.Success {
		return &gateway_bean.GatewayPaymentRefundResp{
			GatewayRefundId: createPaymentRefundContext.Payment.GatewayPaymentId,
			Status:          consts.RefundFailed,
			Type:            consts.RefundTypeGateway,
			Reason:          fmt.Sprintf("Bank131 refund failed: %s", bank131Resp.Error),
		}, nil
	}

	// Save API log
	log.SaveChannelHttpLog("GatewayRefund", refundData, response, nil, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: bank131Resp.Data.RefundId,
		Status:          b.mapBank131RefundStatusToSystemStatus(bank131Resp.Data.Status),
		RefundAmount:    bank131Resp.Data.Amount,
		Currency:        bank131Resp.Data.Currency,
		Type:            consts.RefundTypeGateway,
	}, nil
}

func (b Bank131) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	// Call Bank131 API to query refund details
	response, err := b.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "GET", fmt.Sprintf("/refunds/%s", gatewayRefundId), nil)
	if err != nil {
		return nil, gerror.Newf("Failed to get Bank131 refund detail: %v", err)
	}

	// Parse response
	var bank131Resp Bank131RefundResponse
	err = json.Unmarshal([]byte(response), &bank131Resp)
	if err != nil {
		return nil, gerror.Newf("Failed to parse Bank131 response: %v", err)
	}

	if !bank131Resp.Success {
		return nil, gerror.Newf("Bank131 refund detail failed: %s", bank131Resp.Error)
	}

	// Save API log
	log.SaveChannelHttpLog("GatewayRefundDetail", map[string]interface{}{"refund_id": gatewayRefundId}, response, nil, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)

	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: gatewayRefundId,
		Status:          b.mapBank131RefundStatusToSystemStatus(bank131Resp.Data.Status),
		RefundAmount:    bank131Resp.Data.Amount,
		Currency:        bank131Resp.Data.Currency,
		Type:            consts.RefundTypeGateway,
	}, nil
}

func (b Bank131) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	// Bank131 does not support refund list query, return empty list
	return []*gateway_bean.GatewayPaymentRefundResp{}, nil
}

func (b Bank131) GatewayRefundCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	// Call Bank131 API to cancel refund
	response, err := b.makeAPICall(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", fmt.Sprintf("/refunds/%s/cancel", refund.GatewayRefundId), nil)
	if err != nil {
		return nil, gerror.Newf("Failed to cancel Bank131 refund: %v", err)
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

// The following methods are not supported by Bank131, return empty implementation
func (b Bank131) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	return nil, fmt.Errorf("Bank131 does not support crypto transactions")
}

func (b Bank131) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	return nil, fmt.Errorf("Bank131 does not support payment method management")
}

func (b Bank131) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	return nil, fmt.Errorf("Bank131 does not support payment method management")
}

func (b Bank131) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	return nil, fmt.Errorf("Bank131 does not support payment method management")
}

func (b Bank131) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, currency string, metadata map[string]interface{}) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	return nil, fmt.Errorf("Bank131 does not support payment method management")
}

// Utility methods
func (b Bank131) makeAPICall(ctx context.Context, apiKey, apiSecret, method, path string, data map[string]interface{}) (string, error) {
	// Determine base URL based on environment
	baseURL := "https://proxy.bank131.ru/api/v1"
	if !config.GetConfigInstance().IsProd() {
		baseURL = "https://sandbox.proxy.bank131.ru/api/v1"
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
		return "", gerror.Newf("Bank131 API error: %d - %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}

func (b Bank131) mapBank131StatusToSystemStatus(bank131Status string) consts.PaymentStatusEnum {
	switch strings.ToLower(bank131Status) {
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

func (b Bank131) mapBank131RefundStatusToSystemStatus(bank131Status string) consts.RefundStatusEnum {
	switch strings.ToLower(bank131Status) {
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

func (b Bank131) parseBank131Time(timeStr string) *gtime.Time {
	if timeStr == "" {
		return gtime.Now()
	}

	// Try to parse Bank131 time format
	t, err := time.Parse("2006-01-02T15:04:05Z", timeStr)
	if err != nil {
		// If parsing fails, return current time
		return gtime.Now()
	}

	return gtime.NewFromTime(t)
}
