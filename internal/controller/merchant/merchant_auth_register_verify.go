package merchant

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee-api/utility"

	"encoding/json"
	"unibee-api/api/merchant/auth"
	dao "unibee-api/internal/dao/oversea_pay"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerAuth) RegisterVerify(ctx context.Context, req *auth.RegisterVerifyReq) (res *auth.RegisterVerifyRes, err error) {
	verificationCode, err := g.Redis().Get(ctx, CacheKeyMerchantRegisterPrefix+req.Email+"-verify")
	utility.AssertError(err, "Server Error")
	utility.Assert(verificationCode != nil, "Invalid Code")
	utility.Assert((verificationCode.String()) == req.VerificationCode, "Invalid Code")
	userStr, err := g.Redis().Get(ctx, CacheKeyMerchantRegisterPrefix+req.Email)
	utility.AssertError(err, "Server Error")
	utility.Assert(userStr != nil, "Invalid Code")
	u := struct {
		FirstName, LastName, Email, Password, Phone, Address, UserName string
		MerchantId                                                     uint64
	}{}
	err = json.Unmarshal([]byte(userStr.String()), &u)

	merchantUser := &entity.MerchantUserAccount{
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Email:      u.Email,
		Password:   u.Password,
		MerchantId: u.MerchantId,
		UserName:   u.UserName,
		CreateTime: gtime.Now().Timestamp(),
	}

	// transaction create Merchant
	err = dao.Refund.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
		insert, err := dao.MerchantUserAccount.Ctx(ctx).Data(merchantUser).OmitNil().Insert(merchantUser)
		if err != nil {
			return err
		}
		id, err := insert.LastInsertId()
		if err != nil {
			return err
		}
		merchantUser.Id = uint64(id)

		//create MerchantInfo
		merchantInfo := &entity.MerchantInfo{
			CompanyId: 0,
			UserId:    id,
			ApiKey:    utility.GenerateRandomAlphanumeric(32), //32 bit open api key
		}

		insert, err = dao.MerchantInfo.Ctx(ctx).Data(merchantInfo).OmitNil().Insert(merchantInfo)
		if err != nil {
			return err
		}
		merchantId, err := insert.LastInsertId()
		if err != nil {
			return err
		}
		// bind merchantUserAccount
		_, err = dao.MerchantUserAccount.Ctx(ctx).Data(g.Map{
			dao.MerchantUserAccount.Columns().MerchantId: merchantId,
			dao.MerchantUserAccount.Columns().GmtModify:  gtime.Now(),
		}).Where(dao.MerchantUserAccount.Columns().Id, id).OmitNil().Update()
		if err != nil {
			return err
		}
		return nil
	})

	utility.AssertError(err, "Server Error")
	var newOne *entity.MerchantUserAccount
	newOne = query.GetMerchantUserAccountById(ctx, merchantUser.Id)
	utility.Assert(newOne != nil, "Server Error")
	newOne.Password = ""
	return &auth.RegisterVerifyRes{MerchantUser: newOne}, nil
}
