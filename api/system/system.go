// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package system

import (
	"context"
	
	"go-oversea-pay/api/system/invoice"
	"go-oversea-pay/api/system/subscription"
)

type ISystemInvoice interface {
	BulkChannelSync(ctx context.Context, req *invoice.BulkChannelSyncReq) (res *invoice.BulkChannelSyncRes, err error)
	ChannelSync(ctx context.Context, req *invoice.ChannelSyncReq) (res *invoice.ChannelSyncRes, err error)
}

type ISystemSubscription interface {
	BulkChannelSync(ctx context.Context, req *subscription.BulkChannelSyncReq) (res *subscription.BulkChannelSyncRes, err error)
}


