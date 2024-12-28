package request

type SafeUpdateSSH struct {
	Port   uint `json:"port" form:"port"  validate:"required|min:1|max:65535"`
	Status bool `json:"status" form:"status"`
}

type SafeUpdatePingStatus struct {
	Status bool `json:"status" form:"status"`
}
