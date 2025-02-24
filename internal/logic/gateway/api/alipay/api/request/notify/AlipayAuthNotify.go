package notify

type AlipayAuthNotify struct {
	AlipayNotify
	AuthorizationNotifyType string `json:"authorizationNotifyType,omitempty"`
	AuthClientId            string `json:"authClientId,omitempty"`
	AccessToken             string `json:"accessToken,omitempty"`
	AuthState               string `json:"authState,omitempty"`
	AuthCode                string `json:"authCode,omitempty"`
	Reason                  string `json:"reason,omitempty"`
	UserLoginId             string `json:"userLoginId,omitempty"`
	UserId                  string `json:"userId,omitempty"`
}
