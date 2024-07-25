package merchant

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	member2 "unibee/internal/logic/member"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	_ "unibee/test"
)

func TestMerchantCreateAndDelete(t *testing.T) {
	ctx := context.Background()
	var merchant *entity.Merchant
	var member *entity.MerchantMember
	var err error
	t.Run("Test for Merchant Create|Login|ChangePassword|Delete", func(t *testing.T) {
		merchant, member, err = CreateMerchant(ctx, &CreateMerchantInternalReq{
			FirstName: "test",
			LastName:  "test",
			Email:     "autotest@wowow.io",
			Password:  "test123456",
			Phone:     "123456",
			UserName:  "test",
		})
		require.Nil(t, err)
		require.NotNil(t, merchant)
	})
	t.Run("Test for merchant Login|ChangePassword", func(t *testing.T) {
		merchant = query.GetMerchantById(ctx, merchant.Id)
		require.NotNil(t, merchant)
		member = query.GetMerchantMemberById(ctx, member.Id)
		require.NotNil(t, member)
		one, token := member2.PasswordLogin(ctx, member.Email, "test123456")
		require.NotNil(t, one)
		require.NotNil(t, token)
		member2.ChangeMerchantMemberPassword(ctx, member.Email, "test123456", "test654321")
		one, token = member2.PasswordLogin(ctx, member.Email, "test654321")
		require.NotNil(t, one)
		require.NotNil(t, token)
		member2.ChangeMerchantMemberPasswordWithOutOldVerify(ctx, member.Email, "test123456")
		one, token = member2.PasswordLogin(ctx, member.Email, "test123456")
		require.NotNil(t, one)
		require.NotNil(t, token)
		oldKey := NewOpenApiKey(ctx, one.MerchantId)
		require.NotNil(t, oldKey)
		require.Equal(t, true, len(oldKey) > 0)
		merchant = query.GetMerchantByApiKey(ctx, oldKey)
		require.NotNil(t, merchant)
		newKey := NewOpenApiKey(ctx, one.MerchantId)
		require.NotNil(t, newKey)
		require.Equal(t, true, len(newKey) > 0)
		require.Nil(t, query.GetMerchantByApiKey(ctx, oldKey))
		require.NotNil(t, GetMerchantFromCache(ctx, oldKey))
		require.NotNil(t, query.GetMerchantByApiKey(ctx, newKey))
	})
	t.Run("Test for merchant HardDelete", func(t *testing.T) {
		err = HardDeleteMerchant(ctx, merchant.Id)
		require.Nil(t, err)
		merchant = query.GetMerchantById(ctx, merchant.Id)
		require.Nil(t, merchant)
		member = query.GetMerchantMemberById(ctx, member.Id)
		require.Nil(t, member)
	})
}
