package model

type SettlementInfo struct {
	SettlementCurrency    string                 `json:"settlementCurrency"`
	SettlementBankAccount *SettlementBankAccount `json:"settlementBankAccount"`
}

type SettlementBankAccount struct {
	BankAccountNo        string            `json:"bankAccountNo"`
	AccountHolderName    string            `json:"accountHolderName"`
	SwiftCode            string            `json:"swiftCode"`
	BankRegion           string            `json:"bankRegion"`
	AccountHolderType    AccountHolderType `json:"accountHolderType"`
	RoutingNumber        string            `json:"routingNumber"`
	BranchCode           string            `json:"branchCode"`
	AccountHolderTIN     string            `json:"accountHolderTIN"`
	AccountType          AccountType       `json:"accountType"`
	BankName             string            `json:"bankName"`
	AccountHolderAddress *Address          `json:"accountHolderAddress"`
	Iban                 string            `json:"iban"`
}

type AccountHolderType string

const (
	AccountHolderType_INDIVIDUAL AccountHolderType = "INDIVIDUAL"
	AccountHolderType_ENTERPRISE AccountHolderType = "ENTERPRISE"
)

type AccountType string

const (
	AccountType_CHECKING      AccountType = "CHECKING"
	AccountType_FIXED_DEPOSIT AccountType = "FIXED_DEPOSIT"
)
