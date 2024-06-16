package role

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"unibee/api/bean"
	"unibee/internal/query"
	"unibee/test"
)

func TestMerchantRole(t *testing.T) {
	ctx := context.Background()
	var testRole = "TestRole"
	var err error
	t.Run("Test for merchant role New|Get|Edit|Delete|List", func(t *testing.T) {
		list, _ := MerchantRoleList(ctx, test.TestMerchant.Id)
		require.NotNil(t, list)
		require.Equal(t, true, len(list) == 0)
		err = NewMerchantRole(ctx, &CreateRoleInternalReq{
			MerchantId:     test.TestMerchant.Id,
			Role:           testRole,
			PermissionData: nil,
		})
		require.Nil(t, err)
		require.NotNil(t, query.GetRoleByName(ctx, test.TestMerchant.Id, testRole))
		err = EditMerchantRole(ctx, &CreateRoleInternalReq{
			Id:         query.GetRoleByName(ctx, test.TestMerchant.Id, testRole).Id,
			MerchantId: test.TestMerchant.Id,
			Role:       testRole,
			PermissionData: []*bean.MerchantRolePermission{{
				Group:       "test",
				Permissions: nil,
			}},
		})
		require.Nil(t, err)
		list, _ = MerchantRoleList(ctx, test.TestMerchant.Id)
		require.NotNil(t, list)
		require.Equal(t, true, len(list) > 0)
	})
	t.Run("Test for merchant role HardDelete", func(t *testing.T) {
		err = HardDeleteMerchantRole(ctx, test.TestMerchant.Id, testRole)
		require.Nil(t, err)
	})
}
