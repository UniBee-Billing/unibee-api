package _interface

type PlanCapacity interface {
	QueryCap(accessToken string, userId string, capKey string) int
	IncreaseCap(accessToken string, userId string, capKey string, capUse int) bool
	DecreaseCap(accessToken string, userId string, capKey string, capUse int) bool
}
