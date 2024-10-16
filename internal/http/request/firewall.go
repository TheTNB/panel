package request

type FirewallStatus struct {
	Status bool `json:"status" form:"status"`
}

type FirewallRule struct {
	PortStart uint     `json:"port_start" validate:"required,gte=1,lte=65535"`
	PortEnd   uint     `json:"port_end" validate:"required,gte=1,lte=65535"`
	Protocols []string `json:"protocols" validate:"min=1,dive,oneof=tcp udp"`
	Address   string   `json:"address"`
	Strategy  string   `json:"strategy" validate:"required,oneof=accept drop"`
}
