package examples

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"unibee/internal/logic/analysis/quickbooks/quickbooksdk"
)

func TestReuseToken(t *testing.T) {
	clientId := "<your-client-id>"
	clientSecret := "<your-client-secret>"
	realmId := "<realm-id>"

	token := quickbooksdk.BearerToken{
		RefreshToken: "<saved-refresh-token>",
		AccessToken:  "<saved-access-token>",
	}

	qbClient, err := quickbooksdk.NewClient(clientId, clientSecret, realmId, false, "", &token)
	require.NoError(t, err)

	// Make a request!
	info, err := qbClient.FindCompanyInfo()
	require.NoError(t, err)
	fmt.Println(info)
}
