package user

import (
	"context"
	"fmt"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/user"
	"unibee/internal/logic/vat_gateway"
	"unibee/time"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"unibee/api/user/profile"
	dao "unibee/internal/dao/oversea_pay"
)

func (c *ControllerProfile) Update(ctx context.Context, req *profile.UpdateReq) (res *profile.UpdateRes, err error) {
	// timezone check
	if len(req.TimeZone) > 0 {
		utility.Assert(time.CheckTimeZone(req.TimeZone), fmt.Sprintf("Invalid Timezone:%s", req.TimeZone))
	}
	if req.VATNumber != nil && len(*req.VATNumber) > 0 {
		utility.Assert(vat_gateway.GetDefaultVatGateway(ctx, _interface.GetMerchantId(ctx)) != nil, "Default Vat Gateway Need Setup")
		vatNumberValidate, err := vat_gateway.ValidateVatNumberByDefaultGateway(ctx, _interface.GetMerchantId(ctx), _interface.Context().Get(ctx).User.Id, *req.VATNumber, "")
		utility.AssertError(err, "Update VatNumber error")
		utility.Assert(vatNumberValidate.Valid, "vatNumber invalid")
	}
	if req.CountryCode != nil && len(*req.CountryCode) > 0 {
		utility.Assert(vat_gateway.GetDefaultVatGateway(ctx, _interface.GetMerchantId(ctx)) != nil, "Default Vat Gateway Need Setup")
		user.UpdateUserCountryCode(ctx, _interface.Context().Get(ctx).User.Id, *req.CountryCode)
	}

	utility.Assert(req.Type == 1 || req.Type == 2, "invalid Type, 1-Individual|2-organization")
	_, err = dao.UserAccount.Ctx(ctx).Data(g.Map{
		dao.UserAccount.Columns().Type:            req.Type,
		dao.UserAccount.Columns().LastName:        req.LastName,
		dao.UserAccount.Columns().FirstName:       req.FirstName,
		dao.UserAccount.Columns().Address:         req.Address,
		dao.UserAccount.Columns().CompanyName:     req.CompanyName,
		dao.UserAccount.Columns().VATNumber:       req.VATNumber,
		dao.UserAccount.Columns().Phone:           req.Phone,
		dao.UserAccount.Columns().Telegram:        req.Telegram,
		dao.UserAccount.Columns().WhatsAPP:        req.WhatsApp,
		dao.UserAccount.Columns().WeChat:          req.WeChat,
		dao.UserAccount.Columns().LinkedIn:        req.LinkedIn,
		dao.UserAccount.Columns().Facebook:        req.Facebook,
		dao.UserAccount.Columns().TikTok:          req.TikTok,
		dao.UserAccount.Columns().OtherSocialInfo: req.OtherSocialInfo,
		dao.UserAccount.Columns().TimeZone:        req.TimeZone,
		dao.UserAccount.Columns().GmtModify:       gtime.Now(),
	}).Where(dao.UserAccount.Columns().Id, _interface.Context().Get(ctx).User.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}

	return &profile.UpdateRes{}, nil
}
