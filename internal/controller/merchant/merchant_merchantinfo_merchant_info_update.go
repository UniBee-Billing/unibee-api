package merchant

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/merchant/merchantinfo"
)

func (c *ControllerMerchantinfo) MerchantInfoUpdate(ctx context.Context, req *merchantinfo.MerchantInfoUpdateReq) (res *merchantinfo.MerchantInfoUpdateRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.MerchantId > 0, "MerchantId invalid")
	}
	_, err = dao.MerchantInfo.Ctx(ctx).Data(g.Map{
		dao.MerchantInfo.Columns().Email:       req.Email,
		dao.MerchantInfo.Columns().LastName:    req.LastName,
		dao.MerchantInfo.Columns().FirstName:   req.FirstName,
		dao.MerchantInfo.Columns().Address:     req.Address,
		dao.MerchantInfo.Columns().CompanyName: req.CompanyName,
		dao.MerchantInfo.Columns().CompanyLogo: req.CompanyLogo,
		dao.MerchantInfo.Columns().Phone:       req.Phone,
		dao.MerchantInfo.Columns().GmtModify:   gtime.Now(),
	}).Where(dao.MerchantInfo.Columns().Id, _interface.BizCtx().Get(ctx).MerchantUser.MerchantId).OmitNil().Update()
	if err != nil {
		return nil, err
	}

	return &merchantinfo.MerchantInfoUpdateRes{MerchantInfo: query.GetMerchantInfoById(ctx, int64(_interface.BizCtx().Get(ctx).MerchantUser.MerchantId))}, nil
}