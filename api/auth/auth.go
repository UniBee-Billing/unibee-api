// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package auth

import (
	"context"
	
	"go-oversea-pay/api/auth/v1"
)

type IAuthV1 interface {
	Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error)
	Register(ctx context.Context, req *v1.RegisterReq) (res *v1.RegisterRes, err error)
}


