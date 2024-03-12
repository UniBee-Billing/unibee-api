package bean

type PlanAddonParam struct {
	Quantity    int64  `json:"quantity" dc:"Quantity，Default 1" `
	AddonPlanId uint64 `json:"addonPlanId" dc:"AddonPlanId"`
}

type PlanAddonDetail struct {
	Quantity  int64         `json:"quantity" dc:"Quantity" `
	AddonPlan *PlanSimplify `json:"addonPlan" dc:"addonPlan" `
}
