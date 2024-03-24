package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"unibee/api/bean"
	"unibee/internal/consts"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/test"
)

func TestCryptoPayment(t *testing.T) {
	ctx := context.Background()
	var paymentId string
	var one *entity.Payment
	var err error
	gateway := test.TestCryptoGateway
	t.Run("Test for crypto payment checkout_new|cancel", func(t *testing.T) {
		res, err := GatewayPaymentCreate(ctx, &gateway_bean.GatewayNewPaymentReq{
			CheckoutMode: true,
			Pay: &entity.Payment{
				ExternalPaymentId: uuid.New().String(),
				BizType:           consts.BizTypeOneTime,
				UserId:            test.TestUser.Id,
				GatewayId:         gateway.Id,
				TotalAmount:       100,
				Currency:          "usd",
				CountryCode:       "CN",
				MerchantId:        test.TestMerchant.Id,
				CompanyId:         test.TestMerchant.CompanyId,
				ReturnUrl:         "",
			},
			Gateway:        gateway,
			ExternalUserId: test.TestUser.ExternalUserId,
			Email:          test.TestUser.Email,
			DaysUtilDue:    consts.DEFAULT_DAY_UTIL_DUE,
			PayImmediate:   false,
			Invoice: &bean.InvoiceSimplify{
				TotalAmount:             100,
				Currency:                "usd",
				TotalAmountExcludingTax: 0,
				TaxAmount:               0,
				SendStatus:              consts.InvoiceSendStatusUnnecessary,
				DayUtilDue:              consts.DEFAULT_DAY_UTIL_DUE,
			},
		})
		require.Nil(t, err)
		require.NotNil(t, res)
		paymentId = res.PaymentId
		require.NotNil(t, paymentId)
		require.Equal(t, true, res.Status == consts.PaymentCreated)
		require.Equal(t, true, len(res.Link) > 0)
		one = query.GetPaymentByPaymentId(ctx, paymentId)
		require.NotNil(t, one)
		require.Equal(t, "USD", one.Currency)
		require.Equal(t, int64(100), one.TotalAmount)
		require.Equal(t, true, len(one.InvoiceId) > 0)
		require.Equal(t, true, len(one.CryptoCurrency) > 0)
		require.Equal(t, true, one.CryptoAmount > 0)
		err = PaymentGatewayCancel(ctx, one)
		require.Nil(t, err)
		one = query.GetPaymentByPaymentId(ctx, paymentId)
		require.NotNil(t, one)
		require.Equal(t, true, one.Status == consts.PaymentCancelled)
	})
	t.Run("Test for crypto payment HardDelete", func(t *testing.T) {
		err = HardDeletePayment(ctx, test.TestMerchant.Id, paymentId)
		require.Nil(t, err)
	})

	var refundId string
	t.Run("Test for crypto payment automatic|refund", func(t *testing.T) {
		res, err := GatewayPaymentCreate(ctx, &gateway_bean.GatewayNewPaymentReq{
			CheckoutMode: false,
			Pay: &entity.Payment{
				ExternalPaymentId: uuid.New().String(),
				BizType:           consts.BizTypeOneTime,
				UserId:            test.TestUser.Id,
				GatewayId:         gateway.Id,
				TotalAmount:       100,
				Currency:          "usd",
				CountryCode:       "CN",
				MerchantId:        test.TestMerchant.Id,
				CompanyId:         test.TestMerchant.CompanyId,
				ReturnUrl:         "",
			},
			Gateway:        gateway,
			ExternalUserId: test.TestUser.ExternalUserId,
			Email:          test.TestUser.Email,
			DaysUtilDue:    consts.DEFAULT_DAY_UTIL_DUE,
			PayImmediate:   true,
			Invoice: &bean.InvoiceSimplify{
				TotalAmount:             100,
				Currency:                "usd",
				TotalAmountExcludingTax: 0,
				TaxAmount:               0,
				SendStatus:              consts.InvoiceSendStatusUnnecessary,
				DayUtilDue:              consts.DEFAULT_DAY_UTIL_DUE,
			},
		})
		require.Nil(t, err)
		require.NotNil(t, res)
		paymentId = res.PaymentId

		require.Equal(t, true, res.Status == consts.PaymentSuccess)
		require.Equal(t, true, len(res.Link) > 0)
		one = query.GetPaymentByPaymentId(ctx, paymentId)
		require.NotNil(t, one)
		require.Equal(t, "USD", one.Currency)
		require.Equal(t, int64(100), one.TotalAmount)
		err = PaymentGatewayCancel(ctx, one)
		require.NotNil(t, err)
		refundRes, err := GatewayPaymentRefundCreate(ctx, &NewPaymentRefundInternalReq{
			PaymentId:        one.PaymentId,
			ExternalRefundId: uuid.New().String(),
			RefundAmount:     100,
			Currency:         "usd",
			Reason:           "test",
		})
		require.Nil(t, err)
		require.NotNil(t, refundRes)
		refundId = refundRes.RefundId
		refund := query.GetRefundByRefundId(ctx, refundId)
		require.NotNil(t, refund)
		require.Equal(t, true, refund.Status == consts.RefundCreated)
		require.Equal(t, 1, refund.Type)
		err = PaymentRefundGatewayCancel(ctx, refund)
		require.Nil(t, err)
		err = HardDeleteRefund(ctx, test.TestMerchant.Id, refundId)
		require.Nil(t, err)
		refundRes, err = MarkPaymentRefundCreate(ctx, &NewPaymentRefundInternalReq{
			PaymentId:        one.PaymentId,
			ExternalRefundId: uuid.New().String(),
			RefundAmount:     100,
			Currency:         "usd",
			Reason:           "test",
		})
		require.Nil(t, err)
		require.NotNil(t, refundRes)
		refundId = refundRes.RefundId
		refund = query.GetRefundByRefundId(ctx, refundId)
		require.NotNil(t, refund)
		require.Equal(t, true, refund.Status == consts.RefundSuccess)
		require.Equal(t, 2, refund.Type)
	})
	t.Run("Test for crypto payment HardDelete", func(t *testing.T) {
		err = HardDeletePayment(ctx, test.TestMerchant.Id, paymentId)
		require.Nil(t, err)
		err = HardDeleteRefund(ctx, test.TestMerchant.Id, refundId)
		require.Nil(t, err)
	})
}
