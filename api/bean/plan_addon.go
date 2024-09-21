package bean

type PlanProductParam struct {
	Name        string `json:"name" dc:"Name" `
	Description string `json:"description" dc:"Description"`
}

type PlanAddonParam struct {
	Quantity    int64  `json:"quantity" dc:"Quantity，Default 1" `
	AddonPlanId uint64 `json:"addonPlanId" dc:"AddonPlanId"`
}

type PlanAddonDetail struct {
	Quantity  int64 `json:"quantity" dc:"Quantity" `
	AddonPlan *Plan `json:"addonPlan" dc:"addonPlan" `
}
