package consts

type PermissionType string

const (
	PermissionAccess   PermissionType = "access"
	PermissionRead     PermissionType = "read"
	PermissionWrite    PermissionType = "write"
	PermissionDownload PermissionType = "download"
)

type PermissionTypeGroup string

type Group struct {
	Group            PermissionTypeGroup
	DependencyGroups []PermissionTypeGroup
}

const (
	PermissionGroupPlan           = PermissionTypeGroup("plan")
	PermissionGroupBillableMetric = PermissionTypeGroup("billable-metric")
	PermissionGroupDiscountCode   = PermissionTypeGroup("discount-code")
	PermissionGroupSubscription   = PermissionTypeGroup("subscription")
	PermissionGroupInvoice        = PermissionTypeGroup("invoice")
	PermissionGroupTransaction    = PermissionTypeGroup("transaction")
	PermissionGroupUser           = PermissionTypeGroup("user")
	PermissionGroupAdmin          = PermissionTypeGroup("admin")
	PermissionGroupMyAccount      = PermissionTypeGroup("my-account")
	PermissionGroupReport         = PermissionTypeGroup("report")
	PermissionGroupConfiguration  = PermissionTypeGroup("configuration")
	PermissionGroupActivityLogs   = PermissionTypeGroup("activity-logs")
	PermissionGroupAnalytics      = PermissionTypeGroup("analytics")
)

var PermissionGroupMap map[PermissionTypeGroup]Group = map[PermissionTypeGroup]Group{
	PermissionGroupPlan: {
		Group: PermissionGroupPlan,
		DependencyGroups: []PermissionTypeGroup{
			PermissionGroupBillableMetric}},
	PermissionGroupBillableMetric: {
		Group:            PermissionGroupBillableMetric,
		DependencyGroups: []PermissionTypeGroup{}},
	PermissionGroupDiscountCode: {
		Group: PermissionGroupDiscountCode,
		DependencyGroups: []PermissionTypeGroup{
			PermissionGroupPlan}},
	PermissionGroupSubscription: {
		Group: PermissionGroupSubscription,
		DependencyGroups: []PermissionTypeGroup{
			PermissionGroupPlan,
			PermissionGroupDiscountCode,
			PermissionGroupUser,
			PermissionGroupInvoice,
			PermissionGroupTransaction}},
	PermissionGroupInvoice: {
		Group: PermissionGroupInvoice,
		DependencyGroups: []PermissionTypeGroup{
			PermissionGroupUser,
			PermissionGroupPlan}},
	PermissionGroupTransaction: {
		Group: PermissionGroupTransaction,
		DependencyGroups: []PermissionTypeGroup{
			PermissionGroupUser}},
	PermissionGroupUser: {
		Group:            PermissionGroupUser,
		DependencyGroups: []PermissionTypeGroup{}},
	PermissionGroupAdmin: {
		Group:            PermissionGroupAdmin,
		DependencyGroups: []PermissionTypeGroup{}},
	PermissionGroupMyAccount: {
		Group:            PermissionGroupMyAccount,
		DependencyGroups: []PermissionTypeGroup{}},
	PermissionGroupReport: {
		Group:            PermissionGroupReport,
		DependencyGroups: []PermissionTypeGroup{}},
	PermissionGroupConfiguration: {
		Group:            PermissionGroupConfiguration,
		DependencyGroups: []PermissionTypeGroup{}},
	PermissionGroupActivityLogs: {
		Group:            PermissionGroupActivityLogs,
		DependencyGroups: []PermissionTypeGroup{}},
	PermissionGroupAnalytics: {
		Group:            PermissionGroupAnalytics,
		DependencyGroups: []PermissionTypeGroup{}},
}
