package api

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"unibee/internal/consts"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/default"
)

type Wire struct {
}

func (w Wire) GatewayTest(ctx context.Context, key string, secret string) (icon string, gatewayType int64, err error) {
	return "https://api.unibee.top/oss/file/d6y8q1dfe2owqjnayq.svg", consts.GatewayTypeWireTransfer, gerror.New("Please setup by wire transfer setup api")
}

func (w Wire) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	return nil, gerror.New("not support")
}

func (w Wire) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	return nil, gerror.New("not support")
}

func (w Wire) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	return nil, gerror.New("not support")
}

func (w Wire) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	return nil, gerror.New("not support")
}

func (w Wire) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	return nil, gerror.New("not support")
}

func (w Wire) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	return nil, gerror.New("not support")
}

func (w Wire) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, currency string, metadata map[string]interface{}) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	return nil, gerror.New("not support")
}

func (w Wire) GatewayNewPayment(ctx context.Context, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
	return &gateway_bean.GatewayNewPaymentResp{
		Status:                 consts.PaymentCreated,
		GatewayPaymentId:       createPayContext.Pay.PaymentId,
		GatewayPaymentIntentId: createPayContext.Pay.PaymentId,
		Link:                   "",
	}, nil
}

func (w Wire) GatewayCapture(ctx context.Context, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	return nil, gerror.New("not support")
}

func (w Wire) GatewayCancel(ctx context.Context, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	return &gateway_bean.GatewayPaymentCancelResp{Status: consts.PaymentCancelled}, nil
}

func (w Wire) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	return nil, gerror.New("not support")
}

func (w Wire) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	return []*gateway_bean.GatewayPaymentRo{}, nil
}

func (w Wire) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
	return &gateway_bean.GatewayPaymentRo{
		GatewayPaymentId:     gatewayPaymentId,
		Status:               payment.Status,
		AuthorizeStatus:      payment.AuthorizeStatus,
		AuthorizeReason:      payment.AuthorizeReason,
		CancelReason:         payment.FailureReason,
		PaymentData:          payment.PaymentData,
		TotalAmount:          payment.TotalAmount,
		PaymentAmount:        payment.PaymentAmount,
		GatewayPaymentMethod: payment.GatewayPaymentMethod,
		Currency:             payment.Currency,
		PaidTime:             gtime.NewFromTimeStamp(payment.PaidTime),
		CreateTime:           gtime.NewFromTimeStamp(payment.CreateTime),
		CancelTime:           gtime.NewFromTimeStamp(payment.CancelTime),
	}, nil
}

func (w Wire) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	return []*gateway_bean.GatewayPaymentRefundResp{}, nil
}

func (w Wire) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return &gateway_bean.GatewayPaymentRefundResp{
		MerchantId:       strconv.FormatUint(refund.MerchantId, 10),
		GatewayRefundId:  refund.RefundId,
		GatewayPaymentId: refund.PaymentId,
		Status:           consts.RefundSuccess,
		Reason:           refund.RefundComment,
		RefundAmount:     refund.RefundAmount,
		Currency:         refund.Currency,
		RefundTime:       gtime.NewFromTimeStamp(refund.RefundTime),
	}, nil
}

func (w Wire) GatewayRefund(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: refund.RefundId,
		Status:          consts.RefundCreated,
		Type:            consts.RefundTypeMarked,
	}, nil
}

func (w Wire) GatewayRefundCancel(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
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
