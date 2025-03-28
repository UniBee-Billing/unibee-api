package _interface

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/default"
)

var KEY_MERCHANT_GATEWAY_SORT string = "KEY_MERCHANT_GATEWAY_SORT"

type GatewayPaymentType struct {
	Name        string `json:"name"`
	PaymentType string `json:"paymentType"`
	CountryName string `json:"countryName"`
	AutoCharge  bool   `json:"autoCharge"`
	Category    string `json:"category"`
}
type GatewayInfo struct {
	Name                          string
	Description                   string
	DisplayName                   string
	GatewayWebsiteLink            string
	GatewayWebhookIntegrationLink string
	GatewayLogo                   string
	GatewayIcons                  []string
	GatewayType                   int64
	Sort                          int64
	CurrencyExchangeEnabled       bool
	QueueForRefund                bool
	GatewayPaymentTypes           []*GatewayPaymentType
	IsStaging                     bool
	PublicKeyName                 string
	PrivateSecretName             string
	SubGatewayName                string
	Host                          string
	AutoChargeEnabled             bool
}

type GatewayTestReq struct {
	Key                 string
	Secret              string
	SubGateway          string
	GatewayPaymentTypes []*GatewayPaymentType
}

type GatewayInterface interface {
	GatewayInfo(ctx context.Context) *GatewayInfo
	GatewayTest(ctx context.Context, req *GatewayTestReq) (icon string, gatewayType int64, err error)
	// User
	GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error)
	// Balance
	GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, gatewayUserId string) (res *gateway_bean.GatewayUserDetailQueryResp, err error)
	GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error)
	// Payment
	GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error)
	GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error)
	GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error)
	GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, currency string, metadata map[string]interface{}) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error)
	GatewayNewPayment(ctx context.Context, gateway *entity.MerchantGateway, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error)
	GatewayCapture(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error)
	GatewayCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error)
	GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error)
	GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error)
	GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error)
	GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error)
	GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error)
	GatewayRefund(ctx context.Context, gateway *entity.MerchantGateway, createPaymentRefundContext *gateway_bean.GatewayNewPaymentRefundReq) (res *gateway_bean.GatewayPaymentRefundResp, err error)
	GatewayRefundCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error)
}

type GatewayWebhookInterface interface {
	GatewayCheckAndSetupWebhook(ctx context.Context, gateway *entity.MerchantGateway) (err error)
	GatewayWebhook(r *ghttp.Request, gateway *entity.MerchantGateway)
	GatewayRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayRedirectResp, err error)
	GatewayNewPaymentMethodRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (err error)
}
