package prepare

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
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
	merchantInfo := &entity.Merchant{
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

	return merchantInfo, merchantMasterMember, nil
}
