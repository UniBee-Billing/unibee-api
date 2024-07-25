package metric

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"unibee/api/bean"
	"unibee/internal/query"
	"unibee/test"
)

func TestMerchantMetric(t *testing.T) {
	ctx := context.Background()
	var one *bean.MerchantMetric
	var limit *bean.MerchantMetricPlanLimit
	var err error
	t.Run("Test for merchant metric New|Get|Detail|Edit|Delete|List", func(t *testing.T) {
		list, _ := MerchantMetricList(ctx, test.TestMerchant.Id)
		require.NotNil(t, list)
		require.Equal(t, 0, len(list))
		one, err = NewMerchantMetric(ctx, &NewMerchantMetricInternalReq{
			MerchantId:          test.TestMerchant.Id,
			Code:                "test_metric_ex",
			Name:                "test",
			Description:         "test",
			AggregationType:     1,
			AggregationProperty: "",
		})
		require.Nil(t, err)
		require.NotNil(t, one)
		one = GetMerchantMetricSimplify(ctx, one.Id)
		require.NotNil(t, one)
		one = MerchantMetricDetail(ctx, test.TestMerchant.Id, one.Id)
		require.NotNil(t, one)
		require.Equal(t, "test", one.MetricName)
		require.Equal(t, "test", one.MetricDescription)
		list, _ = MerchantMetricList(ctx, test.TestMerchant.Id)
		require.NotNil(t, list)
		require.Equal(t, 1, len(list))
		one, err = EditMerchantMetric(ctx, test.TestMerchant.Id, one.Id, "test2", "test2")
		require.Nil(t, err)
		require.NotNil(t, one)
		require.Equal(t, "test2", one.MetricName)
		require.Equal(t, "test2", one.MetricDescription)
		err = DeleteMerchantMetric(ctx, test.TestMerchant.Id, one.Id)
		require.Nil(t, err)
		list, _ = MerchantMetricList(ctx, test.TestMerchant.Id)
		require.NotNil(t, list)
		require.Equal(t, 0, len(list))
		oldOne := one
		one, err = NewMerchantMetric(ctx, &NewMerchantMetricInternalReq{
			MerchantId:          test.TestMerchant.Id,
			Code:                "test_metric_ex",
			Name:                "test",
			Description:         "test",
			AggregationType:     1,
			AggregationProperty: "",
		})
		require.Nil(t, err)
		require.NotNil(t, one)
		err = HardDeleteMerchantMetric(ctx, test.TestMerchant.Id, oldOne.Id)
		require.Nil(t, err)
	})
	t.Run("Test for merchant metric limit", func(t *testing.T) {
		list := MerchantMetricPlanLimitCachedList(ctx, test.TestMerchant.Id, test.TestPlan.Id, true)
		require.NotNil(t, list)
		require.Equal(t, 0, len(list))
		limit, err = NewMerchantMetricPlanLimit(ctx, &MerchantMetricPlanLimitInternalReq{
			MerchantId:  test.TestMerchant.Id,
			MetricId:    one.Id,
			PlanId:      test.TestPlan.Id,
			MetricLimit: 1,
		})
		require.Nil(t, err)
		require.NotNil(t, limit)
		entityOne := query.GetMerchantMetricPlanLimit(ctx, limit.Id)
		require.NotNil(t, entityOne)
		require.Equal(t, uint64(1), entityOne.MetricLimit)
		list = MerchantMetricPlanLimitCachedList(ctx, test.TestMerchant.Id, test.TestPlan.Id, false)
		require.NotNil(t, list)
		require.Equal(t, 1, len(list))
		limit, err = EditMerchantMetricPlanLimit(ctx, &MerchantMetricPlanLimitInternalReq{
			MerchantId:        test.TestMerchant.Id,
			MetricId:          one.Id,
			MetricPlanLimitId: limit.Id,
			PlanId:            test.TestPlan.Id,
			MetricLimit:       2,
		})
		require.Nil(t, err)
		require.NotNil(t, limit)
		entityOne = query.GetMerchantMetricPlanLimit(ctx, limit.Id)
		require.NotNil(t, entityOne)
		require.Equal(t, uint64(2), entityOne.MetricLimit)

		err = DeleteMerchantMetricPlanLimit(ctx, test.TestMerchant.Id, limit.Id)
		require.Nil(t, err)
		oldOne := limit
		limit, err = NewMerchantMetricPlanLimit(ctx, &MerchantMetricPlanLimitInternalReq{
			MerchantId:  test.TestMerchant.Id,
			MetricId:    one.Id,
			PlanId:      test.TestPlan.Id,
			MetricLimit: 1,
		})
		require.Nil(t, err)
		require.NotNil(t, limit)
		err = HardDeleteMerchantMetricPlanLimit(ctx, one.MerchantId, oldOne.Id)
		require.Nil(t, err)
	})
	t.Run("Test for merchant metric and limit BulkBind|HardDelete", func(t *testing.T) {
		list := MerchantMetricPlanLimitCachedList(ctx, test.TestMerchant.Id, test.TestPlan.Id, true)
		require.NotNil(t, list)
		require.Equal(t, 1, len(list))
		err = BulkMetricLimitPlanBindingReplace(ctx, test.TestPlan, []*bean.BulkMetricLimitPlanBindingParam{{
			MetricId:    one.Id,
			MetricLimit: 3,
		}})
		require.Nil(t, err)
		list = MerchantMetricPlanLimitCachedList(ctx, test.TestMerchant.Id, test.TestPlan.Id, true)
		require.NotNil(t, list)
		require.Equal(t, 1, len(list))
		for _, limit = range list {
			require.Equal(t, uint64(3), limit.MetricLimit)
			err = HardDeleteMerchantMetricPlanLimit(ctx, one.MerchantId, limit.Id)
			require.Nil(t, err)
		}
		err = HardDeleteMerchantMetric(ctx, one.MerchantId, one.Id)
		require.Nil(t, err)
	})
}
