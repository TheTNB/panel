package firewall

type FireInfo struct {
	Family    string `json:"family"`     // ipv4 ipv6
	Address   string `json:"address"`    // 源地址或目标地址
	PortStart uint   `json:"port_start"` // 1-65535
	PortEnd   uint   `json:"port_end"`   // 1-65535
	Protocol  string `json:"protocol"`   // tcp udp tcp/udp
	Strategy  string `json:"strategy"`   // accept drop reject
	Direction string `json:"direction"`  // in out 入站或出站
}

type FireForwardInfo struct {
	Address    string `json:"address"`
	Port       uint   `json:"port"`     // 1-65535
	Protocol   string `json:"protocol"` // tcp udp tcp/udp
	TargetIP   string `json:"targetIP"`
	TargetPort string `json:"targetPort"` // 1-65535
}

type Forward struct {
	Protocol   string `json:"protocol"`
	Port       uint   `json:"port"` // 1-65535
	TargetIP   string `json:"targetIP"`
	TargetPort uint   `json:"targetPort"` // 1-65535
}
