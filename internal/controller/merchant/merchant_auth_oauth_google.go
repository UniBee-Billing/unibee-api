package merchant

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unibee/api/merchant/auth"
	"unibee/internal/cmd/config"
	"unibee/internal/logic/jwt"
	"unibee/utility"
)

// GoogleTokenResponse represents the response from Google token exchange
type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token,omitempty"`
	Error        string `json:"error,omitempty"`
	ErrorDesc    string `json:"error_description,omitempty"`
}

// GoogleUserResponse represents the response from Google user API
type GoogleUserResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func (c *ControllerAuth) OauthGoogle(ctx context.Context, req *auth.OauthGoogleReq) (res *auth.OauthGoogleRes, err error) {
	// Validate Google code
	utility.Assert(len(req.GoogleCode) > 0, "Google authorization code is required")

	// Get Google OAuth configuration
	utility.Assert(len(config.GetConfigInstance().OAuth.GoogleClientId) > 0, "Google client ID not configured")
	utility.Assert(len(config.GetConfigInstance().OAuth.GoogleClientSecret) > 0, "Google client secret not configured")

	// Exchange authorization code for access token
	accessToken, err := exchangeGoogleCode(ctx, req.GoogleCode, config.GetConfigInstance().OAuth.GoogleClientId, config.GetConfigInstance().OAuth.GoogleClientSecret, req.RedirectUri)
	utility.AssertError(err, "Error exchanging Google code")

	// Get Google user information
	userInfo, err := getGoogleUserInfo(ctx, accessToken)
	utility.AssertError(err, "Error getting Google user info")

	// Convert to API response format
	googleUserInfo := &auth.GoogleUserInfo{
		ID:            userInfo.ID,
		Email:         userInfo.Email,
		VerifiedEmail: userInfo.VerifiedEmail,
		Name:          userInfo.Name,
		GivenName:     userInfo.GivenName,
		FamilyName:    userInfo.FamilyName,
		Picture:       userInfo.Picture,
		Locale:        userInfo.Locale,
	}

	// Return the user info to frontend
	return &auth.OauthGoogleRes{
		UserInfo: googleUserInfo,
		Sign:     jwt.GenerateOAuthSign(googleUserInfo.Email, googleUserInfo.Name),
	}, nil
}

// exchangeGoogleCode exchanges the authorization code for an access token
func exchangeGoogleCode(ctx context.Context, code, clientID, clientSecret, redirectUri string) (string, error) {
	// Prepare the request data
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	// Google OAuth requires the exact redirect_uri that was used in the authorization request
	// For server-side flow, we should not include redirect_uri if it wasn't used in the auth request
	data.Set("redirect_uri", redirectUri) // This might be causing issues

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", "https://oauth2.googleapis.com/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Send request with longer timeout for Google API
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Google API: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Google API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Log successful response for debugging
	fmt.Printf("Google token exchange successful, response length: %d\n", len(string(body)))

	// Parse response
	var tokenResp GoogleTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse token response: %v", err)
	}

	// Check for Google API errors
	if tokenResp.Error != "" {
		return "", fmt.Errorf("Google API error: %s - %s", tokenResp.Error, tokenResp.ErrorDesc)
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("no access token in response: %s", string(body))
	}

	return tokenResp.AccessToken, nil
}

// getGoogleUserInfo retrieves user information from Google API
func getGoogleUserInfo(ctx context.Context, accessToken string) (*GoogleUserResponse, error) {
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	// Send request with longer timeout for Google API
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Google userinfo API: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Google API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var userResp GoogleUserResponse
	if err := json.Unmarshal(body, &userResp); err != nil {
		return nil, fmt.Errorf("failed to parse user response: %v", err)
	}

	return &userResp, nil
}
