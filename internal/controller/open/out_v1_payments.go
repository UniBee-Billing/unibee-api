package open

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/api/open/payment"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/payment/gateway/ro"
	"go-oversea-pay/internal/logic/payment/service"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strings"
)

func (c *ControllerPayment) Payments(ctx context.Context, req *payment.PaymentsReq) (res *payment.PaymentsRes, err error) {
	utility.Assert(req != nil, "request req is nil")
	utility.Assert(req.Amount != nil, "amount is nil")
	utility.Assert(req.Amount.Value > 0, "amount value is nil")
	utility.Assert(len(req.Amount.Currency) > 0, "amount currency is nil")
	//类似日元的小数尾数必须为 0 检查
	currencyNumberCheck(req.Amount)
	utility.Assert(len(req.CountryCode) > 0, "countryCode is nil")
	utility.Assert(req.MerchantId > 0, "merchantId is nil")
	utility.Assert(req.PaymentMethod != nil, "paymentmethod is nil")
	utility.Assert(len(req.PaymentMethod.Channel) > 0, "paymentmethod type is nil")
	utility.Assert(len(req.Reference) > 0, "reference is nil")
	utility.Assert(len(req.ShopperReference) > 0, "shopperReference type is nil")
	utility.Assert(len(req.ShopperEmail) > 0, "shopperEmail is nil")
	utility.Assert(req.LineItems != nil, "lineItems is nil")
	utility.Assert(len(req.Channel) > 0, "outchannel is nil")
	utility.Assert(strings.Contains("WEB，WAP，APP, MINI, INWALLET", req.Channel), "outchannel is invalid, should be WEB，WAP，APP, MINI, INWALLET")

	openApiConfig, merchantInfo := merchantCheck(ctx, req.MerchantId)
	payChannel := query.GetPayChannelByChannel(ctx, req.PaymentMethod.Channel)
	utility.Assert(payChannel != nil, "找不到支付方式 type:"+req.PaymentMethod.Channel)
	//支付方式绑定校验 todo mark

	createPayContext := &ro.CreatePayContext{
		OpenApiId:            int64(openApiConfig.Id),
		PayChannel:           payChannel,
		PaymentBrandAddition: req.PaymentBrandAddition,
		Pay: &entity.OverseaPay{
			BizId:             req.Reference,
			BizType:           consts.PAYMENT_BIZ_TYPE_ORDER,
			ChannelId:         int64(payChannel.Id),
			PaymentFee:        req.Amount.Value,
			Currency:          req.Amount.Currency,
			CountryCode:       req.CountryCode,
			MerchantId:        merchantInfo.Id,
			CompanyId:         merchantInfo.CompanyId,
			NotifyUrl:         req.ReturnUrl,
			CaptureDelayHours: req.CaptureDelayHours,
		},
		Platform:                 req.Channel,
		DeviceType:               req.DeviceType,
		UserId:                   req.ShopperReference,
		ShopperEmail:             req.ShopperEmail,
		ShopperLocale:            req.ShopperLocale,
		Mobile:                   req.TelephoneNumber,
		MediaInfo:                req.Metadata,
		Items:                    req.LineItems,
		BillingDetails:           req.BillingAddress,
		ShippingDetails:          req.DetailAddress,
		ShopperName:              req.ShopperName,
		ShopperInteraction:       req.ShopperInteraction,
		RecurringProcessingModel: req.RecurringProcessingModel,
		StorePaymentMethod:       req.StorePaymentMethod,
		TokenId:                  req.PaymentMethod.TokenId,
		DeviceFingerprint:        req.DeviceFingerprint,
		MerchantOrderReference:   req.MerchantOrderReference,
		DateOfBirth:              gtime.ParseTimeFromContent(req.DateOfBrith, "YYYY-MM-DD"),
	}

	resp, err := service.DoChannelPay(ctx, createPayContext)
	utility.Assert(err == nil, fmt.Sprintf("%+v", err))
	res = &payment.PaymentsRes{
		Status:       "Pending",
		PspReference: resp.PayOrderNo,
		Reference:    req.Reference,
		Action:       resp.Action,
	}
	return
}
