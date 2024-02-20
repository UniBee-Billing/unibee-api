package metric

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee-api/internal/dao/oversea_pay"
	"unibee-api/internal/logic/gateway/ro"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

const (
	MerchantMetricPlanLimitCacheKeyPrefix = "MerchantMetricPlanLimitCacheKeyPrefix_"
	MerchantMetricPlanLimitCacheExpire    = 24 * 60 * 60
)

func MerchantMetricPlanLimitCachedList(ctx context.Context, merchantId int64, planId int64, reloadCache bool) []*ro.MerchantMetricPlanLimitVo {
	utility.Assert(merchantId > 0, "invalid merchantId")
	utility.Assert(planId > 0, "invalid planId")
	var list = make([]*ro.MerchantMetricPlanLimitVo, 0)
	cacheKey := fmt.Sprintf("%s%d%d", MerchantMetricPlanLimitCacheKeyPrefix, merchantId, planId)
	if !reloadCache {
		get, _ := g.Redis().Get(ctx, cacheKey)
		value := get.String()
		if len(value) > 0 {
			_ = utility.UnmarshalFromJsonString(value, &list)
			if len(list) > 0 {
				return list
			}
		}
	}
	if merchantId > 0 {
		var entities []*entity.MerchantMetricPlanLimit
		err := dao.MerchantMetric.Ctx(ctx).
			Where(entity.MerchantMetricPlanLimit{MerchantId: merchantId}).
			Where(entity.MerchantMetricPlanLimit{PlanId: planId}).
			Where(entity.MerchantMetricPlanLimit{IsDeleted: 0}).
			Scan(&entities)
		if err == nil && len(entities) > 0 {
			for _, one := range entities {
				list = append(list, &ro.MerchantMetricPlanLimitVo{
					Id:          one.Id,
					MerchantId:  one.MerchantId,
					MetricId:    one.MetricId,
					Metric:      GetMerchantMetricVo(ctx, one.MetricId),
					PlanId:      one.PlanId,
					MetricLimit: one.MetricLimit,
					UpdateTime:  one.GmtModify.Timestamp(),
					CreateTime:  one.CreateTime,
				})
			}
		}
	}
	if len(list) > 0 {
		_, _ = g.Redis().Set(ctx, cacheKey, utility.MarshalToJsonString(list))
		_, _ = g.Redis().Expire(ctx, cacheKey, MerchantMetricPlanLimitCacheExpire) // one day cache expire time
	}
	return list
}

type MerchantMetricPlanLimitInternalReq struct {
	MerchantId int64  `p:"merchantId" dc:"MerchantId" v:"required"`
	LimitId    uint64 `p:"limitId" dc:"LimitId" `
	MetricId   int64  `p:"metricId" dc:"MetricId" `
	PlanId     int64  `p:"planId" dc:"PlanId" `
	Limit      int64  `p:"limit" dc:"Limit" `
}

func NewMerchantMetricPlanLimit(ctx context.Context, req *MerchantMetricPlanLimitInternalReq) (*ro.MerchantMetricPlanLimitVo, error) {
	utility.Assert(req.MerchantId > 0, "invalid merchantId")
	utility.Assert(req.PlanId > 0, "invalid planId")
	utility.Assert(req.MetricId > 0, "invalid metricId")
	//metric check
	metric := query.GetMerchantMetric(ctx, req.MetricId)
	utility.Assert(metric != nil, "metric not found")
	utility.Assert(metric.Type == MetricTypeLimitMetered, "Metric Not MetricTypeLimitMetered Type")
	//Plan check
	plan := query.GetPlanById(ctx, req.PlanId)
	utility.Assert(plan != nil, "plan not found")
	utility.Assert(plan.MerchantId == req.MerchantId, "plan merchantId not match")

	var one *entity.MerchantMetricPlanLimit
	err := dao.MerchantMetricPlanLimit.Ctx(ctx).
		Where(entity.MerchantMetricPlanLimit{MerchantId: req.MerchantId}).
		Where(entity.MerchantMetricPlanLimit{PlanId: req.PlanId}).
		Where(entity.MerchantMetricPlanLimit{MetricId: req.MetricId}).
		Where(entity.MerchantMetricPlanLimit{IsDeleted: 0}).
		Scan(&one)
	utility.AssertError(err, "server error")
	utility.Assert(one == nil, "metric limit already exist")
	one = &entity.MerchantMetricPlanLimit{
		MerchantId:  req.MerchantId,
		MetricId:    req.MetricId,
		PlanId:      req.PlanId,
		MetricLimit: req.Limit,
		CreateTime:  gtime.Now().Timestamp(),
	}
	result, err := dao.MerchantMetricPlanLimit.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		g.Log().Errorf(ctx, "NewMerchantMetricPlanLimit Insert err:%s", err.Error())
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(id)
	// reload Cache
	MerchantMetricPlanLimitCachedList(ctx, req.MerchantId, req.PlanId, true)
	return &ro.MerchantMetricPlanLimitVo{
		Id:          one.Id,
		MerchantId:  one.MerchantId,
		MetricId:    one.MetricId,
		Metric:      GetMerchantMetricVo(ctx, one.MetricId),
		PlanId:      one.PlanId,
		MetricLimit: one.MetricLimit,
		CreateTime:  one.CreateTime,
	}, nil
}

func EditMerchantMetricPlanLimit(ctx context.Context, req *MerchantMetricPlanLimitInternalReq) (*ro.MerchantMetricPlanLimitVo, error) {
	utility.Assert(req.MerchantId > 0, "invalid merchantId")
	utility.Assert(req.LimitId > 0, "invalid limitId")
	var one *entity.MerchantMetricPlanLimit
	err := dao.MerchantMetricPlanLimit.Ctx(ctx).
		Where(entity.MerchantMetricPlanLimit{MerchantId: req.MerchantId}).
		Where(entity.MerchantMetricPlanLimit{Id: req.LimitId}).
		Where(entity.MerchantMetricPlanLimit{IsDeleted: 0}).
		Scan(&one)
	utility.AssertError(err, "server error")
	utility.Assert(one != nil, "metric limit not found")
	_, err = dao.MerchantMetricPlanLimit.Ctx(ctx).Data(g.Map{
		dao.MerchantMetricPlanLimit.Columns().MetricLimit: req.Limit,
		dao.MerchantMetricPlanLimit.Columns().GmtModify:   gtime.Now(),
	}).Where(dao.MerchantMetricPlanLimit.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		g.Log().Errorf(ctx, "EditMerchantMetricPlanLimit Update err:%s", err.Error())
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	one.MetricLimit = req.Limit
	// reload Cache
	MerchantMetricPlanLimitCachedList(ctx, one.MerchantId, req.PlanId, true)
	return &ro.MerchantMetricPlanLimitVo{
		Id:          one.Id,
		MerchantId:  one.MerchantId,
		MetricId:    one.MetricId,
		Metric:      GetMerchantMetricVo(ctx, one.MetricId),
		PlanId:      one.PlanId,
		MetricLimit: one.MetricLimit,
		UpdateTime:  gtime.Now().Timestamp(),
		CreateTime:  one.CreateTime,
	}, nil
}

func DeleteMerchantMetricPlanLimit(ctx context.Context, merchantId int64, metricId int64) error {
	utility.Assert(merchantId > 0, "invalid merchantId")
	utility.Assert(metricId > 0, "invalid metricId")
	one := query.GetMerchantMetric(ctx, metricId)
	utility.Assert(one != nil, "metric limit not found")
	_, err := dao.MerchantMetricPlanLimit.Ctx(ctx).Data(g.Map{
		dao.MerchantMetricPlanLimit.Columns().IsDeleted: 1,
		dao.MerchantMetricPlanLimit.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantMetricPlanLimit.Columns().Id, one.Id).OmitNil().Update()
	// reload Cache
	MerchantMetricPlanLimitCachedList(ctx, one.MerchantId, metricId, true)
	return err
}
