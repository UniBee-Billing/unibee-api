package bean

type ExternalDiscountParam struct {
	Recurring          *bool             `json:"recurring"       dc:"Discount recurring enable, default false"`
	DiscountAmount     *int64            `json:"discountAmount"     dc:"Amount of discount"`
	DiscountPercentage *int64            `json:"discountPercentage" dc:"Percentage of discount, 100=1%, ignore if discountAmount set"`
	CycleLimit         *int              `json:"cycleLimit"         dc:"the count limitation of subscription recurring cycle, recurring need enable if cycleLimit set"`
	EndTime            *int64            `json:"endTime"            dc:"end of discount available utc time"`
	Metadata           map[string]string `json:"metadata" dc:"Metadata，Map"`
}
