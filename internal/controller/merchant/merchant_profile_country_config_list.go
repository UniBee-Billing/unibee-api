package merchant

import (
	"context"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface"

	"unibee/api/merchant/profile"
)

func (c *ControllerProfile) CountryConfigList(ctx context.Context, req *profile.CountryConfigListReq) (res *profile.CountryConfigListRes, err error) {
	var mainList []*bean.MerchantCountryConfig

	err = dao.MerchantCountryConfig.Ctx(ctx).
		Where(dao.MerchantCountryConfig.Columns().MerchantId, _interface.GetMerchantId(ctx)).OmitEmpty().Scan(&mainList)
	if err != nil {
		return nil, err
	}
	return &profile.CountryConfigListRes{Configs: mainList}, nil
}
