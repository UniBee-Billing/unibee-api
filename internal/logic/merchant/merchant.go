package merchant

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/merchant/cloud"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type CreateMerchantInternalReq struct {
	FirstName, LastName, Email, Password, Phone, UserName string
}

func CreateMerchant(ctx context.Context, req *CreateMerchantInternalReq) (*entity.Merchant, *entity.MerchantMember, error) {
	merchantMasterMember := &entity.MerchantMember{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		Password:   utility.PasswordEncrypt(req.Password),
		UserName:   req.UserName,
		Mobile:     req.Phone,
		CreateTime: gtime.Now().Timestamp(),
	}
	merchantInfo := &entity.Merchant{
		CompanyId: 0,
		Phone:     req.Phone,
		ApiKey:    utility.GenerateRandomAlphanumeric(32), //32 bit open api key
	}
	// transaction create Merchant
	err := dao.Merchant.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
		insert, err := dao.MerchantMember.Ctx(ctx).Data(merchantMasterMember).OmitNil().Insert(merchantMasterMember)
		if err != nil {
			return err
		}
		id, err := insert.LastInsertId()
		if err != nil {
			return err
		}
		merchantMasterMember.Id = uint64(id)

		merchantInfo.UserId = id

		insert, err = dao.Merchant.Ctx(ctx).Data(merchantInfo).OmitNil().Insert(merchantInfo)
		if err != nil {
			return err
		}
		merchantId, err := insert.LastInsertId()
		if err != nil {
			return err
		}
		merchantInfo.Id = uint64(merchantId)
		merchantMasterMember.MerchantId = merchantInfo.Id
		// bind merchantMemberAccount
		_, err = dao.MerchantMember.Ctx(ctx).Data(g.Map{
			dao.MerchantMember.Columns().MerchantId: merchantId,
			dao.MerchantMember.Columns().GmtModify:  gtime.Now(),
		}).Where(dao.MerchantMember.Columns().Id, id).OmitNil().Update()
		if err != nil {
			return err
		}
		return nil
	})

	utility.AssertError(err, "Server Error")
	var newOne *entity.MerchantMember
	newOne = query.GetMerchantMemberById(ctx, merchantMasterMember.Id)
	utility.Assert(newOne != nil, "Server Error")
	err = cloud.MerchantSetupForCloudMode(ctx, merchantInfo.Id)
	return merchantInfo, newOne, err
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
