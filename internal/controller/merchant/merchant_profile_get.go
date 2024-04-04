package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/api/merchant/profile"
	"unibee/internal/cmd/config"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/currency"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/time"
)

func (c *ControllerProfile) Get(ctx context.Context, req *profile.GetReq) (res *profile.GetRes, err error) {
	var member *entity.MerchantMember
	if _interface.Context().Get(ctx).MerchantMember != nil {
		member = query.GetMerchantMemberById(ctx, _interface.Context().Get(ctx).MerchantMember.Id)
	}
	return &profile.GetRes{
		Merchant:       bean.SimplifyMerchant(query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))),
		MerchantMember: detail.ConvertMemberToDetail(ctx, member),
		Currency:       currency.GetMerchantCurrencies(),
		Env:            config.GetConfigInstance().Env,
		IsProd:         config.GetConfigInstance().IsProd(),
		TimeZone:       time.GetTimeZoneList(),
	}, nil
}
