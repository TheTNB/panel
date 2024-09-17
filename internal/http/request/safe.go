package request

type SafeUpdateSSH struct {
	Port   uint `json:"port" form:"port"`
	Status bool `json:"status" form:"status"`
}

type SafeUpdatePingStatus struct {
	Status bool `json:"status" form:"status"`
}
