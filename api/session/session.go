// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package session

import (
	"context"
	
	"unibee-api/api/session/user"
)

type ISessionUser interface {
	New(ctx context.Context, req *user.NewReq) (res *user.NewRes, err error)
}


