package jwt

import (
	"context"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"unibee/test"
)

func TestCreatePortalToken(t *testing.T) {
	ctx := context.Background()
	t.Run("Test for jwt token ", func(t *testing.T) {
		token, err := CreatePortalToken(TOKENTYPEUSER, test.TestMerchant.Id, test.TestUser.Id, test.TestUser.Email)
		require.Nil(t, err)
		require.NotNil(t, token)
		require.Equal(t, true, IsPortalToken(token))
		claim := ParsePortalToken(token)
		require.NotNil(t, claim)
		require.Equal(t, claim.MerchantId, test.TestMerchant.Id)
		require.Equal(t, claim.Id, test.TestUser.Id)
		require.Equal(t, claim.Email, test.TestUser.Email)
		token, err = CreatePortalToken(TOKENTYPEUSER, test.TestMerchant.Id, test.TestMerchantMember.Id, test.TestMerchantMember.Email)
		require.Nil(t, err)
		require.NotNil(t, token)
		require.Equal(t, true, IsPortalToken(token))
		claim = ParsePortalToken(token)
		require.NotNil(t, claim)
		require.Equal(t, claim.MerchantId, test.TestMerchant.Id)
		require.Equal(t, claim.Id, test.TestMerchantMember.Id)
		require.Equal(t, claim.Email, test.TestMerchantMember.Email)
		require.Equal(t, true, PutAuthTokenToCache(ctx, token, strconv.FormatUint(test.TestMerchant.Id, 10)))
		require.Equal(t, true, IsAuthTokenAvailable(ctx, token))
		require.Equal(t, true, ResetAuthTokenTTL(ctx, token))
		require.Equal(t, true, DelAuthToken(ctx, token))
		require.Equal(t, false, IsAuthTokenAvailable(ctx, token))
		require.Equal(t, false, ResetAuthTokenTTL(ctx, token))
	})
}
