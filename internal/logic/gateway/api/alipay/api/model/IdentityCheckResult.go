package model

type IdentityCheckResult string

const (
	// CHECK_PASSED, CHECK_NOT_PASSED;
	IdentityCheckResult_CHECK_PASSED     IdentityCheckResult = "CHECK_PASSED"
	IdentityCheckResult_CHECK_NOT_PASSED IdentityCheckResult = "CHECK_NOT_PASSED"
)

type ClearingChannel string

const (
	//CUP, NUCC, OTHER;
	ClearingChannel_CUP   ClearingChannel = "CUP"
	ClearingChannel_NUCC  ClearingChannel = "NUCC"
	ClearingChannel_OTHER ClearingChannel = "OTHER"
)
