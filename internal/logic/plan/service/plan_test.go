package service

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"unibee/api/merchant/plan"
	"unibee/internal/consts"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/test"
)

func TestPlanCreateAndDelete(t *testing.T) {
	ctx := context.Background()
	var one *entity.Plan
	var err error
	t.Run("Test for Plan Create|Edit|Activate|Delete|HardDelete", func(t *testing.T) {
		one, err = PlanCreate(ctx, &PlanInternalReq{
			PlanName:           "autotest",
			Amount:             100,
			Currency:           "USD",
			IntervalUnit:       "month",
			IntervalCount:      2,
			Description:        "autotest",
			Type:               consts.PlanTypeMain,
			ProductName:        "",
			ProductDescription: "",
			AddonIds:           nil,
			OnetimeAddonIds:    nil,
			MetricLimits:       nil,
			GasPayer:           "user",
			Metadata:           map[string]string{"type": "test"},
			MerchantId:         test.TestMerchant.Id,
		})
		require.Nil(t, err)
		require.NotNil(t, one)
		ones := PlanList(ctx, &SubscriptionPlanListInternalReq{
			MerchantId: test.TestMerchant.Id,
			Type:       []int{consts.PlanTypeMain},
			Status:     []int{consts.PlanStatusActive},
			Page:       0,
			Count:      10,
		})
		require.Equal(t, 1, len(ones))
		one = query.GetPlanById(ctx, one.Id)
		require.NotNil(t, one)
		one, err = PlanEdit(ctx, &PlanInternalReq{
			PlanId:             one.Id,
			PlanName:           "autotest",
			Amount:             200,
			Currency:           "USD",
			IntervalUnit:       "day",
			IntervalCount:      1,
			Description:        "autotest",
			ProductName:        "",
			ProductDescription: "",
			AddonIds:           []int64{int64(test.TestRecurringAddon.Id)},
			OnetimeAddonIds:    []int64{int64(test.TestOneTimeAddon.Id)},
			MetricLimits:       nil,
			GasPayer:           "user",
			Metadata:           map[string]string{"type": "test"},
			MerchantId:         test.TestMerchant.Id,
		})
		require.Nil(t, err)
		require.NotNil(t, one)
		one = query.GetPlanById(ctx, one.Id)
		require.NotNil(t, one)
		require.Equal(t, one.Amount, int64(200))
		detail, err := PlanDetail(ctx, test.TestMerchant.Id, one.Id)
		require.Nil(t, err)
		require.NotNil(t, detail)
		require.NotNil(t, detail.Plan.AddonIds)
		require.NotNil(t, detail.Plan.Addons)
		require.Equal(t, 1, len(detail.Plan.AddonIds))
		require.Equal(t, 1, len(detail.Plan.Addons))
		ones = PlanList(ctx, &SubscriptionPlanListInternalReq{
			MerchantId: test.TestMerchant.Id,
			Type:       []int{consts.PlanTypeMain},
			Status:     []int{consts.PlanStatusActive},
			Page:       0,
			Count:      10,
		})
		require.Equal(t, 1, len(ones))
		_, err = PlanDelete(ctx, one.Id)
		require.Nil(t, err)
		err = PlanActivate(ctx, one.Id)
		require.Nil(t, err)
		err = PlanActivate(ctx, one.Id)
		require.Nil(t, err)
		ones = PlanList(ctx, &SubscriptionPlanListInternalReq{
			MerchantId: test.TestMerchant.Id,
			Type:       []int{consts.PlanTypeMain},
			Status:     []int{consts.PlanStatusActive},
			Page:       0,
			Count:      10,
		})
		require.Equal(t, 2, len(ones))
		one, err = PlanAddonsBinding(ctx, &plan.AddonsBindingReq{
			PlanId:          one.Id,
			Action:          0,
			AddonIds:        nil,
			OnetimeAddonIds: nil,
		})
		detail, err = PlanDetail(ctx, test.TestMerchant.Id, one.Id)
		require.Nil(t, err)
		require.NotNil(t, detail.Plan.AddonIds)
		require.NotNil(t, detail.Plan.Addons)
		require.Equal(t, 0, len(detail.Plan.AddonIds))
		require.Equal(t, 0, len(detail.Plan.Addons))
		one, err = PlanAddonsBinding(ctx, &plan.AddonsBindingReq{
			PlanId:          one.Id,
			Action:          1,
			AddonIds:        []int64{int64(test.TestRecurringAddon.Id)},
			OnetimeAddonIds: []int64{int64(test.TestOneTimeAddon.Id)},
		})
		detail, err = PlanDetail(ctx, test.TestMerchant.Id, one.Id)
		require.Nil(t, err)
		require.NotNil(t, detail.Plan.AddonIds)
		require.NotNil(t, detail.Plan.Addons)
		require.Equal(t, 1, len(detail.Plan.AddonIds))
		require.Equal(t, 1, len(detail.Plan.Addons))
		ones = PlanList(ctx, &SubscriptionPlanListInternalReq{
			MerchantId: test.TestMerchant.Id,
			Type:       []int{consts.PlanTypeMain},
			Status:     []int{consts.PlanStatusActive},
			Page:       0,
			Count:      10,
		})
		require.Equal(t, 2, len(ones))

	})
	t.Run("Test for Plan HardDelete", func(t *testing.T) {
		err = HardDeletePlan(ctx, one.Id)
		require.Nil(t, err)
		one = query.GetPlanById(ctx, one.Id)
		require.Nil(t, one)
	})
	t.Run("Test for Plan Publish|UnPublish|PlanDetail", func(t *testing.T) {
		one = test.TestPlan
		//activate & publish
		publishPlans := PlanList(ctx, &SubscriptionPlanListInternalReq{
			MerchantId:    test.TestMerchant.Id,
			Type:          []int{consts.PlanTypeMain},
			Status:        []int{consts.PlanStatusActive},
			PublishStatus: consts.PlanPublishStatusPublished,
			SortField:     "gmt_create",
			SortType:      "desc",
			Page:          0,
			Count:         10,
		})
		require.Equal(t, 1, len(publishPlans))
		err := PlanUnPublish(ctx, one.Id)
		require.Nil(t, err)
		publishPlans = PlanList(ctx, &SubscriptionPlanListInternalReq{
			MerchantId:    test.TestMerchant.Id,
			Type:          []int{consts.PlanTypeMain},
			Status:        []int{consts.PlanStatusActive},
			PublishStatus: consts.PlanPublishStatusPublished,
			SortField:     "gmt_create",
			Page:          -1,
		})
		require.Equal(t, 0, len(publishPlans))
		err = PlanPublish(ctx, one.Id)
		require.Nil(t, err)
		publishPlans = PlanList(ctx, &SubscriptionPlanListInternalReq{
			MerchantId:    test.TestMerchant.Id,
			Status:        []int{consts.PlanStatusActive},
			PublishStatus: consts.PlanPublishStatusPublished,
			SortField:     "gmt_create",
			SortType:      "desc",
			Page:          0,
			Count:         10,
		})
		require.Equal(t, 3, len(publishPlans))
	})
}