package consts

type UserStatusEnum int

const (
	UserStatusActive  = 0
	UserStatusSuspend = 2

	UserTypeIndividual   = 1
	UserTypeOrganization = 2
)

func (status UserStatusEnum) Description() string {
	switch status {
	case UserStatusActive:
		return "Active"
	case UserStatusSuspend:
		return "Suspend"
	default:
		return "Active"
	}
}

func UserStatusToEnum(status int) UserStatusEnum {
	switch status {
	case UserStatusActive:
		return UserStatusActive
	case UserStatusSuspend:
		return UserStatusSuspend
	default:
		return UserStatusActive
	}
}

type UserTypeEnum int64

func (status UserTypeEnum) Description() string {
	switch status {
	case UserTypeIndividual:
		return "Individual"
	case UserTypeOrganization:
		return "Organization"
	default:
		return "Individual"
	}
}

func UserTypeToEnum(userType int64) UserTypeEnum {
	switch userType {
	case UserTypeIndividual:
		return UserTypeIndividual
	case UserTypeOrganization:
		return UserTypeOrganization
	default:
		return UserTypeIndividual
	}
}
