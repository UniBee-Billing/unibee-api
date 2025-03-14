package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/api/merchant/profile"
	"unibee/internal/cmd/config"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/analysis/segment"
	"unibee/internal/logic/currency"
	"unibee/internal/logic/email"
	"unibee/internal/logic/fiat_exchange"
	"unibee/internal/logic/merchant_config"
	"unibee/internal/logic/vat_gateway"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/time"
	"unibee/utility"
	"unibee/utility/unibee"
)

func (c *ControllerProfile) Get(ctx context.Context, req *profile.GetReq) (res *profile.GetRes, err error) {
	var member *entity.MerchantMember
	var isOwner = false
	var memberRoles = make([]*bean.MerchantRole, 0)
	if _interface.Context().Get(ctx) != nil && _interface.Context().Get(ctx).MerchantMember != nil {
		member = query.GetMerchantMemberById(ctx, _interface.Context().Get(ctx).MerchantMember.Id)
		if member != nil {
			isOwner, memberRoles = detail.ConvertMemberRole(ctx, member)
		}
	}
	merchant := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	utility.Assert(merchant != nil, "merchant not found")
	vatGatewayName, vatGatewayKey := vat_gateway.GetDefaultMerchantVatConfig(ctx, merchant.Id)
	if vatGatewayName != "vatsense" {
		vatGatewayKey = ""
	}
	_, emailData := email.GetDefaultMerchantEmailConfig(ctx, merchant.Id)
	var emailSender *bean.Sender
	one := email.GetMerchantEmailSender(ctx, _interface.GetMerchantId(ctx))
	if one != nil {
		emailSender = &bean.Sender{
			Name:    one.Name,
			Address: one.Address,
		}
	}
	exchangeApiKey := ""
	exchangeApiKeyConfig := merchant_config.GetMerchantConfig(ctx, _interface.GetMerchantId(ctx), fiat_exchange.FiatExchangeApiKey)
	if exchangeApiKeyConfig != nil {
		exchangeApiKey = exchangeApiKeyConfig.ConfigValue
	}
	return &profile.GetRes{
		Merchant:             bean.SimplifyMerchant(merchant),
		MerchantMember:       detail.ConvertMemberToDetail(ctx, member),
		Currency:             currency.GetMerchantCurrencies(),
		Env:                  config.GetConfigInstance().Env,
		IsProd:               config.GetConfigInstance().IsProd(),
		TimeZone:             time.GetTimeZoneList(),
		Gateways:             detail.ConvertGatewayList(ctx, query.GetMerchantGatewayList(ctx, merchant.Id, unibee.Bool(false))),
		ExchangeRateApiKey:   utility.HideStar(exchangeApiKey),
		OpenApiKey:           utility.HideStar(merchant.ApiKey),
		SendGridKey:          utility.HideStar(emailData),
		EmailSender:          emailSender,
		VatSenseKey:          utility.HideStar(vatGatewayKey),
		SegmentServerSideKey: segment.GetMerchantSegmentServerSideConfig(ctx, merchant.Id),
		SegmentUserPortalKey: segment.GetMerchantSegmentUserPortalConfig(ctx, merchant.Id),
		IsOwner:              isOwner,
		MemberRoles:          memberRoles,
	}, nil
}
