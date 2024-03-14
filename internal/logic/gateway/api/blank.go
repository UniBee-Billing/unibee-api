package api

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/oversea_pay"
)

type Blank struct {
}

func (b Blank) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) GatewayRefundCancel(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId int64, data *gjson.Json) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) GatewayTest(ctx context.Context, key string, secret string) (icon string, gatewayType int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res *gateway_bean.GatewayPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) GatewayNewPayment(ctx context.Context, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) GatewayCapture(ctx context.Context, pay *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) GatewayCancel(ctx context.Context, pay *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) GatewayRefund(ctx context.Context, pay *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	//TODO implement me
	panic("implement me")
}
