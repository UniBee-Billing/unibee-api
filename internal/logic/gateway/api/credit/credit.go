package credit

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/credit/payment"
	refund2 "unibee/internal/logic/credit/refund"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/default"
)

type Credit struct {
}

func (c Credit) GatewayInfo(ctx context.Context) *_interface.GatewayInfo {
	return &_interface.GatewayInfo{
		DisplayName:  "Credit",
		GatewayIcons: []string{"https://unibee.dev/wp-content/uploads/2024/05/logo-white.svg?ver=1718007070"},
		GatewayType:  consts.GatewayTypeCredit,
	}
}

func (c Credit) GatewayTest(ctx context.Context, key string, secret string, subGateway string) (icon string, gatewayType int64, err error) {
	return "https://unibee.dev/wp-content/uploads/2024/05/logo-white.svg?ver=1718007070", consts.GatewayTypeCredit, nil
}

func (c Credit) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *gateway_bean.GatewayUserCreateResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Credit) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64) (res *gateway_bean.GatewayUserDetailQueryResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Credit) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayMerchantBalanceQueryResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Credit) GatewayUserAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserAttachPaymentMethodResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Credit) GatewayUserDeAttachPaymentMethodQuery(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, gatewayPaymentMethod string) (res *gateway_bean.GatewayUserDeAttachPaymentMethodResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Credit) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, req *gateway_bean.GatewayUserPaymentMethodReq) (res *gateway_bean.GatewayUserPaymentMethodListResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Credit) GatewayUserCreateAndBindPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64, currency string, metadata map[string]interface{}) (res *gateway_bean.GatewayUserPaymentMethodCreateAndBindResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Credit) GatewayNewPayment(ctx context.Context, gateway *entity.MerchantGateway, createPayContext *gateway_bean.GatewayNewPaymentReq) (res *gateway_bean.GatewayNewPaymentResp, err error) {
	creditPayment, err := payment.NewCreditPayment(ctx, &payment.CreditPaymentInternalReq{
		UserId:                  createPayContext.Pay.UserId,
		MerchantId:              createPayContext.Pay.MerchantId,
		ExternalCreditPaymentId: createPayContext.Pay.PaymentId,
		InvoiceId:               createPayContext.Pay.InvoiceId,
		CurrencyAmount:          createPayContext.Pay.TotalAmount,
		Currency:                createPayContext.Pay.Currency,
		CreditType:              consts.CreditAccountTypeMain,
		Name:                    "",
		Description:             createPayContext.Invoice.InvoiceName,
	})
	if err != nil {
		return nil, err
	}
	return &gateway_bean.GatewayNewPaymentResp{
		Payment:                createPayContext.Pay,
		Status:                 consts.PaymentSuccess,
		GatewayPaymentId:       creditPayment.CreditPayment.CreditPaymentId,
		GatewayPaymentIntentId: creditPayment.CreditPayment.CreditPaymentId,
		GatewayPaymentMethod:   "",
		Link:                   "",
	}, nil
}

func (c Credit) GatewayCapture(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCaptureResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Credit) GatewayCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment) (res *gateway_bean.GatewayPaymentCancelResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Credit) GatewayCryptoFiatTrans(ctx context.Context, from *gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq) (to *gateway_bean.GatewayCryptoToCurrencyAmountDetailRes, err error) {
	return nil, gerror.New("Not Support")
}

func (c Credit) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *gateway_bean.GatewayPaymentListReq) (res []*gateway_bean.GatewayPaymentRo, err error) {
	return nil, gerror.New("Not Support")
}

func (c Credit) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string, payment *entity.Payment) (res *gateway_bean.GatewayPaymentRo, err error) {
	return nil, gerror.New("Not Support")
}

func (c Credit) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Credit) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("Not Support")
}

func (c Credit) GatewayRefund(ctx context.Context, gateway *entity.MerchantGateway, createPaymentRefundContext *gateway_bean.GatewayNewPaymentRefundReq) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	creditRefund, err := refund2.NewCreditRefund(ctx, &refund2.CreditRefundInternalReq{
		UserId:                 createPaymentRefundContext.Payment.UserId,
		MerchantId:             createPaymentRefundContext.Payment.MerchantId,
		CreditPaymentId:        createPaymentRefundContext.Payment.GatewayPaymentId,
		ExternalCreditRefundId: createPaymentRefundContext.Refund.RefundId,
		InvoiceId:              createPaymentRefundContext.Refund.InvoiceId,
		RefundAmount:           createPaymentRefundContext.Refund.RefundAmount,
		Currency:               createPaymentRefundContext.Refund.Currency,
		Name:                   createPaymentRefundContext.Refund.RefundComment,
		Description:            createPaymentRefundContext.Refund.RefundCommentExplain,
	})
	if err != nil {
		return nil, err
	}
	return &gateway_bean.GatewayPaymentRefundResp{
		GatewayRefundId: creditRefund.CreditRefund.CreditRefundId,
		Status:          consts.RefundSuccess,
		Type:            consts.RefundTypeGateway,
	}, nil
}

func (c Credit) GatewayRefundCancel(ctx context.Context, gateway *entity.MerchantGateway, payment *entity.Payment, refund *entity.Refund) (res *gateway_bean.GatewayPaymentRefundResp, err error) {
	return nil, gerror.New("Not Support")
}
