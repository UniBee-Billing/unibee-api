package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"unibee/api/bean"
	dao "unibee/internal/dao/oversea_pay"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/member"
	"unibee/internal/query"
	"unibee/time"
	"unibee/utility"

	"unibee/api/merchant/profile"
)

func (c *ControllerProfile) Update(ctx context.Context, req *profile.UpdateReq) (res *profile.UpdateRes, err error) {
	if len(req.TimeZone) > 0 {
		utility.Assert(time.CheckTimeZone(req.TimeZone), fmt.Sprintf("Invalid Timezone:%s", req.TimeZone))
	}
	merchant := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	utility.Assert(merchant != nil, "merchantInfo not found")
	var companyLogo = merchant.CompanyLogo
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
	}).Where(dao.Merchant.Columns().Id, _interface.Context().Get(ctx).MerchantMember.MerchantId).OmitEmpty().Update()
	if err != nil {
		return nil, err
	}
	member.AppendOptLog(ctx, &member.OptLogRequest{
		Target:         fmt.Sprintf("Profile"),
		Content:        "Update",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)

	return &profile.UpdateRes{Merchant: bean.SimplifyMerchant(query.GetMerchantById(ctx, _interface.Context().Get(ctx).MerchantMember.MerchantId))}, nil
}
