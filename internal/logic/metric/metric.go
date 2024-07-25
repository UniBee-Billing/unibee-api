package metric

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/operation_log"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

const (
	MetricTypeLimitMetered           = 1
	MetricTypeChargeMetered          = 2
	MetricTypeChargeRecurring        = 3
	MetricAggregationTypeCount       = 1
	MetricAggregationTypeCountUnique = 2
	MetricAggregationTypeLatest      = 3
	MetricAggregationTypeMax         = 4
	MetricAggregationTypeSum         = 5
)

func GetMerchantMetricSimplify(ctx context.Context, id uint64) *bean.MerchantMetricSimplify {
	one := query.GetMerchantMetric(ctx, id)
	if one != nil {
		return bean.SimplifyMerchantMetric(one)
	}
	return nil
}

func MerchantMetricDetail(ctx context.Context, merchantId uint64, merchantMetricId uint64) *bean.MerchantMetricSimplify {
	utility.Assert(merchantMetricId > 0, "invalid merchantMetricId")
	if merchantMetricId > 0 {
		var one *entity.MerchantMetric
		err := dao.MerchantMetric.Ctx(ctx).
			Where(dao.MerchantMetric.Columns().Id, merchantMetricId).
			Scan(&one)
		if err == nil && one != nil {
			utility.Assert(one.MerchantId == merchantId, "wrong merchant account")
			return bean.SimplifyMerchantMetric(one)
		}
	}
	return nil
}

type NewMerchantMetricInternalReq struct {
	MerchantId          uint64 `json:"merchantId" dc:"MerchantId" v:"required"`
	Code                string `json:"code" dc:"Code" v:"required"`
	Name                string `json:"name" dc:"Name" v:"required"`
	Description         string `json:"description" dc:"Description"`
	AggregationType     int    `json:"aggregationType" dc:"AggregationType,1-count，2-count unique, 3-latest, 4-max, 5-sum"`
	AggregationProperty string `json:"aggregationProperty" dc:"AggregationProperty, Will Needed When AggregationType != count"`
}

func NewMerchantMetric(ctx context.Context, req *NewMerchantMetricInternalReq) (*bean.MerchantMetricSimplify, error) {
	utility.Assert(req.MerchantId > 0, "invalid merchantId")
	utility.Assert(len(req.Code) > 0, "code is nil")
	utility.Assert(req.AggregationType > 0 && req.AggregationType < 6, "aggregationType should be one of 1-count，2-count unique, 3-latest, 4-max, 5-sum")
	if req.AggregationType > 1 {
		//check property should contain
		utility.Assert(len(req.AggregationProperty) > 0, "aggregationProperty should be set when aggregationType not count Type")
	}

	one := query.GetMerchantMetricByCode(ctx, req.Code)
	utility.Assert(one == nil, "metric already exist")
	one = &entity.MerchantMetric{
		MerchantId:          req.MerchantId,
		Code:                req.Code,
		MetricName:          req.Name,
		MetricDescription:   req.Description,
		Type:                MetricTypeLimitMetered, // other type not support now
		AggregationType:     req.AggregationType,
		AggregationProperty: req.AggregationProperty,
		CreateTime:          gtime.Now().Timestamp(),
	}
	result, err := dao.MerchantMetric.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		g.Log().Errorf(ctx, "NewMerchantMetric Insert err:%s", err.Error())
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(id)
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Metric(%v)", one.Id),
		Content:        "New",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return bean.SimplifyMerchantMetric(one), nil
}

func EditMerchantMetric(ctx context.Context, merchantId uint64, metricId uint64, name string, description string) (*bean.MerchantMetricSimplify, error) {
	utility.Assert(merchantId > 0, "invalid merchantId")
	utility.Assert(metricId > 0, "invalid metricId")
	one := query.GetMerchantMetric(ctx, metricId)
	utility.Assert(one != nil, "endpoint not found")
	_, err := dao.MerchantMetric.Ctx(ctx).Data(g.Map{
		dao.MerchantMetric.Columns().MetricName:        name,
		dao.MerchantMetric.Columns().MetricDescription: description,
		dao.MerchantMetric.Columns().GmtModify:         gtime.Now(),
	}).Where(dao.MerchantMetric.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		g.Log().Errorf(ctx, "EditMerchantMetric Update err:%s", err.Error())
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	one.MetricName = name
	one.MetricDescription = description
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Metric(%v)", one.Id),
		Content:        "Edit",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return bean.SimplifyMerchantMetric(one), nil
}

func DeleteMerchantMetric(ctx context.Context, merchantId uint64, metricId uint64) error {
	utility.Assert(merchantId > 0, "invalid merchantId")
	utility.Assert(metricId > 0, "invalid metricId")
	one := query.GetMerchantMetric(ctx, metricId)
	utility.Assert(one != nil, "endpoint not found")
	_, err := dao.MerchantMetric.Ctx(ctx).Data(g.Map{
		dao.MerchantMetric.Columns().IsDeleted: gtime.Now().Timestamp(),
		dao.MerchantMetric.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantMetric.Columns().Id, one.Id).OmitNil().Update()
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Metric(%v)", one.Id),
		Content:        "Delete",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return err
}

func HardDeleteMerchantMetric(ctx context.Context, merchantId uint64, metricId uint64) error {
	utility.Assert(merchantId > 0, "invalid merchantId")
	utility.Assert(metricId > 0, "invalid metricId")
	_, err := dao.MerchantMetric.Ctx(ctx).Where(dao.MerchantMetric.Columns().Id, metricId).Delete()
	return err
}
