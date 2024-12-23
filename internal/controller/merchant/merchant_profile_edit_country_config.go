package merchant

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/profile"
)

func (c *ControllerProfile) EditCountryConfig(ctx context.Context, req *profile.EditCountryConfigReq) (res *profile.EditCountryConfigRes, err error) {
	utility.Assert(len(req.CountryCode) > 0, "Invalid CountryCode")
	utility.Assert(len(req.Name) > 0, "Invalid Name")

	one := query.GetMerchantCountryConfig(ctx, _interface.GetMerchantId(ctx), req.CountryCode)
	if one == nil {
		one = &entity.MerchantCountryConfig{
			MerchantId:  _interface.GetMerchantId(ctx),
			CountryCode: req.CountryCode,
			Name:        req.Name,
			CreateTime:  gtime.Now().Timestamp(),
		}
		_, err = dao.MerchantCountryConfig.Ctx(ctx).Data(one).OmitNil().Insert(one)
		utility.AssertError(err, "Server Error")
	} else {
		_, err = dao.MerchantCountryConfig.Ctx(ctx).Data(g.Map{
			dao.MerchantCountryConfig.Columns().Name: req.Name,
			dao.Merchant.Columns().GmtModify:         gtime.Now(),
		}).Where(dao.MerchantCountryConfig.Columns().Id, one.Id).OmitEmpty().Update()
		utility.AssertError(err, "Server Error")
	}
	return &profile.EditCountryConfigRes{}, nil
}
