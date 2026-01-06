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

// GitHubTokenResponse represents the response from GitHub token exchange
type GitHubTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// GitHubUserResponse represents the response from GitHub user API
type GitHubUserResponse struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Company   string `json:"company"`
	Location  string `json:"location"`
	Bio       string `json:"bio"`
}

func (c *ControllerAuth) OauthGithub(ctx context.Context, req *auth.OauthGithubReq) (res *auth.OauthGithubRes, err error) {
	// Validate GitHub code
	utility.Assert(len(req.GithubCode) > 0, "GitHub authorization code is required")

	// Get GitHub OAuth configuration
	//config := config.GetConfigInstance()
	utility.Assert(len(config.GetConfigInstance().OAuth.GithubClientId) > 0, "GitHub client ID not configured")
	utility.Assert(len(config.GetConfigInstance().OAuth.GithubClientSecret) > 0, "GitHub client secret not configured")

	// Exchange authorization code for access token
	accessToken, err := exchangeGitHubCode(ctx, req.GithubCode, config.GetConfigInstance().OAuth.GithubClientId, config.GetConfigInstance().OAuth.GithubClientSecret)
	utility.AssertError(err, "Error exchanging GitHub code")

	// Get GitHub user information
	userInfo, err := getGitHubUserInfo(ctx, accessToken)
	utility.AssertError(err, "Error getting GitHub user info")

	if len(userInfo.Name) == 0 {
		userInfo.Name = userInfo.Login
	}
	// Convert to API response format
	githubUserInfo := &auth.GitHubUserInfo{
		ID:        userInfo.ID,
		Login:     userInfo.Login,
		Email:     userInfo.Email,
		Name:      userInfo.Name,
		AvatarURL: userInfo.AvatarURL,
		Company:   userInfo.Company,
		Location:  userInfo.Location,
		Bio:       userInfo.Bio,
	}

	// Return the access token and user info to frontend
	return &auth.OauthGithubRes{
		//GithubAccessToken: accessToken,
		UserInfo: githubUserInfo,
		Sign:     jwt.GenerateOAuthSign(githubUserInfo.Email, githubUserInfo.Name),
	}, nil
}

// exchangeGitHubCode exchanges the authorization code for an access token
func exchangeGitHubCode(ctx context.Context, code, clientID, clientSecret string) (string, error) {
	// Prepare the request data
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var tokenResp GitHubTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse token response: %v", err)
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("no access token in response: %s", string(body))
	}

	return tokenResp.AccessToken, nil
}

// getGitHubUserInfo retrieves user information from GitHub API
func getGitHubUserInfo(ctx context.Context, accessToken string) (*GitHubUserResponse, error) {
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "Unibee-OAuth-Client")

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var userResp GitHubUserResponse
	if err := json.Unmarshal(body, &userResp); err != nil {
		return nil, fmt.Errorf("failed to parse user response: %v", err)
	}

	return &userResp, nil
}
