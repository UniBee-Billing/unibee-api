package prepare

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func CreateTestMerchantAccount(ctx context.Context) (*entity.Merchant, *entity.MerchantMember, error) {
	merchantMasterMember := &entity.MerchantMember{
		FirstName:  "test",
		LastName:   "test",
		Email:      "test@wowow.io",
		Password:   utility.PasswordEncrypt("test123456"),
		UserName:   "test",
		Mobile:     "123456",
		CreateTime: gtime.Now().Timestamp(),
	}
	merchant := &entity.Merchant{
		CompanyId:   0,
		Phone:       "123456",
		Host:        "autotest.unibee.top",
		CompanyName: "Unibee_AutoTest",
		CompanyLogo: "http://unibee.top/files/invoice/cm/czgryizwt00wwofira.png",
		Email:       "jack.fu@wowow.io",
		ApiKey:      utility.GenerateRandomAlphanumeric(32), //32 bit open api key
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
		}).Where(dao.Merchant.Columns().Id, id).OmitNil().Update()
		if err != nil {
			return err
		}
		return nil
	})

	utility.AssertError(err, "Server Error")
	var newOne *entity.MerchantMember
	newOne = query.GetMerchantMemberById(ctx, merchantMasterMember.Id)
	utility.Assert(newOne != nil, "Server Error")

	return merchant, merchantMasterMember, nil
}
