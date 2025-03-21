package user

import (
	"context"
	"encoding/json"
	"fmt"
	"unibee/api/user/auth"
	"unibee/internal/cmd/i18n"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/email"
	"unibee/internal/logic/user"
	"unibee/internal/logic/vat_gateway"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"

	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
)

const CacheKeyUserRegisterPrefix = "CacheKeyUserRegisterPrefix-"

func (c *ControllerAuth) Register(ctx context.Context, req *auth.RegisterReq) (res *auth.RegisterRes, err error) {
	utility.Assert(len(req.Email) > 0, "Email Needed")
	utility.Assert(utility.IsEmailValid(req.Email), "Invalid Email")

	redisKey := fmt.Sprintf("UserAuth-Regist-Email:%s", req.Email)

	if !utility.TryLock(ctx, redisKey, 10) {
		utility.Assert(false, i18n.LocalizationFormat(ctx, "{#ClickTooFast}"))
	}

	var newOne *entity.UserAccount
	newOne = query.GetUserAccountByEmail(ctx, _interface.GetMerchantId(ctx), req.Email) //Id(ctx, user.Id)
	utility.Assert(newOne == nil, "Email already existed")

	var vatNumber = ""
	if req.VATNumber != nil && len(*req.VATNumber) > 0 {
		utility.Assert(vat_gateway.GetDefaultVatGateway(ctx, _interface.GetMerchantId(ctx)).VatRatesEnabled(), "Vat Gateway Need Setup")
		vatNumberValidate, err := vat_gateway.ValidateVatNumberByDefaultGateway(ctx, _interface.GetMerchantId(ctx), 0, *req.VATNumber, "")
		utility.AssertError(err, "Update VAT number error")
		utility.Assert(vatNumberValidate.Valid, i18n.LocalizationFormat(ctx, "{#VatValidateError}", *req.VATNumber))
		if len(req.CountryCode) > 0 {
			utility.Assert(req.CountryCode == vatNumberValidate.CountryCode, "Your country from vat number is "+vatNumberValidate.CountryCode)
		}
		vatNumber = *req.VATNumber
	}

	userStr, err := json.Marshal(
		&user.NewUserInternalReq{
			Email:       req.Email,
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Password:    req.Password,
			Phone:       req.Phone,
			Address:     req.Address,
			UserName:    req.UserName,
			CountryCode: req.CountryCode,
			Type:        req.Type,
			CompanyName: req.CompanyName,
			VATNumber:   vatNumber,
			City:        req.City,
			ZipCode:     req.ZipCode,
		},
	)
	utility.AssertError(err, "Server Error")
	_, err = g.Redis().Set(ctx, CacheKeyUserRegisterPrefix+req.Email, userStr)
	utility.AssertError(err, "Server Error")
	_, err = g.Redis().Expire(ctx, CacheKeyUserRegisterPrefix+req.Email, 3*60)
	utility.AssertError(err, "Server Error")
	verificationCode := utility.GenerateRandomCode(6)
	fmt.Printf("verification %s", verificationCode)
	_, err = g.Redis().Set(ctx, CacheKeyUserRegisterPrefix+req.Email+"-verify", verificationCode)
	utility.AssertError(err, "Server Error")
	_, err = g.Redis().Expire(ctx, CacheKeyUserRegisterPrefix+req.Email+"-verify", 3*60)
	utility.AssertError(err, "Server Error")

	err = email.SendTemplateEmail(ctx, _interface.GetMerchantId(ctx), req.Email, "", "", email.TemplateUserRegistrationCodeVerify, "", &email.TemplateVariable{
		CodeExpireMinute: "3",
		Code:             verificationCode,
	})
	utility.AssertError(err, "Server Error")

	return &auth.RegisterRes{}, nil
}
