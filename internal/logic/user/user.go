package user

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	"strings"
	"unibee/internal/cmd/i18n"
	redismq2 "unibee/internal/cmd/redismq"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/analysis/segment"
	"unibee/internal/logic/jwt"
	metric2 "unibee/internal/logic/metric"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/subscription/service"
	"unibee/internal/logic/vat_gateway"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func ChangeUserEmail(ctx context.Context, userId uint64, newEmail string) {
	newEmail = strings.TrimSpace(newEmail)
	one := query.GetUserAccountById(ctx, userId)
	utility.Assert(one != nil, "user not found")
	if one.Email == newEmail {
		return
	}
	exist := query.GetUserAccountByEmail(ctx, one.MerchantId, newEmail)
	utility.Assert(exist == nil, "email has been used by another user")
	_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
		dao.UserAccount.Columns().Email:     newEmail,
		dao.UserAccount.Columns().GmtModify: gtime.Now(),
	}).Where(dao.UserAccount.Columns().Id, one.Id).OmitNil().Update()
	utility.AssertError(err, "update user email err")
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("User(%v)", one.Id),
		Content:        fmt.Sprintf("ChangeToNewEmail(%s)", newEmail),
		UserId:         one.Id,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	if len(one.SubscriptionId) > 0 {
		// change user's sub gateway immediately
		_, _ = redismq.Send(&redismq.Message{
			Topic: redismq2.TopicUserMetricUpdate.Topic,
			Tag:   redismq2.TopicUserMetricUpdate.Tag,
			Body: utility.MarshalToJsonString(&metric2.UserMetricUpdateMessage{
				UserId:         one.Id,
				SubscriptionId: one.SubscriptionId,
				Description:    "ChangeUserEmail",
			}),
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
		_, _ = redismq.Send(&redismq.Message{
			Topic:      redismq2.TopicUserAccountUpdate.Topic,
			Tag:        redismq2.TopicUserAccountUpdate.Tag,
			Body:       fmt.Sprintf("%d", one.Id),
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
	}
}

func ChangeUserPasswordWithOutOldVerify(ctx context.Context, merchantId uint64, email string, newPassword string) {
	email = strings.TrimSpace(email)
	one := query.GetUserAccountByEmail(ctx, merchantId, email)
	utility.Assert(one != nil, "user not found")
	_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
		dao.UserAccount.Columns().Password:  utility.PasswordEncrypt(newPassword),
		dao.UserAccount.Columns().GmtModify: gtime.Now(),
	}).Where(dao.UserAccount.Columns().Id, one.Id).OmitNil().Update()
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("User(%v)", one.Id),
		Content:        "ChangePassword",
		UserId:         one.Id,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	utility.AssertError(err, "server error")
}

func ChangeUserPassword(ctx context.Context, merchantId uint64, email string, oldPassword string, newPassword string) {
	email = strings.TrimSpace(email)
	one := query.GetUserAccountByEmail(ctx, merchantId, email)
	utility.Assert(one != nil, "user not found")
	utility.Assert(utility.ComparePasswords(one.Password, oldPassword), "password not match")
	_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
		dao.UserAccount.Columns().Password:  utility.PasswordEncrypt(newPassword),
		dao.UserAccount.Columns().GmtModify: gtime.Now(),
	}).Where(dao.UserAccount.Columns().Id, one.Id).OmitNil().Update()
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("User(%v)", one.Id),
		Content:        "ChangePassword",
		UserId:         one.Id,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	utility.AssertError(err, "server error")
}

func FrozenUser(ctx context.Context, userId int64) {
	one := query.GetUserAccountById(ctx, uint64(userId))
	utility.Assert(one != nil, "user not found")
	utility.Assert(one.Status != 2, "user already suspend")
	_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
		dao.UserAccount.Columns().Status:    2,
		dao.UserAccount.Columns().GmtModify: gtime.Now(),
	}).Where(dao.UserAccount.Columns().Id, one.Id).OmitNil().Update()
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("User(%v)", one.Id),
		Content:        "Suspend",
		UserId:         one.Id,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	utility.AssertError(err, "server error")
	subs := query.GetLatestActiveOrIncompleteOrCreateSubscriptionsByUserId(ctx, one.Id, one.MerchantId)
	for _, sub := range subs {
		if sub != nil {
			err = service.SubscriptionCancel(ctx, sub.SubscriptionId, false, false, "CancelledByAdminSuspendUser")
			utility.AssertError(err, "server error")
		}
	}
	if len(one.SubscriptionId) > 0 {
		// change user's sub gateway immediately
		_, _ = redismq.Send(&redismq.Message{
			Topic: redismq2.TopicUserMetricUpdate.Topic,
			Tag:   redismq2.TopicUserMetricUpdate.Tag,
			Body: utility.MarshalToJsonString(&metric2.UserMetricUpdateMessage{
				UserId:         one.Id,
				SubscriptionId: one.SubscriptionId,
				Description:    "SuspendUser",
			}),
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
		_, _ = redismq.Send(&redismq.Message{
			Topic:      redismq2.TopicUserAccountUpdate.Topic,
			Tag:        redismq2.TopicUserAccountUpdate.Tag,
			Body:       fmt.Sprintf("%d", one.Id),
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
	}
}

func ReleaseUser(ctx context.Context, userId int64) {
	one := query.GetUserAccountById(ctx, uint64(userId))
	utility.Assert(one != nil, "user not found")
	_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
		dao.UserAccount.Columns().Status:    0,
		dao.UserAccount.Columns().GmtModify: gtime.Now(),
	}).Where(dao.UserAccount.Columns().Id, one.Id).OmitNil().Update()
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("User(%v)", one.Id),
		Content:        "Resume",
		UserId:         one.Id,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	utility.AssertError(err, "server error")
	if len(one.SubscriptionId) > 0 {
		// change user's sub gateway immediately
		_, _ = redismq.Send(&redismq.Message{
			Topic: redismq2.TopicUserMetricUpdate.Topic,
			Tag:   redismq2.TopicUserMetricUpdate.Tag,
			Body: utility.MarshalToJsonString(&metric2.UserMetricUpdateMessage{
				UserId:         one.Id,
				SubscriptionId: one.SubscriptionId,
				Description:    "Resume",
			}),
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
		_, _ = redismq.Send(&redismq.Message{
			Topic:      redismq2.TopicUserAccountUpdate.Topic,
			Tag:        redismq2.TopicUserAccountUpdate.Tag,
			Body:       fmt.Sprintf("%d", one.Id),
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
	}
}

func PasswordLogin(ctx context.Context, merchantId uint64, email string, password string) (one *entity.UserAccount, token string) {
	email = strings.TrimSpace(email)
	one = query.GetUserAccountByEmail(ctx, merchantId, email)
	utility.Assert(one != nil, "Email Not Found")
	utility.Assert(one.Status == 0, "Account Status Abnormal")
	utility.Assert(utility.ComparePasswords(one.Password, password), "Login Failed, Password Not Match")

	token, err := jwt.CreatePortalToken(jwt.TOKENTYPEUSER, one.MerchantId, one.Id, one.Email, one.Language)
	fmt.Println("logged-in, save email/id in token: ", one.Email, "/", one.Id)
	utility.AssertError(err, "Server Error")
	utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("User#%d", one.Id)), "Cache Error")
	return one, token
}

func CreateUser(ctx context.Context, req *NewUserInternalReq) (one *entity.UserAccount, err error) {
	utility.Assert(req.MerchantId > 0, "merchantId invalid")
	utility.Assert(req != nil, "Server Error")
	req.Email = strings.TrimSpace(req.Email)
	if len(req.ExternalUserId) > 0 {
		one = query.GetUserAccountByExternalUserId(ctx, req.MerchantId, req.ExternalUserId)
	}
	if one == nil {
		one = query.GetUserAccountByEmail(ctx, req.MerchantId, req.Email)
	}
	utility.Assert(one == nil, "email or externalUserId exist")
	emailOne := query.GetUserAccountByEmail(ctx, req.MerchantId, req.Email)
	utility.Assert(emailOne == nil, "email exist")

	var taxPercentage int64 = 0
	var countryName = ""
	if len(req.CountryCode) > 0 && vat_gateway.GetDefaultVatGateway(ctx, req.MerchantId).VatRatesEnabled() {
		//utility.Assert(vat_gateway.GetDefaultVatGateway(ctx, req.MerchantId).VatRatesEnabled(), "vat gateway need setup while countryCode is not blank")
		if len(req.VATNumber) > 0 {
			vatNumberValidate, err := vat_gateway.ValidateVatNumberByDefaultGateway(ctx, _interface.GetMerchantId(ctx), 0, req.VATNumber, "")
			utility.AssertError(err, "Validate vatNumber error")
			utility.Assert(vatNumberValidate.Valid, i18n.LocalizationFormat(ctx, "{#VatValidateError}", req.VATNumber))
			utility.Assert(req.CountryCode == vatNumberValidate.CountryCode, i18n.LocalizationFormat(ctx, "{#CountryCodeVatNumberNotMatch}", vatNumberValidate.CountryCode))
		}
		taxPercentage, countryName = vat_gateway.ComputeMerchantVatPercentage(ctx, req.MerchantId, req.CountryCode, 0, req.VATNumber)
	} else if len(req.VATNumber) > 0 {
		utility.Assert(false, "countryCode is blank while vatNumber provided")
	}

	one = &entity.UserAccount{
		FirstName:          req.FirstName,
		LastName:           req.LastName,
		Password:           utility.PasswordEncrypt(req.Password),
		Email:              req.Email,
		Phone:              req.Phone,
		Address:            req.Address,
		ExternalUserId:     req.ExternalUserId,
		CountryCode:        req.CountryCode,
		CountryName:        countryName,
		UserName:           req.UserName,
		MerchantId:         req.MerchantId,
		Type:               req.Type,
		CompanyName:        req.CompanyName,
		VATNumber:          req.VATNumber,
		City:               req.City,
		ZipCode:            req.ZipCode,
		Custom:             req.Custom,
		TaxPercentage:      taxPercentage,
		Language:           req.Language,
		RegistrationNumber: req.RegistrationNumber,
		CreateTime:         gtime.Now().Timestamp(),
	}
	result, err := dao.UserAccount.Ctx(ctx).Data(one).OmitNil().Insert(one)
	utility.AssertError(err, "Server Error")
	id, err := result.LastInsertId()
	utility.AssertError(err, "Server Error")
	one.Id = uint64(id)
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("User(%v)", one.Id),
		Content:        "New",
		UserId:         one.Id,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	_, _ = redismq.Send(&redismq.Message{
		Topic:      redismq2.TopicUserAccountCreate.Topic,
		Tag:        redismq2.TopicUserAccountCreate.Tag,
		Body:       fmt.Sprintf("%d", one.Id),
		CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
	})
	segment.RegisterSegmentUserBackground(ctx, one.MerchantId, one)
	return one, nil
}

func QueryOrCreateUser(ctx context.Context, req *NewUserInternalReq) (one *entity.UserAccount, err error) {
	utility.Assert(req.MerchantId > 0, "merchantId invalid")
	utility.Assert(req != nil, "Server Error")
	req.Email = strings.TrimSpace(req.Email)
	if len(req.ExternalUserId) > 0 {
		one = query.GetUserAccountByExternalUserId(ctx, req.MerchantId, req.ExternalUserId)
	}
	if one == nil {
		one = query.GetUserAccountByEmail(ctx, req.MerchantId, req.Email)
	}
	if one == nil {
		// check email not exist
		one, err = CreateUser(ctx, req)
		utility.AssertError(err, "Server Error")
	} else {
		if strings.Compare(one.Email, req.Email) != 0 {
			//email changed, update email
			emailOne := query.GetUserAccountByEmail(ctx, req.MerchantId, req.Email)
			utility.Assert(emailOne == nil || emailOne.Id == one.Id, "email of other externalUserId exist")
			_, err = dao.UserAccount.Ctx(ctx).Data(g.Map{
				dao.UserAccount.Columns().Email:     req.Email,
				dao.UserAccount.Columns().GmtModify: gtime.Now(),
			}).Where(dao.UserAccount.Columns().Id, one.Id).OmitEmpty().Update()
			operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
				MerchantId:     one.MerchantId,
				Target:         fmt.Sprintf("User(%v)", one.Id),
				Content:        fmt.Sprintf("ChangeToNewEmail(%s)", req.Email),
				UserId:         one.Id,
				SubscriptionId: "",
				InvoiceId:      "",
				PlanId:         0,
				DiscountCode:   "",
			}, err)
			utility.AssertError(err, "Server Error")
		}
		if strings.Compare(one.ExternalUserId, req.ExternalUserId) != 0 {
			//externalUserId not match, update externalUserId
			otherOne := query.GetUserAccountByExternalUserId(ctx, req.MerchantId, req.ExternalUserId)
			utility.Assert(otherOne == nil || otherOne.Id == one.Id, "externalUserId of other email exist")
			_, err = dao.UserAccount.Ctx(ctx).Data(g.Map{
				dao.UserAccount.Columns().ExternalUserId: req.ExternalUserId,
				dao.UserAccount.Columns().GmtModify:      gtime.Now(),
			}).Where(dao.UserAccount.Columns().Id, one.Id).OmitEmpty().Update()
			operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
				MerchantId:     one.MerchantId,
				Target:         fmt.Sprintf("User(%v)", one.Id),
				Content:        fmt.Sprintf("UpdateExternalUserId(%s)", req.ExternalUserId),
				UserId:         one.Id,
				SubscriptionId: "",
				InvoiceId:      "",
				PlanId:         0,
				DiscountCode:   "",
			}, err)
			utility.AssertError(err, "Server Error")
		}
		utility.Assert(one.Status == 0, "account status abnormal")
		_, err = dao.UserAccount.Ctx(ctx).Data(g.Map{
			dao.UserAccount.Columns().Address:            req.Address,
			dao.UserAccount.Columns().Phone:              req.Phone,
			dao.UserAccount.Columns().FirstName:          req.FirstName,
			dao.UserAccount.Columns().LastName:           req.LastName,
			dao.UserAccount.Columns().City:               req.City,
			dao.UserAccount.Columns().Type:               req.Type,
			dao.UserAccount.Columns().ZipCode:            req.ZipCode,
			dao.UserAccount.Columns().Language:           req.Language,
			dao.UserAccount.Columns().Address:            req.Address,
			dao.UserAccount.Columns().CompanyName:        req.CompanyName,
			dao.UserAccount.Columns().RegistrationNumber: req.RegistrationNumber,
			dao.UserAccount.Columns().GmtModify:          gtime.Now(),
		}).Where(dao.UserAccount.Columns().Id, one.Id).OmitEmpty().Update()
		utility.AssertError(err, "Server Error")
		one = query.GetUserAccountById(ctx, one.Id)
		_, _ = redismq.Send(&redismq.Message{
			Topic:      redismq2.TopicUserAccountUpdate.Topic,
			Tag:        redismq2.TopicUserAccountUpdate.Tag,
			Body:       fmt.Sprintf("%d", one.Id),
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
	}
	return
}

func HardDeleteUser(ctx context.Context, userId uint64) error {
	_, err := dao.UserAccount.Ctx(ctx).Where(dao.UserAccount.Columns().Id, userId).Delete()
	if err != nil {
		return err
	}
	return nil
}
