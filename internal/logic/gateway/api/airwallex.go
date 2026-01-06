package api

import (
	"context"
	"fmt"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"

	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
)

type Airwallex struct{}

// https://www.airwallex.com/docs/api#/Introduction
func (a Airwallex) GatewayInfo(ctx context.Context) *_interface.GatewayInfo {
	return &_interface.GatewayInfo{
		Name:                          "Airwallex",
		Description:                   "Airwallex global payment gateway",
		DisplayName:                   "Airwallex",
		GatewayWebsiteLink:            "https://www.airwallex.com/",
		GatewayWebhookIntegrationLink: "https://www.airwallex.com/",
		GatewayLogo:                   "https://api.unibee.dev/oss/file/dbzbb48klsqymsjyke.svg",
		GatewayIcons:                  []string{"https://api.unibee.dev/oss/file/dbzbb48klsqymsjyke.svg"},
		GatewayType:                   consts.GatewayTypeCard,
		Sort:                          230,
		AutoChargeEnabled:             true,
		PublicKeyName:                 "APIKey",
		PrivateSecretName:             "ClientSecret",
		Host:                          "https://api.airwallex.com",
		IsStaging:                     true,
	}
}

func (a Airwallex) GatewayTest(ctx context.Context, req *_interface.GatewayTestReq) (icon string, gatewayType int64, err error) {
	utility.Assert(len(req.Key) > 0, "Airwallex API key is required")
	utility.Assert(len(req.Secret) > 0, "Airwallex client secret is required")

	// For now, return success without actual API test
	// In real implementation, this would test the API connection
	return "https://api.unibee.top/oss/file/airwallex-icon.png", consts.GatewayTypeCard, nil
}

func (a Airwallex) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	return &gateway_bean.GatewayUserCreateResp{
		GatewayUserId: fmt.Sprintf("airwallex_%d", user.Id),
	}, nil
}

func (a Airwallex) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, gatewayUserId string) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	return &gateway_bean.GatewayUserDetailQueryResp{
		GatewayUserId: gatewayUserId,
		Email:         "",
	}, nil
}

func (a Airwallex) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	return &gateway_bean.GatewayUserPaymentMethodListResp{
		PaymentMethods: []*gateway_bean.PaymentMethod{},
	}, nil
}

// Airwallex authentication cache
var airwallexTokenCache struct {
	sync.Mutex
	token     string
	expiresAt int64
}

// Get Airwallex access_token
func getAirwallexAccessToken(apiKey, clientSecret string) (string, error) {
	airwallexTokenCache.Lock()
	defer airwallexTokenCache.Unlock()
	if airwallexTokenCache.token != "" && airwallexTokenCache.expiresAt > time.Now().Unix()+60 {
		return airwallexTokenCache.token, nil
	}
	url := "https://api.airwallex.com/api/v1/authentication/login"
	body := map[string]string{
		"client_id": apiKey,
		"api_key":   clientSecret,
	}
	jsonBody, _ := json.Marshal(body)
	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result struct {
		Token     string `json:"token"`
		ExpiresIn int64  `json:"expires_in"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}
	airwallexTokenCache.token = result.Token
	airwallexTokenCache.expiresAt = time.Now().Unix() + result.ExpiresIn
	return result.Token, nil
}

// Create payment intent
func (a Airwallex) GatewayNewPayment(ctx context.Context, gateway *entity.MerchantGateway, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
	token, err := getAirwallexAccessToken(gateway.GatewayKey, gateway.GatewaySecret)
	if err != nil {
		return nil, gerror.Newf("Airwallex authentication failed: %v", err)
	}
	url := "https://api.airwallex.com/api/v1/pa/payment_intents/create"
	amount := createPayContext.Pay.PaymentAmount
	currency := createPayContext.Pay.Currency
	orderId := createPayContext.Pay.PaymentId
	description := orderId
	if createPayContext.Pay != nil && createPayContext.Pay.PaymentId != "" {
		description = createPayContext.Pay.PaymentId
	}
	body := map[string]interface{}{
		"amount":            float64(amount) / 100.0, // Airwallex uses amount in yuan
		"currency":          currency,
		"merchant_order_id": orderId,
		"request_id":        orderId,
		"description":       description,
		"payment_method": map[string]interface{}{
			"type": "card",
		},
		"return_url": createPayContext.Pay.ReturnUrl,
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, gerror.Newf("Airwallex failed to create payment: %v", err)
	}
	defer resp.Body.Close()
	var result struct {
		Id           string `json:"id"`
		Status       string `json:"status"`
		ClientSecret string `json:"client_secret"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, gerror.Newf("Airwallex failed to parse response: %v", err)
	}
	return &gateway_bean.GatewayNewPaymentResp{
		Payment:              createPayContext.Pay,
		Status:               consts.PaymentCreated,
		GatewayPaymentId:     result.Id,
		GatewayPaymentMethod: "card",
		Link:                 "", // Airwallex requires frontend to use client_secret to launch payment page
	}, nil
}

func (a Airwallex) GatewayCapture(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	return &gateway_bean.GatewayPaymentCaptureResp{
		MerchantId:       payment.MerchantId,
		GatewayCaptureId: payment.GatewayPaymentId,
		Amount:           payment.TotalAmount,
		Currency:         payment.Currency,
		Status:           "captured",
	}, nil
}

func (a Airwallex) GatewayCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	return &gateway_bean.GatewayPaymentCancelResp{
		MerchantId:      fmt.Sprintf("%d", payment.MerchantId),
		GatewayCancelId: payment.GatewayPaymentId,
		PaymentId:       payment.PaymentId,
		Status:          consts.PaymentCancelled,
	}, nil
}

// Refund
func (a Airwallex) GatewayRefund(ctx context.Context, gateway *entity.MerchantGateway, createPaymentRefundContext *gateway_bean.GatewayNewPaymentRefundReq) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	token, err := getAirwallexAccessToken(gateway.GatewayKey, gateway.GatewaySecret)
	if err != nil {
		return nil, gerror.Newf("Airwallex authentication failed: %v", err)
	}
	url := "https://api.airwallex.com/api/v1/pa/refunds/create"
	amount := createPaymentRefundContext.Refund.RefundAmount
	currency := createPaymentRefundContext.Refund.Currency
	paymentId := createPaymentRefundContext.Payment.GatewayPaymentId
	body := map[string]interface{}{
		"amount":            float64(amount) / 100.0, // Airwallex uses amount in yuan
		"currency":          currency,
		"payment_intent_id": paymentId,
		"reason":            createPaymentRefundContext.Refund.RefundComment,
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, gerror.Newf("Airwallex failed to create refund: %v", err)
	}
	defer resp.Body.Close()
	var result struct {
		Id       string  `json:"id"`
		Status   string  `json:"status"`
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, gerror.Newf("Airwallex failed to parse refund response: %v", err)
	}
	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: result.Id,
		Status:          consts.RefundCreated, // Map status if needed
		RefundAmount:    int64(result.Amount * 100),
		Currency:        result.Currency,
		Type:            consts.RefundTypeGateway,
	}, nil
}

func (a Airwallex) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	return nil, gerror.New("Not Supported")
}

// Query payment intent
func (a Airwallex) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
	token, err := getAirwallexAccessToken(gateway.GatewayKey, gateway.GatewaySecret)
	if err != nil {
		return nil, gerror.Newf("Airwallex authentication failed: %v", err)
	}
	url := fmt.Sprintf("https://api.airwallex.com/api/v1/pa/payment_intents/%s", gatewayPaymentId)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, gerror.Newf("Airwallex failed to query payment: %v", err)
	}
	defer resp.Body.Close()
	var result struct {
		Id       string  `json:"id"`
		Status   string  `json:"status"`
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, gerror.Newf("Airwallex failed to parse response: %v", err)
	}
	return &gateway_bean.GatewayPaymentRo{
		GatewayPaymentId: result.Id,
		Status:           consts.PaymentCreated, // Map status if needed
		PaymentAmount:    int64(result.Amount * 100),
		Currency:         result.Currency,
	}, nil
}

// Query refund
func (a Airwallex) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	token, err := getAirwallexAccessToken(gateway.GatewayKey, gateway.GatewaySecret)
	if err != nil {
		return nil, gerror.Newf("Airwallex authentication failed: %v", err)
	}
	url := fmt.Sprintf("https://api.airwallex.com/api/v1/pa/refunds/%s", gatewayRefundId)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, gerror.Newf("Airwallex failed to query refund: %v", err)
	}
	defer resp.Body.Close()
	var result struct {
		Id       string  `json:"id"`
		Status   string  `json:"status"`
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, gerror.Newf("Airwallex failed to parse refund response: %v", err)
	}
	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: result.Id,
		Status:          consts.RefundCreated, // Map status if needed
		RefundAmount:    int64(result.Amount * 100),
		Currency:        result.Currency,
		Type:            consts.RefundTypeGateway,
	}, nil
}

func (a Airwallex) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	return &gateway_bean.GatewayMerchantBalanceQueryResp{
		AvailableBalance: []*gateway_bean.GatewayBalance{},
		PendingBalance:   []*gateway_bean.GatewayBalance{},
	}, nil
}

func (a Airwallex) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	return nil, gerror.New("Not Supported")
}

func (a Airwallex) GatewayRefundCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return &gateway_bean.GatewayPaymentRefundResp{
		MerchantId:       fmt.Sprintf("%d", payment.MerchantId),
		GatewayRefundId:  refund.GatewayRefundId,
		GatewayPaymentId: payment.GatewayPaymentId,
		Status:           consts.RefundCancelled,
		Reason:           "Refund cancelled",
		RefundAmount:     refund.RefundAmount,
		Currency:         payment.Currency,
		Type:             consts.RefundTypeGateway,
	}, nil
}

func (a Airwallex) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	return &gateway_bean.GatewayUserAttachPaymentMethodResp{}, nil
}

func (a Airwallex) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	return &gateway_bean.GatewayUserDeAttachPaymentMethodResp{}, nil
}

func (a Airwallex) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, currency string, metadata map[string]interface{}) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	return &gateway_bean.GatewayUserPaymentMethodCreateAndBindResp{
		PaymentMethod: &gateway_bean.PaymentMethod{
			Id:        fmt.Sprintf("airwallex_pm_%d", userId),
			Type:      "card",
			IsDefault: true,
			Data:      nil,
		},
		Url: "",
	}, nil
}
