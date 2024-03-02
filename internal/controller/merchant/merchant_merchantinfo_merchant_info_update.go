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
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.MerchantId > 0, "MerchantId invalid")
	}
	if len(req.TimeZone) > 0 {
		utility.Assert(time.CheckTimeZone(req.TimeZone), fmt.Sprintf("Invalid Timezone:%s", req.TimeZone))
	}
	info := query.GetMerchantById(ctx, _interface.BizCtx().Get(ctx).MerchantMember.MerchantId)
	utility.Assert(info != nil, "merchantInfo not found")
	var companyLogo = info.CompanyLogo
	if len(req.CompanyLogo) > 0 {
		utility.Assert(strings.HasPrefix(req.CompanyLogo, "http://") || strings.HasPrefix(req.CompanyLogo, "https://"), "companyLogo Invalid, should has http:// or https:// prefix")
		companyLogo = req.CompanyLogo
	}
	_, err = dao.Merchant.Ctx(ctx).Data(g.Map{
		dao.Merchant.Columns().Email:       req.Email,
		dao.Merchant.Columns().Address:     req.Address,
		dao.Merchant.Columns().CompanyName: req.CompanyName,
		dao.Merchant.Columns().CompanyLogo: companyLogo,
		dao.Merchant.Columns().Phone:       req.Phone,
		dao.Merchant.Columns().TimeZone:    req.TimeZone,
		dao.Merchant.Columns().Host:        req.Host,
		dao.Merchant.Columns().GmtModify:   gtime.Now(),
	}).Where(dao.Merchant.Columns().Id, _interface.BizCtx().Get(ctx).MerchantMember.MerchantId).OmitEmpty().Update()
	if err != nil {
		return nil, err
	}

	return &merchantinfo.MerchantInfoUpdateRes{MerchantInfo: query.GetMerchantById(ctx, _interface.BizCtx().Get(ctx).MerchantMember.MerchantId)}, nil
}
