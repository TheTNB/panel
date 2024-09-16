package biz

type SafeRepo interface {
	GetSSH() (uint, bool, error)
	UpdateSSH(port uint, status bool) error
	GetPingStatus() (bool, error)
	UpdatePingStatus(status bool) error
}
