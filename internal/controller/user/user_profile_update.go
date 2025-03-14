package user

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"unibee/api/bean/detail"
	"unibee/internal/cmd/i18n"
	"unibee/internal/consts"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/user/sub_update"
	"unibee/internal/logic/vat_gateway"
	"unibee/internal/query"
	"unibee/time"
	"unibee/utility"

	"unibee/api/user/profile"
	dao "unibee/internal/dao/default"
)

func (c *ControllerProfile) Update(ctx context.Context, req *profile.UpdateReq) (res *profile.UpdateRes, err error) {
	// timezone check
	if len(req.TimeZone) > 0 {
		utility.Assert(time.CheckTimeZone(req.TimeZone), fmt.Sprintf("Invalid Timezone:%s", req.TimeZone))
	}

	if req.GatewayId != nil && *req.GatewayId > 0 {
		one := query.GetUserAccountById(ctx, _interface.Context().Get(ctx).User.Id)
		utility.Assert(one != nil, "user not found")
		if len(one.GatewayId) > 0 {
			oldGatewayId, err := strconv.ParseUint(one.GatewayId, 10, 64)
			if err == nil {
				gateway := query.GetGatewayById(ctx, oldGatewayId)
				newGateway := query.GetGatewayById(ctx, *req.GatewayId)
				if oldGatewayId != *req.GatewayId {
					utility.Assert(gateway.GatewayType != consts.GatewayTypeWireTransfer, "Can't change gateway from wire transfer to other, Please contact billing admin")
					utility.Assert(newGateway.GatewayType != consts.GatewayTypeWireTransfer, "Can't change gateway to wire transfer, Please contact billing admin")
				}
			}
		} else {
			newGateway := query.GetGatewayById(ctx, *req.GatewayId)
			utility.Assert(newGateway.GatewayType != consts.GatewayTypeWireTransfer, "Can't change gateway to wire transfer, Please contact billing admin")
		}
		var paymentMethodId = ""
		if req.PaymentMethodId != nil {
			paymentMethodId = *req.PaymentMethodId
		}
		var paymentType = ""
		if req.GatewayPaymentType != nil {
			paymentType = *req.GatewayPaymentType
		}
		sub_update.UpdateUserDefaultGatewayPaymentMethod(ctx, _interface.Context().Get(ctx).User.Id, *req.GatewayId, paymentMethodId, paymentType)
	}
	one := query.GetUserAccountById(ctx, _interface.Context().Get(ctx).User.Id)
	var vatNumber = one.VATNumber
	if req.VATNumber != nil {
		if len(*req.VATNumber) > 0 {
			utility.Assert(vat_gateway.GetDefaultVatGateway(ctx, _interface.GetMerchantId(ctx)).VatRatesEnabled(), "Vat Gateway Need Setup")
			vatNumberValidate, err := vat_gateway.ValidateVatNumberByDefaultGateway(ctx, _interface.GetMerchantId(ctx), _interface.Context().Get(ctx).User.Id, *req.VATNumber, "")
			utility.AssertError(err, "Update VAT number error")
			utility.Assert(vatNumberValidate.Valid, i18n.LocalizationFormat(ctx, "{#VatValidateError}", *req.VATNumber))
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
			utility.Assert(gateway.VatRatesEnabled(), "Vat Gateway Need Setup")
			vatNumberValidate, err := vat_gateway.ValidateVatNumberByDefaultGateway(ctx, _interface.GetMerchantId(ctx), _interface.Context().Get(ctx).User.Id, vatNumber, "")
			utility.AssertError(err, "Update VAT number error")
			utility.Assert(vatNumberValidate.Valid, i18n.LocalizationFormat(ctx, "{#VatValidateError}", vatNumber))
			utility.Assert(vatNumberValidate.CountryCode == *req.CountryCode, "Your country from vat number is "+vatNumberValidate.CountryCode)
		}
		if one.CountryCode != *req.CountryCode {
			sub_update.UpdateUserCountryCode(ctx, _interface.Context().Get(ctx).User.Id, *req.CountryCode)
		}
	}

	if req.Type != nil {
		utility.Assert(*req.Type == 1 || *req.Type == 2, "invalid Type, 1-Individual|2-organization")
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
		dao.UserAccount.Columns().TimeZone:        req.TimeZone,
		dao.UserAccount.Columns().City:            req.City,
		dao.UserAccount.Columns().Language:        req.Language,
		dao.UserAccount.Columns().ZipCode:         req.ZipCode,
		//dao.UserAccount.Columns().ReMark:             req.GatewayPaymentType,
		dao.UserAccount.Columns().RegistrationNumber: req.RegistrationNumber,
		dao.UserAccount.Columns().GmtModify:          gtime.Now(),
	}).Where(dao.UserAccount.Columns().Id, _interface.Context().Get(ctx).User.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}
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
	one = query.GetUserAccountById(ctx, _interface.Context().Get(ctx).User.Id)
	return &profile.UpdateRes{User: detail.ConvertUserAccountToDetail(ctx, one)}, nil
}
