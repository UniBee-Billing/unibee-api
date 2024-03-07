package api

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"unibee/internal/logic/gateway/ro"
	entity "unibee/internal/model/entity/oversea_pay"
)

type Invalid struct{}

func (i Invalid) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId int64, data *gjson.Json) (res *ro.GatewayUserPaymentMethodCreateAndBindInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayTest(ctx context.Context, key string, secret string) (err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64, gatewayPaymentMethod string) (res *ro.GatewayUserAttachPaymentMethodInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64, gatewayPaymentMethod string) (res *ro.GatewayUserDeAttachPaymentMethodInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64) (res *ro.GatewayUserPaymentMethodListInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *ro.GatewayUserCreateInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *ro.GatewayPaymentListReq) (res []*ro.GatewayPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res *ro.GatewayPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *ro.GatewayMerchantBalanceQueryInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64) (res *ro.GatewayUserDetailQueryInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayNewPayment(ctx context.Context, createPayContext *ro.NewPaymentInternalReq) (res *ro.NewPaymentInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayCapture(ctx context.Context, payment *entity.Payment) (res *ro.OutPayCaptureRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayCancel(ctx context.Context, payment *entity.Payment) (res *ro.OutPayCancelRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayPayStatusCheck(ctx context.Context, payment *entity.Payment) (res *ro.GatewayPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayRefundStatusCheck(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayRefund(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}
