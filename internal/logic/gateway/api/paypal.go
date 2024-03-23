package api

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/plutov/paypal/v4"
	"net/http"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/oversea_pay"
)

// link：https://developer.paypal.com/docs/api/payments/v1/#payment_create
// https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_transactions
// clientId ATaWQ8G9oJNFyle9YCt59
// Secret EHUy5GALkYr1Qp0n6MepJY8LnUwYCBIWElG4Iv_DO3mdYcbB2l6zwJxk99OrPhbdNRLk7GkHEqb5RHEA

// Other ClientId AXy9orp-CDaHhBZ9C78QHW2BKZpACgroqo85_NIOa9mIfJ9QnSVKzY-X_rivR_fTUUr6aLjcJsj6sDur
// Other Secret EBoIiUSkCKeSk49hHSgTem1qnjzzJgRQHDEHvGpzlLEf_nIoJd91xu8rPOBDCdR_UYNKVxJE-UgS2iCw

// Other 2 ClientId AT-HU_WUeHCis_uqkU2Y8-0f54qq_QkoNXJeBj1-4S01__m1OLQn1jXnG9F86bcaH5TbcYiFed7UBRGH
// Other 2 Secret  EL2TLXWp_6XyZEtYqeRjLLVb9S_uYjwZOrBUiqhHhw96-50VisMsQvBDA09qMVntXrPf6TukiyfRCkG0

//APIBaseSandBox = "https://api-m.sandbox.paypal.com"
//APIBaseLive = "https://api-m.paypal.com"

type Paypal struct {
}

func (p Paypal) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayRefundCancel(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, data *gjson.Json) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayTest(ctx context.Context, key string, secret string) (icon string, gatewayType int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res *gateway_bean.GatewayPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	//TODO implement me
	panic("implement me")
}

func init() {
	// gateway_webhook_entry
}

// todo mark 确认改造成单例是否可行，不用每次都去获取 accessToken
func NewClient(clientID string, secret string, APIBase string) (*paypal.Client, error) {
	if clientID == "" || secret == "" || APIBase == "" {
		return nil, errors.New("ClientID, Secret and APIBase are required to create a Client")
	}

	return &paypal.Client{
		Client:   &http.Client{},
		ClientID: clientID,
		Secret:   secret,
		APIBase:  APIBase,
	}, nil
}

func (p Paypal) GatewayNewPayment(ctx context.Context, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayCapture(ctx context.Context, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayCancel(ctx context.Context, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayRefund(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	//TODO implement me
	panic("implement me")
}
