package metric

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee-api/internal/dao/oversea_pay"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

const (
	MetricTypeLimitMetered    = 1
	MetricTypeChargeMetered   = 2
	MetricTypeChargeRecurring = 3
)

type MerchantMetricVo struct {
	Id                  uint64 `json:"id"            description:"id"`                                                                                // id
	MerchantId          int64  `json:"merchantId"          description:"merchantId"`                                                                  // merchantId
	Code                string `json:"code"                description:"code"`                                                                        // code
	MetricName          string `json:"metricName"          description:"metric name"`                                                                 // metric name
	MetricDescription   string `json:"metricDescription"   description:"metric description"`                                                          // metric description
	Type                int    `json:"type"                description:"1-limit_metered，2-charge_metered(come later),3-charge_recurring(come later)"` // 1-limit_metered，2-charge_metered(come later),3-charge_recurring(come later)
	AggregationType     int    `json:"aggregationType"     description:"0-count，1-count unique, 2-latest, 3-max, 4-sum"`                              // 0-count，1-count unique, 2-latest, 3-max, 4-sum
	AggregationProperty string `json:"aggregationProperty" description:"aggregation property"`
	UpdateTime          int64  `json:"gmtModify"     description:"update time"`     // update time
	CreateTime          int64  `json:"createTime"    description:"create utc time"` // create utc time
}

func MerchantMetricList(ctx context.Context, merchantId int64) []*MerchantMetricVo {
	utility.Assert(merchantId > 0, "invalid merchantId")
	var list = make([]*MerchantMetricVo, 0)
	if merchantId > 0 {
		var entities []*entity.MerchantMetric
		err := dao.MerchantMetric.Ctx(ctx).
			Where(entity.MerchantMetric{MerchantId: merchantId}).
			Where(entity.MerchantMetric{IsDeleted: 0}).
			Scan(&entities)
		if err == nil && len(entities) > 0 {
			for _, one := range entities {
				list = append(list, &MerchantMetricVo{
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

func NewMerchantMetric(ctx context.Context, req *NewMerchantMetricInternalReq) (*MerchantMetricVo, error) {
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

	return &MerchantMetricVo{
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

func EditMerchantMetric(ctx context.Context, merchantId int64, metricId int64, name string, description string) (*MerchantMetricVo, error) {
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

	return &MerchantMetricVo{
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
