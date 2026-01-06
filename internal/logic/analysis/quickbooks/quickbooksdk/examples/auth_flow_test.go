package examples

import (
	"fmt"
	"testing"
	"unibee/internal/logic/analysis/quickbooks/quickbooksdk"

	"github.com/stretchr/testify/require"
)

func TestAuthorizationFlow(t *testing.T) {
	clientId := "<your-client-id>"
	clientSecret := "<your-client-secret>"
	realmId := "<realm-id>"

	qbClient, err := quickbooksdk.NewClient(clientId, clientSecret, realmId, false, "", nil)
	require.NoError(t, err)

	// To do first when you receive the authorization code from quickbooks callback
	authorizationCode := "<received-from-callback>"
	redirectURI := "https://developer.intuit.com/v2/OAuth2Playground/RedirectUrl"
	bearerToken, err := qbClient.RetrieveBearerToken(authorizationCode, redirectURI)
	require.NoError(t, err)
	// Save the bearer token inside a db

	// When the token expire, you can use the following function
	bearerToken, err = qbClient.RefreshToken(bearerToken.RefreshToken)
	require.NoError(t, err)

	// Make a request!
	info, err := qbClient.FindCompanyInfo()
	require.NoError(t, err)
	fmt.Println(info)

	// Revoke the token, this should be done only if a user unsubscribe from your app
	require.NoError(t, qbClient.RevokeToken(bearerToken.RefreshToken))
}
