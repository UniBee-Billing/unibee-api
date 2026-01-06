package quickbooksdk

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAttachable(t *testing.T) {
	jsonFile, err := os.Open("data/testing/attachable.json")
	require.NoError(t, err)
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)

	var r struct {
		Attachable Attachable
		Time       Date
	}

	require.NoError(t, json.Unmarshal(byteValue, &r))

	assert.Equal(t, "0", r.Attachable.SyncToken)
	assert.False(t, r.Attachable.AttachableRef[0].IncludeOnSend)
	assert.Equal(t, "95", r.Attachable.AttachableRef[0].EntityRef.Value)
	assert.Equal(t, "This is an attached note.", r.Attachable.Note)
	assert.Equal(t, "200900000000000008541", r.Attachable.Id)
	assert.Equal(t, "2015-11-17T11:05:15-08:00", r.Attachable.MetaData.CreateTime.String())
	assert.Equal(t, "2015-11-17T11:05:15-08:00", r.Attachable.MetaData.LastUpdatedTime.String())
}
