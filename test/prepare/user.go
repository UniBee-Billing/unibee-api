package prepare

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

func CreateTestUser(ctx context.Context, merchantId uint64) (one *entity.UserAccount, err error) {
	one = &entity.UserAccount{
		ExternalUserId: "auto_x",
		Email:          "testuser@wowow.io",
		FirstName:      "test",
		LastName:       "test",
		Password:       utility.PasswordEncrypt("test123456"),
		Phone:          "test",
		Address:        "test",
		UserName:       "test",
		CountryCode:    "CN",
		CountryName:    "China",
		MerchantId:     merchantId,
		CreateTime:     gtime.Now().Timestamp(),
	}
	result, err := dao.UserAccount.Ctx(ctx).Data(one).OmitNil().Insert(one)
	utility.AssertError(err, "Server Error")
	id, err := result.LastInsertId()
	utility.AssertError(err, "Server Error")
	one.Id = uint64(id)
	return one, nil
}
