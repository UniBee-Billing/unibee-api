package _interface

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/ghttp"
	"unibee/internal/logic/gateway/ro"
	entity "unibee/internal/model/entity/oversea_pay"
)

type GatewayInterface interface {
	GatewayTest(ctx context.Context, key string, secret string) (err error)
	// User
	GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *ro.GatewayUserCreateInternalResp, err error)
	// Balance
	GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64) (res *ro.GatewayUserDetailQueryInternalResp, err error)
	GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *ro.GatewayMerchantBalanceQueryInternalResp, err error)
	// Payment
	GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64, gatewayPaymentMethod string) (res *ro.GatewayUserAttachPaymentMethodInternalResp, err error)
	GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64, gatewayPaymentMethod string) (res *ro.GatewayUserDeAttachPaymentMethodInternalResp, err error)
	GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64) (res *ro.GatewayUserPaymentMethodListInternalResp, err error)
	GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId int64, data *gjson.Json) (res *ro.GatewayUserPaymentMethodCreateAndBindInternalResp, err error)
	GatewayNewPayment(ctx context.Context, createPayContext *ro.NewPaymentInternalReq) (res *ro.NewPaymentInternalResp, err error)
	GatewayCapture(ctx context.Context, payment *entity.Payment) (res *ro.OutPayCaptureRo, err error)
	GatewayCancel(ctx context.Context, payment *entity.Payment) (res *ro.OutPayCancelRo, err error)
	GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *ro.GatewayPaymentListReq) (res []*ro.GatewayPaymentRo, err error)
	GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res *ro.GatewayPaymentRo, err error)
	GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*ro.OutPayRefundRo, err error)
	GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string) (res *ro.OutPayRefundRo, err error)
	GatewayRefund(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error)
	GatewayRefundCancel(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error)
}

type GatewayWebhookInterface interface {
	GatewayCheckAndSetupWebhook(ctx context.Context, gateway *entity.MerchantGateway) (err error)
	GatewayWebhook(r *ghttp.Request, gateway *entity.MerchantGateway)
	GatewayRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (res *ro.GatewayRedirectInternalResp, err error)
}
