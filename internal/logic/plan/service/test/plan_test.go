package test

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"unibee/internal/consts"
	"unibee/internal/logic/plan/service"
	"unibee/internal/query"
	"unibee/test"
)

func TestPlanCreateAndDelete(t *testing.T) {
	ctx := context.Background()
	t.Run("Test for Plan Create|Edit|Publish|UnPublish|Delete", func(t *testing.T) {
		one, err := service.PlanCreate(ctx, &service.PlanInternalReq{
			PlanName:           "autotest",
			Amount:             100,
			Currency:           "USD",
			IntervalUnit:       "day",
			IntervalCount:      1,
			Description:        "autotest",
			Type:               consts.PlanTypeMain,
			ProductName:        "autotest",
			ProductDescription: "autotest",
			ImageUrl:           "http://api.unibee.top",
			HomeUrl:            "http://api.unibee.top",
			AddonIds:           nil,
			OnetimeAddonIds:    nil,
			MetricLimits:       nil,
			GasPayer:           "",
			Metadata:           map[string]string{"type": "test"},
			MerchantId:         test.TestMerchant.Id,
		})
		require.Nil(t, err)
		require.NotNil(t, one)
		one = query.GetPlanById(ctx, one.Id)
		require.NotNil(t, one)
		one, err = service.PlanEdit(ctx, &service.PlanInternalReq{
			PlanId:             one.Id,
			PlanName:           "autotest",
			Amount:             200,
			Currency:           "USD",
			IntervalUnit:       "day",
			IntervalCount:      1,
			Description:        "autotest",
			ProductName:        "autotest",
			ProductDescription: "autotest",
			ImageUrl:           "http://api.unibee.top",
			HomeUrl:            "http://api.unibee.top",
			AddonIds:           nil,
			OnetimeAddonIds:    nil,
			MetricLimits:       nil,
			GasPayer:           "",
			Metadata:           map[string]string{"type": "test"},
			MerchantId:         test.TestMerchant.Id,
		})
		require.Nil(t, err)
		require.NotNil(t, one)
		one = query.GetPlanById(ctx, one.Id)
		require.NotNil(t, one)
		require.Equal(t, one.Amount, int64(200))
		//activate & publish
		publishPlans := service.SubscriptionPlanList(ctx, &service.SubscriptionPlanListInternalReq{
			MerchantId:    test.TestMerchant.Id,
			Status:        []int{consts.PlanStatusActive},
			PublishStatus: consts.PlanPublishStatusPublished,
			Page:          0,
			Count:         10,
		})
		require.Equal(t, len(publishPlans), 0)
		err = service.PlanActivate(ctx, one.Id)
		require.Nil(t, err)
		err = service.PlanPublish(ctx, one.Id)
		require.Nil(t, err)
		publishPlans = service.SubscriptionPlanList(ctx, &service.SubscriptionPlanListInternalReq{
			MerchantId:    test.TestMerchant.Id,
			Status:        []int{consts.PlanStatusActive},
			PublishStatus: consts.PlanPublishStatusPublished,
			Page:          0,
			Count:         10,
		})
		require.Equal(t, len(publishPlans), 1)
		err = service.HardDeletePlan(ctx, one.Id)
		require.Nil(t, err)
		one = query.GetPlanById(ctx, one.Id)
		require.Nil(t, one)
	})
}
