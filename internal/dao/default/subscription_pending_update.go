// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"unibee/internal/dao/default/internal"
)

// internalSubscriptionPendingUpdateDao is internal type for wrapping internal DAO implements.
type internalSubscriptionPendingUpdateDao = *internal.SubscriptionPendingUpdateDao

// subscriptionPendingUpdateDao is the data access object for table subscription_pending_update.
// You can define custom methods on it to extend its functionality as you wish.
type subscriptionPendingUpdateDao struct {
	internalSubscriptionPendingUpdateDao
}

var (
	// SubscriptionPendingUpdate is globally public accessible object for table subscription_pending_update operations.
	SubscriptionPendingUpdate = subscriptionPendingUpdateDao{
		internal.NewSubscriptionPendingUpdateDao(),
	}
)

// Fill with you ideas below.