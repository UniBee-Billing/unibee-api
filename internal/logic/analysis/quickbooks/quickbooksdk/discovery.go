package quickbooksdk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type DiscoveryAPI struct {
	Issuer                string `json:"issuer"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	UserinfoEndpoint      string `json:"userinfo_endpoint"`
	RevocationEndpoint    string `json:"revocation_endpoint"`
	JwksUri               string `json:"jwks_uri"`
}

// CallDiscoveryAPI
// See https://developer.intuit.com/app/developer/qbo/docs/develop/authentication-and-authorization/openid-connect#discovery-document
func CallDiscoveryAPI(discoveryEndpoint EndpointUrl) (*DiscoveryAPI, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", string(discoveryEndpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create req: %v", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make req: %v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %v", err)
	}

	respData := DiscoveryAPI{}
	if err = json.Unmarshal(body, &respData); err != nil {
		return nil, fmt.Errorf("error getting DiscoveryAPIResponse: %v", err)
	}

	return &respData, nil
}
