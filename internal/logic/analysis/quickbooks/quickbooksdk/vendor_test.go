package quickbooksdk

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVendor(t *testing.T) {
	jsonFile, err := os.Open("data/testing/vendor.json")
	require.NoError(t, err)
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)

	var resp struct {
		Vendor Vendor
		Time   Date
	}

	require.NoError(t, json.Unmarshal(byteValue, &resp))
	assert.NotNil(t, resp.Vendor.PrimaryEmailAddr)
	assert.False(t, resp.Vendor.Vendor1099)
	assert.Equal(t, "Bessie", resp.Vendor.GivenName)
	assert.Equal(t, "Books by Bessie", resp.Vendor.DisplayName)
	assert.NotNil(t, resp.Vendor.BillAddr)
	assert.Equal(t, "0", resp.Vendor.SyncToken)
	assert.Equal(t, "Books by Bessie", resp.Vendor.PrintOnCheckName)
	assert.Equal(t, "Williams", resp.Vendor.FamilyName)
	assert.NotNil(t, resp.Vendor.PrimaryPhone)
	assert.Equal(t, "1345", resp.Vendor.AcctNum)
	assert.Equal(t, "Books by Bessie", resp.Vendor.CompanyName)
	assert.NotNil(t, resp.Vendor.WebAddr)
	assert.True(t, resp.Vendor.Active)
	assert.Equal(t, "0", resp.Vendor.Balance.String())
	assert.Equal(t, "30", resp.Vendor.Id)
	assert.Equal(t, "2014-09-12T10:07:56-07:00", resp.Vendor.MetaData.CreateTime.String())
	assert.Equal(t, "2014-09-17T11:13:46-07:00", resp.Vendor.MetaData.LastUpdatedTime.String())
}
