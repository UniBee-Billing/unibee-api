// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package xin

import (
	"context"

	"go-oversea-pay/api/xin/v1"
)

type IXinV1 interface {
	Insert(ctx context.Context, req *v1.InsertReq) (res *v1.InsertRes, err error)
	Get(ctx context.Context, req *v1.GetReq) (res *v1.GetRes, err error)
}
