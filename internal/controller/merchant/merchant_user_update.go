package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	"unibee/api/bean/detail"
	"unibee/api/merchant/user"
	"unibee/internal/cmd/i18n"
	redismq2 "unibee/internal/cmd/redismq"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/user/sub_update"
	"unibee/internal/logic/vat_gateway"
	"unibee/internal/query"
	"unibee/utility"
	"unibee/utility/unibee"
)

func (c *ControllerUser) Update(ctx context.Context, req *user.UpdateReq) (res *user.UpdateRes, err error) {
	if req.UserId == nil {
		utility.Assert(req.Email != nil && len(*req.Email) > 0, "either Email or UserId needed")
		one := query.GetUserAccountByEmail(ctx, _interface.GetMerchantId(ctx), *req.Email)
		utility.Assert(one != nil, "user not found")
		req.UserId = unibee.Uint64(one.Id)
	}
	utility.Assert(req.UserId != nil, "either Email or UserId needed")
	if req.ExternalUserId != nil && len(*req.ExternalUserId) > 0 {
		//update externalUserId
		one := query.GetUserAccountByExternalUserId(ctx, _interface.GetMerchantId(ctx), *req.ExternalUserId)
		if one != nil {
			utility.Assert(one.Id == *req.UserId, fmt.Sprintf("ExternalUserId has bean used by another user:%v email:%s", one.Id, one.Email))
		}
		_, err = dao.UserAccount.Ctx(ctx).Data(g.Map{
			dao.UserAccount.Columns().ExternalUserId: req.ExternalUserId,
		}).Where(dao.UserAccount.Columns().Id, req.UserId).OmitNil().Update()
	}

	if req.GatewayId != nil && *req.GatewayId > 0 {
		var paymentMethodId = ""
		if req.PaymentMethodId != nil {
			paymentMethodId = *req.PaymentMethodId
		}
		var paymentType = ""
		if req.GatewayPaymentType != nil {
			paymentType = *req.GatewayPaymentType
		}
		sub_update.UpdateUserDefaultGatewayPaymentMethod(ctx, *req.UserId, *req.GatewayId, paymentMethodId, paymentType)
	}
	one := query.GetUserAccountById(ctx, *req.UserId)
	var vatNumber = one.VATNumber
	if req.VATNumber != nil {
		if len(*req.VATNumber) > 0 {
			utility.Assert(vat_gateway.GetDefaultVatGateway(ctx, _interface.GetMerchantId(ctx)).VatRatesEnabled(), "Vat Gateway Need Setup")
			vatNumberValidate, err := vat_gateway.ValidateVatNumberByDefaultGateway(ctx, _interface.GetMerchantId(ctx), *req.UserId, *req.VATNumber, "")
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
			vatNumberValidate, err := vat_gateway.ValidateVatNumberByDefaultGateway(ctx, _interface.GetMerchantId(ctx), *req.UserId, vatNumber, "")
			utility.AssertError(err, "Update VAT number error")
			utility.Assert(vatNumberValidate.Valid, i18n.LocalizationFormat(ctx, "{#VatValidateError}", vatNumber))
			utility.Assert(vatNumberValidate.CountryCode == *req.CountryCode, "Your country from vat number is "+vatNumberValidate.CountryCode)
		}
		if one.CountryCode != *req.CountryCode {
			sub_update.UpdateUserCountryCode(ctx, *req.UserId, *req.CountryCode)
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
		dao.UserAccount.Columns().City:            req.City,
		dao.UserAccount.Columns().ZipCode:         req.ZipCode,
		dao.UserAccount.Columns().Language:        req.Language,
		//dao.UserAccount.Columns().ReMark:             req.GatewayPaymentType,
		dao.UserAccount.Columns().RegistrationNumber: req.RegistrationNumber,
		dao.UserAccount.Columns().GmtModify:          gtime.Now(),
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
	if req.Metadata != nil {
		var metadata = make(map[string]interface{})
		if len(one.MetaData) > 0 {
			_ = gjson.Unmarshal([]byte(one.MetaData), &metadata)
		}
		for k, v := range *req.Metadata {
			metadata[k] = v
		}
		_, _ = dao.UserAccount.Ctx(ctx).Data(g.Map{
			dao.UserAccount.Columns().MetaData: utility.MarshalToJsonString(metadata),
		}).Where(dao.UserAccount.Columns().Id, req.UserId).OmitNil().Update()
	}
	one = query.GetUserAccountById(ctx, *req.UserId)
	_, _ = redismq.Send(&redismq.Message{
		Topic:      redismq2.TopicUserAccountUpdate.Topic,
		Tag:        redismq2.TopicUserAccountUpdate.Tag,
		Body:       fmt.Sprintf("%d", one.Id),
		CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
	})

	return &user.UpdateRes{User: detail.ConvertUserAccountToDetail(ctx, one)}, nil
}
