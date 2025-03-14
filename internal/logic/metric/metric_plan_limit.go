package metric

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/operation_log"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

const (
	MerchantMetricPlanLimitCacheKeyPrefix = "MerchantMetricPlanLimitCacheKeyPrefix_"
	MerchantMetricPlanLimitCacheExpire    = 24 * 60 * 60
)

type MerchantMetricPlanLimitInternalReq struct {
	MerchantId        uint64 `json:"merchantId" dc:"MerchantId" v:"required"`
	MetricId          uint64 `json:"metricId" dc:"MetricId" `
	MetricPlanLimitId uint64 `json:"metricPlanLimitId" dc:"MetricPlanLimitId,use for edit" `
	PlanId            uint64 `json:"planId" dc:"PlanId" `
	MetricLimit       uint64 `json:"metricLimit" dc:"MetricLimit" `
}

func NewMerchantMetricPlanLimit(ctx context.Context, req *MerchantMetricPlanLimitInternalReq) (*detail.MerchantMetricPlanLimitDetail, error) {
	utility.Assert(req.MerchantId > 0, "invalid merchantId")
	utility.Assert(req.PlanId > 0, "invalid planId")
	utility.Assert(req.MetricId > 0, "invalid metricId")
	utility.Assert(req.MetricPlanLimitId == 0, "invalid MetricPlanLimitId, should not enter in")
	utility.Assert(req.MetricLimit > 0, "invalid MetricLimit")
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
		Where(dao.MerchantMetricPlanLimit.Columns().MerchantId, req.MerchantId).
		Where(dao.MerchantMetricPlanLimit.Columns().PlanId, req.PlanId).
		Where(dao.MerchantMetricPlanLimit.Columns().MetricId, req.MetricId).
		Where(dao.MerchantMetricPlanLimit.Columns().IsDeleted, 0).
		Scan(&one)
	utility.AssertError(err, "server error")
	utility.Assert(one == nil, "metric limit already exist")
	one = &entity.MerchantMetricPlanLimit{
		MerchantId:  req.MerchantId,
		MetricId:    req.MetricId,
		PlanId:      req.PlanId,
		MetricLimit: req.MetricLimit,
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
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Metric(%v)", one.Id),
		Content:        "NewPlanLimit",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return &detail.MerchantMetricPlanLimitDetail{
		Id:          one.Id,
		MerchantId:  one.MerchantId,
		MetricId:    one.MetricId,
		Metric:      GetMerchantMetricSimplify(ctx, one.MetricId),
		PlanId:      one.PlanId,
		MetricLimit: one.MetricLimit,
		CreateTime:  one.CreateTime,
	}, nil
}

func EditMerchantMetricPlanLimit(ctx context.Context, req *MerchantMetricPlanLimitInternalReq) (*detail.MerchantMetricPlanLimitDetail, error) {
	utility.Assert(req.MerchantId > 0, "invalid merchantId")
	utility.Assert(req.MetricPlanLimitId > 0, "invalid MetricPlanLimitId")
	utility.Assert(req.MetricLimit > 0, "invalid MetricLimit")
	var one *entity.MerchantMetricPlanLimit
	err := dao.MerchantMetricPlanLimit.Ctx(ctx).
		Where(dao.MerchantMetricPlanLimit.Columns().MerchantId, req.MerchantId).
		Where(dao.MerchantMetricPlanLimit.Columns().Id, req.MetricPlanLimitId).
		Where(dao.MerchantMetricPlanLimit.Columns().IsDeleted, 0).
		Scan(&one)
	utility.AssertError(err, "server error")
	utility.Assert(one != nil, "metric limit not found")
	_, err = dao.MerchantMetricPlanLimit.Ctx(ctx).Data(g.Map{
		dao.MerchantMetricPlanLimit.Columns().MetricLimit: req.MetricLimit,
		dao.MerchantMetricPlanLimit.Columns().GmtModify:   gtime.Now(),
	}).Where(dao.MerchantMetricPlanLimit.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		g.Log().Errorf(ctx, "EditMerchantMetricPlanLimit Update err:%s", err.Error())
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	one.MetricLimit = req.MetricLimit
	// reload Cache
	MerchantMetricPlanLimitCachedList(ctx, one.MerchantId, req.PlanId, true)
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Metric(%v)", one.Id),
		Content:        "EditPlanLimit",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return &detail.MerchantMetricPlanLimitDetail{
		Id:          one.Id,
		MerchantId:  one.MerchantId,
		MetricId:    one.MetricId,
		Metric:      GetMerchantMetricSimplify(ctx, one.MetricId),
		PlanId:      one.PlanId,
		MetricLimit: one.MetricLimit,
		UpdateTime:  gtime.Now().Timestamp(),
		CreateTime:  one.CreateTime,
	}, nil
}

func DeleteMerchantMetricPlanLimit(ctx context.Context, merchantId uint64, metricPlanLimitId uint64) error {
	utility.Assert(merchantId > 0, "invalid merchantId")
	utility.Assert(metricPlanLimitId > 0, "invalid metricPlanLimitId")
	one := query.GetMerchantMetricPlanLimit(ctx, metricPlanLimitId)
	utility.Assert(one != nil, "metric limit not found")
	_, err := dao.MerchantMetricPlanLimit.Ctx(ctx).Data(g.Map{
		dao.MerchantMetricPlanLimit.Columns().IsDeleted: gtime.Now().Timestamp(),
		dao.MerchantMetricPlanLimit.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantMetricPlanLimit.Columns().Id, one.Id).OmitNil().Update()
	// reload Cache
	MerchantMetricPlanLimitCachedList(ctx, one.MerchantId, one.PlanId, true)
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Metric(%v)", one.Id),
		Content:        "DeletePlanLimit",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return err
}

func HardDeleteMerchantMetricPlanLimit(ctx context.Context, merchantId uint64, metricPlanLimitId uint64) error {
	utility.Assert(merchantId > 0, "invalid merchantId")
	utility.Assert(metricPlanLimitId > 0, "invalid metricPlanLimitId")
	_, err := dao.MerchantMetricPlanLimit.Ctx(ctx).Where(dao.MerchantMetricPlanLimit.Columns().Id, metricPlanLimitId).Delete()
	return err
}

func BulkMetricLimitPlanBindingReplace(ctx context.Context, plan *entity.Plan, params []*bean.PlanMetricLimitParam) error {
	utility.Assert(plan != nil, "invalid plan")
	if len(params) > 0 {
		var oldList []*entity.MerchantMetricPlanLimit
		_ = dao.MerchantMetricPlanLimit.Ctx(ctx).
			Where(dao.MerchantMetricPlanLimit.Columns().MerchantId, plan.MerchantId).
			Where(dao.MerchantMetricPlanLimit.Columns().PlanId, plan.Id).
			Where(dao.MerchantMetricPlanLimit.Columns().IsDeleted, 0).
			Scan(&oldList)
		var oldMap = make(map[uint64]*entity.MerchantMetricPlanLimit)
		for _, old := range oldList {
			oldMap[old.MetricId] = old
		}
		var newMap = make(map[uint64]uint64)
		for _, ml := range params {
			utility.Assert(ml.MetricId > 0, "invalid metricId")
			utility.Assert(ml.MetricLimit > 0, "invalid MetricLimit")
			me := query.GetMerchantMetric(ctx, ml.MetricId)
			utility.Assert(me != nil, "metric not found")
			utility.Assert(me.Type == MetricTypeLimitMetered, "metric type invalid")

			if old, ok := oldMap[ml.MetricId]; ok {
				//edit
				delete(oldMap, ml.MetricId)
				if old.MetricLimit != ml.MetricLimit {
					//need update
					_, err := dao.MerchantMetricPlanLimit.Ctx(ctx).Data(g.Map{
						dao.MerchantMetricPlanLimit.Columns().MetricLimit: ml.MetricLimit,
						dao.MerchantMetricPlanLimit.Columns().GmtModify:   gtime.Now(),
					}).Where(dao.MerchantMetricPlanLimit.Columns().Id, old.Id).Update()
					if err == nil {
						newMap[old.Id] = ml.MetricLimit
					}
				}
			} else {
				//create
				one := &entity.MerchantMetricPlanLimit{
					MerchantId:  plan.MerchantId,
					MetricId:    ml.MetricId,
					PlanId:      plan.Id,
					MetricLimit: ml.MetricLimit,
					CreateTime:  gtime.Now().Timestamp(),
				}
				result, _ := dao.MerchantMetricPlanLimit.Ctx(ctx).Data(one).OmitNil().Insert(one)
				if result != nil {
					id, _ := result.LastInsertId()
					one.Id = uint64(uint(id))
					newMap[one.Id] = ml.MetricLimit
				}
			}
		}
		// delete other all
		for _, other := range oldMap {
			_, _ = dao.MerchantMetricPlanLimit.Ctx(ctx).Data(g.Map{
				dao.MerchantMetricPlanLimit.Columns().IsDeleted: gtime.Now().Timestamp(),
				dao.MerchantMetricPlanLimit.Columns().GmtModify: gtime.Now(),
			}).Where(dao.MerchantMetricPlanLimit.Columns().Id, other.Id).Update()
		}

		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     plan.MerchantId,
			Target:         fmt.Sprintf("MetricLimitJson(%s)", utility.MarshalToJsonString(newMap)),
			Content:        "OverrideMetricLimit",
			UserId:         0,
			SubscriptionId: "",
			InvoiceId:      "",
			PlanId:         plan.Id,
			DiscountCode:   "",
		}, nil)
		// reload Cache
		MerchantMetricPlanLimitCachedList(ctx, plan.MerchantId, plan.Id, true)
	}
	return nil
}
