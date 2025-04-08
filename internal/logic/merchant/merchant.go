package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	"unibee/internal/cmd/config"
	redismq2 "unibee/internal/cmd/redismq"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/member"
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

func GetMerchantByOpenApiKeyFromCache(ctx context.Context, openApiKey string) *entity.Merchant {
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
	ReloadAllMerchantsCacheForSDKAuthBackground()
	return apiKey
}

func QueryOrCreateMerchant(ctx context.Context, req *CreateMerchantInternalReq) (*entity.Merchant, *entity.MerchantMember, error) {
	one := query.GetMerchantMemberByEmail(ctx, req.Email)
	if one == nil {
		return CreateMerchant(ctx, req)
	} else {
		merchant := query.GetMerchantById(ctx, one.MerchantId)
		utility.Assert(merchant != nil, "Merchant Not Found")
		if one.Role == "Owner" {
			_, err := dao.Merchant.Ctx(ctx).Data(g.Map{
				dao.Merchant.Columns().Phone:     req.Phone,
				dao.Merchant.Columns().GmtModify: gtime.Now(),
			}).Where(dao.Merchant.Columns().Id, one.MerchantId).OmitEmpty().Update()
			if err != nil {
				g.Log().Errorf(ctx, "QueryOrCreateMerchant UpdateMerchant error:%s", err.Error())
			}
		}
		_, err := dao.MerchantMember.Ctx(ctx).Data(g.Map{
			dao.MerchantMember.Columns().FirstName: req.FirstName,
			dao.MerchantMember.Columns().LastName:  req.LastName,
		}).Where(dao.MerchantMember.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			g.Log().Errorf(ctx, "QueryOrCreateMerchant UpdateOwnerMember error:%s", err.Error())
		}
		merchant = query.GetMerchantById(ctx, one.MerchantId)
		one = query.GetMerchantMemberByEmail(ctx, req.Email)
		return merchant, one, nil
	}
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
	ReloadAllMerchantsCacheForSDKAuthBackground()
	member.ReloadMemberCacheForSdkAuthBackground(merchantMasterMember.Id)
	_, _ = redismq.Send(&redismq.Message{
		Topic:      redismq2.TopicMerchantCreatedWebhook.Topic,
		Tag:        redismq2.TopicMerchantCreatedWebhook.Tag,
		Body:       fmt.Sprintf("%d", merchant.Id),
		CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
	})
	_, _ = redismq.Send(&redismq.Message{
		Topic:      redismq2.TopicMerchantMemberCreatedWebhook.Topic,
		Tag:        redismq2.TopicMerchantMemberCreatedWebhook.Tag,
		Body:       fmt.Sprintf("%d", merchantMasterMember.Id),
		CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
	})
	return merchant, newOne, err
}

func SendMerchantRegisterEmail(ctx context.Context, req *CreateMerchantInternalReq, verificationCode string) {
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
