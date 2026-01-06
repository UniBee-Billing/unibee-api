package quickbooksdk

import (
	"context"
	"testing"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"
)

func TestToken(t *testing.T) {
	ctx := context.Background()
	qbClient, err := NewClient("", "", "9341454732969707", false, "", nil)

	// To do first when you receive the authorization code from quickbooks callback
	bearerToken, err := qbClient.RetrieveBearerToken("", "https://api.unibee.top/integrate/quickbooks/auth_back")
	if err != nil {
		g.Log().Errorf(ctx, "RetrieveBearerToken Error: %s", err.Error())
	} else {
		g.Log().Infof(ctx, "RetrieveBearerToken Success:%s", utility.MarshalToJsonString(bearerToken))
	}
}
