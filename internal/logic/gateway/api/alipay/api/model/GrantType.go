package model

type GrantType string

const (
	GrantTypeAUTHORIZATION_CODE GrantType = "AUTHORIZATION_CODE"
	GrantTypeREFRESH_TOKEN      GrantType = "REFRESH_TOKEN"
)
