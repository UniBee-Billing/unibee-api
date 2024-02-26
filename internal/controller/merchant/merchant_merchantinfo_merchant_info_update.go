package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	_interface "unibee/internal/interface"
	"unibee/internal/query"
	"unibee/time"
	"unibee/utility"

	"unibee/api/merchant/merchantinfo"
)

func (c *ControllerMerchantinfo) MerchantInfoUpdate(ctx context.Context, req *merchantinfo.MerchantInfoUpdateReq) (res *merchantinfo.MerchantInfoUpdateRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.MerchantId > 0, "MerchantId invalid")
	}
	if len(req.TimeZone) > 0 {
		utility.Assert(time.CheckTimeZone(req.TimeZone), fmt.Sprintf("Invalid Timezone:%s", req.TimeZone))
	}
	info := query.GetMerchantInfoById(ctx, _interface.BizCtx().Get(ctx).MerchantUser.MerchantId)
	utility.Assert(info != nil, "merchantInfo not found")
	var companyLogo = info.CompanyLogo
	if len(req.CompanyLogo) > 0 {
		utility.Assert(strings.HasPrefix(req.CompanyLogo, "http://") || strings.HasPrefix(req.CompanyLogo, "https://"), "companyLogo Invalid, should has http:// or https:// prefix")
		companyLogo = req.CompanyLogo
	}
	_, err = dao.MerchantInfo.Ctx(ctx).Data(g.Map{
		dao.MerchantInfo.Columns().Email:       req.Email,
		dao.MerchantInfo.Columns().Address:     req.Address,
		dao.MerchantInfo.Columns().CompanyName: req.CompanyName,
		dao.MerchantInfo.Columns().CompanyLogo: companyLogo,
		dao.MerchantInfo.Columns().Phone:       req.Phone,
		dao.MerchantInfo.Columns().TimeZone:    req.TimeZone,
		dao.MerchantInfo.Columns().Host:        req.Host,
		dao.MerchantInfo.Columns().GmtModify:   gtime.Now(),
	}).Where(dao.MerchantInfo.Columns().Id, _interface.BizCtx().Get(ctx).MerchantUser.MerchantId).OmitEmpty().Update()
	if err != nil {
		return nil, err
	}

	return &merchantinfo.MerchantInfoUpdateRes{MerchantInfo: query.GetMerchantInfoById(ctx, _interface.BizCtx().Get(ctx).MerchantUser.MerchantId)}, nil
}
