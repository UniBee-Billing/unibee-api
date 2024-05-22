package merchant

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/merchant/user"
	dao "unibee/internal/dao/oversea_pay"
	_interface "unibee/internal/interface"
	user2 "unibee/internal/logic/user"
	"unibee/internal/logic/vat_gateway"
	"unibee/utility"
)

func (c *ControllerUser) Update(ctx context.Context, req *user.UpdateReq) (res *user.UpdateRes, err error) {
	if req.VATNumber != nil && len(*req.VATNumber) > 0 {
		gateway := vat_gateway.GetDefaultVatGateway(ctx, _interface.GetMerchantId(ctx))
		utility.Assert(gateway != nil, "Default Vat Gateway Need Setup")
		vatNumberValidate, err := vat_gateway.ValidateVatNumberByDefaultGateway(ctx, _interface.GetMerchantId(ctx), _interface.Context().Get(ctx).User.Id, *req.VATNumber, "")
		utility.AssertError(err, "Update VatNumber error")
		utility.Assert(vatNumberValidate.Valid, "vatNumber invalid")
	}
	if req.CountryCode != nil && len(*req.CountryCode) > 0 {
		utility.Assert(vat_gateway.GetDefaultVatGateway(ctx, _interface.GetMerchantId(ctx)) != nil, "Default Vat Gateway Need Setup")
		user2.UpdateUserCountryCode(ctx, _interface.Context().Get(ctx).User.Id, *req.CountryCode)
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
		dao.UserAccount.Columns().GmtModify:       gtime.Now(),
	}).Where(dao.UserAccount.Columns().Id, req.UserId).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	return &user.UpdateRes{}, nil
}
