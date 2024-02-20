package metric

import (
	"context"
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
	MetricTypeLimitMetered    = 1
	MetricTypeChargeMetered   = 2
	MetricTypeChargeRecurring = 3
)

func GetMerchantMetricVo(ctx context.Context, id int64) *ro.MerchantMetricVo {
	one := query.GetMerchantMetric(ctx, id)
	if one != nil {
		return &ro.MerchantMetricVo{
			Id:                  one.Id,
			MerchantId:          one.MerchantId,
			Code:                one.Code,
			MetricName:          one.MetricName,
			MetricDescription:   one.MetricDescription,
			Type:                one.Type,
			AggregationType:     one.AggregationType,
			AggregationProperty: one.AggregationProperty,
			UpdateTime:          one.GmtModify.Timestamp(),
			CreateTime:          one.CreateTime,
		}
	}
	return nil
}

func MerchantMetricList(ctx context.Context, merchantId int64) []*ro.MerchantMetricVo {
	utility.Assert(merchantId > 0, "invalid merchantId")
	var list = make([]*ro.MerchantMetricVo, 0)
	if merchantId > 0 {
		var entities []*entity.MerchantMetric
		err := dao.MerchantMetric.Ctx(ctx).
			Where(entity.MerchantMetric{MerchantId: merchantId}).
			Where(entity.MerchantMetric{IsDeleted: 0}).
			Scan(&entities)
		if err == nil && len(entities) > 0 {
			for _, one := range entities {
				list = append(list, &ro.MerchantMetricVo{
					Id:                  one.Id,
					MerchantId:          one.MerchantId,
					Code:                one.Code,
					MetricName:          one.MetricName,
					MetricDescription:   one.MetricDescription,
					Type:                one.Type,
					AggregationType:     one.AggregationType,
					AggregationProperty: one.AggregationProperty,
					UpdateTime:          one.GmtModify.Timestamp(),
					CreateTime:          one.CreateTime,
				})
			}
		}
	}
	return list
}

type NewMerchantMetricInternalReq struct {
	MerchantId          int64  `p:"merchantId" dc:"MerchantId" v:"required"`
	Code                string `p:"code" dc:"Code" v:"required"`
	Name                string `p:"name" dc:"Name" v:"required"`
	Description         string `p:"description" dc:"Description"`
	AggregationType     int    `p:"aggregationType" dc:"AggregationType,0-count，1-count unique, 2-latest, 3-max, 4-sum"`
	AggregationProperty string `p:"aggregationProperty" dc:"AggregationProperty, Will Needed When AggregationType != count"`
}

func NewMerchantMetric(ctx context.Context, req *NewMerchantMetricInternalReq) (*ro.MerchantMetricVo, error) {
	utility.Assert(req.MerchantId > 0, "invalid merchantId")
	utility.Assert(len(req.Code) > 0, "code is nil")

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

	return &ro.MerchantMetricVo{
		Id:                  one.Id,
		MerchantId:          one.MerchantId,
		Code:                one.Code,
		MetricName:          one.MetricName,
		MetricDescription:   one.MetricDescription,
		Type:                one.Type,
		AggregationType:     one.AggregationType,
		AggregationProperty: one.AggregationProperty,
		CreateTime:          one.CreateTime,
	}, nil
}

func EditMerchantMetric(ctx context.Context, merchantId int64, metricId int64, name string, description string) (*ro.MerchantMetricVo, error) {
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

	return &ro.MerchantMetricVo{
		Id:                  one.Id,
		MerchantId:          one.MerchantId,
		Code:                one.Code,
		MetricName:          one.MetricName,
		MetricDescription:   one.MetricDescription,
		Type:                one.Type,
		AggregationType:     one.AggregationType,
		AggregationProperty: one.AggregationProperty,
		UpdateTime:          gtime.Now().Timestamp(),
		CreateTime:          one.CreateTime,
	}, nil
}

func DeleteMerchantMetric(ctx context.Context, merchantId int64, metricId int64) error {
	utility.Assert(merchantId > 0, "invalid merchantId")
	utility.Assert(metricId > 0, "invalid metricId")
	one := query.GetMerchantMetric(ctx, metricId)
	utility.Assert(one != nil, "endpoint not found")
	_, err := dao.MerchantMetric.Ctx(ctx).Data(g.Map{
		dao.MerchantMetric.Columns().IsDeleted: 1,
		dao.MerchantMetric.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantMetric.Columns().Id, one.Id).OmitNil().Update()
	return err
}
