package auth

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/auth/v1"
)

func (c *ControllerV1) Register(ctx context.Context, req *v1.RegisterReq) (res *v1.RegisterRes, err error) {
	utility.Assert(len(req.Phone) > 0, "phone not null")
	user := &entity.UserAccount{
		UserName: req.Email,
		Email:    req.Email,
		Mobile:   req.Phone,
		Gender:   req.Gender,
	}
	result, err := dao.UserAccount.Ctx(ctx).Data(user).OmitEmpty().Insert(user)
	if err != nil {
		err = gerror.Newf(`record insert failure %s`, err)
		return
	}
	id, _ := result.LastInsertId()
	user.Id = uint64(id)
	var newOne *entity.UserAccount
	err = dao.UserAccount.Ctx(ctx).Where(entity.UserAccount{Id: user.Id}).OmitEmpty().Scan(&newOne)
	if err != nil {
		return nil, err
	}

	return &v1.RegisterRes{User: newOne}, nil
}
