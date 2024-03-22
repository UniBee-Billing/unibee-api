package test

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	merchant2 "unibee/internal/logic/merchant"
	"unibee/internal/query"
	_ "unibee/test"
)

func TestMerchantCreateAndDelete(t *testing.T) {
	ctx := context.Background()
	t.Run("Test for Merchant Create|Login|ChangePassword|Delete", func(t *testing.T) {
		merchant, member, err := merchant2.CreateMerchant(ctx, &merchant2.CreateMerchantInternalReq{
			FirstName: "test",
			LastName:  "test",
			Email:     "autotest@wowow.io",
			Password:  "test123456",
			Phone:     "123456",
			UserName:  "test",
		})
		require.Nil(t, err)
		require.NotNil(t, merchant)
		merchant = query.GetMerchantById(ctx, merchant.Id)
		require.NotNil(t, merchant)
		member = query.GetMerchantMemberById(ctx, member.Id)
		require.NotNil(t, member)
		one, token := merchant2.PasswordLogin(ctx, member.Email, "test123456")
		require.NotNil(t, one)
		require.NotNil(t, token)
		merchant2.ChangeMerchantMemberPassword(ctx, member.Email, "test123456", "test654321")
		one, token = merchant2.PasswordLogin(ctx, member.Email, "test654321")
		require.NotNil(t, one)
		require.NotNil(t, token)
		merchant2.ChangeMerchantMemberPasswordWithOutOldVerify(ctx, member.Email, "test123456")
		one, token = merchant2.PasswordLogin(ctx, member.Email, "test123456")
		require.NotNil(t, one)
		require.NotNil(t, token)
		err = merchant2.HardDeleteMerchant(ctx, merchant.Id)
		require.Nil(t, err)
		merchant = query.GetMerchantById(ctx, merchant.Id)
		require.Nil(t, merchant)
		member = query.GetMerchantMemberById(ctx, member.Id)
		require.Nil(t, member)
	})
}
