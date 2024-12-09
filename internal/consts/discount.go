package consts

type DiscountStatusEnum int
type DiscountBillingTypeEnum int
type DiscountTypeEnum int

const (
	DiscountStatusEditable       = 1
	DiscountStatusActive         = 2
	DiscountStatusDeactivate     = 3
	DiscountStatusExpired        = 4
	DiscountStatusArchived       = 10
	DiscountBillingTypeOnetime   = 1
	DiscountBillingTypeRecurring = 2
	DiscountTypePercentage       = 1
	DiscountTypeFixedAmount      = 2
)

func (status DiscountStatusEnum) Description() string {
	switch status {
	case DiscountStatusEditable:
		return "Editable"
	case DiscountStatusActive:
		return "Active"
	case DiscountStatusDeactivate:
		return "Deactivate"
	case DiscountStatusExpired:
		return "Expired"
	case DiscountStatusArchived:
		return "Archived"
	default:
		return "Active"
	}
}

func DiscountStatusToEnum(status int) DiscountStatusEnum {
	switch status {
	case DiscountStatusEditable:
		return DiscountStatusEditable
	case DiscountStatusActive:
		return DiscountStatusActive
	case DiscountStatusDeactivate:
		return DiscountStatusDeactivate
	case DiscountStatusExpired:
		return DiscountStatusExpired
	case DiscountStatusArchived:
		return DiscountStatusArchived
	default:
		return DiscountStatusActive
	}
}

func (status DiscountTypeEnum) Description() string {
	switch status {
	case DiscountTypePercentage:
		return "Percentage"
	case DiscountTypeFixedAmount:
		return "FixedAmount"
	default:
		return "Percentage"
	}
}

func DiscountTypeToEnum(one int) DiscountTypeEnum {
	switch one {
	case DiscountTypePercentage:
		return DiscountTypePercentage
	case DiscountTypeFixedAmount:
		return DiscountTypeFixedAmount
	default:
		return DiscountTypePercentage
	}
}

func (status DiscountBillingTypeEnum) Description() string {
	switch status {
	case DiscountBillingTypeOnetime:
		return "OneTime"
	case DiscountBillingTypeRecurring:
		return "Recurring"
	default:
		return "OneTime"
	}
}

func DiscountBillingTypeToEnum(one int) DiscountBillingTypeEnum {
	switch one {
	case DiscountBillingTypeOnetime:
		return DiscountBillingTypeOnetime
	case DiscountBillingTypeRecurring:
		return DiscountBillingTypeRecurring
	default:
		return DiscountBillingTypeOnetime
	}
}
