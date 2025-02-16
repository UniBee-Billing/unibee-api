// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"unibee/internal/dao/default/internal"
)

// internalMerchantRoleDao is internal type for wrapping internal DAO implements.
type internalMerchantRoleDao = *internal.MerchantRoleDao

// merchantRoleDao is the data access object for table merchant_role.
// You can define custom methods on it to extend its functionality as you wish.
type merchantRoleDao struct {
	internalMerchantRoleDao
}

var (
	// MerchantRole is globally public accessible object for table merchant_role operations.
	MerchantRole = merchantRoleDao{
		internal.NewMerchantRoleDao(),
	}
)

// Fill with you ideas below.
