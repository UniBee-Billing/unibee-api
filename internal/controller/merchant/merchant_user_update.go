package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean/detail"
	"unibee/api/merchant/user"
	dao "unibee/internal/dao/oversea_pay"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/operation_log"
	user2 "unibee/internal/logic/user"
	"unibee/internal/logic/vat_gateway"
	"unibee/internal/query"
	"unibee/utility"
	"unibee/utility/unibee"
)

func (c *ControllerUser) Update(ctx context.Context, req *user.UpdateReq) (res *user.UpdateRes, err error) {
	if req.GatewayId != nil && *req.GatewayId > 0 {
		var paymentMethodId = ""
		if req.PaymentMethodId != nil {
			paymentMethodId = *req.PaymentMethodId
		}
		user2.UpdateUserDefaultGatewayPaymentMethod(ctx, req.UserId, *req.GatewayId, paymentMethodId)
	}
	one := query.GetUserAccountById(ctx, req.UserId)
	var vatNumber = one.VATNumber
	if req.VATNumber != nil {
		if len(*req.VATNumber) > 0 {
			utility.Assert(vat_gateway.GetDefaultVatGateway(ctx, _interface.GetMerchantId(ctx)) != nil, "Default Vat Gateway Need Setup")
			vatNumberValidate, err := vat_gateway.ValidateVatNumberByDefaultGateway(ctx, _interface.GetMerchantId(ctx), req.UserId, *req.VATNumber, "")
			utility.AssertError(err, "Update VAT number error")
			utility.Assert(vatNumberValidate.Valid, "VAT number invalid")
			if req.CountryCode != nil {
				utility.Assert(*req.CountryCode == vatNumberValidate.CountryCode, "Your country from vat number is "+vatNumberValidate.CountryCode)
			} else {
				utility.Assert(one.CountryCode == vatNumberValidate.CountryCode, "Your country from vat number is "+vatNumberValidate.CountryCode)
			}
		}
		vatNumber = *req.VATNumber
	}

	if req.CountryCode != nil && len(*req.CountryCode) > 0 {
		if len(vatNumber) > 0 {
			gateway := vat_gateway.GetDefaultVatGateway(ctx, _interface.GetMerchantId(ctx))
			utility.Assert(gateway != nil, "Default Vat Gateway Need Setup")
			vatNumberValidate, err := vat_gateway.ValidateVatNumberByDefaultGateway(ctx, _interface.GetMerchantId(ctx), req.UserId, vatNumber, "")
			utility.AssertError(err, "Update VAT number error")
			utility.Assert(vatNumberValidate.Valid, "VAT number invalid")
			utility.Assert(vatNumberValidate.CountryCode == *req.CountryCode, "Your country from vat number is "+vatNumberValidate.CountryCode)
		}
		if one.CountryCode != *req.CountryCode {
			utility.Assert(vat_gateway.GetDefaultVatGateway(ctx, _interface.GetMerchantId(ctx)) != nil, "Default Vat Gateway Need Setup")
			user2.UpdateUserCountryCode(ctx, _interface.Context().Get(ctx).User.Id, *req.CountryCode)
		}
	}

	if req.Type != nil {
		utility.Assert(*req.Type == 1 || *req.Type == 2, "invalid Type, 1-Individual|2-organization")
	} else {
		req.Type = unibee.Int64(1)
	}
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
		dao.UserAccount.Columns().City:            req.City,
		dao.UserAccount.Columns().ZipCode:         req.ZipCode,
		dao.UserAccount.Columns().GmtModify:       gtime.Now(),
	}).Where(dao.UserAccount.Columns().Id, req.UserId).OmitNil().Update()
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("User(%v)", one.Id),
		Content:        "Update",
		UserId:         one.Id,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	if err != nil {
		return nil, err
	}
	one = query.GetUserAccountById(ctx, req.UserId)
	return &user.UpdateRes{User: detail.ConvertUserAccountToDetail(ctx, one)}, nil
}
