package prepare

import (
	"context"
	"unibee/internal/logic/merchant"
	entity "unibee/internal/model/entity/oversea_pay"
)

func CreateTestMerchantAccount(ctx context.Context) (*entity.Merchant, *entity.MerchantMember, error) {
	return merchant.CreateMerchant(ctx, &merchant.CreateMerchantInternalReq{
		FirstName: "test",
		LastName:  "test",
		Email:     "test@wowow.io",
		Password:  "test123456",
		Phone:     "123456",
		UserName:  "test",
	})
}
