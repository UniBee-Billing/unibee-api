package model

type TransitType string

const (
	FLIGHT TransitType = "FLIGHT"
	TRAIN  TransitType = "TRAIN"
	CRUISE TransitType = "CRUISE"
	COACH  TransitType = "COACH"
)
