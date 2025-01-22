package api

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type AutoTestCrypto struct {
}

func (a AutoTestCrypto) GatewayInfo(ctx context.Context) *_interface.GatewayInfo {
	return &_interface.GatewayInfo{
		DisplayName:  "AutoTestCrypto",
		GatewayLogo:  "https://unibee.dev/wp-content/uploads/2024/05/logo-white.svg?ver=1718007070",
		GatewayIcons: []string{"https://unibee.dev/wp-content/uploads/2024/05/logo-white.svg?ver=1718007070"},
		GatewayType:  consts.GatewayTypeCard,
	}
}

func (a AutoTestCrypto) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	return &gateway_bean.GatewayCryptoToCurrencyAmountDetailRes{
		Amount:         from.Amount,
		Currency:       from.Currency,
		CountryCode:    from.CountryCode,
		CryptoAmount:   utility.RoundUp(float64(from.Amount) / 1),
		CryptoCurrency: "USD",
		Rate:           1,
	}, nil
}

func (a AutoTestCrypto) GatewayRefundCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
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

func (a AutoTestCrypto) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, currency string, metadata map[string]interface{}) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	return &gateway_bean.GatewayUserPaymentMethodCreateAndBindResp{PaymentMethod: &gateway_bean.PaymentMethod{
		Id:   strconv.FormatUint(userId, 10),
		Type: "card",
		Data: gjson.New(""),
	}}, nil
}

func (a AutoTestCrypto) GatewayTest(ctx context.Context, key string, secret string) (icon string, gatewayType int64, err error) {
	return "http://autotest_crypto.unibee.com", consts.GatewayTypeCrypto, nil
}

func (a AutoTestCrypto) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	return &gateway_bean.GatewayUserAttachPaymentMethodResp{}, nil
}

func (a AutoTestCrypto) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	return &gateway_bean.GatewayUserDeAttachPaymentMethodResp{}, nil
}

func (a AutoTestCrypto) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	return &gateway_bean.GatewayUserPaymentMethodListResp{}, nil
}

func (a AutoTestCrypto) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	return &gateway_bean.GatewayUserCreateResp{GatewayUserId: strconv.FormatUint(user.Id, 10)}, nil
}

func (a AutoTestCrypto) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	return []*gateway_bean.GatewayPaymentRo{}, nil
}

func (a AutoTestCrypto) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	return []*gateway_bean.GatewayPaymentRefundResp{}, nil
}

func (a AutoTestCrypto) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
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

func (a AutoTestCrypto) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
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

func (a AutoTestCrypto) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	return &gateway_bean.GatewayMerchantBalanceQueryResp{
		AvailableBalance:       []*gateway_bean.GatewayBalance{},
		ConnectReservedBalance: []*gateway_bean.GatewayBalance{},
		PendingBalance:         []*gateway_bean.GatewayBalance{},
	}, nil
}

func (a AutoTestCrypto) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
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

func (a AutoTestCrypto) GatewayNewPayment(ctx context.Context, gateway *entity.MerchantGateway, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
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

func (a AutoTestCrypto) GatewayCapture(ctx context.Context, gateway *entity.MerchantGateway, pay *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTestCrypto) GatewayCancel(ctx context.Context, gateway *entity.MerchantGateway, pay *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	return &gateway_bean.GatewayPaymentCancelResp{
		MerchantId:      strconv.FormatUint(pay.MerchantId, 10),
		GatewayCancelId: pay.PaymentId,
		PaymentId:       pay.PaymentId,
		Status:          consts.PaymentCancelled,
	}, nil
}

func (a AutoTestCrypto) GatewayRefund(ctx context.Context, gateway *entity.MerchantGateway, pay *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: refund.RefundId,
		Status:          consts.RefundCreated,
		Type:            consts.RefundTypeMarked,
	}, nil
}
