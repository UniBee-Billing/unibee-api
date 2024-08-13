package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/cmd/config"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/platform"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type CreateMerchantInternalReq struct {
	FirstName, LastName, Email, Password, Phone, UserName string
}

func GetOpenApiKeyRedisKey(token string) string {
	return fmt.Sprintf("openApiKey#%s#%s", config.GetConfigInstance().Env, token)
}

func GetMerchantFromCache(ctx context.Context, openApiKey string) *entity.Merchant {
	get, err := g.Redis().Get(ctx, GetOpenApiKeyRedisKey(openApiKey))
	if err != nil {
		return nil
	}
	if get != nil && len(get.String()) > 0 {
		one := query.GetMerchantById(ctx, get.Uint64())
		return one
	}
	return nil
}

func PutOpenApiKeyToCache(ctx context.Context, openApiKey string, merchantId uint64) bool {
	err := g.Redis().SetEX(ctx, GetOpenApiKeyRedisKey(openApiKey), merchantId, 24*3600)
	if err != nil {
		return false
	}
	return true
}

func NewOpenApiKey(ctx context.Context, merchantId uint64) string {
	one := query.GetMerchantById(ctx, merchantId)
	utility.Assert(one != nil, "Merchant Not Found")
	oldApikey := one.ApiKey
	utility.Assert(PutOpenApiKeyToCache(ctx, oldApikey, merchantId), "Server Error")
	apiKey := utility.GenerateRandomAlphanumeric(32)
	_, err := dao.Merchant.Ctx(ctx).Data(g.Map{
		dao.Merchant.Columns().ApiKey:    apiKey,
		dao.Merchant.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Merchant.Columns().Id, merchantId).Update()
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     merchantId,
		Target:         fmt.Sprintf("ApiKey(%v)", one.ApiKey),
		Content:        fmt.Sprintf("NewApiKey(%v)", apiKey),
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	utility.AssertError(err, "Server Error")
	return apiKey
}

func CreateMerchant(ctx context.Context, req *CreateMerchantInternalReq) (*entity.Merchant, *entity.MerchantMember, error) {
	merchantMasterMember := &entity.MerchantMember{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		Password:   utility.PasswordEncrypt(req.Password),
		UserName:   req.UserName,
		Mobile:     req.Phone,
		Role:       "Owner",
		CreateTime: gtime.Now().Timestamp(),
	}
	merchant := &entity.Merchant{
		Phone:  req.Phone,
		Email:  req.Email,
		ApiKey: utility.GenerateRandomAlphanumeric(32), //32 bit open api key
	}
	// transaction create Merchant
	err := dao.Merchant.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
		insert, err := dao.Merchant.Ctx(ctx).Data(merchant).OmitNil().Insert(merchant)
		if err != nil {
			return err
		}
		merchantId, err := insert.LastInsertId()
		if err != nil {
			return err
		}
		merchant.Id = uint64(merchantId)
		merchantMasterMember.MerchantId = merchant.Id

		insert, err = dao.MerchantMember.Ctx(ctx).Data(merchantMasterMember).OmitNil().Insert(merchantMasterMember)
		if err != nil {
			return err
		}
		id, err := insert.LastInsertId()
		if err != nil {
			return err
		}
		merchantMasterMember.Id = uint64(id)
		merchant.UserId = id

		// bind merchantMemberAccount
		_, err = dao.Merchant.Ctx(ctx).Data(g.Map{
			dao.Merchant.Columns().UserId:    merchant.UserId,
			dao.Merchant.Columns().GmtModify: gtime.Now(),
		}).Where(dao.Merchant.Columns().Id, merchant.Id).Update()
		if err != nil {
			return err
		}
		return nil
	})
	utility.AssertError(err, "Server Error")
	var newOne *entity.MerchantMember
	newOne = query.GetMerchantMemberById(ctx, merchantMasterMember.Id)
	utility.Assert(newOne != nil, "Server Error")
	err = SetupForCloudMode(ctx, merchant.Id)
	return merchant, newOne, err
}

func SendMerchantRegisterEmail(ctx context.Context, req *CreateMerchantInternalReq, verificationCode string) {
	//if config.GetConfigInstance().Mode == "cloud" {
	//	err := email.SendTemplateEmail(ctx, consts.CloudModeManagerMerchantId, req.Email, "", email.TemplateMerchantRegistrationCodeVerify, "", &email.TemplateVariable{
	//		CodeExpireMinute: "3",
	//		Code:             verificationCode,
	//	})
	//	utility.AssertError(err, "Server Error")
	//} else {
	//	utility.Assert(true, "not support")
	//}
	utility.Assert(req != nil, "Server Error,nil")
	publicIp := utility.GetPublicIP()
	if len(publicIp) == 0 {
		publicIp = "0.0.0.0"
	}
	err := platform.SentPlatformMerchantRegisterEmail(map[string]string{
		"ownerEmail": req.Email,
		"ip":         publicIp,
		"firstName":  req.FirstName,
		"lastName":   req.LastName,
		"Phone":      req.Phone,
		"userName":   req.UserName,
		"code":       verificationCode,
	})
	utility.AssertError(err, "Server Error")
}

func HardDeleteMerchant(ctx context.Context, merchantId uint64) error {
	_, err := dao.MerchantMember.Ctx(ctx).Where(dao.MerchantMember.Columns().MerchantId, merchantId).Delete()
	if err != nil {
		return err
	}
	_, err = dao.Merchant.Ctx(ctx).Where(dao.Merchant.Columns().Id, merchantId).Delete()
	if err != nil {
		return err
	}
	return nil
}
