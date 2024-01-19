package xin_service

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	dao "go-oversea-pay/internal/dao/xin"
	entity "go-oversea-pay/internal/model/entity/xin"
)

func QueryTest(ctx context.Context) (out []*entity.Test, err error) {
	err = dao.Test.Ctx(ctx).Scan(&out)
	if err != nil {
		return
	}
	if out == nil {
		err = gerror.Newf(`record not found`)
	}
	return
}

func InsertTest(ctx context.Context, name string) (out *entity.Test, err error) {
	test := &entity.Test{
		Name: name,
	}
	result, err := dao.Test.Ctx(ctx).Data(test).OmitNil().Insert(test)
	if err != nil {
		err = gerror.Newf(`record insert failure %s`, err)
		return
	}
	id, _ := result.LastInsertId()
	test.Id = uint(id)
	out = test
	g.Log().Infof(ctx, "insert success result: %s", result)
	return
}
