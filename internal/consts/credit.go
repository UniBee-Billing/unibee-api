package consts

const (
	CreditAccountTypeMain  = 1
	CreditAccountTypePromo = 2
)

const (
	CreditRechargeCreated = 10
	CreditRechargeSuccess = 20
	CreditRechargeFailed  = 30
)

type TransactionTypeEnum int

const (
	// Transaction type。1-recharge income，2-payment out，3-refund income，4-withdraw out，5-withdraw failed income, 6-admin change，7-recharge refund out
	CreditTransactionRechargeIncome       = 1
	CreditTransactionPayout               = 2
	CreditTransactionRefundIncome         = 3
	CreditTransactionWithdrawOut          = 4
	CreditTransactionWithdrawFailedIncome = 5
	CreditTransactionAdminChange          = 6
	CreditTransactionRechargeRefundOut    = 7
)

func (transactionType TransactionTypeEnum) Description() string {
	switch transactionType {
	case CreditTransactionRechargeIncome:
		return "RechargeIncome"
	case CreditTransactionPayout:
		return "Payout"
	case CreditTransactionRefundIncome:
		return "RefundIncome"
	case CreditTransactionWithdrawOut:
		return "WithdrawOut"
	case CreditTransactionWithdrawFailedIncome:
		return "WithdrawFailedIncome"
	case CreditTransactionRechargeRefundOut:
		return "RechargeRefundOut"
	default:
		return "AdminChange"
	}
}

func (transactionType TransactionTypeEnum) ExportDescription(amount int64) string {
	switch transactionType {
	case CreditTransactionRechargeIncome:
		return "RechargeIncome"
	case CreditTransactionPayout:
		return "Applied to an invoice"
	case CreditTransactionRefundIncome:
		return "From refund"
	case CreditTransactionWithdrawOut:
		return "WithdrawOut"
	case CreditTransactionWithdrawFailedIncome:
		return "WithdrawFailedIncome"
	case CreditTransactionRechargeRefundOut:
		return "RechargeRefundOut"
	default:
		if amount > 0 {
			return "Added by admin"
		} else {
			return "Reduced by admin"
		}
	}
}

func CreditTransactionTypeToEnum(transactionType int) TransactionTypeEnum {
	switch transactionType {
	case CreditTransactionRechargeIncome:
		return CreditTransactionRechargeIncome
	case CreditTransactionPayout:
		return CreditTransactionPayout
	case CreditTransactionRefundIncome:
		return CreditTransactionRefundIncome
	case CreditTransactionWithdrawOut:
		return CreditTransactionWithdrawOut
	case CreditTransactionWithdrawFailedIncome:
		return CreditTransactionWithdrawFailedIncome
	case CreditTransactionRechargeRefundOut:
		return CreditTransactionRechargeRefundOut
	default:
		return CreditTransactionAdminChange
	}
}
