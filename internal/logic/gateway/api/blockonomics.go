package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unibee/internal/cmd/config"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	gateway2 "unibee/internal/logic/gateway"
	"unibee/internal/logic/gateway/api/log"
	"unibee/internal/logic/gateway/gateway_bean"
	"unibee/internal/logic/gateway/util"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
	"unibee/utility/liberr"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Blockonomics Gateway Implementation
type Blockonomics struct {
}

// getUniqueAddress gets a unique address from Blockonomics API with Redis caching and balance verification
func getUniqueAddress(ctx context.Context, gateway *entity.MerchantGateway, crypto string) (string, int64, error) {
	const maxRetries = 5
	const cacheExpireDays = 3

	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Call Blockonomics API to get new address
		param := map[string]interface{}{
			"match_callback": gateway2.GetPaymentWebhookEntranceUrl(gateway.Id),
		}
		responseJson, err := SendBlockonomicsRequest(ctx, gateway.GatewaySecret, "POST", fmt.Sprintf("/api/new_address?crypto=%s", crypto), param)
		log.SaveChannelHttpLog("GetUniqueAddress", param, responseJson, err, "BlockonomicsGetAddress", nil, gateway)

		if err != nil {
			return "", 0, fmt.Errorf("failed to get address from Blockonomics API: %v", err)
		}

		address := responseJson.Get("address").String()
		if address == "" {
			return "", 0, gerror.New("invalid request, address not found")
		}

		// Check if address already exists in Redis (skip for USDT)
		if strings.ToUpper(crypto) != "USDT" {
			cacheKey := fmt.Sprintf("blockonomics_address:%s", address)
			exists, err := g.Redis().Exists(ctx, cacheKey)
			liberr.ErrIsNil(ctx, err, "Redis check address existence failure")

			if exists > 0 {
				g.Log().Warningf(ctx, "Address %s already exists in cache, retrying... (attempt %d/%d)", address, attempt, maxRetries)
				continue
			}
		} else {
			g.Log().Infof(ctx, "Skipping duplicate address check for USDT")
		}

		// Check address balance to ensure unconfirmed is 0 (only for BTC)
		var confirmedBalance int64
		if strings.ToUpper(crypto) == "BTC" {
			var err error
			confirmedBalance, err = checkAddressBalance(ctx, gateway, address)
			if err != nil {
				g.Log().Warningf(ctx, "Failed to check balance for address %s: %v, retrying... (attempt %d/%d)", address, err, attempt, maxRetries)
				continue
			}
		} else {
			// For non-BTC currencies (like USDT), balance check is not supported
			confirmedBalance = 0
			g.Log().Infof(ctx, "Skipping balance check for %s (not supported), using confirmedBalance = 0", crypto)
		}

		// Address is unique and has no unconfirmed balance, cache it and return (skip caching for USDT)
		if strings.ToUpper(crypto) != "USDT" {
			cacheKey := fmt.Sprintf("blockonomics_address:%s", address)
			_, err = g.Redis().Set(ctx, cacheKey, "used")
			liberr.ErrIsNil(ctx, err, "Redis cache address failure")

			// Set expiration to 3 days
			_, err = g.Redis().Expire(ctx, cacheKey, int64(cacheExpireDays*24*3600))
			liberr.ErrIsNil(ctx, err, "Redis set address expiration failure")
		} else {
			g.Log().Infof(ctx, "Skipping address caching for USDT")
		}

		g.Log().Infof(ctx, "Successfully got unique address: %s with confirmed balance: %d satoshis (attempt %d)", address, confirmedBalance, attempt)
		return address, confirmedBalance, nil
	}

	return "", 0, gerror.New("failed to get unique address")
}

// checkAddressBalance checks the balance of an address using Blockonomics balance API
func checkAddressBalance(ctx context.Context, gateway *entity.MerchantGateway, address string) (int64, error) {
	// Call Blockonomics balance API
	urlPath := fmt.Sprintf("/api/balance?addr=%s", address)
	responseJson, err := SendBlockonomicsRequest(ctx, gateway.GatewaySecret, "GET", urlPath, nil)
	log.SaveChannelHttpLog("CheckAddressBalance", nil, responseJson, err, "BlockonomicsCheckBalance", nil, gateway)

	if err != nil {
		return 0, fmt.Errorf("failed to check address balance: %v", err)
	}

	// Parse response from array format
	// Response format: {"response": [{"addr": "...", "confirmed": 0, "unconfirmed": 0}]}
	var confirmed int64 = 0
	var unconfirmed int64 = 0

	if responseJson.Contains("response") {
		responseArray := responseJson.Get("response").Array()
		if len(responseArray) > 0 {
			// Get the first address data (should match our address)
			// Convert map[string]interface{} to gjson.Json
			firstAddrMap := responseArray[0].(map[string]interface{})
			firstAddrJson, err := gjson.Encode(firstAddrMap)
			if err != nil {
				return 0, fmt.Errorf("failed to encode first address data: %v", err)
			}
			firstAddr := gjson.New(firstAddrJson)
			confirmed = firstAddr.Get("confirmed").Int64()
			unconfirmed = firstAddr.Get("unconfirmed").Int64()
		}
	}

	g.Log().Debugf(ctx, "Address %s balance - confirmed: %d, unconfirmed: %d", address, confirmed, unconfirmed)

	// Check if unconfirmed balance is 0
	if unconfirmed != 0 {
		return 0, fmt.Errorf("address has unconfirmed balance: %d satoshis", unconfirmed)
	}

	return confirmed, nil
}

// getCryptocurrencyPrecision returns the decimal precision for different cryptocurrencies
func getCryptocurrencyPrecision(crypto string) (int, int64) {
	cryptoPrecision := map[string]int{
		"BTC":  8,
		"USDT": 6,
		"ETH":  18,
		"LTC":  8,
		"BCH":  8,
		"XRP":  6,
		"ADA":  6,
		"DOT":  10,
		"LINK": 18,
		"UNI":  18,
	}

	precision, exists := cryptoPrecision[strings.ToUpper(crypto)]
	if !exists {
		precision = 8 // Default precision for unknown cryptocurrencies
	}
	precisionFactor := int64(1)
	for i := 0; i < precision; i++ {
		precisionFactor *= 10
	}
	return precision, precisionFactor
}

// getCryptocurrencyPrice fetches cryptocurrency price from Blockonomics API
func getCryptocurrencyPrice(ctx context.Context, crypto, currency string) (float64, error) {
	// Build API request URL
	apiURL := fmt.Sprintf("https://www.blockonomics.co/api/price?crypto=%s&currency=%s",
		strings.ToUpper(crypto), strings.ToUpper(currency))

	g.Log().Debugf(ctx, "Fetching cryptocurrency price from: %s", apiURL)

	// Send HTTP request
	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, fmt.Errorf("failed to call price API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("price API returned status: %d", resp.StatusCode)
	}

	// Parse response
	var priceResponse struct {
		Price float64 `json:"price"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&priceResponse); err != nil {
		return 0, fmt.Errorf("failed to decode price response: %v", err)
	}

	if priceResponse.Price <= 0 {
		return 0, fmt.Errorf("invalid price received: %f", priceResponse.Price)
	}

	g.Log().Debugf(ctx, "Cryptocurrency price: %s = %f %s", crypto, priceResponse.Price, currency)
	return priceResponse.Price, nil
}

func fetchBlockonomicsPaymentTypes(ctx context.Context) []*_interface.GatewayPaymentType {
	//filename := "alipay_plus_payment_types.json"
	//data, err := os.ReadFile(filename)
	//if err != nil {
	//	g.Log().Errorf(ctx, "Read payment type file: %s", err.Error())
	//}
	//
	//jsonString := string(data)
	jsonString := "[\n  {\"name\": \"Bitcoin\", \"paymentType\": \"BTC\", \"countryName\": \"Global\", \"autoCharge\": false, \"category\": \"Online Bitcoin\"},\n  {\"name\": \"USDT(ERC-20)\", \"paymentType\": \"USDT\", \"countryName\": \"ETH ERC-20\", \"autoCharge\": false, \"category\": \"Online ETH ERC-20\"}\n]"
	if !gjson.Valid(jsonString) {
		g.Log().Errorf(ctx, "Parse payment type file error, invalid json file")
	}

	var list = make([]*_interface.GatewayPaymentType, 0)
	err := utility.UnmarshalFromJsonString(jsonString, &list)
	if err != nil {
		g.Log().Errorf(ctx, "UnmarshalFromJsonString file error: %s", err.Error())
	}

	return list
}

func (b Blockonomics) GatewayInfo(ctx context.Context) *_interface.GatewayInfo {
	return &_interface.GatewayInfo{
		Name:                          "Blockonomics",
		Description:                   "Use API Key to securely process Bitcoin payments",
		DisplayName:                   "Blockonomics",
		GatewayWebsiteLink:            "https://www.blockonomics.co/",
		GatewayWebhookIntegrationLink: "https://www.blockonomics.co/dashboard#/store",
		GatewayLogo:                   "https://s3.amazonaws.com/cdn.freshdesk.com/data/helpdesk/attachments/production/33000207544/logo/zwUwDQK4ElmALBs7MZNUijfCmfpCy7bhdQ.png",
		GatewayIcons:                  []string{"https://api.unibee.top/oss/file/d6yhnz0wty7w6m7zhd.svg", "https://s3.amazonaws.com/cdn.freshdesk.com/data/helpdesk/attachments/production/33000207544/logo/zwUwDQK4ElmALBs7MZNUijfCmfpCy7bhdQ.png"},
		GatewayType:                   consts.GatewayTypeCrypto,
		GatewayPaymentTypes:           fetchBlockonomicsPaymentTypes(ctx),
		PublicKeyName:                 "API Key",
		PrivateSecretName:             "Secret Key(Same with API Key)",
		Sort:                          85,
		AutoChargeEnabled:             false,
		IsStaging:                     false,
	}
}

func (b Blockonomics) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	// 1. Get cryptocurrency price
	cryptoPrice, err := getCryptocurrencyPrice(ctx, from.CryptoCurrency, from.Currency)
	if err != nil {
		return nil, fmt.Errorf("failed to get cryptocurrency price: %v", err)
	}

	// 2. Calculate cryptocurrency amount
	// TotalAmount is fiat amount in cents, convert to dollars
	fiatAmount := utility.ConvertCentToDollarFloat(from.Amount, from.Currency)
	cryptoAmountFloat := fiatAmount / cryptoPrice

	// 3. Format cryptocurrency amount with dynamic precision
	_, precisionFactor := getCryptocurrencyPrecision(from.CryptoCurrency)
	cryptoAmount := int64(cryptoAmountFloat * float64(precisionFactor))

	return &gateway_bean.GatewayCryptoToCurrencyAmountDetailRes{
		Amount:         from.Amount,
		Currency:       from.Currency,
		CountryCode:    from.CountryCode,
		CryptoAmount:   cryptoAmount,
		CryptoCurrency: "USDT",
		Rate:           cryptoPrice * float64(precisionFactor),
	}, nil
}

func (b Blockonomics) GatewayTest(ctx context.Context, req *_interface.GatewayTestReq) (icon string, gatewayType int64, err error) {
	urlPath := "/api/currencies"
	param := map[string]interface{}{}
	responseJson, err := SendBlockonomicsRequest(ctx, req.Key, "GET", urlPath, param)
	utility.Assert(err == nil, fmt.Sprintf("invalid key, call blockonomics error %s", err))
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())

	// Check if response contains valid data
	utility.Assert(responseJson.Contains("currencies") || responseJson.Contains("BTC") || len(responseJson.Array()) > 0, "invalid response, no currency data found")

	return "https://api.unibee.top/oss/file/blockonomics-test.png", consts.GatewayTypeCrypto, nil
}

func (b Blockonomics) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	return &gateway_bean.GatewayUserCreateResp{
		GatewayUserId: strconv.FormatUint(user.Id, 10),
	}, nil
}

func (b Blockonomics) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, gatewayUserId string) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	return &gateway_bean.GatewayUserDetailQueryResp{
		GatewayUserId: gatewayUserId,
	}, nil
}

func (b Blockonomics) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	urlPath := "/api/balance"
	param := map[string]interface{}{}
	responseJson, err := SendBlockonomicsRequest(ctx, gateway.GatewaySecret, "GET", urlPath, param)
	if err != nil {
		return nil, err
	}

	var availableBalances []*gateway_bean.GatewayBalance
	if responseJson.Contains("BTC") {
		btcBalance := responseJson.Get("BTC").Float64()
		availableBalances = append(availableBalances, &gateway_bean.GatewayBalance{
			Amount:   int64(btcBalance * 100000000), // Convert BTC to satoshis
			Currency: "BTC",
		})
	}

	return &gateway_bean.GatewayMerchantBalanceQueryResp{
		AvailableBalance:       availableBalances,
		ConnectReservedBalance: []*gateway_bean.GatewayBalance{},
		PendingBalance:         []*gateway_bean.GatewayBalance{},
	}, nil
}

func (b Blockonomics) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	return nil, gerror.New("Not Support")
}

func (b Blockonomics) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	return nil, gerror.New("Not Support")
}

func (b Blockonomics) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	paymentMethods := []*gateway_bean.PaymentMethod{
		{
			Id: "BTC",
		},
	}
	return &gateway_bean.GatewayUserPaymentMethodListResp{
		PaymentMethods: paymentMethods,
	}, nil
}

func (b Blockonomics) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, currency string, metadata map[string]interface{}) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	return nil, gerror.New("Not Support")
}

func (b Blockonomics) GatewayNewPayment(ctx context.Context, gateway *entity.MerchantGateway, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
	if len(createPayContext.Gateway.BrandData) > 0 {
		gatewayPaymentTypes := utility.SplitToArray(createPayContext.Gateway.BrandData)
		if gatewayPaymentTypes != nil && len(gatewayPaymentTypes) == 1 {
			createPayContext.GatewayPaymentType = gatewayPaymentTypes[0]
		} else if len(createPayContext.GatewayPaymentType) == 0 && gatewayPaymentTypes != nil && len(gatewayPaymentTypes) > 0 {
			createPayContext.GatewayPaymentType = gatewayPaymentTypes[0]
		}
	}
	utility.Assert(len(createPayContext.GatewayPaymentType) > 0, "invalid Gateway PaymentType")
	// 1. Get cryptocurrency price
	cryptoPrice, err := getCryptocurrencyPrice(ctx, createPayContext.GatewayPaymentType, createPayContext.Pay.Currency)
	if err != nil {
		return nil, fmt.Errorf("failed to get cryptocurrency price: %v", err)
	}

	// 2. Calculate cryptocurrency amount
	// TotalAmount is fiat amount in cents, convert to dollars
	fiatAmount := utility.ConvertCentToDollarFloat(createPayContext.Pay.TotalAmount, createPayContext.Pay.Currency)
	cryptoAmountFloat := fiatAmount / cryptoPrice

	// 3. Convert to int64 for CryptoAmount field (multiply by precision factor to remove decimals)
	_, precisionFactor := getCryptocurrencyPrecision(createPayContext.GatewayPaymentType)
	cryptoAmountInt64 := int64(cryptoAmountFloat * float64(precisionFactor))

	// Get unique address with Redis caching and balance verification
	address, confirmedBalance, err := getUniqueAddress(ctx, gateway, createPayContext.GatewayPaymentType)
	if err != nil {
		return nil, err
	}
	paymentUrl := fmt.Sprintf("%s/embedded/blockonomics?paymentId=%s", config.GetConfigInstance().Server.GetServerPath(), createPayContext.Pay.PaymentId)
	name, description := createPayContext.GetInvoiceSingleProductNameAndDescription()
	// Fill action data for frontend integration
	action := gjson.New("")
	if !config.GetConfigInstance().IsProd() {
		_ = action.Set("testnet", 1)
	} else {
		_ = action.Set("testnet", 0)
	}
	_ = action.Set("blockonomicsReturnUrl", gateway2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, true))
	_ = action.Set("blockonomicsName", name)
	_ = action.Set("blockonomicsDescription", description)
	_ = action.Set("blockonomicsAddress", address)
	_ = action.Set("originalConfirmed", confirmedBalance) // Store original confirmed balance for payment completion check
	_ = action.Set("returnUrl", util.GetPaymentRedirectUrl(ctx, createPayContext.Pay, "true"))
	_ = action.Set("cancelUrl", util.GetPaymentRedirectUrl(ctx, createPayContext.Pay, "false"))

	return &gateway_bean.GatewayNewPaymentResp{
		Status:                 consts.PaymentCreated,
		GatewayPaymentId:       address,
		GatewayPaymentIntentId: address,
		Link:                   paymentUrl,
		Action:                 action,
		CryptoAmount:           cryptoAmountInt64,
		CryptoCurrency:         createPayContext.GatewayPaymentType,
	}, nil
}

func (b Blockonomics) GatewayCapture(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	return nil, gerror.New("Not Support")
}

func (b Blockonomics) GatewayCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	return &gateway_bean.GatewayPaymentCancelResp{
		Status: consts.PaymentCancelled,
	}, nil
}

func (b Blockonomics) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	return nil, gerror.New("Not Support")
}

func (b Blockonomics) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
	// Choose different processing methods based on cryptocurrency type
	cryptoCurrency := strings.ToUpper(payment.CryptoCurrency)

	switch cryptoCurrency {
	case "BTC":
		return b.handleBTCPaymentDetail(ctx, gateway, gatewayPaymentId, payment)
	case "USDT":
		return b.handleUSDTPaymentDetail(ctx, gateway, payment.Code, payment)
	default:
		// Other cryptocurrencies return unsuccessful status
		return b.handleOtherCryptoPaymentDetail(ctx, gateway, gatewayPaymentId, payment)
	}
}

// handleBTCPaymentDetail handles BTC payment detail query
func (b Blockonomics) handleBTCPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
	// Use /api/balance endpoint to query balance
	urlPath := "/api/balance"
	param := map[string]interface{}{
		"addr": gatewayPaymentId,
	}

	responseJson, err := SendBlockonomicsRequest(ctx, gateway.GatewaySecret, "GET", urlPath, param)
	log.SaveChannelHttpLog("GatewayPaymentDetail_BTC", param, responseJson, err, "BlockonomicsPaymentDetail_BTC", nil, gateway)
	if err != nil {
		return nil, err
	}

	g.Log().Debugf(ctx, "BTC Payment Detail responseJson: %s", responseJson.String())

	// Parse balance information from response array
	// Response format: {"response": [{"addr": "...", "confirmed": 0, "unconfirmed": 0}]}
	var confirmed int64 = 0
	var unconfirmed int64 = 0 // Not used for now, reserved for future extension

	if responseJson.Contains("response") {
		responseArray := responseJson.Get("response").Array()
		if len(responseArray) > 0 {
			// Get the first address data (should match our gatewayPaymentId)
			// Convert map[string]interface{} to gjson.Json
			firstAddrMap := responseArray[0].(map[string]interface{})
			firstAddrJson, err := gjson.Encode(firstAddrMap)
			if err != nil {
				return nil, fmt.Errorf("failed to encode first address data: %v", err)
			}
			firstAddr := gjson.New(firstAddrJson)
			confirmed = firstAddr.Get("confirmed").Int64()
			unconfirmed = firstAddr.Get("unconfirmed").Int64()
		}
	}

	// Get original confirmed balance from payment's Action
	originalConfirmed := int64(0)
	if payment.PaymentData != "" {
		paymentDataJson, err := gjson.LoadJson(payment.PaymentData)
		if err == nil && paymentDataJson.Contains("action") {
			actionValue := paymentDataJson.Get("action")
			if actionValue != nil {
				actionJson, err := gjson.LoadJson(actionValue.String())
				if err == nil && actionJson.Contains("originalConfirmed") {
					originalConfirmed = actionJson.Get("originalConfirmed").Int64()
				}
			}
		}
	}

	// Calculate new confirmed balance
	newConfirmed := confirmed - originalConfirmed

	// Determine payment success: new confirmed balance >= required crypto amount
	var status = consts.PaymentCreated
	var authorizeStatus = consts.WaitingAuthorized
	var paymentAmount int64 = 0
	var paidTime *gtime.Time
	var lastErr = "Pending"

	if newConfirmed >= payment.CryptoAmount {
		status = consts.PaymentSuccess
		authorizeStatus = consts.Authorized
		paymentAmount = payment.CryptoAmount
		paidTime = gtime.Now()
		lastErr = "Confirmed"
	} else if newConfirmed > 0 {
		lastErr = "Partially Confirmed"
	} else if unconfirmed > 0 {
		lastErr = "Unconfirmed"
	}

	return &gateway_bean.GatewayPaymentRo{
		GatewayPaymentId:     gatewayPaymentId,
		Status:               status,
		AuthorizeStatus:      authorizeStatus,
		AuthorizeReason:      "",
		CancelReason:         "",
		PaymentData:          responseJson.String(),
		TotalAmount:          payment.TotalAmount,
		PaymentAmount:        paymentAmount,
		GatewayPaymentMethod: "BTC",
		Currency:             payment.Currency,
		PaidTime:             paidTime,
		LastError:            lastErr,
	}, nil
}

// handleUSDTPaymentDetail handles USDT payment detail query
func (b Blockonomics) handleUSDTPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
	// Use /api/monitor_tx endpoint to monitor USDT transaction
	urlPath := "/api/monitor_tx"
	param := map[string]interface{}{
		"txhash":         gatewayPaymentId,
		"crypto":         "USDT",
		"match_callback": gateway2.GetPaymentWebhookEntranceUrl(gateway.Id),
	}
	if !config.GetConfigInstance().IsProd() {
		param["testnet"] = 1
	}

	responseJson, err := SendBlockonomicsRequest(ctx, gateway.GatewaySecret, "POST", urlPath, param)
	log.SaveChannelHttpLog("GatewayPaymentDetail_USDT", param, responseJson, err, "BlockonomicsPaymentDetail_USDT", nil, gateway)
	if err != nil {
		return nil, err
	}

	g.Log().Debugf(ctx, "USDT Payment Detail responseJson: %s", responseJson.String())

	//// According to documentation, USDT monitoring endpoint response may contain status information
	//// Since documentation is incomplete, we make reasonable guesses based on common response formats
	//var status = consts.PaymentCreated
	//var authorizeStatus = consts.WaitingAuthorized
	//var paymentAmount int64 = 0
	//var paidTime *gtime.Time
	//var lastErr = "Pending"
	//var txId = ""
	//
	//// Check if response contains success status
	//// According to Blockonomics documentation, possible fields include status, txid, addr, value, etc.
	//// The response might be in array format like other Blockonomics APIs
	//if responseJson.Contains("response") {
	//	// Handle array response format
	//	responseArray := responseJson.Get("response").Array()
	//	if len(responseArray) > 0 {
	//		// Convert map[string]interface{} to gjson.Json
	//		firstItemMap := responseArray[0].(map[string]interface{})
	//		firstItemJson, err := gjson.Encode(firstItemMap)
	//		if err != nil {
	//			return nil, fmt.Errorf("failed to encode first item data: %v", err)
	//		}
	//		firstItem := gjson.New(firstItemJson)
	//		txId = firstItem.Get("txid").String()
	//		if firstItem.Contains("status") {
	//			responseStatus := firstItem.Get("status").Int()
	//			switch responseStatus {
	//			case 2: // 2 confirmations - payment success
	//				status = consts.PaymentSuccess
	//				authorizeStatus = consts.Authorized
	//				paymentAmount = payment.CryptoAmount
	//				paidTime = gtime.Now()
	//				lastErr = "Confirmed"
	//			case 1: // 1 confirmation - partial confirmation
	//				status = consts.PaymentCreated
	//				authorizeStatus = consts.WaitingAuthorized
	//				lastErr = "Partially Confirmed"
	//			case 0: // 0 confirmations - unconfirmed
	//				status = consts.PaymentCreated
	//				authorizeStatus = consts.WaitingAuthorized
	//				lastErr = "Unconfirmed"
	//			case -1: // -1 - transaction failed or reverted
	//				status = consts.PaymentFailed
	//				authorizeStatus = consts.WaitingAuthorized
	//				lastErr = "Reverted"
	//			}
	//		} else if firstItem.Contains("txid") {
	//			// If transaction ID is returned, transaction is submitted but status is unknown
	//			status = consts.PaymentCreated
	//			authorizeStatus = consts.WaitingAuthorized
	//		}
	//	}
	//	//} else if responseJson.Contains("status") {
	//	//	// Handle direct response format (fallback)
	//	//	responseStatus := responseJson.Get("status").Int()
	//	//	switch responseStatus {
	//	//	case 2: // 2 confirmations - payment success
	//	//		status = consts.PaymentSuccess
	//	//		authorizeStatus = consts.Authorized
	//	//		paymentAmount = payment.CryptoAmount
	//	//		paidTime = gtime.Now()
	//	//	case 1: // 1 confirmation - partial confirmation
	//	//		status = consts.PaymentCreated
	//	//		authorizeStatus = consts.WaitingAuthorized
	//	//	case 0: // 0 confirmations - unconfirmed
	//	//		status = consts.PaymentCreated
	//	//		authorizeStatus = consts.WaitingAuthorized
	//	//	case -1: // -1 - transaction failed or reverted
	//	//		status = consts.PaymentFailed
	//	//		authorizeStatus = consts.WaitingAuthorized
	//	//	}
	//	//} else if responseJson.Contains("txid") {
	//	//	// If transaction ID is returned, transaction is submitted but status is unknown
	//	//	status = consts.PaymentCreated
	//	//	authorizeStatus = consts.WaitingAuthorized
	//}

	return &gateway_bean.GatewayPaymentRo{
		GatewayPaymentId: payment.GatewayPaymentId,
		Status:           payment.Status,
		AuthorizeStatus:  payment.AuthorizeStatus,
		AuthorizeReason:  "",
		CancelReason:     "",
		PaymentData:      responseJson.String(),
		TotalAmount:      payment.TotalAmount,
		PaymentAmount:    payment.PaymentAmount,
		Currency:         payment.Currency,
		PaidTime:         gtime.NewFromTimeStamp(payment.PaidTime),
		LastError:        payment.LastError,
	}, nil
}

// handleOtherCryptoPaymentDetail handles other cryptocurrency payment details
func (b Blockonomics) handleOtherCryptoPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
	// Other cryptocurrencies are not supported yet, return unsuccessful status
	return &gateway_bean.GatewayPaymentRo{
		GatewayPaymentId:     gatewayPaymentId,
		Status:               consts.PaymentCreated,
		AuthorizeStatus:      consts.WaitingAuthorized,
		AuthorizeReason:      "",
		CancelReason:         "",
		PaymentData:          "",
		TotalAmount:          payment.TotalAmount,
		PaymentAmount:        0,
		GatewayPaymentMethod: payment.CryptoCurrency,
		Currency:             payment.Currency,
		PaidTime:             nil,
	}, nil
}

func (b Blockonomics) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("Not Support")
}

func (b Blockonomics) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("Not Support")
}

func (b Blockonomics) GatewayRefund(ctx context.Context, gateway *entity.MerchantGateway, createPaymentRefundContext *gateway_bean.GatewayNewPaymentRefundReq) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: createPaymentRefundContext.Refund.RefundId,
		Status:          consts.RefundCreated,
		Type:            consts.RefundTypeMarked,
	}, nil
}

func (b Blockonomics) GatewayRefundCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return &gateway_bean.GatewayPaymentRefundResp{
		MerchantId:       strconv.FormatUint(payment.MerchantId, 10),
		GatewayRefundId:  refund.GatewayRefundId,
		GatewayPaymentId: payment.GatewayPaymentId,
		Status:           consts.RefundCancelled,
		Reason:           refund.RefundComment,
		RefundAmount:     refund.RefundAmount,
		Currency:         refund.Currency,
		RefundTime:       gtime.Now(),
	}, nil
}

func parseBlockonomicsPayment(item *gjson.Json) *gateway_bean.GatewayPaymentRo {
	var status = consts.PaymentCreated
	var authorizeStatus = consts.WaitingAuthorized

	balance := item.Get("balance").Float64()
	unconfirmedBalance := item.Get("unconfirmed_balance").Float64()

	if balance > 0 || unconfirmedBalance > 0 {
		status = consts.PaymentSuccess
		authorizeStatus = consts.Authorized
	}

	var paymentAmount int64 = 0
	if balance > 0 {
		paymentAmount = int64(balance * 100000000)
	}

	var paidTime *gtime.Time
	if balance > 0 {
		paidTime = gtime.Now()
	}

	return &gateway_bean.GatewayPaymentRo{
		GatewayPaymentId:     item.Get("addr").String(),
		Status:               status,
		AuthorizeStatus:      authorizeStatus,
		AuthorizeReason:      "",
		CancelReason:         "",
		PaymentData:          item.String(),
		TotalAmount:          paymentAmount,
		PaymentAmount:        paymentAmount,
		GatewayPaymentMethod: "BTC",
		PaidTime:             paidTime,
	}
}

func SendBlockonomicsRequest(ctx context.Context, apiKey string, method string, urlPath string, param map[string]interface{}) (res *gjson.Json, err error) {
	utility.Assert(len(apiKey) > 0, "apiKey is nil")

	baseUrl := "https://www.blockonomics.co"
	fullUrl := baseUrl + urlPath

	var body []byte
	if method == "POST" && param != nil {
		jsonData, err := gjson.Marshal(param)
		utility.Assert(err == nil, fmt.Sprintf("json format error %s param %s", err, param))
		body = []byte(jsonData)
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + apiKey,
	}

	g.Log().Debugf(ctx, "\nBlockonomics_Start %s %s %s\n", method, fullUrl, apiKey)

	response, err := utility.SendRequest(fullUrl, method, body, headers)
	g.Log().Debugf(ctx, "\nBlockonomics_End %s %s response: %s error %s\n", method, fullUrl, response, err)

	if err != nil {
		return nil, err
	}

	responseJson, err := gjson.LoadJson(string(response))
	if err != nil {
		return nil, err
	}

	return responseJson, nil
}
