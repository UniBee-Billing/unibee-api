package model

type ChallengeActionForm struct {
	ChallengeType        ChallengeType              `json:"challengeType,omitempty"`
	ChallengeRenderValue string                     `json:"challengeRenderValue,omitempty"`
	TriggerSource        ChallengeTriggerSourceType `json:"triggerSource,omitempty"`
	ExtendInfo           string                     `json:"extendInfo,omitempty"`
}

type ChallengeType string

const (
	SMS_OTP           ChallengeType = "SMS_OTP"
	PLAINTEXT_CARD_NO ChallengeType = "PLAINTEXT_CARD_NO"
	CARD_EXPIRE_DATE  ChallengeType = "CARD_EXPIRE_DATE"
)

type ChallengeTriggerSourceType string

const (
	AMS     ChallengeTriggerSourceType = "AMS"
	CHANNEL ChallengeTriggerSourceType = "CHANNEL"
)
