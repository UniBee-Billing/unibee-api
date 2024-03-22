package prepare

import (
	"context"
	"unibee/internal/logic/auth"
	entity "unibee/internal/model/entity/oversea_pay"
)

func CreateTestUser(ctx context.Context, merchantId uint64) (one *entity.UserAccount, err error) {
	return auth.CreateUser(ctx, &auth.NewReq{
		ExternalUserId: "auto_x",
		Email:          "testuser@wowow.io",
		FirstName:      "test",
		LastName:       "test",
		Password:       "test123456",
		Phone:          "test",
		Address:        "test",
		UserName:       "test",
		CountryCode:    "CN",
		CountryName:    "China",
		MerchantId:     merchantId,
	})
}
