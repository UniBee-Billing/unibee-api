// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PlanDao is the data access object for table plan.
type PlanDao struct {
	table   string      // table is the underlying table name of the DAO.
	group   string      // group is the database configuration group name of current DAO.
	columns PlanColumns // columns contains all the column names of Table for convenient usage.
}

// PlanColumns defines and stores column names for table plan.
type PlanColumns struct {
	Id                        string //
	GmtCreate                 string // create time
	GmtModify                 string // update time
	CompanyId                 string // company id
	MerchantId                string // merchant id
	PlanName                  string // PlanName
	Amount                    string // amount, cent, without tax
	Currency                  string // currency
	IntervalUnit              string // period unit,day|month|year|week
	IntervalCount             string // period unit count
	Description               string // description
	ImageUrl                  string // image_url
	HomeUrl                   string // home_url
	GatewayProductName        string // gateway product name
	GatewayProductDescription string // gateway product description
	TaxPercentage             string // taxPercentage 1000 = 10%
	TaxInclusive              string // deperated
	Type                      string // type，1-main plan，2-recurring addon plan 3-onetime addon plan
	Status                    string // status，1-editing，2-active，3-inactive，4-expired
	IsDeleted                 string // 0-UnDeleted，1-Deleted
	BindingAddonIds           string // binded recurring addon planIds，split with ,
	BindingOnetimeAddonIds    string // binded onetime addon planIds，split with ,
	PublishStatus             string // 1-UnPublish,2-Publish, Use For Display Plan At UserPortal
	CreateTime                string // create utc time
	ExtraMetricData           string //
	MetaData                  string // meta_data(json)
	GasPayer                  string // who pay the gas, merchant|user
	TrialAmount               string // amount of trial, 0 for free
	TrialDurationTime         string // duration of trial
	TrialDemand               string //
	CancelAtTrialEnd          string // whether cancel at subscripiton first trial end，0-false | 1-true, will pass to cancelAtPeriodEnd of subscription
	ExternalPlanId            string // external_user_id
	ProductId                 string // product id
	DisableAutoCharge         string // disable auto-charge, 0-false,1-true
	MetricCharge              string // metric charge(json)
}

// planColumns holds the columns for table plan.
var planColumns = PlanColumns{
	Id:                        "id",
	GmtCreate:                 "gmt_create",
	GmtModify:                 "gmt_modify",
	CompanyId:                 "company_id",
	MerchantId:                "merchant_id",
	PlanName:                  "plan_name",
	Amount:                    "amount",
	Currency:                  "currency",
	IntervalUnit:              "interval_unit",
	IntervalCount:             "interval_count",
	Description:               "description",
	ImageUrl:                  "image_url",
	HomeUrl:                   "home_url",
	GatewayProductName:        "gateway_product_name",
	GatewayProductDescription: "gateway_product_description",
	TaxPercentage:             "tax_percentage",
	TaxInclusive:              "tax_inclusive",
	Type:                      "type",
	Status:                    "status",
	IsDeleted:                 "is_deleted",
	BindingAddonIds:           "binding_addon_ids",
	BindingOnetimeAddonIds:    "binding_onetime_addon_ids",
	PublishStatus:             "publish_status",
	CreateTime:                "create_time",
	ExtraMetricData:           "extra_metric_data",
	MetaData:                  "meta_data",
	GasPayer:                  "gas_payer",
	TrialAmount:               "trial_amount",
	TrialDurationTime:         "trial_duration_time",
	TrialDemand:               "trial_demand",
	CancelAtTrialEnd:          "cancel_at_trial_end",
	ExternalPlanId:            "external_plan_id",
	ProductId:                 "product_id",
	DisableAutoCharge:         "disable_auto_charge",
	MetricCharge:              "metric_charge",
}

// NewPlanDao creates and returns a new DAO object for table data access.
func NewPlanDao() *PlanDao {
	return &PlanDao{
		group:   "default",
		table:   "plan",
		columns: planColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PlanDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PlanDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PlanDao) Columns() PlanColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PlanDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PlanDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PlanDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
