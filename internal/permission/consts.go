package permission

type PermissionGroup string

const (
	PermissionGroupConfig          = "App Config"
	PermissionGroupEmailTemplate   = "Email Template"
	PermissionGroupInvoiceTemplate = "Invoice Template"
	PermissionGroupPlan            = "Plan"
	PermissionGroupSubscription    = "Subscription"
	PermissionGroupInvoice         = "Invoice"
	PermissionGroupCustomer        = "Customer"
	PermissionGroupAnalytics       = "Analytics"
)

var PermissionGroupList = []PermissionGroup{
	PermissionGroupConfig,
	PermissionGroupEmailTemplate,
	PermissionGroupInvoiceTemplate,
	PermissionGroupPlan,
	PermissionGroupSubscription,
	PermissionGroupInvoice,
	PermissionGroupCustomer,
	PermissionGroupAnalytics,
}

func IsValidGroup(target string) bool {
	if len(target) <= 0 {
		return false
	}
	for _, event := range PermissionGroupList {
		if event == PermissionGroup(target) {
			return true
		}
	}
	return false
}

type PermissionType string

const (
	PermissionTypeRead   = "Read"
	PermissionTypeWrite  = "Write"
	PermissionTypeExport = "Export"
)
