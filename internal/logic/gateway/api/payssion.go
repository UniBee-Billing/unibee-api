package api

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
	"strings"
	"unibee/internal/cmd/config"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	webhook2 "unibee/internal/logic/gateway"
	"unibee/internal/logic/gateway/api/log"
	"unibee/internal/logic/gateway/gateway_bean"
	"unibee/internal/logic/gateway/util"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
	"unibee/utility/unibee"
)

//https://payssion.com/cn/docs/#api-reference-payment-request
// todo mark auto-charge
// todo mark 3ds check

type Payssion struct {
}

func fetchPayssionPaymentTypes(ctx context.Context) []*_interface.GatewayPaymentType {
	//filename := "payssion_payment_types.json"
	//data, err := os.ReadFile(filename)
	//if err != nil {
	//	g.Log().Errorf(ctx, "Read payment type file: %s", err.Error())
	//}
	//
	//jsonString := string(data)
	jsonString := "[\n  { \"name\": \"FPX\", \"paymentType\": \"fpx_my\", \"countryName\": \"Malaysia\" },\n  { \"name\": \"eNets\", \"paymentType\": \"enets_sg\", \"countryName\": \"Singapore\" },\n  { \"name\": \"PayNow\", \"paymentType\": \"paynow_sg\", \"countryName\": \"Singapore\" },\n  { \"name\": \"E-banking\", \"paymentType\": \"ebanking_th\", \"countryName\": \"Thailand\" },\n  { \"name\": \"Doku\", \"paymentType\": \"doku_id\", \"countryName\": \"Indonesia\" },\n  { \"name\": \"ATM\", \"paymentType\": \"atm_id\", \"countryName\": \"Indonesia\" },\n  { \"name\": \"Alfamart\", \"paymentType\": \"alfamart_id\", \"countryName\": \"Indonesia\" },\n  { \"name\": \"Dragonpay\", \"paymentType\": \"dragonpay_ph\", \"countryName\": \"Philippines\" },\n  { \"name\": \"Globe Gcash\", \"paymentType\": \"gcash_ph\", \"countryName\": \"Philippines\" },\n  { \"name\": \"CherryCredits\", \"paymentType\": \"cherrycredits\", \"countryName\": \"Global including South East\" },\n  { \"name\": \"MOLPoints\", \"paymentType\": \"molpoints\", \"countryName\": \"Global including South East\" },\n  { \"name\": \"MOLPoints card\", \"paymentType\": \"molpointscard\", \"countryName\": \"Global including South East\" },\n  { \"name\": \"Alipay\", \"paymentType\": \"alipay_cn\", \"countryName\": \"China\" },\n  { \"name\": \"Tenpay\", \"paymentType\": \"tenpay_cn\", \"countryName\": \"China\" },\n  { \"name\": \"Unionpay\", \"paymentType\": \"unionpay_cn\", \"countryName\": \"China\" },\n  { \"name\": \"Gash\", \"paymentType\": \"gash_tw\", \"countryName\": \"Taiwan\" },\n  { \"name\": \"UPI\", \"paymentType\": \"upi_in\", \"countryName\": \"India\" },\n  { \"name\": \"Indian Wallets\", \"paymentType\": \"wallet_in\", \"countryName\": \"India\" },\n  { \"name\": \"India Netbanking\", \"paymentType\": \"ebanking_in\", \"countryName\": \"India\" },\n  { \"name\": \"India Credit/Debit Card\", \"paymentType\": \"bankcard_in\", \"countryName\": \"India\" },\n  { \"name\": \"South Korea Credit Card\", \"paymentType\": \"creditcard_kr\", \"countryName\": \"South Korea\" },\n  { \"name\": \"South Korea Internet Banking\", \"paymentType\": \"ebanking_kr\", \"countryName\": \"South Korea\" },\n  { \"name\": \"KakaoPay\", \"paymentType\": \"kakaopay_kr\", \"countryName\": \"South Korea\" },\n  { \"name\": \"PAYCO\", \"paymentType\": \"payco_kr\", \"countryName\": \"South Korea\" },\n  { \"name\": \"SSG Pay\", \"paymentType\": \"ssgpay_kr\", \"countryName\": \"South Korea\" },\n  { \"name\": \"Samsung Pay\", \"paymentType\": \"samsungpay_kr\", \"countryName\": \"South Korea\" },\n  { \"name\": \"onecard\", \"paymentType\": \"onecard\", \"countryName\": \"Middle East & North Africa\" },\n  { \"name\": \"Fawry\", \"paymentType\": \"fawry_eg\", \"countryName\": \"Egypt\" },\n  { \"name\": \"Santander Rio\", \"paymentType\": \"santander_ar\", \"countryName\": \"Argentina\" },\n  { \"name\": \"Pago Fácil\", \"paymentType\": \"pagofacil_ar\", \"countryName\": \"Argentina\" },\n  { \"name\": \"Rapi Pago\", \"paymentType\": \"rapipago_ar\", \"countryName\": \"Argentina\" },\n  { \"name\": \"bancodobrasil\", \"paymentType\": \"bancodobrasil_br\", \"countryName\": \"Brazil\" },\n  { \"name\": \"itau\", \"paymentType\": \"itau_br\", \"countryName\": \"Brazil\" },\n  { \"name\": \"Boleto\", \"paymentType\": \"boleto_br\", \"countryName\": \"Brazil\" },\n  { \"name\": \"bradesco\", \"paymentType\": \"bradesco_br\", \"countryName\": \"Brazil\" },\n  { \"name\": \"caixa\", \"paymentType\": \"caixa_br\", \"countryName\": \"Brazil\" },\n  { \"name\": \"Santander\", \"paymentType\": \"santander_br\", \"countryName\": \"Brazil\" },\n  { \"name\": \"BBVA Bancomer\", \"paymentType\": \"bancomer_mx\", \"countryName\": \"Mexico\" },\n  { \"name\": \"Santander\", \"paymentType\": \"santander_mx\", \"countryName\": \"Mexico\" },\n  { \"name\": \"oxxo\", \"paymentType\": \"oxxo_mx\", \"countryName\": \"Mexico\" },\n  { \"name\": \"SPEI\", \"paymentType\": \"spei_mx\", \"countryName\": \"Mexico\" },\n  { \"name\": \"redpagos\", \"paymentType\": \"redpagos_uy\", \"countryName\": \"Uruguay\" },\n  { \"name\": \"Abitab\", \"paymentType\": \"abitab_uy\", \"countryName\": \"Uruguay\" },\n  { \"name\": \"Banco de Chile\", \"paymentType\": \"bancochile_cl\", \"countryName\": \"Chile\" },\n  { \"name\": \"RedCompra\", \"paymentType\": \"redcompra_cl\", \"countryName\": \"Chile\" },\n  { \"name\": \"WebPay plus\", \"paymentType\": \"webpay_cl\", \"countryName\": \"Chile\" },\n  { \"name\": \"Servipag\", \"paymentType\": \"servipag_cl\", \"countryName\": \"Chile\" },\n  { \"name\": \"Santander\", \"paymentType\": \"santander_cl\", \"countryName\": \"Chile\" },\n  { \"name\": \"Efecty\", \"paymentType\": \"efecty_co\", \"countryName\": \"Colombia\" },\n  { \"name\": \"PSE\", \"paymentType\": \"pse_co\", \"countryName\": \"Colombia\" },\n  { \"name\": \"BCP\", \"paymentType\": \"bcp_pe\", \"countryName\": \"Peru\" },\n  { \"name\": \"Interbank\", \"paymentType\": \"interbank_pe\", \"countryName\": \"Peru\" },\n  { \"name\": \"BBVA\", \"paymentType\": \"bbva_pe\", \"countryName\": \"Peru\" },\n  { \"name\": \"Pago Efectivo\", \"paymentType\": \"pagoefectivo_pe\", \"countryName\": \"Peru\" },\n  { \"name\": \"BoaCompra\", \"paymentType\": \"boacompra\", \"countryName\": \"Latin America\" },\n  { \"name\": \"QIWI\", \"paymentType\": \"qiwi\", \"countryName\": \"CIS countries\" },\n  { \"name\": \"Yandex.Money\", \"paymentType\": \"yamoney\", \"countryName\": \"CIS countries\" },\n  { \"name\": \"Webmoney\", \"paymentType\": \"webmoney\", \"countryName\": \"CIS countries\" },\n  { \"name\": \"Bank Card (Yandex.Money)\", \"paymentType\": \"yamoneyac\", \"countryName\": \"CIS countries\" },\n  { \"name\": \"Cash (Yandex.Money)\", \"paymentType\": \"yamoneygp\", \"countryName\": \"Russia\" },\n  { \"name\": \"Moneta\", \"paymentType\": \"moneta_ru\", \"countryName\": \"Russia\" },\n  { \"name\": \"Alfa-Click\", \"paymentType\": \"alfaclick_ru\", \"countryName\": \"Russia\" },\n  { \"name\": \"Promsvyazbank\", \"paymentType\": \"promsvyazbank_ru\", \"countryName\": \"Russia\" },\n  { \"name\": \"Faktura\", \"paymentType\": \"faktura_ru\", \"countryName\": \"Russia\" },\n  { \"name\": \"Russia Bank transfer\", \"paymentType\": \"banktransfer_ru\", \"countryName\": \"Russia\" },\n  { \"name\": \"Turkish Credit/Bank Card\", \"paymentType\": \"bankcard_tr\", \"countryName\": \"Turkey\" },\n  { \"name\": \"ininal\", \"paymentType\": \"ininal_tr\", \"countryName\": \"Turkey\" },\n  { \"name\": \"bkmexpress\", \"paymentType\": \"bkmexpress_tr\", \"countryName\": \"Turkey\" },\n  { \"name\": \"Turkish Bank Transfer\", \"paymentType\": \"banktransfer_tr\", \"countryName\": \"Turkey\" },\n  { \"name\": \"Paysafecard\", \"paymentType\": \"paysafecard\", \"countryName\": \"Global\" },\n  { \"name\": \"Sofort\", \"paymentType\": \"sofort\", \"countryName\": \"Europe\" },\n  { \"name\": \"Giropay\", \"paymentType\": \"giropay_de\", \"countryName\": \"Germany\" },\n  { \"name\": \"EPS\", \"paymentType\": \"eps_at\", \"countryName\": \"Austria\" },\n  { \"name\": \"Bancontact/Mistercash\", \"paymentType\": \"bancontact_be\", \"countryName\": \"Belgium\" },\n  { \"name\": \"Dotpay\", \"paymentType\": \"dotpay_pl\", \"countryName\": \"Poland\" },\n  { \"name\": \"P24\", \"paymentType\": \"p24_pl\", \"countryName\": \"Poland\" },\n  { \"name\": \"PayU\", \"paymentType\": \"payu_pl\", \"countryName\": \"Poland\" },\n  { \"name\": \"PayU\", \"paymentType\": \"payu_cz\", \"countryName\": \"Czech Republic\" },\n  { \"name\": \"iDeal\", \"paymentType\": \"ideal_nl\", \"countryName\": \"Netherlands\" },\n  { \"name\": \"Multibanco\", \"paymentType\": \"multibanco_pt\", \"countryName\": \"Portugal\" },\n  { \"name\": \"Neosurf\", \"paymentType\": \"neosurf\", \"countryName\": \"France\" },\n  { \"name\": \"Polipayment\", \"paymentType\": \"polipayment\", \"countryName\": \"Australia & New Zealand\" }\n]\n"
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

func (c Payssion) GatewayInfo(ctx context.Context) *_interface.GatewayInfo {
	return &_interface.GatewayInfo{
		Name:                          "Payssion",
		Description:                   "Use App Key and Secret Key to secure the payment",
		DisplayName:                   "Payssion",
		GatewayWebsiteLink:            "https://payssion.com",
		GatewayWebhookIntegrationLink: "https://www.payssion.com/account/app",
		GatewayLogo:                   "https://api.unibee.top/oss/file/d76q4s98dnw7x1yzzg.png",
		GatewayIcons:                  []string{"https://api.unibee.top/oss/file/d76q4s98dnw7x1yzzg.png"},
		GatewayType:                   consts.GatewayTypeCard,
		QueueForRefund:                true,
		GatewayPaymentTypes:           fetchPayssionPaymentTypes(ctx),
		CurrencyExchangeEnabled:       true,
		Sort:                          60,
		SubGatewayName:                "PM ID",
	}
}

func (c Payssion) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	return nil, gerror.New("not support")
}

func (c Payssion) GatewayTest(ctx context.Context, req *_interface.GatewayTestReq) (icon string, gatewayType int64, err error) {
	urlPath := "/api/v1/payments"
	pmID := "payssion_test"
	if config.GetConfigInstance().IsProd() {
		pmID = req.SubGateway
	}
	param := map[string]interface{}{
		"currency":    "EUR",
		"pm_id":       pmID,
		"amount":      100,
		"description": "test payment description",
		"order_id":    uuid.New().String(),
		"payer_email": "jack.fu@wowow.io",
	}
	if !config.GetConfigInstance().IsProd() {
		param["pm_id"] = "payssion_test"
	}
	param["api_key"] = req.Key
	param["api_sig"] = utility.MD5(fmt.Sprintf("%v|%v|%v|%v|%v|%v", param["api_key"], param["pm_id"], param["amount"], param["currency"], param["order_id"], req.Secret))
	responseJson, err := SendPayssionPaymentRequest(ctx, req.Key, req.Secret, "POST", urlPath, param)
	utility.Assert(err == nil, fmt.Sprintf("invalid keys,  call error %s", err))
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	utility.Assert(responseJson.Contains("transaction.transaction_id"), "invalid keys, transaction_id is nil")
	return "http://unibee.top/files/invoice/changelly.png", consts.GatewayTypeCrypto, nil
}

func (c Payssion) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Payssion) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, gatewayUserId string) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Payssion) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Payssion) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Payssion) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Payssion) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Payssion) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, currency string, metadata map[string]interface{}) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Payssion) GatewayNewPayment(ctx context.Context, gateway *entity.MerchantGateway, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
	urlPath := "/api/v1/payments"
	//var name = ""
	var description = createPayContext.Invoice.ProductName
	if len(createPayContext.Invoice.Lines) > 0 {
		var line = createPayContext.Invoice.Lines[0]
		if len(line.Name) > 0 {
			description = line.Name
		} else if len(line.Description) > 0 {
			description = line.Description
		}
	}
	pmID := "alipay_cn"
	if len(gateway.SubGateway) > 0 {
		pmID = gateway.SubGateway
	}
	if len(createPayContext.GatewayPaymentType) > 0 {
		pmID = createPayContext.GatewayPaymentType
	}
	var currency = createPayContext.Pay.Currency
	var totalAmount = createPayContext.Pay.TotalAmount
	{
		// Currency Exchange
		if createPayContext.GatewayCurrencyExchange != nil && createPayContext.ExchangeAmount > 0 && len(createPayContext.ExchangeCurrency) > 0 {
			currency = createPayContext.ExchangeCurrency
			totalAmount = createPayContext.ExchangeAmount
		}
	}
	param := map[string]interface{}{
		"currency": currency,
		"amount":   utility.ConvertCentToDollarStr(totalAmount, currency),
		"pm_id":    pmID,
		//"title":               name,
		"description": description,
		"order_id":    createPayContext.Pay.PaymentId,
		//"customer_id": strconv.FormatUint(createPayContext.Pay.UserId, 10),
		"payer_email": createPayContext.Email,
		"return_url":  webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, true),
		"notify_url":  webhook2.GetPaymentWebhookEntranceUrl(gateway.Id),
		//"backUrl":     webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, false),
		//"payment_data": createPayContext.Metadata,
		//"pending_deadline_at": time.Unix(createPayContext.Pay.ExpireTime, 0).Format("2006-01-02T15:04:05.876Z"),
	}
	if !config.GetConfigInstance().IsProd() {
		param["pm_id"] = "payssion_test"
	}
	param["api_key"] = gateway.GatewayKey
	param["api_sig"] = utility.MD5(fmt.Sprintf("%v|%v|%v|%v|%v|%v", param["api_key"], param["pm_id"], param["amount"], param["currency"], param["order_id"], gateway.GatewaySecret))
	responseJson, err := SendPayssionPaymentRequest(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", urlPath, param)
	log.SaveChannelHttpLog("GatewayNewPayment", param, responseJson, err, "PayssionNewPayment", nil, gateway)
	if err != nil {
		return nil, err
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	if !responseJson.Contains("result_code") || responseJson.Get("result_code").Int() != 200 {
		return nil, gerror.New("invalid request, result_code is nil or not 200")
	}
	if !responseJson.Contains("transaction") {
		return nil, gerror.New("invalid request, transaction is nil")
	}
	var status consts.PaymentStatusEnum = consts.PaymentCreated
	//transaction := responseJson.Get("transaction")

	return &gateway_bean.GatewayNewPaymentResp{
		Status:                 status,
		GatewayPaymentId:       responseJson.Get("transaction.transaction_id").String(),
		GatewayPaymentIntentId: responseJson.Get("transaction.transaction_id").String(),
		Link:                   responseJson.Get("redirect_url").String(),
	}, nil
}

func (c Payssion) GatewayCapture(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Payssion) GatewayCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	return &gateway_bean.GatewayPaymentCancelResp{Status: consts.PaymentCancelled}, nil
}

func (c Payssion) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	return nil, gerror.New("Not Support")
}

func (c Payssion) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
	urlPath := "/api/v1/payment/getDetail"
	param := map[string]interface{}{}
	param["transaction_id"] = gatewayPaymentId
	param["order_id"] = payment.PaymentId
	param["api_key"] = gateway.GatewayKey
	param["api_sig"] = utility.MD5(fmt.Sprintf("%v|%v|%v|%v", param["api_key"], param["transaction_id"], param["order_id"], gateway.GatewaySecret))
	responseJson, err := SendPayssionPaymentRequest(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", urlPath, param)
	log.SaveChannelHttpLog("GatewayPaymentDetail", param, responseJson, err, "PayssionPaymentDetail", nil, gateway)
	if err != nil {
		return nil, err
	}
	if !responseJson.Contains("result_code") || responseJson.Get("result_code").Int() != 200 {
		return nil, gerror.New("invalid request, result_code is nil or not 200")
	}
	if !responseJson.Contains("transaction") {
		return nil, gerror.New("invalid request, transaction is nil")
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())

	return parsePayssionPayment(responseJson), nil
}

func (c Payssion) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("Not Support")
}

// https://payssion.com/cn/docs/#api-reference-payment-details
func (c Payssion) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	payment := util.GetPaymentByPaymentId(ctx, refund.PaymentId)
	if payment == nil {
		return nil, gerror.New("payment not found")
	}
	detail, err := c.GatewayPaymentDetail(ctx, gateway, payment.GatewayPaymentId, payment)
	if err != nil {
		return nil, err
	}
	if detail.RefundSequence > int64(refund.RefundGatewaySequence) {
		return &gateway_bean.GatewayPaymentRefundResp{
			GatewayRefundId: gatewayRefundId,
			Status:          consts.RefundSuccess,
			Reason:          refund.RefundComment,
			RefundAmount:    refund.RefundAmount,
			Currency:        strings.ToUpper(refund.Currency),
			RefundTime:      gtime.Now(),
		}, nil
	} else {
		return &gateway_bean.GatewayPaymentRefundResp{
			GatewayRefundId: gatewayRefundId,
			Status:          consts.RefundCreated,
			Reason:          refund.RefundComment,
			RefundAmount:    refund.RefundAmount,
			Currency:        strings.ToUpper(refund.Currency),
			RefundTime:      gtime.Now(),
		}, nil
	}
}

func (c Payssion) GatewayRefund(ctx context.Context, gateway *entity.MerchantGateway, createPaymentRefundContext *gateway_bean.GatewayNewPaymentRefundReq) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	detail, err := c.GatewayPaymentDetail(ctx, gateway, createPaymentRefundContext.Payment.GatewayPaymentId, createPaymentRefundContext.Payment)
	if err != nil {
		return &gateway_bean.GatewayPaymentRefundResp{
			GatewayRefundId: createPaymentRefundContext.Payment.GatewayPaymentId,
			Status:          consts.RefundFailed,
			Type:            consts.RefundTypeGateway,
			Reason:          fmt.Sprintf("Get Gateway Refund Sequence Failed:%s", err.Error()),
		}, nil
	}

	urlPath := "/api/v1/refunds"
	param := map[string]interface{}{}
	param["transaction_id"] = createPaymentRefundContext.Payment.GatewayPaymentId
	if createPaymentRefundContext.GatewayCurrencyExchange != nil && createPaymentRefundContext.ExchangeRefundAmount > 0 && len(createPaymentRefundContext.ExchangeRefundCurrency) > 0 {
		param["amount"] = utility.ConvertCentToDollarStr(createPaymentRefundContext.ExchangeRefundAmount, createPaymentRefundContext.ExchangeRefundCurrency)
		param["currency"] = strings.ToUpper(createPaymentRefundContext.ExchangeRefundCurrency)
	} else {
		param["amount"] = utility.ConvertCentToDollarStr(createPaymentRefundContext.Refund.RefundAmount, createPaymentRefundContext.Refund.Currency)
		param["currency"] = strings.ToUpper(createPaymentRefundContext.Refund.Currency)
	}
	param["track_id"] = createPaymentRefundContext.Refund.RefundId
	param["description"] = createPaymentRefundContext.Refund.RefundComment
	param["api_key"] = gateway.GatewayKey
	param["api_sig"] = utility.MD5(fmt.Sprintf("%v|%v|%v|%v|%v", param["api_key"], param["transaction_id"], param["amount"], param["currency"], gateway.GatewaySecret))
	responseJson, err := SendPayssionPaymentRequest(ctx, gateway.GatewayKey, gateway.GatewaySecret, "POST", urlPath, param)
	log.SaveChannelHttpLog("GatewayRefund", param, responseJson, err, fmt.Sprintf("%s-%d", gateway.GatewayName, gateway.Id), nil, gateway)
	if err != nil {
		return nil, err
	}
	if !responseJson.Contains("result_code") || responseJson.Get("result_code").Int() != 200 {
		return &gateway_bean.GatewayPaymentRefundResp{
			GatewayRefundId: createPaymentRefundContext.Payment.GatewayPaymentId,
			Status:          consts.RefundFailed,
			Type:            consts.RefundTypeGateway,
			Reason:          fmt.Sprintf("invalid request, %s", responseJson.Get("description").String()),
		}, nil
	}
	if !responseJson.Contains("refund") {
		return &gateway_bean.GatewayPaymentRefundResp{
			GatewayRefundId: createPaymentRefundContext.Payment.GatewayPaymentId,
			Status:          consts.RefundFailed,
			Type:            consts.RefundTypeGateway,
			Reason:          fmt.Sprintf("invalid request, %s", responseJson.Get("description").String()),
		}, nil
	}
	g.Log().Debugf(ctx, "responseJson :%s", responseJson.String())
	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: responseJson.Get("refund.transaction_id").String(),
		Status:          consts.RefundCreated,
		Type:            consts.RefundTypeGateway,
		RefundSequence:  unibee.Int64(detail.RefundSequence),
	}, nil
}

func (c Payssion) GatewayRefundCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("Not Support")
}

func parsePayssionPayment(item *gjson.Json) *gateway_bean.GatewayPaymentRo {
	var status = consts.PaymentCreated
	var authorizeStatus = consts.WaitingAuthorized
	if strings.Compare(item.Get("transaction.state").String(), "pending") == 0 {
		authorizeStatus = consts.Authorized
	} else if strings.Compare(item.Get("transaction.state").String(), "completed") == 0 {
		status = consts.PaymentSuccess
	} else if strings.Compare(item.Get("transaction.state").String(), "cancelled") == 0 {
		status = consts.PaymentCancelled
	} else if strings.Compare(item.Get("transaction.state").String(), "failed") == 0 {
		status = consts.PaymentFailed
	}

	var authorizeReason = ""
	var paymentAmount = utility.ConvertDollarStrToCent(item.GetJson("transaction.amount").String(), item.Get("transaction.currency").String())
	var paymentMethod = ""
	if item.Contains("selected_payment_method") && item.GetJson("selected_payment_method").Contains("payins") {
		//paymentAmount = utility.ConvertDollarStrToCent(item.GetJson("orderSum").String(), item.Get("orderCurrency").String())
		//paymentMethod = item.Get("payin_currency").String() + "|" + item.Get("payin_network").String()
	}
	var paidTime *gtime.Time
	if item.Contains("transaction.updated") {
		if t, err := gtime.StrToTime(item.Get("transaction.updated").String()); err == nil {
			paidTime = t
		}
	}
	var refundSequence int64 = 0
	if item.Contains("transaction.refund") {
		refundSequence = int64(item.Get("transaction.refund").Float64() * 100)
	}
	return &gateway_bean.GatewayPaymentRo{
		GatewayPaymentId:     item.Get("transaction.transaction_id").String(),
		Status:               status,
		AuthorizeStatus:      authorizeStatus,
		AuthorizeReason:      authorizeReason,
		CancelReason:         "",
		PaymentData:          item.String(),
		TotalAmount:          utility.ConvertDollarStrToCent(item.Get("transaction.amount").String(), item.Get("transaction.currency").String()),
		PaymentAmount:        paymentAmount,
		GatewayPaymentMethod: paymentMethod,
		PaidTime:             paidTime,
		RefundSequence:       refundSequence,
	}
}

func SendPayssionPaymentRequest(ctx context.Context, publicKey string, privateKey string, method string, urlPath string, param map[string]interface{}) (res *gjson.Json, err error) {
	utility.Assert(param != nil, "param is nil")
	domain := "http://www.payssion.com"
	if !config.GetConfigInstance().IsProd() {
		domain = "http://sandbox.payssion.com"
		param["pm_id"] = "payssion_test"
	}
	jsonData, err := gjson.Marshal(param)
	jsonString := string(jsonData)
	utility.Assert(err == nil, fmt.Sprintf("json format error %s param %s", err, param))
	g.Log().Debugf(ctx, "\nPayssion_Start %s %s %s %s\n", method, urlPath, publicKey, jsonString)
	body := []byte(jsonString)
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	response, err := utility.SendRequest(domain+urlPath, method, body, headers)
	g.Log().Debugf(ctx, "\nPayssion_End %s %s response: %s error %s\n", method, urlPath, response, err)
	if err != nil {
		return nil, err
	}
	responseJson, err := gjson.LoadJson(string(response))
	if err != nil {
		return nil, err
	}
	return responseJson, nil
}
