package prepare

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

func GetPlanByName(ctx context.Context, name string) (one *entity.Plan) {
	if len(name) <= 0 {
		return nil
	}
	err := dao.Plan.Ctx(ctx).Where(dao.Plan.Columns().PlanName, name).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func CreateTestPlan(ctx context.Context, merchantId uint64) (one *entity.Plan, err error) {
	one = &entity.Plan{
		MerchantId:                merchantId,
		PlanName:                  "autotest_x",
		Amount:                    100,
		Currency:                  "USD",
		IntervalUnit:              "day",
		IntervalCount:             1,
		Status:                    consts.PlanStatusActive,
		Description:               "autotest_x",
		Type:                      consts.PlanTypeMain,
		GatewayProductName:        "autotest_x",
		GatewayProductDescription: "autotest_x",
		ImageUrl:                  "http://api.unibee.top",
		HomeUrl:                   "http://api.unibee.top",
		BindingAddonIds:           "",
		BindingOnetimeAddonIds:    "",
		ExtraMetricData:           "",
		GasPayer:                  "",
		PublishStatus:             consts.PlanPublishStatusPublished,
		MetaData:                  utility.MarshalToJsonString(map[string]string{"type": "test"}),
		CreateTime:                gtime.Now().Timestamp(),
	}
	result, err := dao.Plan.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		return nil, gerror.Newf(`PlanCreate record insert failure %s`, err)
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))
	return one, nil
}

func CreateTestAddon(ctx context.Context, merchantId uint64, name string, addonType int) (one *entity.Plan, err error) {
	one = &entity.Plan{
		MerchantId:                merchantId,
		PlanName:                  name,
		Amount:                    100,
		Currency:                  "USD",
		IntervalUnit:              "day",
		IntervalCount:             1,
		Status:                    consts.PlanStatusActive,
		Description:               "autotest_x",
		Type:                      addonType,
		GatewayProductName:        "autotest_x",
		GatewayProductDescription: "autotest_x",
		ImageUrl:                  "http://api.unibee.top",
		HomeUrl:                   "http://api.unibee.top",
		BindingAddonIds:           "",
		BindingOnetimeAddonIds:    "",
		ExtraMetricData:           "",
		GasPayer:                  "",
		PublishStatus:             consts.PlanPublishStatusPublished,
		MetaData:                  utility.MarshalToJsonString(map[string]string{"type": "test"}),
		CreateTime:                gtime.Now().Timestamp(),
	}
	result, err := dao.Plan.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		return nil, gerror.Newf(`PlanCreate record insert failure %s`, err)
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))
	return one, nil
}
