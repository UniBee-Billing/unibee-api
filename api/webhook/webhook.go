// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package webhook

import (
	"context"
	
	"unibee-api/api/webhook/setup"
)

type IWebhookSetup interface {
	New(ctx context.Context, req *setup.NewReq) (res *setup.NewRes, err error)
	Update(ctx context.Context, req *setup.UpdateReq) (res *setup.UpdateRes, err error)
}


