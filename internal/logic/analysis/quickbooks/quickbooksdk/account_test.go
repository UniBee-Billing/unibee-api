package quickbooksdk

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	jsonFile, err := os.Open("data/testing/account.json")
	require.NoError(t, err)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)

	var r struct {
		Account Account
		Time    Date
	}
	err = json.Unmarshal(byteValue, &r)
	require.NoError(t, err)

	assert.Equal(t, "MyJobs", r.Account.FullyQualifiedName)
	assert.Equal(t, "MyJobs", r.Account.Name)
	assert.Equal(t, "Asset", r.Account.Classification)
	assert.Equal(t, "AccountsReceivable", r.Account.AccountSubType)
	assert.Equal(t, json.Number("0"), r.Account.CurrentBalanceWithSubAccounts)
	assert.Equal(t, "2014-12-31T09:29:05-08:00", r.Account.MetaData.CreateTime.String())
	assert.Equal(t, "2014-12-31T09:29:05-08:00", r.Account.MetaData.LastUpdatedTime.String())
	assert.Equal(t, AccountsReceivableAccountType, r.Account.AccountType)
	assert.Equal(t, json.Number("0"), r.Account.CurrentBalance)
	assert.True(t, r.Account.Active)
	assert.Equal(t, "0", r.Account.SyncToken)
	assert.Equal(t, "94", r.Account.Id)
	assert.False(t, r.Account.SubAccount)
}
