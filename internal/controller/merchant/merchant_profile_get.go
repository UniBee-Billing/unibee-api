package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/api/merchant/profile"
	"unibee/internal/cmd/config"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/analysis/segment"
	"unibee/internal/logic/currency"
	"unibee/internal/logic/email"
	"unibee/internal/logic/vat_gateway"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/time"
	"unibee/utility"
)

func (c *ControllerProfile) Get(ctx context.Context, req *profile.GetReq) (res *profile.GetRes, err error) {
	var member *entity.MerchantMember
	var isOwner = false
	var memberRoles = make([]*bean.MerchantRoleSimplify, 0)
	if _interface.Context().Get(ctx) != nil && _interface.Context().Get(ctx).MerchantMember != nil {
		member = query.GetMerchantMemberById(ctx, _interface.Context().Get(ctx).MerchantMember.Id)
		if member != nil {
			isOwner, memberRoles = detail.ConvertMemberRole(ctx, member)
		}
	}
	merchant := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	utility.Assert(merchant != nil, "merchant not found")
	_, vatData := vat_gateway.GetDefaultMerchantVatConfig(ctx, merchant.Id)
	_, emailData := email.GetDefaultMerchantEmailConfig(ctx, merchant.Id)
	return &profile.GetRes{
		Merchant:             bean.SimplifyMerchant(merchant),
		MerchantMember:       detail.ConvertMemberToDetail(ctx, member),
		Currency:             currency.GetMerchantCurrencies(),
		Env:                  config.GetConfigInstance().Env,
		IsProd:               config.GetConfigInstance().IsProd(),
		TimeZone:             time.GetTimeZoneList(),
		Gateways:             bean.SimplifyGatewayList(query.GetMerchantGatewayList(ctx, merchant.Id)),
		OpenApiKey:           utility.HideStar(merchant.ApiKey),
		SendGridKey:          utility.HideStar(emailData),
		VatSenseKey:          utility.HideStar(vatData),
		SegmentServerSideKey: segment.GetMerchantSegmentServerSideConfig(ctx, merchant.Id),
		SegmentUserPortalKey: segment.GetMerchantSegmentUserPortalConfig(ctx, merchant.Id),
		IsOwner:              isOwner,
		MemberRoles:          memberRoles,
	}, nil
}
