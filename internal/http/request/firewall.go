package request

type FirewallStatus struct {
	Status bool `json:"status" form:"status"`
}

type FirewallRule struct {
	Type      string `json:"type"`
	Family    string `json:"family" validate:"required|in:ipv4,ipv6"`
	PortStart uint   `json:"port_start" validate:"required|min:1|max:65535"`
	PortEnd   uint   `json:"port_end" validate:"required|min:1|max:65535"`
	Protocol  string `json:"protocol" validate:"required|in:tcp,udp,tcp/udp"`
	Address   string `json:"address"`
	Strategy  string `json:"strategy" validate:"required|in:accept,drop,reject"`
	Direction string `json:"direction" validate:"required|in:in,out"`
}

type FirewallIPRule struct {
	Family    string `json:"family" validate:"required|in:ipv4,ipv6"`
	Protocol  string `json:"protocol" validate:"required|in:tcp,udp,tcp/udp"`
	Address   string `json:"address"`
	Strategy  string `json:"strategy" validate:"required|in:accept,drop,reject"`
	Direction string `json:"direction" validate:"required|in:in,out"`
}

type FirewallForward struct {
	Protocol   string `json:"protocol" validate:"required|in:tcp,udp,tcp/udp"`
	Port       uint   `json:"port" validate:"required|min:1|max:65535"`
	TargetIP   string `json:"target_ip" validate:"required"`
	TargetPort uint   `json:"target_port" validate:"required|min:1|max:65535"`
}
