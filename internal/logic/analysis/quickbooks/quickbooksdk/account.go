package quickbooksdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"unibee/api/bean/detail"
)

const (
	BankAccountType                  = "Bank"
	OtherCurrentAssetAccountType     = "Other Current Asset"
	FixedAssetAccountType            = "Fixed Asset"
	OtherAssetAccountType            = "Other Asset"
	AccountsReceivableAccountType    = "Accounts Receivable"
	EquityAccountType                = "Equity"
	ExpenseAccountType               = "Expense"
	OtherExpenseAccountType          = "Other Expense"
	CostOfGoodsSoldAccountType       = "Cost of Goods Sold"
	AccountsPayableAccountType       = "Accounts Payable"
	CreditCardAccountType            = "Credit Card"
	LongTermLiabilityAccountType     = "Long Term Liability"
	OtherCurrentLiabilityAccountType = "Other Current Liability"
	IncomeAccountType                = "Income"
	OtherIncomeAccountType           = "Other Income"
)

type Account struct {
	Id                            string         `json:"Id,omitempty"`
	Name                          string         `json:"Name,omitempty"`
	SyncToken                     string         `json:"SyncToken,omitempty"`
	AcctNum                       string         `json:"AcctNum,omitempty"`
	CurrencyRef                   *ReferenceType `json:"CurrencyRef,omitempty"`
	ParentRef                     *ReferenceType `json:"ParentRef,omitempty"`
	Description                   string         `json:"Description,omitempty"`
	Active                        bool           `json:"Active,omitempty"`
	MetaData                      *MetaData      `json:"MetaData,omitempty"`
	SubAccount                    bool           `json:"SubAccount,omitempty"`
	Classification                string         `json:"Classification,omitempty"`
	FullyQualifiedName            string         `json:"FullyQualifiedName,omitempty"`
	TxnLocationType               string         `json:"TxnLocationType,omitempty"`
	AccountType                   string         `json:"AccountType,omitempty"`
	CurrentBalanceWithSubAccounts json.Number    `json:"CurrentBalanceWithSubAccounts,omitempty"`
	AccountAlias                  string         `json:"AccountAlias,omitempty"`
	TaxCodeRef                    *ReferenceType `json:"TaxCodeRef,omitempty"`
	AccountSubType                string         `json:"AccountSubType,omitempty"`
	CurrentBalance                json.Number    `json:"CurrentBalance,omitempty"`
}

// CreateAccount creates the given account within QuickBooks
func (c *Client) CreateAccount(account *Account) (*Account, error) {
	var resp struct {
		Account Account
		Time    Date
	}

	if err := c.post("account", account, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Account, nil
}

// FindAccounts gets the full list of Accounts in the QuickBooks account.
func (c *Client) FindAccounts() ([]Account, error) {
	var resp struct {
		QueryResponse struct {
			Accounts      []Account `json:"Account"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM Account", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no accounts could be found")
	}

	accounts := make([]Account, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM Account ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.Accounts == nil {
			return nil, errors.New("no accounts could be found")
		}

		accounts = append(accounts, resp.QueryResponse.Accounts...)
	}

	return accounts, nil
}

// FindAccountById returns an account with a given Id.
func (c *Client) FindAccountById(id string) (*Account, error) {
	var resp struct {
		Account Account
		Time    Date
	}

	if err := c.get("account/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Account, nil
}

// QueryAccounts accepts an SQL query and returns all accounts found using it
func (c *Client) QueryAccounts(query string) ([]Account, error) {
	var resp struct {
		QueryResponse struct {
			Accounts      []Account `json:"Account"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.Accounts == nil {
		return make([]Account, 0), nil
	}

	return resp.QueryResponse.Accounts, nil
}

// UpdateAccount updates the account
func (c *Client) UpdateAccount(account *Account) (*Account, error) {
	if account.Id == "" {
		return nil, errors.New("missing account id")
	}

	existingAccount, err := c.FindAccountById(account.Id)
	if err != nil {
		return nil, err
	}

	account.SyncToken = existingAccount.SyncToken

	payload := struct {
		*Account
		Sparse bool `json:"sparse"`
	}{
		Account: account,
		Sparse:  true,
	}

	var accountData struct {
		Account Account
		Time    Date
	}

	if err = c.post("account", payload, &accountData, nil); err != nil {
		return nil, err
	}

	return &accountData.Account, err
}

func (c *Client) FindOrCreateBankAccount(gateway *detail.Gateway) (*Account, error) {
	var name = fmt.Sprintf("UniBee %s %d", gateway.Name, gateway.Id)

	query := "SELECT * FROM Account WHERE Name = '" + name + "' AND AccountType = '" + BankAccountType + "'"
	existingAccounts, err := c.QueryAccounts(query)
	if err != nil {
		return nil, err
	}

	if len(existingAccounts) > 0 {
		return &existingAccounts[0], nil
	}

	newAccount, err := c.CreateAccount(&Account{
		Name:           name,
		AcctNum:        fmt.Sprintf("%d", gateway.Id),
		Description:    "Used For Receive UniBee Payments",
		AccountType:    BankAccountType,
		AccountSubType: "Checking",
	})
	if err != nil {
		return nil, err
	}

	return newAccount, nil
}
