package firewall

type FireInfo struct {
	Family   string `json:"family"` // ipv4 ipv6
	Address  string `json:"address"`
	Port     uint   `json:"port"`     // 1-65535
	Protocol string `json:"protocol"` // tcp udp tcp/udp
	Strategy string `json:"strategy"` // accept drop

	Num        string `json:"num"`
	TargetIP   string `json:"targetIP"`
	TargetPort string `json:"targetPort"` // 1-65535

	UsedStatus  string `json:"usedStatus"`
	Description string `json:"description"`
}

type Forward struct {
	Num        string `json:"num"`
	Protocol   string `json:"protocol"`
	Port       uint   `json:"port"` // 1-65535
	TargetIP   string `json:"targetIP"`
	TargetPort uint   `json:"targetPort"` // 1-65535
}
