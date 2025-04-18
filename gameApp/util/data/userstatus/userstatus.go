package userstatus

type UserStatus uint8

//go:generate stringer -type=UserStatus
const (
	Online UserStatus = iota + 1
	Offline
)
