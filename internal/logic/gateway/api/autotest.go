package api

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"unibee/api/bean"
	"unibee/internal/consts"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/default"
)

type AutoTest struct {
}

func (a AutoTest) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	return nil, gerror.New("not support")
}

func (a AutoTest) GatewayRefundCancel(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
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

func (a AutoTest) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, currency string, metadata map[string]interface{}) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	return &gateway_bean.GatewayUserPaymentMethodCreateAndBindResp{PaymentMethod: &bean.PaymentMethod{
		Id:   strconv.FormatUint(userId, 10),
		Type: "card",
		Data: gjson.New(""),
	}}, nil
}

func (a AutoTest) GatewayTest(ctx context.Context, key string, secret string) (icon string, gatewayType int64, err error) {
	return "http://autotest.unibee.com", consts.GatewayTypeCard, nil
}

func (a AutoTest) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	return &gateway_bean.GatewayUserAttachPaymentMethodResp{}, nil
}

func (a AutoTest) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	return &gateway_bean.GatewayUserDeAttachPaymentMethodResp{}, nil
}

func (a AutoTest) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	return &gateway_bean.GatewayUserPaymentMethodListResp{}, nil
}

func (a AutoTest) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	return &gateway_bean.GatewayUserCreateResp{GatewayUserId: strconv.FormatUint(user.Id, 10)}, nil
}

func (a AutoTest) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	return []*gateway_bean.GatewayPaymentRo{}, nil
}

func (a AutoTest) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	return []*gateway_bean.GatewayPaymentRefundResp{}, nil
}

func (a AutoTest) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
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

func (a AutoTest) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
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

func (a AutoTest) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	return &gateway_bean.GatewayMerchantBalanceQueryResp{
		AvailableBalance:       []*gateway_bean.GatewayBalance{},
		ConnectReservedBalance: []*gateway_bean.GatewayBalance{},
		PendingBalance:         []*gateway_bean.GatewayBalance{},
	}, nil
}

func (a AutoTest) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	return &gateway_bean.GatewayUserDetailQueryResp{
		GatewayUserId:        strconv.FormatUint(userId, 10),
		DefaultPaymentMethod: "",
		Balance:              nil,
		CashBalance:          []*gateway_bean.GatewayBalance{},
		InvoiceCreditBalance: []*gateway_bean.GatewayBalance{},
		Description:          "",
		Email:                "",
	}, nil
}

func (a AutoTest) GatewayNewPayment(ctx context.Context, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
	if createPayContext.CheckoutMode || !createPayContext.PayImmediate {
		return &gateway_bean.GatewayNewPaymentResp{
			Status:                 consts.PaymentCreated,
			GatewayPaymentId:       createPayContext.Pay.PaymentId,
			GatewayPaymentIntentId: createPayContext.Pay.PaymentId,
			Link:                   "http://unibee.top",
		}, nil
	} else {
		return &gateway_bean.GatewayNewPaymentResp{
			Status:                 consts.PaymentSuccess,
			GatewayPaymentId:       createPayContext.Pay.PaymentId,
			GatewayPaymentIntentId: createPayContext.Pay.PaymentId,
			Link:                   "http://unibee.top",
		}, nil
	}
}

func (a AutoTest) GatewayCapture(ctx context.Context, pay *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	return nil, gerror.New("not support")
}

func (a AutoTest) GatewayCancel(ctx context.Context, pay *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	return &gateway_bean.GatewayPaymentCancelResp{
		MerchantId:      strconv.FormatUint(pay.MerchantId, 10),
		GatewayCancelId: pay.PaymentId,
		PaymentId:       pay.PaymentId,
		Status:          consts.PaymentCancelled,
	}, nil
}

func (a AutoTest) GatewayRefund(ctx context.Context, pay *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: refund.RefundId,
		Status:          consts.RefundSuccess,
		Type:            consts.RefundTypeMarked,
	}, nil
}
