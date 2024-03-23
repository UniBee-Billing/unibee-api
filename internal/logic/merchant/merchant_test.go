package merchant

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"unibee/internal/query"
	_ "unibee/test"
)

func TestMerchantCreateAndDelete(t *testing.T) {
	ctx := context.Background()
	t.Run("Test for Merchant Create|Login|ChangePassword|Delete", func(t *testing.T) {
		merchant, member, err := CreateMerchant(ctx, &CreateMerchantInternalReq{
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
		one, token := PasswordLogin(ctx, member.Email, "test123456")
		require.NotNil(t, one)
		require.NotNil(t, token)
		ChangeMerchantMemberPassword(ctx, member.Email, "test123456", "test654321")
		one, token = PasswordLogin(ctx, member.Email, "test654321")
		require.NotNil(t, one)
		require.NotNil(t, token)
		ChangeMerchantMemberPasswordWithOutOldVerify(ctx, member.Email, "test123456")
		one, token = PasswordLogin(ctx, member.Email, "test123456")
		require.NotNil(t, one)
		require.NotNil(t, token)
		err = HardDeleteMerchant(ctx, merchant.Id)
		require.Nil(t, err)
		merchant = query.GetMerchantById(ctx, merchant.Id)
		require.Nil(t, merchant)
		member = query.GetMerchantMemberById(ctx, member.Id)
		require.Nil(t, member)
	})
}
