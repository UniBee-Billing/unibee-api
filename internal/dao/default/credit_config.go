// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"unibee/internal/dao/default/internal"
)

// internalCreditConfigDao is internal type for wrapping internal DAO implements.
type internalCreditConfigDao = *internal.CreditConfigDao

// creditConfigDao is the data access object for table credit_config.
// You can define custom methods on it to extend its functionality as you wish.
type creditConfigDao struct {
	internalCreditConfigDao
}

var (
	// CreditConfig is globally public accessible object for table credit_config operations.
	CreditConfig = creditConfigDao{
		internal.NewCreditConfigDao(),
	}
)

// Fill with you ideas below.