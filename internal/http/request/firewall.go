package request

type FirewallStatus struct {
	Status bool `json:"status" form:"status"`
}

type FirewallRule struct {
	Type      string `json:"type"`
	Family    string `json:"family" validate:"required,oneof=ipv4 ipv6"`
	PortStart uint   `json:"port_start" validate:"required,gte=1,lte=65535"`
	PortEnd   uint   `json:"port_end" validate:"required,gte=1,lte=65535"`
	Protocol  string `json:"protocol" validate:"min=1,oneof=tcp udp tcp/udp"`
	Address   string `json:"address"`
	Strategy  string `json:"strategy" validate:"required,oneof=accept drop reject"`
	Direction string `json:"direction"`
}

type FirewallIPRule struct {
	Family    string `json:"family" validate:"required,oneof=ipv4 ipv6"`
	Protocol  string `json:"protocol" validate:"min=1,oneof=tcp udp tcp/udp"`
	Address   string `json:"address"`
	Strategy  string `json:"strategy" validate:"required,oneof=accept drop reject"`
	Direction string `json:"direction"`
}

type FirewallForward struct {
	Protocol   string `json:"protocol" validate:"min=1,oneof=tcp udp tcp/udp"`
	Port       uint   `json:"port" validate:"required,gte=1,lte=65535"`
	TargetIP   string `json:"target_ip" validate:"required"`
	TargetPort uint   `json:"target_port" validate:"required,gte=1,lte=65535"`
}
