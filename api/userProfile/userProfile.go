// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package userProfile

import (
	"context"

	"go-oversea-pay/api/userProfile/v1"
)

type IUserProfileV1 interface {
	Profile(ctx context.Context, req *v1.ProfileReq) (res *v1.ProfileRes, err error)
}
