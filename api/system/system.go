// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package system

import (
	"context"
	
	"go-oversea-pay/api/system/subscription"
)

type ISystemSubscription interface {
	BulkChannelSync(ctx context.Context, req *subscription.BulkChannelSyncReq) (res *subscription.BulkChannelSyncRes, err error)
}


