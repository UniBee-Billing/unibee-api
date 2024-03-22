package prepare

import (
	"context"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/plan/service"
	entity "unibee/internal/model/entity/oversea_pay"
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
	return service.PlanCreate(ctx, &service.PlanInternalReq{
		PlanName:           "autotest_x",
		Amount:             100,
		Currency:           "USD",
		IntervalUnit:       "day",
		IntervalCount:      1,
		Description:        "autotest_x",
		Type:               consts.PlanTypeMain,
		ProductName:        "autotest_x",
		ProductDescription: "autotest_x",
		ImageUrl:           "http://api.unibee.top",
		HomeUrl:            "http://api.unibee.top",
		AddonIds:           nil,
		OnetimeAddonIds:    nil,
		MetricLimits:       nil,
		GasPayer:           "",
		Metadata:           map[string]string{"type": "test"},
		MerchantId:         merchantId,
	})
}

func CreateTestAddon(ctx context.Context, merchantId uint64) (one *entity.Plan, err error) {
	return service.PlanCreate(ctx, &service.PlanInternalReq{
		PlanName:           "autotest_addon_x",
		Amount:             100,
		Currency:           "USD",
		IntervalUnit:       "day",
		IntervalCount:      1,
		Description:        "autotest_x",
		Type:               consts.PlanTypeRecurringAddon,
		ProductName:        "autotest_x",
		ProductDescription: "autotest_x",
		ImageUrl:           "http://api.unibee.top",
		HomeUrl:            "http://api.unibee.top",
		AddonIds:           nil,
		OnetimeAddonIds:    nil,
		MetricLimits:       nil,
		GasPayer:           "",
		Metadata:           map[string]string{"type": "test"},
		MerchantId:         merchantId,
	})
}
