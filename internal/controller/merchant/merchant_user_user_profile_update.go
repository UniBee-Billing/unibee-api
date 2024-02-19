package merchant

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee-api/internal/dao/oversea_pay"
	_interface "unibee-api/internal/interface"
	"unibee-api/utility"

	"unibee-api/api/merchant/user"
)

func (c *ControllerUser) UserProfileUpdate(ctx context.Context, req *user.UserProfileUpdateReq) (res *user.UserProfileUpdateRes, err error) {
	utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "Merchant User Not Found")
	utility.Assert(len(_interface.BizCtx().Get(ctx).MerchantUser.Token) > 0, "Merchant Token Not Found")

	_, err = dao.UserAccount.Ctx(ctx).Data(g.Map{
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
		dao.UserAccount.Columns().PaymentMethod:   req.PaymentMethod,
		dao.UserAccount.Columns().CountryCode:     req.CountryCode,
		dao.UserAccount.Columns().CountryName:     req.CountryName,
		dao.UserAccount.Columns().GmtModify:       gtime.Now(),
	}).Where(dao.UserAccount.Columns().Id, req.UserId).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	return &user.UserProfileUpdateRes{}, nil
}
